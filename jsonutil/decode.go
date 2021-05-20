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
	"fmt"
	"io"
	"os"
	"strconv"
)

// Decode decodes the result into data, requiring fields to match
// strictly and checking for trailing text. The reader is read to
// completion, even on error, so that HTTP response bodies are properly
// closed and connections can be reused.
func Decode(r io.Reader, v interface{}) error {
	return decode(r, v, true, true)
}

// DecodeAllowUnknownFields decodes the result into data, checking for
// trailing text. The reader is read to completion, even on error, so
// that HTTP response bodies are properly closed and connections can be
// reused.
func DecodeAllowUnknownFields(r io.Reader, v interface{}) error {
	return decode(r, v, false, true)
}

// DecodeFile opens the given file and decodes the result into data,
// requiring fields to match strictly and checking for trailing text.
func DecodeFile(filename string, v interface{}) error {
	return decodeFile(filename, v, true)
}

// DecodeFileAllowUnknownFields opens the given file and decodes the
// result into data, checking for trailing text.
func DecodeFileAllowUnknownFields(filename string, v interface{}) error {
	return decodeFile(filename, v, false)
}

func decodeFile(filename string, v interface{}, strict bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return decode(f, v, strict, false)
}

func decode(r io.Reader, v interface{}, strict, readAll bool) error {
	d := json.NewDecoder(r)
	if strict {
		d.DisallowUnknownFields()
	}
	if err := d.Decode(v); err != nil {
		if readAll {
			_, _ = io.Copy(io.Discard, r)
		}
		return err
	}

	br := io.MultiReader(d.Buffered(), r)
	var buf [4096]byte
	for {
		n, err := br.Read(buf[:])
		if err == io.EOF {
			return nil
		}
		for _, b := range buf[:n] {
			if !isSpace(b) {
				if readAll {
					_, _ = io.Copy(io.Discard, r)
				}
				return fmt.Errorf("json: invalid trailing character: %q", b)
			}
		}
	}
}

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
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
