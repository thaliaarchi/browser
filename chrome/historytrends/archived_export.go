// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"strconv"
	"time"

	"github.com/andrewarchi/browser/chrome"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

/*
	Archived Export ("Transfer History" / "Auto Backup")

	An archived export is a tab-delimited file with the fields listed
	below. It is created by using the "Transfer History" or "Auto Backup"
	features on the Options page.

	0: URL              visited URL
	1: Visit Time       visit time in milliseconds            i.e. U1384634958041.754 (Unix epoch: 1970-01-01)
	2: Transition Type  how the browser navigated to the URL  i.e. 1
	3: Page Title*      page title of the visited URL
	* column can be blank

	Format docs: chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/export_details.html
*/

// readArchivedVisit reads a single visit in an archived export.
func (r *ExportReader) readArchivedVisit() (*Visit, error) {
	record, err := r.readRecord()
	if err != nil {
		return nil, err
	}
	t, err := parseEpochTime(record[1])
	if err != nil {
		return nil, err
	}
	transition, err := strconv.ParseUint(record[2], 10, 32)
	if err != nil {
		return nil, err
	}
	return &Visit{
		URL:        record[0],
		VisitTime:  t,
		Transition: chrome.PageTransition(transition),
		PageTitle:  normalizeTitle(record[3]),
	}, nil
}

// parseEpochTime parses a fractional millisecond timestamp relative to
// the Unix or Windows epoch.
func parseEpochTime(msec string) (time.Time, error) {
	if msec == "" {
		return time.Time{}, nil
	}
	epoch := timefmt.Windows
	if msec[0] == 'U' { // >= v1.4.1
		epoch = timefmt.Unix
		msec = msec[1:]
	}
	return timefmt.Parse(msec, timefmt.Milli, epoch)
}
