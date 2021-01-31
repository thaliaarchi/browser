// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

// fullURL represents an entry in the SQL urls table, defined in
// chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/js/calcTrendsWebSql.js,
// line 25.
type fullURL struct {
	URLID      int    // primary key, sequential
	URL        string // unique
	Host       string // derived from URL
	RootDomain string // derived from Host
	Title      string // /[\t\r\n]+/ is normalized to ' ' then /\s\s+/ to ' '
}

// fullVisit represents an entry in the SQL visits table, defined in
// chrome-extension://pnmchffiealhkdloeffcdnbgdnedheme/js/calcTrendsWebSql.js,
// line 37.
type fullVisit struct {
	VisitID        int // primary key, sequential
	URLID          int // unique together with VisitTime
	VisitTime      int // unique together with URLID
	VisitDate      int
	Year           int
	Month          int
	MonthDay       int
	WeekDay        int
	Hour           int
	TransitionType string
}
