// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import "fmt"

// UnknownObj represents a json object for which the full type
// information is not known. Any unmarshal of a value that is not null
// or {} will raise an error. This is to ensure no data loss until all
// types have been determined.
type UnknownObj struct{}

// UnmarshalJSON implements the json.Unmarshaler interface. Any
// unmarshal of a value that is not null will raise an error. This is to
// ensure no data loss until all types have been determined.
func (u *UnknownObj) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "{}" {
		return nil
	}
	return fmt.Errorf("jsonutil: unmarshal of unknown object type: %s", data)
}

// UnknownType represents a json value for which the full type
// information is not known. Any unmarshal of a value that is not null
// will raise an error. This is to ensure no data loss until all types
// have been determined.
type UnknownType struct{}

// UnmarshalJSON implements the json.Unmarshaler interface. Any
// unmarshal of a value that is not null will raise an error. This is to
// ensure no data loss until all types have been determined.
func (u *UnknownType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	return fmt.Errorf("jsonutil: unmarshal of unknown type: %s", data)
}
