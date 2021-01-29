package jsonutil

import (
	"encoding/json"
	"os"
)

// Decode opens the given file and decodes the result into data,
// requiring fields to match strictly.
func Decode(filename string, data interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	return d.Decode(data)
}
