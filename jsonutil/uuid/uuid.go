// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package uuid encodes and decodes UUIDs.
package uuid

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/andrewarchi/browser/jsonutil"
)

// Format is a UUID formatting style.
type Format uint8

// UUID formats:
const (
	Normal Format = iota // "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	Braced               // "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}"
)

// UUID is a UUID that is formatted as
// "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" or
// "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}".
type UUID [16]byte

// Encode encodes a UUID in the given format.
func (uuid *UUID) Encode(format Format) []byte {
	if uuid == nil {
		return nil
	}
	if format == Braced {
		var buf [38]byte
		buf[0] = '{'
		encode(buf[1:37], *uuid)
		buf[37] = '}'
		return buf[:]
	}
	var buf [36]byte
	encode(buf[:], *uuid)
	return buf[:]
}

func encode(dst []byte, uuid [16]byte) {
	hex.Encode(dst[:8], uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], uuid[10:16])
}

// Decode decodes a UUID.
func Decode(uuid []byte) (*UUID, error) {
	if len(uuid) == 38 {
		if uuid[0] != '{' || uuid[37] != '}' {
			return nil, fmt.Errorf("uuid: invalid braced UUID: %q", uuid)
		}
		uuid = uuid[1:37]
	}
	if len(uuid) != 36 || uuid[8] != '-' ||
		uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return nil, fmt.Errorf("uuid: invalid UUID: %q", uuid)
	}
	var dst UUID
	if _, err := hex.Decode(dst[:4], []byte(uuid[:8])); err != nil {
		return nil, err
	}
	if _, err := hex.Decode(dst[4:6], []byte(uuid[9:13])); err != nil {
		return nil, err
	}
	if _, err := hex.Decode(dst[6:8], []byte(uuid[14:18])); err != nil {
		return nil, err
	}
	if _, err := hex.Decode(dst[8:10], []byte(uuid[19:23])); err != nil {
		return nil, err
	}
	if _, err := hex.Decode(dst[10:16], []byte(uuid[24:36])); err != nil {
		return nil, err
	}
	return &dst, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (uuid *UUID) MarshalText() ([]byte, error) {
	return uuid.Encode(Normal), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (uuid *UUID) UnmarshalText(data []byte) error {
	u, err := Decode(data)
	if err != nil {
		return err
	}
	*uuid = *u
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (uuid *UUID) MarshalJSON() ([]byte, error) {
	if uuid == nil {
		return []byte("null"), nil
	}
	return []byte(strconv.Quote(string(uuid.Encode(Normal)))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (uuid *UUID) UnmarshalJSON(data []byte) error {
	return jsonutil.QuotedUnmarshal(data, uuid.UnmarshalText)
}

func (uuid *UUID) String() string {
	return string(uuid.Encode(Normal))
}
