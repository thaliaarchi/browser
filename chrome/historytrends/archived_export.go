// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/andrewarchi/browser/chrome"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// ArchivedExportReader reads the browsing history visits in an
// exported_archived_history_{date}.{tsv|txt} or
// history_autobackup_{date}_{full|incremental}.{tsv|txt|zip} file.
type ArchivedExportReader struct {
	time time.Time // export time
	exportReader
}

var archivedExportPattern = regexp.MustCompile(
	`^(?:exported_archived_history_(\d{8})` +
		`|(?:history_autobackup_(\d{8})_(?:full|incremental)))\.(?:tsv|txt|zip)$`)

// OpenArchivedExport opens an archived export for reading.
func OpenArchivedExport(filename string) (*ArchivedExportReader, error) {
	r, err := newExportReader(filename, 4)
	if err != nil {
		return nil, err
	}
	matches := archivedExportPattern.FindStringSubmatch(r.filename)
	if len(matches) != 3 {
		return nil, fmt.Errorf("historytrends: filename is not an archived export: %q", r.filename)
	}
	exportTime := matches[1]
	if exportTime == "" {
		exportTime = matches[2]
	}
	t, err := time.Parse("20060102", exportTime)
	if err != nil {
		return nil, err
	}
	return &ArchivedExportReader{time: t, exportReader: *r}, nil
}

// Time returns the time of export in UTC.
func (r *ArchivedExportReader) Time() time.Time { return r.time }

// Close closes the underlying reader.
func (r *ArchivedExportReader) Close() error { return r.r.Close() }

// ReadAll reads all visits in the export.
func (r *ArchivedExportReader) ReadAll() (*Export, error) {
	var visits []Visit
	for {
		visit, err := r.Read()
		if err == io.EOF {
			return &Export{r.time, visits}, nil
		}
		if err != nil {
			return nil, err
		}
		visits = append(visits, *visit)
	}
}

// Read reads a single visit in the export.
func (r *ArchivedExportReader) Read() (*Visit, error) {
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

func parseEpochTime(msec string) (time.Time, error) {
	if msec == "" {
		return time.Time{}, nil
	}
	epoch := timefmt.Windows
	if msec[0] == 'U' {
		epoch = timefmt.Unix
		msec = msec[1:]
	}
	return timefmt.Parse(msec, timefmt.Milli, epoch)
}
