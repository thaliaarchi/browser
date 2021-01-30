// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonutil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4/v4"
)

const mozLz4Magic = 0x6d6f7a4c7a343000 // "mozLz40\x00"

func DecompressMozLz4(b []byte) ([]byte, error) {
	if len(b) < 12 {
		return nil, errors.New("mozlz4: missing header")
	}
	magic := binary.BigEndian.Uint64(b)
	if magic != mozLz4Magic {
		return nil, fmt.Errorf("mozlz4: invalid magic number: %08x", magic)
	}
	size := binary.LittleEndian.Uint32(b[8:])

	data := make([]byte, size)
	n, err := lz4.UncompressBlock(b[12:], data)
	if err != nil {
		return nil, fmt.Errorf("mozlz4: decompress: %w", err)
	}
	if n != int(size) {
		return nil, fmt.Errorf("mozlz4: header size %d and decompressed size %d differ", size, n)
	}
	return data, nil
}

func UnmarshalMozLz4(b []byte, v interface{}) error {
	data, err := DecompressMozLz4(b)
	if err != nil {
		return err
	}
	return Decode(bytes.NewReader(data), v)
}

func DecodeMozLz4(r io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := UnmarshalMozLz4(b, v); err != nil {
		return err
	}
	return nil
}

func DecodeMozLz4File(filename string, v interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := UnmarshalMozLz4(b, v); err != nil {
		return err
	}
	return nil
}
