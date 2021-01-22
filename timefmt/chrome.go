package timefmt

import (
	"strconv"
	"time"
)

// FromChromeMicro constructs a time.Time from a Chrome timestamp.
// For some fields, Chrome uses a format of microseconds since
// 1 Jan 1601.
func FromChromeMicro(usec int64) time.Time {
	// http://fileformats.archiveteam.org/wiki/Chrome_bookmarks#Date_format
	sec := usec / 1e6
	nsec := (usec % 1e6) * 1000
	return time.Date(1601, 1, 1, 0, 0, int(sec), int(nsec), time.UTC)
}

// ChromeMicro handles parsing of Chrome timestamps in JSON.
type ChromeMicro struct {
	time.Time
}

func (t ChromeMicro) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("0"), nil
	}
	usec := t.Sub(time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)) / time.Microsecond
	return []byte(strconv.FormatInt(int64(usec), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a number representing a Chrome timestamp
// in microseconds since 1 Jan 1601.
func (t *ChromeMicro) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "0" {
		return nil
	}
	usec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = ChromeMicro{FromChromeMicro(usec)}
	return nil
}
