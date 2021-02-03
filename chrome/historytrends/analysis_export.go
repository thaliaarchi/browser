// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/andrewarchi/browser/chrome"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
	"golang.org/x/net/publicsuffix"
)

/*
	Analysis Export ("Export These Results")

	An analysis export is a tab-delimited file with the fields listed
	below. It is created by clicking "Export These Results" on the Trends
	or Search pages.

	0: URL                  visited URL
	1: Host*                hostname of visited URL
	2: Domain*              public suffix of visited URL
	3: Visit Time (ms)      visit time in milliseconds since 1970-01-01  i.e. 1384634958041.754
	4: Visit Time (string)  visit time in local time                     i.e. 2013-11-16 14:49:18.041
	5: Day of Week          day of the week for the visit time           0 for Sunday
	6: Transition Type      how the browser navigated to the URL         i.e. link
	7: Page Title*          page title of visited URL
	* column can be blank

	Several fields are redundant: host and domain are derived from the
	URL; visit time (string) and day of week are less precise than visit
	time (ms). These fields are validated for consistency, then discarded.

	The string-formatted visit time is in local time, at the time of the
	export, so the timezone of the export is known.

	Format docs: chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/export_details.html
*/

// readAnalysisVisit reads a single visit in an analysis export.
func (r *ExportReader) readAnalysisVisit() (*Visit, error) {
	record, err := r.readRecord()
	if err != nil {
		return nil, err
	}

	if err := checkURL(record[0], record[1], record[2]); err != nil {
		return nil, r.err(err)
	}

	t, offset, err := parseTimes(record[3], record[4], record[5])
	if err != nil {
		return nil, r.err(err)
	}
	if r.record == 1 {
		// Retrieve timezone offset from first record.
		d := time.Duration(offset) * time.Second
		zone := time.FixedZone("", offset)
		r.time = r.time.Add(-d).In(zone)
		r.tz = offset
	} else if offset != r.tz {
		// Check that all visits have same timezone offset.
		return nil, r.err(fmt.Errorf("%s differs from timezone offset %s",
			time.Duration(offset)*time.Second, time.Duration(r.tz)*time.Second))
	}

	// The page transition string only contains the core type, so
	// qualifiers are lost.
	transition, err := chrome.PageTransitionFromString(record[6])
	if err != nil {
		return nil, r.err(err)
	}

	return &Visit{
		URL:        record[0],
		VisitTime:  t,
		Transition: transition,
		PageTitle:  normalizeTitle(record[7]),
	}, nil
}

func parseTimes(timeMsec, timeLocal, weekday string) (time.Time, int, error) {
	msec, err := timefmt.Parse(timeMsec, timefmt.Milli, timefmt.Unix)
	if err != nil {
		return time.Time{}, 0, err
	}
	local, err := time.Parse("2006-01-02 15:04:05.000", timeLocal)

	// timeMsec and timeLocal both represent the same time. timeMsec is in
	// UTC with sub-millisecond precision. timeLocal is in the local
	// timezone at the time of export and has truncated millisecond
	// precision.
	diff := msec.Truncate(time.Millisecond).Sub(local)
	if diff%time.Second != 0 {
		return time.Time{}, 0, fmt.Errorf("time difference is fractional: %s", diff)
	}
	offset := int(diff / time.Second)

	day, err := strconv.Atoi(weekday)
	if err != nil {
		return time.Time{}, 0, err
	}
	if d := local.Weekday(); d != time.Weekday(day) {
		return time.Time{}, 0, fmt.Errorf("inconsistent weekday: %s and %s", time.Weekday(day), d)
	}
	return msec, offset, nil
}

func checkURL(rawURL, host, domain string) error {
	if host != "" {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		if h := u.Hostname(); h != host {
			// When the URL path contains @, utils.extractHost in utils.js
			// will erroneously return the segment after the @.
			// For example, extractHost reports that the host of
			// https://web.archive.org/save/https://medium.com/@user/article
			// is "user" instead of "web.archive.org".
			if strings.IndexByte(u.Path, '@') == -1 {
				return fmt.Errorf("%q differs from computed host %q", host, h)
			}
		}
	}
	if domain != "" {
		tld1, err := publicsuffix.EffectiveTLDPlusOne(host)
		if err != nil {
			return err
		}
		if tld1 != domain {
			return fmt.Errorf("%q differs from computed eTLD+1 %q", domain, tld1)
		}
	}
	return nil
}
