package jsonutil

import (
	"encoding/base64"
	"errors"
)

// Base64 is a byte slice that is formatted in json as a base64 string.
type Base64 []byte

// MarshalJSON implements the json.Marshaler interface.
func (b Base64) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte("null"), nil
	}
	n := base64.StdEncoding.EncodedLen(len(b))
	b64 := make([]byte, n+2)
	base64.StdEncoding.Encode(b64[1:n+1], b)
	b64[0] = '"'
	b64[n+1] = '"'
	return b64, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Base64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("base64 data is not a quoted string")
	}
	data = data[1 : len(data)-1]
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		return err
	}
	*b = buf[:n]
	return nil
}

func (b Base64) String() string {
	return base64.StdEncoding.EncodeToString(b)
}
