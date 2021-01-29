package jsonutil

import "fmt"

// UnknownObj represents a json object for which the full type
// information is not known. Any unmarshal of a value that is not null
// or {} will raise an error. This is to ensure no data loss until all
// types have been determined.
type UnknownObj struct{}

// UnmarshalJSON implements the json.Unmarshaler interface. Any
// unmarshal of a value that is not null will raise an error. This is to
// ensure no data loss until all types have been determined.
func (u *UnknownObj) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "{}" {
		return nil
	}
	return fmt.Errorf("unmarshal of unknown object type: %s", data)
}

// UnknownType represents a json value for which the full type
// information is not known. Any unmarshal of a value that is not null
// will raise an error. This is to ensure no data loss until all types
// have been determined.
type UnknownType struct{}

// UnmarshalJSON implements the json.Unmarshaler interface. Any
// unmarshal of a value that is not null will raise an error. This is to
// ensure no data loss until all types have been determined.
func (u *UnknownType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	return fmt.Errorf("unmarshal of unknown type: %s", data)
}
