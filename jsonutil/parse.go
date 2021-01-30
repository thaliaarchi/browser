// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import (
	"encoding/json"
	"io"
	"os"
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
