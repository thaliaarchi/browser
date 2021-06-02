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
	1: Visit Time       visit time in UTC                     e.g. U1384634958041.754 (Unix milliseconds), 13149893660345543 (Windows microseconds)
	2: Transition Type  how the browser navigated to the URL  e.g. 1
	3: Page Title*      page title of the visited URL
	* column can be blank

	Format docs: chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/export_details.html
*/

// readArchivedVisit reads a single visit in an archived export.
func (r *Reader) readArchivedVisit(rawURL, timeMsec, transition, title string) (*Visit, error) {
	t, err := parseEpochTime(timeMsec)
	if err != nil {
		return nil, err
	}
	typ, err := strconv.ParseInt(transition, 10, 32)
	if err != nil {
		return nil, err
	}
	return &Visit{
		URL:        rawURL,
		VisitTime:  t,
		Transition: chrome.PageTransition(typ),
		PageTitle:  normalizeTitle(title),
	}, nil
}

// parseEpochTime parses a fractional millisecond timestamp relative to
// the Unix or Windows epoch.
func parseEpochTime(t string) (time.Time, error) {
	if t == "" {
		return time.Time{}, nil
	}
	if t[0] == 'U' { // >= v1.4.1
		return timefmt.Parse(t[1:], timefmt.Milli, timefmt.Unix)
	}
	return timefmt.Parse(t, timefmt.Micro, timefmt.Windows)
}

// writeArchivedVisit writes a single visit in an archived export.
func (w *Writer) writeArchivedVisit(v *Visit) []string {
	return []string{
		v.URL,
		"U" + timefmt.Format(v.VisitTime, timefmt.Milli, timefmt.Unix),
		strconv.Itoa(int(v.Transition)),
		v.PageTitle,
	}
}
