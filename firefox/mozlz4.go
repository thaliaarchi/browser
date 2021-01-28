package firefox

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

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

func UnmarshalMozLz4Json(b []byte, v interface{}) error {
	data, err := DecompressMozLz4(b)
	if err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
