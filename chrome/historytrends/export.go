// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package historytrends parses Chromium browsing history exports from
// the History Trends Unlimited extension.
package historytrends

import (
	"time"

	"github.com/andrewarchi/browser/chrome"
)

// Export contains browsing history exported from History Trends
// Unlimited.
type Export struct {
	Filename string
	Type     ExportType
	Time     time.Time // time of export (analysis: local, exported: UTC)
	Visits   []Visit
}

// Visit is a page visit in browsing history. URL and visit time
// combined are unique; no two visits have the same URL and visit time.
type Visit struct {
	URL        string
	VisitTime  time.Time // UTC
	Transition chrome.PageTransition
	PageTitle  string
}

type ExportType uint8

// Values for ExportType:
const (
	_ ExportType = iota
	AnalysisExport
	ArchivedExport
	IncrementalExport
)
