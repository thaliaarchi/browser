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
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/andrewarchi/browser/archive"
	"github.com/andrewarchi/browser/chrome"
)

// Export contains browsing history exported from History Trends
// Unlimited.
type Export struct {
	Time   time.Time // time of export (analysis: local, exported: UTC)
	Visits []Visit
}

// Visit is a page visit in browsing history.
type Visit struct {
	URL            string
	VisitTime      time.Time // UTC
	TransitionType chrome.TransitionType
	PageTitle      string
}

type exportReader struct {
	cr       *csv.Reader
	r        io.ReadCloser
	filename string // filename of tsv in zip or given filename
	record   int    // index of record
}

func newExportReader(filename string, fieldsPerRecord int) (*exportReader, error) {
	r, name, err := openExport(filename)
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = fieldsPerRecord
	cr.LazyQuotes = true
	return &exportReader{cr, r, name, 0}, nil
}

func openExport(filename string) (io.ReadCloser, string, error) {
	switch ext := filepath.Ext(filename); ext {
	case ".tsv", ".txt":
		f, err := os.Open(filename)
		return f, filename, err
	case ".zip":
		return archive.OpenSingleFileZip(filename)
	default:
		return nil, "", fmt.Errorf("historytrends: bad file extension: %q", ext)
	}
}

func (r *exportReader) readRecord() ([]string, error) {
	r.record++
	return r.cr.Read()
}

func (r *exportReader) Close() error {
	return r.r.Close()
}

func (r *exportReader) err(err error) error {
	return fmt.Errorf("historytrends: record %d: %w", r.record, err)
}

var spacePattern = regexp.MustCompile(`\p{Z}+`)

// normalizeTitle fixes incompletely normalized spaces.
func normalizeTitle(title string) string {
	// utils.formatTitle in utils.js replaces /[\t\r\n]/g, then
	// /\s\s+/g with ' ', which overlooks non-repeated Unicode spaces
	// (JavaScript \s matches Unicode spaces, unlike Go).
	return spacePattern.ReplaceAllString(title, " ")
}
