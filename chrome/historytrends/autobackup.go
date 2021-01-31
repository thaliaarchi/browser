// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package historytrends

import (
	"time"

	"github.com/andrewarchi/browser/chrome"
)

type AutoBackup struct {
	Date       time.Time // YYYYMMDD
	BackupType string    // "full" or "incremental"
	Visits     []AutoBackupVisit
}

type AutoBackupVisit struct {
	URL        string
	VisitTime  time.Time             // i.e. "U1384634958041.754"
	Transition chrome.TransitionType // i.e. 1
	PageTitle  string
}
