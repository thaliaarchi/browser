// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import "encoding/base64"

// Base64 is a byte slice that is formatted in json as a base64 string.
type Base64 []byte

// MarshalText implements the encoding.TextMarshaler interface.
func (b Base64) MarshalText() ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(buf, b)
	return buf, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *Base64) UnmarshalText(data []byte) error {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		return err
	}
	*b = buf[:n]
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (b Base64) MarshalJSON() ([]byte, error) {
	if b == nil {
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
	return QuotedUnmarshal(data, b)
}

func (b Base64) String() string {
	return base64.StdEncoding.EncodeToString(b)
}
