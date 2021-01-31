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
func (uuid UUID) Encode(format Format) string {
	if format == Braced {
		var buf [38]byte
		buf[0] = '{'
		encode(buf[1:37], uuid)
		buf[37] = '}'
		return string(buf[:])
	}
	var buf [36]byte
	encode(buf[:], uuid)
	return string(buf[:])
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
func Decode(uuid string) (*UUID, error) {
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

// MarshalJSON implements the json.Marshaler interface.
func (uuid UUID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(uuid.Encode(Normal))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (uuid *UUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	b, err := jsonutil.UnquoteSimple(data)
	if err != nil {
		return err
	}
	u, err := Decode(string(b))
	if err != nil {
		return err
	}
	*uuid = *u
	return nil
}

func (uuid UUID) String() string {
	return uuid.Encode(Normal)
}

// FirefoxID is an ID or UUID and is used by Firefox addons. ID is
// preferred for display.
type FirefoxID struct {
	ID   string // i.e. "addon@example.com"
	UUID *UUID  // i.e. "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}"
}

// MarshalJSON implements the json.Marshaler interface.
func (id FirefoxID) MarshalJSON() ([]byte, error) {
	if id.ID != "" {
		return []byte(strconv.Quote(id.ID)), nil
	}
	if id.UUID != nil {
		return id.UUID.MarshalJSON()
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *FirefoxID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	b, err := jsonutil.UnquoteSimple(data)
	if err != nil {
		return err
	}
	var i FirefoxID
	if b[0] == '{' && b[len(b)-1] == '}' {
		uuid, err := Decode(string(b))
		if err != nil {
			return err
		}
		i.UUID = uuid
	} else {
		i.ID = string(b)
	}
	*id = i
	return nil
}

func (id FirefoxID) String() string {
	if id.ID != "" {
		return id.ID
	}
	if id.UUID != nil {
		return id.UUID.Encode(Normal)
	}
	return ""
}
