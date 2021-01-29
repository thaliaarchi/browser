package timefmt

import (
	"strconv"
	"time"
)

// Chrome represents times internally as microseconds (s/1,000,000)
// since the Windows epoch (1601-01-01 00:00:00 UTC) using the Time
// class found at base/time/time.h. Values for Time may skew and jump
// around as the operating system makes adjustments to synchronize.
//
// https://source.chromium.org/chromium/chromium/src/+/master:base/time/time.h

// FromChrome returns the time corresponding to the given internal
// Chrome time in microseconds since 1601-01-01 00:00:00 UTC.
func FromChrome(usec int64) time.Time {
	if usec == 0 {
		return time.Time{}
	}
	sec := usec / 1e6
	nsec := (usec % 1e6) * 1000
	return time.Date(1601, 1, 1, 0, 0, int(sec), int(nsec), time.UTC)
}

var chromeEpoch = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)

// ToChrome returns the internal Chrome time in microseconds since
// 1601-01-01 00:00:00 UTC corresponding to the given time.
func ToChrome(t time.Time) (usec int64) {
	if t.IsZero() {
		return 0
	}
	return int64(t.Sub(chromeEpoch) / time.Microsecond)
}

// Chrome is a time that is formatted in json as an integer representing
// an internal Chrome time in microseconds since
// 1601-01-01 00:00:00 UTC.
type Chrome struct{ time.Time }

// MarshalJSON implements the json.Marshaler interface.
func (t Chrome) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ToChrome(t.Time), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Chrome) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	usec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Chrome{FromChrome(usec)}
	return nil
}

// QuotedChrome is a time that is formatted in json as a quoted integer
// representing an internal Chrome time in microseconds since
// 1601-01-01 00:00:00 UTC.
type QuotedChrome struct{ time.Time }

// MarshalJSON implements the json.Marshaler interface.
func (t QuotedChrome) MarshalJSON() ([]byte, error) {
	var buf []byte
	buf = append(buf, '"')
	strconv.AppendInt(buf, ToChrome(t.Time), 10)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *QuotedChrome) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	q, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	usec, err := strconv.ParseInt(q, 10, 64)
	if err != nil {
		return err
	}
	*t = QuotedChrome{FromChrome(usec)}
	return nil
}
