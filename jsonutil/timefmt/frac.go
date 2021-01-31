// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package timefmt

import (
	"strconv"
	"strings"
	"time"
)

func ParseMilliFrac(s string) (time.Time, error) {
	msec, nsec, err := splitFrac(s, 6)
	if err != nil {
		return time.Time{}, err
	}
	sec := msec / 1e3
	nsec += (msec % 1e3) * 1e6
	return time.Unix(sec, nsec).UTC(), nil
}

type MilliFrac struct{ time.Time }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *MilliFrac) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	m, err := ParseMilliFrac(string(data))
	if err != nil {
		return err
	}
	*t = MilliFrac{m}
	return nil
}

func splitFrac(num string, places int) (msec, nsec int64, err error) {
	if i := strings.IndexByte(num, '.'); i != -1 {
		frac := (num[i+1:] + "000000000")[:places]
		nsec, err = strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return
		}
		num = num[:i]
	}
	msec, err = strconv.ParseInt(num, 10, 64)
	return
}
