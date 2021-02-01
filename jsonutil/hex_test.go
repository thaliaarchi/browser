// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"
)

type encDecTest struct {
	enc string
	dec Hex
}

var marshalTextTests = []encDecTest{
	{"", Hex{}},
	{"0001020304050607", Hex{0, 1, 2, 3, 4, 5, 6, 7}},
	{"08090a0b0c0d0e0f", Hex{8, 9, 10, 11, 12, 13, 14, 15}},
	{"f0f1f2f3f4f5f6f7", Hex{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"f8f9fafbfcfdfeff", Hex{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
	{"67", Hex{'g'}},
	{"e3a1", Hex{0xe3, 0xa1}},
}

var unmarshalTextTests []encDecTest
var marshalJSONTests []encDecTest
var unmarshalJSONTests []encDecTest

func init() {
	for _, test := range marshalTextTests {
		test.enc = strconv.Quote(test.enc)
		marshalJSONTests = append(marshalJSONTests, test)
	}
	marshalJSONTests = append(marshalJSONTests, encDecTest{"null", nil})

	unmarshalTextTests = append(marshalTextTests, encDecTest{"F8F9FAFBFCFDFEFF", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}})
	unmarshalJSONTests = append(marshalJSONTests, encDecTest{`"F8F9FAFBFCFDFEFF"`, []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}})
}

func TestHexMarshalText(t *testing.T) {
	for i, test := range marshalTextTests {
		enc, err := test.dec.MarshalText()
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		if string(enc) != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, enc, test.enc)
		}
	}
}

func TestHexUnmarshalText(t *testing.T) {
	for i, test := range unmarshalTextTests {
		var dec Hex
		if err := dec.UnmarshalText([]byte(test.enc)); err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		if !bytes.Equal(dec, test.dec) {
			t.Errorf("#%d: got: %#v want: %#v", i, dec, test.dec)
		}
	}
}

func TestHexMarshalJSON(t *testing.T) {
	for i, test := range marshalJSONTests {
		enc, err := json.Marshal(test.dec)
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		if string(enc) != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, enc, test.enc)
		}
	}
}

func TestHexUnmarshalJSON(t *testing.T) {
	for i, test := range unmarshalJSONTests {
		var dec Hex
		if err := json.Unmarshal([]byte(test.enc), &dec); err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		if !bytes.Equal(dec, test.dec) {
			t.Errorf("#%d: got: %#v want: %#v", i, dec, test.dec)
		}
	}
}

func TestHexString(t *testing.T) {
	for i, test := range marshalTextTests {
		if enc := test.dec.String(); enc != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, enc, test.enc)
		}
	}
}
