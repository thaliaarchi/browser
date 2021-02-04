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

// Reader reads a History Trends Unlimited browsing history export.
type Reader struct {
	cr       *csv.Reader
	typ      ExportType
	filename string    // filename of tsv within zip or as given
	time     time.Time // export time
	tz       int       // timezone offset in seconds (analysis exports-only)
	record   int       // index of record
}

// ReadCloser reads and closes a History Trends Unlimited browsing
// history export.
type ReadCloser struct {
	Reader
	rc io.ReadCloser
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader, exportTime time.Time) *Reader {
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.LazyQuotes = true
	return &Reader{
		cr:   cr,
		typ:  0, // detect on first record
		time: exportTime,
	}
}

// OpenReader opens a History Trends Unlimited browsing history export
// for reading.
func OpenReader(filename string) (*ReadCloser, error) {
	r, name, err := openExport(filename)
	if err != nil {
		return nil, err
	}
	// Use the filename inside of the zip, when possible, to recover the
	// original name for renamed files.
	typ, exportTime, err := ParseExportFilename(name)
	if err != nil {
		return nil, err
	}

	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.LazyQuotes = true
	rc := &ReadCloser{
		Reader: Reader{
			cr:       cr,
			typ:      typ,
			filename: filepath.Base(name),
			time:     exportTime,
		},
		rc: r,
	}
	return rc, nil
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

// filenamePattern matches allowed export filenames. More combinations
// are allowed here than actually exported. An optional suffix is also
// permitted.
//
// As of v1.6, these are the filename formats that have been used:
//
// exported_analysis_history_{date:20060102_150405}.tsv (>= v1.5.2)
// exported_analysis_history_{date:20060102}.tsv (>= v1.4.3)
// exported_analysis_history_{date:20060102}.txt (< v1.4.3)
//
// exported_archived_history_{date:20060102}.tsv (>= v1.4.3)
// exported_archived_history_{date:20060102}.txt (< v1.4.3)
//
// history_autobackup_{date:20060102}_{full|incremental}.{tsv|zip} (>= 1.4.3)
// history_autobackup_{date:20060102}_{full|incremental}.{txt|zip} (>= 1.4.1)
var filenamePattern = regexp.MustCompile(
	`^(?:exported_(analysis|archived)_history_(\d{8}(?:_\d{6})?)` +
		`|(?:history_autobackup_(\d{8}(?:_\d{6})?)_(?:full|incremental)))` +
		`(?:[^\d].*)?` + // suffix
		`\.(?:tsv|txt|zip)$`)

// ParseExportFilename extracts the type and time of export from the
// given filename. A suffix like "(1)" or "copy", for example, is
// permitted.
func ParseExportFilename(filename string) (ExportType, time.Time, error) {
	base := filepath.Base(filename)
	matches := filenamePattern.FindStringSubmatch(base)
	if len(matches) != 4 {
		return 0, time.Time{}, fmt.Errorf("historytrends: not an export: filename %q does not match pattern", base)
	}

	exportTime := matches[2]
	if exportTime == "" {
		exportTime = matches[3]
	}
	t, err := time.Parse("20060102_150405"[:len(exportTime)], exportTime)
	if err != nil {
		return 0, time.Time{}, err
	}

	typ := ArchivedExport
	if matches[1] == "analysis" {
		typ = AnalysisExport
	}
	return typ, t, nil
}

// Read reads a single visit in an export.
func (r *Reader) Read() (*Visit, error) {
	v, err := r.read()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("historytrends: record %d: %w", r.record, err)
	}
	return v, nil
}

func (r *Reader) read() (*Visit, error) {
	r.record++
	record, err := r.cr.Read()
	if err != nil {
		return nil, err
	}

	var typ ExportType
	switch len(record) {
	case 8:
		typ = AnalysisExport
	case 4, 3:
		typ = ArchivedExport
	default:
		return nil, fmt.Errorf("record with %d fields is not an export", len(record))
	}
	if r.typ == 0 { // infer export type
		r.typ = typ
	} else if r.typ != typ {
		return nil, fmt.Errorf("record with %d fields cannot be %s export", len(record), r.typ)
	}

	if typ == AnalysisExport {
		return r.readAnalysisVisit(record[0], record[1], record[2], record[3],
			record[4], record[5], record[6], record[7])
	} else if len(record) == 3 {
		return r.readArchivedVisit(record[0], record[1], record[2], "")
	}
	return r.readArchivedVisit(record[0], record[1], record[2], record[3])
}

// ReadAll reads all visits in an export.
func (r *Reader) ReadAll() (*Export, error) {
	var visits []Visit
	for {
		visit, err := r.Read()
		if err == io.EOF {
			return &Export{r.filename, r.typ, r.time, visits}, nil
		}
		if err != nil {
			return nil, err
		}
		visits = append(visits, *visit)
	}
}

// ExportTime returns the time of export. For analysis exports, the
// timezone is initially UTC, then is determined upon reading the first
// record. For archived exports, the timezone is always UTC.
func (r *Reader) ExportTime() time.Time { return r.time }

// Close closes the underlying reader.
func (r *ReadCloser) Close() error { return r.rc.Close() }

var spacePattern = regexp.MustCompile(`\p{Z}+`)

// normalizeTitle fixes incompletely normalized spaces.
func normalizeTitle(title string) string {
	// utils.formatTitle in utils.js replaces /[\t\r\n]/g, then
	// /\s\s+/g with ' ', which overlooks non-repeated Unicode spaces
	// (JavaScript \s matches Unicode spaces, unlike Go).
	return spacePattern.ReplaceAllString(title, " ")
}
