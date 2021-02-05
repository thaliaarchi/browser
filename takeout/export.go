// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package takeout traverses and parses Google Takeout exports.
package takeout

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/andrewarchi/archive"
)

// Export contains the paths to each part in a Takeout export and the
// time of export. Zip exports are significantly faster to traverse than
// tgz and should be preferred.
type Export struct {
	Time      time.Time // time of export from filename
	Timestamp string    // raw timestamp
	Parts     []string  // paths to multi-part archives
}

var exportPattern = regexp.MustCompile(`^takeout-(\d{8}T\d{6}Z)-(\d{3})\.(tgz|zip)$`)

// NewExport opens a Takeout export, given the path to the first archive
// in a multi-part export.
func NewExport(filename string) (*Export, error) {
	dir, base := filepath.Split(filename)
	match := exportPattern.FindStringSubmatch(base)
	if len(match) != 4 {
		return nil, fmt.Errorf("takeout: path is not an export: %q", base)
	}
	timestamp, seq, ext := match[1], match[2], match[3]
	if seq != "001" {
		return nil, fmt.Errorf("takeout: archive not first in sequence: %q", seq)
	}
	t, err := time.Parse("20060102T150405Z", timestamp)
	if err != nil {
		return nil, fmt.Errorf("takeout: export timestamp: %w", err)
	}
	glob := filepath.Join(dir, "takeout-"+timestamp+"-???."+ext)
	parts, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	return &Export{t, timestamp, parts}, nil
}

// Walk traverses a Takeout export and executes the given walk function
// on each file.
func (ex *Export) Walk(walk archive.WalkFunc) error {
	for _, part := range ex.Parts {
		if err := archive.Walk(part, walk); err != nil {
			return err
		}
	}
	return nil
}
