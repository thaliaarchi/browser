// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package uuid

import "testing"

var uuid = UUID{
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

var tests = []struct {
	uuid   UUID
	str    string
	format Format
}{
	{uuid, "01234567-89ab-cdef-0123-456789abcdef", Normal},
	{uuid, "{01234567-89ab-cdef-0123-456789abcdef}", Braced},
}

func TestEncode(t *testing.T) {
	for _, test := range tests {
		if s := test.uuid.Encode(test.format); string(s) != test.str {
			t.Errorf("%q.Encode(%d) = %q, want %q", test.uuid, test.format, s, test.str)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, test := range tests {
		uuid, err := Decode([]byte(test.str))
		if err != nil {
			t.Error(err)
			continue
		}
		if uuid == nil || *uuid != test.uuid {
			t.Errorf("Decode(%q) = %q, want %q", test.str, uuid, test.uuid)
		}
	}
}
