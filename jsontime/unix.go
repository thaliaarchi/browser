package jsontime

import (
	"strconv"
	"time"
)

// UnixSec is a time formatted as a Unix timestamp.
type UnixSec struct{ time.Time }

// MarshalJSON implements the json.Marshaler interface.
// The time is a number representing a Unix timestamp.
func (t UnixSec) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a number representing a Unix timestamp.
func (t *UnixSec) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "0" {
		return nil
	}
	sec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixSec{time.Unix(sec, 0)}
	return nil
}

type UnixMicro struct{ time.Time }

func (t UnixMicro) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.UnixNano()/1000, 10)), nil
}

func (t *UnixMicro) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "0" {
		return nil
	}
	usec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixMicro{time.Unix(usec/1e6, (usec%1e6)*1e3)}
	return nil
}
