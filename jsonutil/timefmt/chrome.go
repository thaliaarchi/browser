// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package timefmt

import (
	"time"

	"github.com/andrewarchi/browser/jsonutil"
)

// Chrome represents times internally as microseconds (s/1,000,000)
// since the Windows epoch (1601-01-01 00:00:00 UTC) using the Time
// class found at base/time/time.h. Values for Time may skew and jump
// around as the operating system makes adjustments to synchronize.
//
// https://source.chromium.org/chromium/chromium/src/+/master:base/time/time.h

// Chrome is a time that is formatted in json as an integer representing
// an internal Chrome time in microseconds since
// 1601-01-01 00:00:00 UTC.
type Chrome = WindowsMicro

// QuotedChrome is a time that is formatted in json as a quoted integer
// representing an internal Chrome time in microseconds since
// 1601-01-01 00:00:00 UTC.
type QuotedChrome struct{ time.Time }

// MarshalText implements the text.Marshaler interface.
func (t QuotedChrome) MarshalText() ([]byte, error) {
	return FormatBytes(t.Time, Micro, Windows), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t QuotedChrome) MarshalJSON() ([]byte, error) {
	var buf []byte
	buf = append(buf, '"')
	Append(buf, t.Time, Micro, Windows)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *QuotedChrome) UnmarshalText(data []byte) error {
	t0, err := Parse(string(data), Milli, Windows)
	if err != nil {
		return err
	}
	*t = QuotedChrome{t0}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *QuotedChrome) UnmarshalJSON(data []byte) error {
	return jsonutil.QuotedUnmarshal(data, t.UnmarshalText)
}
