// Package timefmt provides types for representing time formats in json.
//
package timefmt

import (
	"strconv"
	"time"
)

// FromUnix returns the time corresponding to the given Unix time
// in seconds since 1970-01-01 00:00:00 UTC.
func FromUnix(sec int64) time.Time {
	if sec == 0 {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}

// ToUnix returns the Unix time in seconds since 1970-01-01 00:00:00 UTC
// corresponding to the given time.
func ToUnix(t time.Time) (sec int64) {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

// FromUnixMilli returns the time corresponding to the given Unix time
// in milliseconds since 1970-01-01 00:00:00 UTC.
func FromUnixMilli(msec int64) time.Time {
	if msec == 0 {
		return time.Time{}
	}
	return time.Unix(msec/1e3, (msec%1e3)*1e6).UTC()
}

// ToUnixMilli returns Unix time in milliseconds since
// 1970-01-01 00:00:00 UTC corresponding to the given time.
func ToUnixMilli(t time.Time) (msec int64) {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1e6
}

// FromUnixMicro returns the time corresponding to the given Unix time
// in microseconds since 1970-01-01 00:00:00 UTC.
func FromUnixMicro(usec int64) time.Time {
	if usec == 0 {
		return time.Time{}
	}
	return time.Unix(usec/1e6, (usec%1e6)*1e3).UTC()
}

// ToUnixMicro returns Unix time in microseconds since
// 1970-01-01 00:00:00 UTC corresponding to the given time.
func ToUnixMicro(t time.Time) (msec int64) {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1e3
}

// Unix is a time that is formatted in json as an integer representing a
// Unix time in seconds since 1970-01-01 00:00:00 UTC.
type Unix struct{ time.Time }

// MarshalJSON implements the json.Marshaler interface.
func (t Unix) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ToUnix(t.Time), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Unix) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	sec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Unix{FromUnix(sec)}
	return nil
}

// UnixMilli is a time that is formatted in json as an integer
// representing a Unix time in milliseconds since
// 1970-01-01 00:00:00 UTC.
type UnixMilli struct{ time.Time }

// MarshalJSON implements the json.Marshaler interface.
func (t UnixMilli) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ToUnixMilli(t.Time), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixMilli) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	msec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixMilli{FromUnixMilli(msec)}
	return nil
}

// UnixMicro is a time that is formatted in json as an integer
// representing a Unix time in microseconds since
// 1970-01-01 00:00:00 UTC.
type UnixMicro struct{ time.Time }

// MarshalJSON implements the json.Marshaler interface.
func (t UnixMicro) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ToUnixMicro(t.Time), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixMicro) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	usec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixMicro{FromUnixMicro(usec)}
	return nil
}
