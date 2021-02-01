// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package timefmt

import (
	"encoding"
	"time"
)

// UnixSec is a time that is formatted as an integer representing
// a Unix time in seconds since 1970-01-01 00:00:00 UTC.
type UnixSec struct{ time.Time }

// UnixMilli is a time that is formatted as an integer representing
// a Unix time in milliseconds since 1970-01-01 00:00:00 UTC.
type UnixMilli struct{ time.Time }

// UnixMicro is a time that is formatted as an integer representing
// a Unix time in microseconds since 1970-01-01 00:00:00 UTC.
type UnixMicro struct{ time.Time }

// UnixNano is a time that is formatted as an integer representing
// a Unix time in nanoseconds since 1970-01-01 00:00:00 UTC.
type UnixNano struct{ time.Time }

// WindowsSec is a time that is formatted as an integer representing
// a Windows time in seconds since 1601-01-01 00:00:00 UTC.
type WindowsSec struct{ time.Time }

// WindowsMilli is a time that is formatted as an integer representing
// a Windows time in milliseconds since 1601-01-01 00:00:00 UTC.
type WindowsMilli struct{ time.Time }

// WindowsMicro is a time that is formatted as an integer representing
// a Windows time in microseconds since 1601-01-01 00:00:00 UTC.
type WindowsMicro struct{ time.Time }

// WindowsNano is a time that is formatted as an integer representing
// a Windows time in nanoseconds since 1601-01-01 00:00:00 UTC.
type WindowsNano struct{ time.Time }

// Epoch returns the epoch that times are relative to. Always Unix.
func (t UnixSec) Epoch() Epoch { return Unix }

// Epoch returns the epoch that times are relative to. Always Unix.
func (t UnixMilli) Epoch() Epoch { return Unix }

// Epoch returns the epoch that times are relative to. Always Unix.
func (t UnixMicro) Epoch() Epoch { return Unix }

// Epoch returns the epoch that times are relative to. Always Unix.
func (t UnixNano) Epoch() Epoch { return Unix }

// Epoch returns the epoch that times are relative to. Always Windows.
func (t WindowsSec) Epoch() Epoch { return Windows }

// Epoch returns the epoch that times are relative to. Always Windows.
func (t WindowsMilli) Epoch() Epoch { return Windows }

// Epoch returns the epoch that times are relative to. Always Windows.
func (t WindowsMicro) Epoch() Epoch { return Windows }

// Epoch returns the epoch that times are relative to. Always Windows.
func (t WindowsNano) Epoch() Epoch { return Windows }

// Unit returns the unit that times are measured in. Always Sec.
func (t UnixSec) Unit() Unit { return Sec }

// Unit returns the unit that times are measured in. Always Milli.
func (t UnixMilli) Unit() Unit { return Milli }

// Unit returns the unit that times are measured in. Always Micro.
func (t UnixMicro) Unit() Unit { return Micro }

// Unit returns the unit that times are measured in. Always Nano.
func (t UnixNano) Unit() Unit { return Nano }

// Unit returns the unit that times are measured in. Always Sec.
func (t WindowsSec) Unit() Unit { return Sec }

// Unit returns the unit that times are measured in. Always Milli.
func (t WindowsMilli) Unit() Unit { return Milli }

// Unit returns the unit that times are measured in. Always Micro.
func (t WindowsMicro) Unit() Unit { return Micro }

// Unit returns the unit that times are measured in. Always Nano.
func (t WindowsNano) Unit() Unit { return Nano }

// MarshalText implements the text.Marshaler interface.
func (t UnixSec) MarshalText() ([]byte, error) { return Format(t.Time, Sec, Unix), nil }

// MarshalText implements the text.Marshaler interface.
func (t UnixMilli) MarshalText() ([]byte, error) { return Format(t.Time, Milli, Unix), nil }

// MarshalText implements the text.Marshaler interface.
func (t UnixMicro) MarshalText() ([]byte, error) { return Format(t.Time, Micro, Unix), nil }

// MarshalText implements the text.Marshaler interface.
func (t UnixNano) MarshalText() ([]byte, error) { return Format(t.Time, Nano, Unix), nil }

// MarshalText implements the text.Marshaler interface.
func (t WindowsSec) MarshalText() ([]byte, error) { return Format(t.Time, Sec, Windows), nil }

// MarshalText implements the text.Marshaler interface.
func (t WindowsMilli) MarshalText() ([]byte, error) { return Format(t.Time, Milli, Windows), nil }

// MarshalText implements the text.Marshaler interface.
func (t WindowsMicro) MarshalText() ([]byte, error) { return Format(t.Time, Micro, Windows), nil }

// MarshalText implements the text.Marshaler interface.
func (t WindowsNano) MarshalText() ([]byte, error) { return Format(t.Time, Nano, Windows), nil }

// MarshalJSON implements the json.Marshaler interface.
func (t UnixSec) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t UnixMilli) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t UnixMicro) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t UnixNano) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t WindowsSec) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t WindowsMilli) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t WindowsMicro) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// MarshalJSON implements the json.Marshaler interface.
func (t WindowsNano) MarshalJSON() ([]byte, error) { return t.MarshalText() }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixSec) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixMilli) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixMicro) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixNano) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *WindowsSec) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *WindowsMilli) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *WindowsMicro) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *WindowsNano) UnmarshalJSON(data []byte) error { return unmarshalJSON(t, data) }

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *UnixSec) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Sec, Unix)
	if err != nil {
		return err
	}
	*t = UnixSec{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *UnixMilli) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Milli, Unix)
	if err != nil {
		return err
	}
	*t = UnixMilli{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *UnixMicro) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Micro, Unix)
	if err != nil {
		return err
	}
	*t = UnixMicro{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *UnixNano) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Nano, Unix)
	if err != nil {
		return err
	}
	*t = UnixNano{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *WindowsSec) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Sec, Windows)
	if err != nil {
		return err
	}
	*t = WindowsSec{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *WindowsMilli) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Milli, Windows)
	if err != nil {
		return err
	}
	*t = WindowsMilli{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *WindowsMicro) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Micro, Windows)
	if err != nil {
		return err
	}
	*t = WindowsMicro{t0}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *WindowsNano) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Nano, Windows)
	if err != nil {
		return err
	}
	*t = WindowsNano{t0}
	return nil
}

func unmarshalJSON(v encoding.TextUnmarshaler, data []byte) error {
	if string(data) == "null" {
		return nil
	}
	return v.UnmarshalText(data)
}
