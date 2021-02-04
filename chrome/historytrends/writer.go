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
	"time"
)

// Writer writes a History Trends Unlimited browsing history export.
type Writer struct {
	cw  *csv.Writer
	typ ExportType
	loc *time.Location
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer, typ ExportType, exportTime time.Time) (*Writer, error) {
	if typ != AnalysisExport && typ != ArchivedExport {
		return nil, fmt.Errorf("historytrends: illegal export type: %s", typ)
	}
	return &Writer{
		cw:  csv.NewWriter(w),
		typ: typ,
		loc: exportTime.Location(),
	}, nil
}

// Write writes a single visit in an export.
func (w *Writer) Write(v *Visit) error {
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
	return w.cw.Write(record)
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

// Write writes the export to w.
func (ex *Export) Write(w io.Writer) error {
	ew, err := NewWriter(w, ex.Type, ex.ExportTime)
	if err != nil {
		return err
	}
	return ew.WriteAll(ex.Visits)
}
