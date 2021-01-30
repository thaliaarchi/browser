// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package firefox

import (
	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// Times contains installation times in times.json.
type Times struct {
	Created  timefmt.UnixMilli `json:"created"`
	FirstUse timefmt.UnixMilli `json:"firstUse"`
}

// ParseTimes parses times.json in a Firefox profile.
func ParseTimes(filename string) (*Times, error) {
	var times Times
	if err := jsonutil.DecodeFile(filename, &times); err != nil {
		return nil, err
	}
	return &times, nil
}
