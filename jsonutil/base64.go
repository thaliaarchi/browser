// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import "encoding/base64"

// Base64 is a byte slice that is formatted in json as a base64 string.
type Base64 []byte

// MarshalJSON implements the json.Marshaler interface.
func (b Base64) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte("null"), nil
	}
	n := base64.StdEncoding.EncodedLen(len(b))
	buf := make([]byte, n+2)
	base64.StdEncoding.Encode(buf[1:n+1], b)
	buf[0] = '"'
	buf[n+1] = '"'
	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Base64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	q, err := UnquoteSimple(data)
	if err != nil {
		return err
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(q)))
	n, err := base64.StdEncoding.Decode(buf, q)
	if err != nil {
		return err
	}
	*b = buf[:n]
	return nil
}

func (b Base64) String() string {
	return base64.StdEncoding.EncodeToString(b)
}
