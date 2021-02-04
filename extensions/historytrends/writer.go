// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

// Writer writes a History Trends Unlimited browsing history export.
type Writer struct {
	w      *bufio.Writer
	typ    ExportType
	loc    *time.Location
	record int
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer, typ ExportType, exportTime time.Time) (*Writer, error) {
	if typ != AnalysisExport && typ != ArchivedExport {
		return nil, fmt.Errorf("historytrends: illegal export type: %s", typ)
	}
	return &Writer{
		w:   bufio.NewWriter(w),
		typ: typ,
		loc: exportTime.Location(),
	}, nil
}

// Write writes a single visit in an export.
func (w *Writer) Write(v *Visit) error {
	err := w.write(v)
	if err == io.EOF {
		return err
	}
	if err != nil {
		return fmt.Errorf("historytrends: writing record %d: %w", w.record, err)
	}
	return nil
}

func (w *Writer) write(v *Visit) error {
	w.record++
	var record []string
	var err error
	if w.typ == AnalysisExport {
		record, err = w.writeAnalysisVisit(v)
	} else {
		record = w.writeArchivedVisit(v)
	}
	if err != nil {
		return err
	}
	return w.writeRecord(record)
}

// writeRecord writes a tab-separated record, but does not quote fields,
// unlike csv.Writer.
func (w *Writer) writeRecord(record []string) error {
	for i, field := range record {
		if i != 0 {
			if err := w.w.WriteByte('\t'); err != nil {
				return err
			}
		}
		if _, err := w.w.WriteString(field); err != nil {
			return err
		}
	}
	_, err := w.w.WriteString("\r\n")
	return err
}

// WriteAll writes all visits in an export.
func (w *Writer) WriteAll(visits []Visit) error {
	for i := range visits {
		if err := w.Write(&visits[i]); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error { return w.w.Flush() }

// Write writes the export to w.
func (ex *Export) Write(w io.Writer) error {
	ew, err := NewWriter(w, ex.Type, ex.ExportTime)
	if err != nil {
		return err
	}
	return ew.WriteAll(ex.Visits)
}
