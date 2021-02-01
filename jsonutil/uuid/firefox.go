// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package uuid

import (
	"strconv"

	"github.com/andrewarchi/browser/jsonutil"
)

// Firefox is an ID or UUID and is used by Firefox addons. ID is
// preferred for display.
type Firefox struct {
	ID   string // i.e. "addon@example.com"
	UUID *UUID  // i.e. "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}"
}

// MarshalText implements the encoding.TextMarshaler interface.
func (id Firefox) MarshalText() ([]byte, error) {
	if id.ID != "" {
		return []byte(id.ID), nil
	}
	if id.UUID != nil {
		return id.UUID.MarshalText()
	}
	return nil, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *Firefox) UnmarshalText(data []byte) error {
	var i Firefox
	if data[0] == '{' && data[len(data)-1] == '}' {
		uuid, err := Decode(data)
		if err != nil {
			return err
		}
		i.UUID = uuid
	} else {
		i.ID = string(data)
	}
	*id = i
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (id Firefox) MarshalJSON() ([]byte, error) {
	if id.ID != "" {
		return []byte(strconv.Quote(id.ID)), nil
	}
	if id.UUID != nil {
		return id.UUID.MarshalJSON()
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *Firefox) UnmarshalJSON(data []byte) error {
	return jsonutil.QuotedUnmarshal(data, id.UnmarshalText)
}

func (id Firefox) String() string {
	if id.ID != "" {
		return id.ID
	}
	if id.UUID != nil {
		return string(id.UUID.Encode(Braced))
	}
	return ""
}
