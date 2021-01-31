// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package uuid

import "testing"

var guid = [16]byte{
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

const guidStr = "01234567-89ab-cdef-0123-456789abcdef"
const bracedGUIDStr = "{01234567-89ab-cdef-0123-456789abcdef}"

func TestEncodeGUID(t *testing.T) {
	if s := EncodeGUID(guid); s != guidStr {
		t.Errorf("EncodeGUID: got %q, want %q", s, guidStr)
	}
}

func TestEncodeBracedGUID(t *testing.T) {
	if s := EncodeBracedGUID(guid); s != bracedGUIDStr {
		t.Errorf("EncodeBracedGUID: got %q, want %q", s, bracedGUIDStr)
	}
}
