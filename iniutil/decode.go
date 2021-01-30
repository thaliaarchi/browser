// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package iniutil provides utilities for parsing INI files.
package iniutil

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/ini.v1"
)

// DecodeINIStrict decodes an INI section into a struct and checks for
// unknown fields.
func DecodeINIStrict(section *ini.Section, v interface{}) error {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr {
		return errors.New("ini: not a pointer to a struct")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		return errors.New("ini: not a pointer to a struct")
	}

	for _, key := range section.KeyStrings() {
		f, ok := typ.FieldByName(key)
		if !ok || f.Tag.Get("ini") == "-" {
			return fmt.Errorf("ini: section %s has unknown key: %s", section.Name(), key)
		}
	}
	return section.StrictMapTo(v)
}
