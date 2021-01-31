// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/andrewarchi/browser/chrome"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
	"golang.org/x/net/publicsuffix"
)

type AnalysisExport struct {
	Date   time.Time // date of export
	Visits []AnalysisExportVisit
}

type AnalysisExportVisit struct {
	URL            string
	VisitTime      time.Time
	TransitionType chrome.TransitionType
	PageTitle      string
}

func ParseAnalysisExport(filename string) (*AnalysisExport, error) {
	// exported_analysis_history_20060102_150405.tsv
	var r io.Reader
	switch ext := filepath.Ext(filename); ext {
	case ".tsv":
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	case ".zip":
		zr, err := zip.OpenReader(filename)
		if err != nil {
			return nil, err
		}
		defer zr.Close()
		panic("unimplemented")
	default:
		return nil, fmt.Errorf("historytrends: bad file extension: %s", ext)
	}

	var export AnalysisExport
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = 8
	cr.LazyQuotes = true
	for line := 1; ; line++ {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		visit, err := parseExportVisit(record)
		if err != nil {
			return nil, fmt.Errorf("historytrends: export on line %d: %w", line, err)
		}
		export.Visits = append(export.Visits, *visit)
	}
	return &export, nil
}

func parseExportVisit(record []string) (*AnalysisExportVisit, error) {
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

	rawURL, host, domain := record[0], record[1], record[2]
	if host != "" {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		if h := u.Hostname(); h != host {
			// When the URL path contains @, utils.extractHost in utils.js
			// will erroneously return the segment after the @.
			// For example, extractHost reports that the host of
			// https://web.archive.org/save/https://medium.com/@user/article
			// is "user" instead of "web.archive.org".
			if strings.IndexByte(u.Path, '@') == -1 {
				return nil, fmt.Errorf("%q differs from computed host %q", host, h)
			}
		}
	}
	if domain != "" {
		tld1, err := publicsuffix.EffectiveTLDPlusOne(host)
		if err != nil {
			return nil, err
		}
		if tld1 != domain {
			return nil, fmt.Errorf("%q differs from computed eTLD+1 %q", domain, tld1)
		}
	}

	timeMsec, err := timefmt.ParseMilliFrac(record[3])
	if err != nil {
		return nil, err
	}
	timeLocal, err := time.Parse("2006-01-02 15:04:05.000", record[4])
	// timeMsec and timeLocal both represent the same time. timeMsec is
	// in UTC with sub-millisecond precision. timeLocal is local, at the
	// time of export, and has truncated millisecond precision.
	diff := (timeMsec.Sub(timeLocal) / time.Millisecond) * time.Millisecond
	// TODO handle timezone
	_ = diff
	// return nil, fmt.Errorf("inconsistent visit times: %s and %s", timeMsec, timeLocal)

	day, err := strconv.Atoi(record[5])
	if err != nil {
		return nil, err
	}
	if d := timeLocal.Weekday(); d != time.Weekday(day) {
		return nil, fmt.Errorf("inconsistent weekday: %s and %s", time.Weekday(day), d)
	}

	transition, err := chrome.ParseTransitionType(record[6])
	if err != nil {
		return nil, err
	}

	return &AnalysisExportVisit{
		URL:            rawURL,
		VisitTime:      timeMsec,
		TransitionType: transition,
		PageTitle:      record[7],
	}, nil
}
