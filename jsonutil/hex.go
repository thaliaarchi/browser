// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import (
	"encoding/hex"
	"strconv"
)

// Hex is a byte slice that is formatted in json as a hexadecimal
// string.
type Hex []byte

// MarshalJSON implements the json.Marshaler interface.
func (h Hex) MarshalJSON() ([]byte, error) {
	if len(h) == 0 {
		return []byte("null"), nil
	}
	n := hex.EncodedLen(len(h))
	buf := make([]byte, n+2)
	hex.Encode(buf[1:n+1], h)
	buf[0] = '"'
	buf[n+1] = '"'
	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *Hex) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	q, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	buf, err := hex.DecodeString(q)
	if err != nil {
		return err
	}
	*h = buf
	return nil
}

func (h Hex) String() string {
	return hex.EncodeToString(h)
}
