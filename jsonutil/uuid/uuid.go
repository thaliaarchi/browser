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

// EncodeGUID formats a GUID as "01234567-89ab-cdef-0123-456789abcdef".
func EncodeGUID(guid [16]byte) string {
	var buf [36]byte
	encodeGUID(buf[:], guid)
	return string(buf[:])
}

// EncodeBracedGUID formats a GUID as
// "{01234567-89ab-cdef-0123-456789abcdef}".
func EncodeBracedGUID(guid [16]byte) string {
	var buf [38]byte
	buf[0] = '{'
	encodeGUID(buf[1:37], guid)
	buf[37] = '}'
	return string(buf[:])
}

func encodeGUID(dst []byte, guid [16]byte) {
	hex.Encode(dst[:8], guid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], guid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], guid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], guid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], guid[10:16])
}

// DecodeGUID decodes a GUID formatted as
// "01234567-89ab-cdef-0123-456789abcdef".
func DecodeGUID(guid string) (*[16]byte, error) {
	var buf [16]byte
	if err := decodeGUID(buf[:], guid); err != nil {
		return nil, err
	}
	return &buf, nil
}

// DecodeBracedGUID decodes a GUID formatted as
// "{01234567-89ab-cdef-0123-456789abcdef}".
func DecodeBracedGUID(guid string) (*[16]byte, error) {
	if len(guid) != 38 || guid[0] != '{' || guid[37] != '}' {
		return nil, fmt.Errorf("jsonutil: invalid braced GUID: %q", guid)
	}
	var buf [16]byte
	if err := decodeGUID(buf[:], guid[1:37]); err != nil {
		return nil, err
	}
	return &buf, nil
}

func decodeGUID(dst []byte, guid string) error {
	if len(guid) != 36 || guid[8] != '-' ||
		guid[13] != '-' || guid[18] != '-' || guid[23] != '-' {
		return fmt.Errorf("jsonutil: invalid GUID: %q", guid)
	}
	if _, err := hex.Decode(dst[:4], []byte(guid[:8])); err != nil {
		return err
	}
	if _, err := hex.Decode(dst[4:6], []byte(guid[9:13])); err != nil {
		return err
	}
	if _, err := hex.Decode(dst[6:8], []byte(guid[14:18])); err != nil {
		return err
	}
	if _, err := hex.Decode(dst[8:10], []byte(guid[19:23])); err != nil {
		return err
	}
	if _, err := hex.Decode(dst[10:16], []byte(guid[24:36])); err != nil {
		return err
	}
	return nil
}

// GUID is a GUID that is formatted as
// "01234567-89ab-cdef-0123-456789abcdef".
type GUID [16]byte

// MarshalJSON implements the json.Marshaler interface.
func (g GUID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(EncodeGUID(g))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (g *GUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	b, err := jsonutil.UnquoteSimple(data)
	if err != nil {
		return err
	}
	guid, err := DecodeGUID(string(b))
	if err != nil {
		return err
	}
	*g = *guid
	return nil
}

func (g GUID) String() string {
	return EncodeGUID(g)
}

// BracedGUID is a GUID that is formatted as
// "{01234567-89ab-cdef-0123-456789abcdef}".
type BracedGUID GUID

// MarshalJSON implements the json.Marshaler interface.
func (g BracedGUID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(EncodeBracedGUID(g))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (g *BracedGUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	b, err := jsonutil.UnquoteSimple(data)
	if err != nil {
		return err
	}
	guid, err := DecodeBracedGUID(string(b))
	if err != nil {
		return err
	}
	*g = *guid
	return nil
}

func (g BracedGUID) String() string {
	return EncodeBracedGUID(g)
}

// FirefoxID is an ID or GUID and is used by Firefox addons. ID is
// preferred for display.
type FirefoxID struct {
	ID   string      // i.e. "addon@example.com"
	GUID *BracedGUID // i.e. "{01234567-89ab-cdef-0123-456789abcdef}"
}

// MarshalJSON implements the json.Marshaler interface.
func (id FirefoxID) MarshalJSON() ([]byte, error) {
	if id.ID != "" {
		return []byte(strconv.Quote(id.ID)), nil
	}
	if id.GUID != nil {
		return id.GUID.MarshalJSON()
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
	if b[0] == '{' {
		guid, err := DecodeBracedGUID(string(b))
		if err != nil {
			return err
		}
		i.GUID = (*BracedGUID)(guid)
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
	if id.GUID != nil {
		return id.GUID.String()
	}
	return ""
}
