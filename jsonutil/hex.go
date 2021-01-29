package jsonutil

import (
	"encoding/hex"
	"fmt"
)

// Hex is a byte slice that is formatted in json as a hexadecimal
// string.
type Hex []byte

// MarshalJSON implements the json.Marshaler interface.
func (h Hex) MarshalJSON() ([]byte, error) {
	if len(h) == 0 {
		return []byte("null"), nil
	}
	n := hex.EncodedLen(len(h))
	buf := make([]byte, n+2)
	hex.Encode(buf[1:n+1], h)
	buf[0] = '"'
	buf[n+1] = '"'
	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *Hex) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("hex data is not a quoted string: %q", string(data))
	}
	data = data[1 : len(data)-1]
	buf := make([]byte, hex.DecodedLen(len(data)))
	n, err := hex.Decode(buf, data)
	if err != nil {
		return err
	}
	*h = buf[:n]
	return nil
}

func (h Hex) String() string {
	return hex.EncodeToString(h)
}
