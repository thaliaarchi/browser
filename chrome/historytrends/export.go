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
	"time"

	"github.com/andrewarchi/browser/chrome"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
	"golang.org/x/net/publicsuffix"
)

type Export struct {
	Date   time.Time // date of export
	Visits []ExportVisit
}

type ExportVisit struct {
	URL            string
	VisitTime      time.Time
	TransitionType chrome.TransitionType
	PageTitle      string
}

func ParseExport(filename string) (*Export, error) {
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

	var export Export
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = 8
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

func parseExportVisit(record []string) (*ExportVisit, error) {
	/*
		"Export These Results" format
		chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/export_details.html

		0: URL              visited URL
		1: Host*            hostname of visited URL
		2: Domain*          public suffix of visited URL
		3: Visit Time       visit time in milliseconds since 1970-01-01 i.e. 1384634958041.754
		4: Visit Time       visit time in local time                    i.e. 2006-01-02 15:04:05.000
		5: Day of Week      day of the week for the visit time          0 for Sunday
		6: Transition Type  how the browser navigated to the URL        i.e. link
		7: Page Title*      page title of visited URL
		* optional
	*/

	rawURL, host, domain := record[0], record[1], record[2]
	if host != "" {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		if u.Host != host {
			return nil, fmt.Errorf("%q differs from computed host %q", host, u.Host)
		}
	}
	if domain != "" {
		if d, _ := publicsuffix.PublicSuffix(rawURL); d != domain {
			return nil, fmt.Errorf("%q differs from computed domain %q", domain, d)
		}
	}

	// TODO handle local time
	timeMs, err := strconv.ParseInt(record[3], 10, 64)
	if err != nil {
		return nil, err
	}
	t1 := timefmt.FromUnixMilli(timeMs)
	t2, err := time.Parse("2006-01-02 15:04:05.000", record[4])
	if t1 != t2 {
		return nil, fmt.Errorf("inconsistent visit times: %s and %s", t1, t2)
	}

	day, err := strconv.Atoi(record[5])
	if err != nil {
		return nil, err
	}
	if d := t1.Weekday(); d != time.Weekday(day) {
		return nil, fmt.Errorf("inconsistent weekday: %s and %s", time.Weekday(day), d)
	}

	transition, err := chrome.ParseTransitionType(record[6])
	if err != nil {
		return nil, err
	}

	return &ExportVisit{
		URL:            rawURL,
		VisitTime:      t1,
		TransitionType: transition,
		PageTitle:      record[7],
	}, nil
}
