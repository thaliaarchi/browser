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

type MilliFrac struct{ time.Time }

func (t *MilliFrac) UnmarshalJSON(data []byte) error {
	msec, nsec, err := splitFrac(string(data), 6)
	if err != nil {
		return err
	}
	sec := msec / 1e3
	nsec += (msec % 1e3) * 1e6
	*t = MilliFrac{time.Unix(sec, nsec).UTC()}
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
