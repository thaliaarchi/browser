// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andrewarchi/browser/archive"
	"github.com/andrewarchi/browser/chrome"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
	"golang.org/x/net/publicsuffix"
)

/*
	"Export These Results" format
	chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/export_details.html

	0: URL                  visited URL
	1: Host*                hostname of visited URL
	2: Domain*              public suffix of visited URL
	3: Visit Time (ms)      visit time in milliseconds since 1970-01-01 i.e. 1384634958041.754
	4: Visit Time (string)  visit time in local time                    i.e. 2013-11-16 14:49:18.041
	5: Day of Week          day of the week for the visit time          0 for Sunday
	6: Transition Type      how the browser navigated to the URL        i.e. link
	7: Page Title*          page title of visited URL
	* optional
*/

type AnalysisExport struct {
	Time   time.Time // time of export
	Visits []AnalysisExportVisit
}

type AnalysisExportVisit struct {
	URL            string
	VisitTime      time.Time
	TransitionType chrome.TransitionType
	PageTitle      string
}

var analysisExportPattern = regexp.MustCompile(`^exported_analysis_history_(\d{8}_\d{6})\.(?:tsv|zip)$`)

// ParseAnalysisExport parses the history visits in an
// exported_analysis_history_{date}_{time}.tsv file.
func ParseAnalysisExport(filename string) (*AnalysisExport, error) {
	var export AnalysisExport
	r, name, err := openExport(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	base := filepath.Base(name)
	matches := analysisExportPattern.FindStringSubmatch(base)
	if len(matches) != 2 {
		return nil, fmt.Errorf("historytrends: filename is not an analysis export: %q", base)
	}
	t, err := time.Parse("20060102_150405", matches[1])
	if err != nil {
		return nil, err
	}
	export.Time = t

	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = 8
	cr.LazyQuotes = true
	var tzOffset int
	for line := 1; ; line++ {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if err := checkURL(record[0], record[1], record[2]); err != nil {
			return nil, lineErr(line, err)
		}

		t, offset, err := parseTimes(record[3], record[4], record[5])
		if err != nil {
			return nil, lineErr(line, err)
		}
		if line == 1 {
			// Attach timezone offset to export time.
			d := time.Duration(-offset) * time.Second
			zone := time.FixedZone("", offset)
			export.Time = export.Time.Add(d).In(zone)
			tzOffset = offset
		} else if offset != tzOffset {
			// Check that all visits have same timezone offset.
			return nil, lineErr(line, fmt.Errorf("%s differs from timezone offset %s",
				time.Duration(offset)*time.Second, time.Duration(tzOffset)*time.Second))
		}

		transition, err := chrome.ParseTransitionType(record[6])
		if err != nil {
			return nil, lineErr(line, err)
		}

		export.Visits = append(export.Visits, AnalysisExportVisit{
			URL:            record[0],
			VisitTime:      t,
			TransitionType: transition,
			PageTitle:      normalizeTitle(record[7]),
		})
	}
	return &export, nil
}

func openExport(filename string) (io.ReadCloser, string, error) {
	switch ext := filepath.Ext(filename); ext {
	case ".tsv":
		f, err := os.Open(filename)
		return f, filename, err
	case ".zip":
		return archive.OpenSingleFileZip(filename)
	default:
		return nil, "", fmt.Errorf("historytrends: bad file extension: %q", ext)
	}
}

func parseTimes(timeMsec, timeLocal, weekday string) (time.Time, int, error) {
	tMsec, err := timefmt.ParseMilliFrac(timeMsec)
	if err != nil {
		return time.Time{}, 0, err
	}
	tLocal, err := time.Parse("2006-01-02 15:04:05.000", timeLocal)

	// tMsec and tLocal both represent the same time. tMsec is in UTC with
	// sub-millisecond precision. tLocal is local, at the time of export,
	// and has truncated millisecond precision.
	diff := tMsec.Truncate(time.Millisecond).Sub(tLocal)
	if diff%time.Second != 0 {
		return time.Time{}, 0, fmt.Errorf("time difference is fractional: %s", diff)
	}
	offset := int(diff / time.Second)

	day, err := strconv.Atoi(weekday)
	if err != nil {
		return time.Time{}, 0, err
	}
	if d := tLocal.Weekday(); d != time.Weekday(day) {
		return time.Time{}, 0, fmt.Errorf("inconsistent weekday: %s and %s", time.Weekday(day), d)
	}
	return tMsec, offset, nil
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

var spacePattern = regexp.MustCompile(`\p{Z}+`)

// normalizeTitle fixes incompletely normalized spaces.
func normalizeTitle(title string) string {
	// utils.formatTitle in utils.js replaces /[\t\r\n]/g, then
	// /\s\s+/g with ' ', which overlooks non-repeated Unicode spaces
	// (JavaScript \s matches Unicode spaces, unlike Go).
	return spacePattern.ReplaceAllString(title, " ")
}

func lineErr(line int, err error) error {
	return fmt.Errorf("historytrends: export on line %d: %w", line, err)
}
