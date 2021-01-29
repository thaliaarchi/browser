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

// Chrome is a time that is formatted in json as an internal Chrome
// timestamp in microseconds since 1601-01-01 00:00:00 UTC.
type Chrome struct{ time.Time }

var chromeEpoch = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)

// MarshalJSON implements the json.Marshaler interface.
func (t Chrome) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("0"), nil
	}
	usec := t.Sub(chromeEpoch) / time.Microsecond
	return []byte(strconv.FormatInt(int64(usec), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Chrome) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "0" {
		return nil
	}
	usec, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Chrome{FromChrome(usec)}
	return nil
}
