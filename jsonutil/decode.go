// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package jsonutil provides utilities for parsing JSON files.
package jsonutil

import (
	"encoding"
	"encoding/json"
	"io"
	"os"
	"strconv"
)

// Decode decodes the result into data, requiring fields to match
// strictly.
func Decode(r io.Reader, data interface{}) error {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	return d.Decode(data)
}

// DecodeFile opens the given file and decodes the result into data,
// requiring fields to match strictly.
func DecodeFile(filename string, data interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return Decode(f, data)
}

// QuotedUnmarshal removes quotes then unmarshals the data. Escape
// sequences are not checked.
func QuotedUnmarshal(data []byte, v encoding.TextUnmarshaler) error {
	if string(data) == "null" {
		return nil
	}
	q, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return v.UnmarshalText([]byte(q))
}
