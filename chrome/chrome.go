// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package chrome

import (
	"os"
	"path/filepath"
	"time"
)

// GetFirstRun retrieves the time that Chrome was first ran from
// "First Run" in the Chrome root.
func GetFirstRun(chromeDir string) (time.Time, error) {
	fi, err := os.Stat(filepath.Join(chromeDir, "First Run"))
	if err != nil {
		return time.Time{}, err
	}
	return fi.ModTime(), nil
}
