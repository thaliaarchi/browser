package firefox

import (
	"encoding/binary"
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

	dst := make([]byte, size)
	n, err := lz4.UncompressBlock(b[12:], dst)
	if err != nil {
		return nil, fmt.Errorf("mozlz4: decompress: %w", err)
	}
	if n != int(size) {
		return nil, fmt.Errorf("mozlz4: header size %d and decompressed size %d differ", size, n)
	}
	return dst, nil
}
