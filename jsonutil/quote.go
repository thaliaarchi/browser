// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import "fmt"

// unquoteSimple removes quotes from a byte slice. Escape sequences are
// not checked.
func unquoteSimple(b []byte) ([]byte, error) {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return nil, fmt.Errorf("jsonutil: not a quoted string: %q", b)
	}
	return b[1 : len(b)-1], nil
}

// QuotedUnmarshal removes quotes then unmarshals the data. Escape
// sequences are not checked.
func QuotedUnmarshal(data []byte, unmarshal func([]byte) error) error {
	if string(data) == "null" {
		return nil
	}
	q, err := unquoteSimple(data)
	if err != nil {
		return err
	}
	return unmarshal(q)
}
