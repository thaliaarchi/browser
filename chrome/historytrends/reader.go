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
)

// ExportReader reads a History Trends Unlimited browsing history
// export.
type ExportReader struct {
	cr       *csv.Reader
	r        io.ReadCloser
	typ      ExportType
	filename string    // filename of tsv within zip or bare
	time     time.Time // export time
	tz       int       // timezone offset in seconds (analysis exports-only)
	record   int       // index of record
}

// filenamePattern matches allowed export filenames. More combinations
// are allowed here than actually exported. An optional suffix like
// " (1)" is also permitted.
//
// As of v1.6, these are the filename formats that have been used:
//
// exported_analysis_history_{date:20060102_150405}.tsv (>= v1.5.2)
// exported_analysis_history_{date:20060102}.tsv (< v1.5.2)
// exported_analysis_history_{date:20060102}.txt (< v1.4.3)
//
// exported_archived_history_{date:20060102}.tsv (>= v1.4.3)
// exported_archived_history_{date:20060102}.txt (< v1.4.3)
//
// history_autobackup_{date:20060102}_{full|incremental}.{tsv|zip} (>= 1.5.2)
// history_autobackup_{date:20060102}_{full|incremental}.{txt|zip} (>= 1.4.1)
var filenamePattern = regexp.MustCompile(
	`^(?:exported_(analysis|archived)_history_(\d{8}(?:_\d{6})?)` +
		`|(?:history_autobackup_(\d{8}(?:_\d{6})?)_(full|incremental)))` +
		`(?:[^\d].*)?` + // suffix
		`\.(?:tsv|txt|zip)$`)

// OpenExport opens a History Trends Unlimited browsing history export
// for reading.
func OpenExport(filename string) (*ExportReader, error) {
	r, name, err := openExport(filename)
	if err != nil {
		return nil, err
	}
	// Use the filename inside of the zip, when possible, to recover the
	// original name for renamed files.
	base := filepath.Base(name)
	matches := filenamePattern.FindStringSubmatch(base)
	if len(matches) != 5 {
		return nil, fmt.Errorf("historytrends: filename is not an export: %q", base)
	}

	exportTime := matches[2]
	if exportTime == "" {
		exportTime = matches[3]
	}
	t, err := time.Parse("20060102_150405"[:len(exportTime)], exportTime)
	if err != nil {
		return nil, err
	}

	typ := ArchivedExport
	fields := 4
	if matches[1] == "analysis" {
		typ = AnalysisExport
		fields = 8
	} else if matches[4] == "incremental" {
		typ = ArchivedExport
	}

	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = fields
	cr.LazyQuotes = true
	return &ExportReader{
		cr:       cr,
		r:        r,
		typ:      typ,
		filename: base,
		time:     t,
	}, nil
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

func (r *ExportReader) readRecord() ([]string, error) {
	r.record++
	return r.cr.Read()
}

// Time returns the time of export. For analysis exports, the timezone
// is initially UTC, then is determined upon reading the first record.
// For archived exports, the timezone is always UTC.
func (r *ExportReader) Time() time.Time { return r.time }

// Close closes the underlying reader.
func (r *ExportReader) Close() error { return r.r.Close() }

// ReadAll reads all visits in the export.
func (r *ExportReader) ReadAll() (*Export, error) {
	var visits []Visit
	for {
		var visit *Visit
		var err error
		if r.typ == AnalysisExport {
			visit, err = r.readAnalysisVisit()
		} else {
			visit, err = r.readArchivedVisit()
		}
		if err == io.EOF {
			return &Export{r.filename, r.typ, r.time, visits}, nil
		}
		if err != nil {
			return nil, err
		}
		visits = append(visits, *visit)
	}
}

func (r *ExportReader) err(err error) error {
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
