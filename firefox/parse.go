package firefox

import (
	"encoding/json"
	"fmt"
	"os"
)

type unknownObj struct{}

func (u *unknownObj) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "{}" {
		return nil
	}
	return fmt.Errorf("unmarshal of unknown object: %s", data)
}

type unknownType struct{}

func (u *unknownType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	return fmt.Errorf("unmarshal of unknown type: %s", data)
}

func parseJSON(filename string, data interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	return d.Decode(data)
}
