// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package takeout traverses and parses Google Takeout exports.
package takeout

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/andrewarchi/browser/archive"
	"github.com/andrewarchi/browser/bookmark"
	"github.com/andrewarchi/browser/jsonutil"
)

// Export contains the paths to each part in a Takeout export and the
// time of export. Zip exports are significantly faster to traverse than
// tgz and should be preferred.
type Export struct {
	Time  time.Time // time of export from filename
	Ext   string    // zip or tgz
	Parts []string  // paths to multi-part archives
}

var exportPattern = regexp.MustCompile(`^takeout-(\d{8}T\d{6}Z)-(\d{3})\.(tgz|zip)$`)

// NewExport opens a Takeout export, given the path to the first archive
// in a multi-part export.
func NewExport(filename string) (*Export, error) {
	base := filepath.Base(filename)
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
	parts, err := filepath.Glob("takeout-" + timestamp + "-???." + ext)
	if err != nil {
		return nil, err
	}
	return &Export{t, ext, parts}, nil
}

// Walk traverses a Takeout export and executes the given walk function
// on each file.
func (ex *Export) Walk(walk archive.WalkFunc) error {
	var walker func(string, archive.WalkFunc) error
	switch ex.Ext {
	case "zip":
		walker = archive.WalkZip
	case "tgz":
		walker = archive.WalkTgz
	default:
		return fmt.Errorf("takeout: illegal extension: %q", ex.Ext)
	}
	for _, part := range ex.Parts {
		if err := walker(part, walk); err != nil {
			return err
		}
	}
	return nil
}

// ParseChrome parses Chrome data in a Takeout export.
func ParseChrome(filename string) (*Chrome, error) {
	ex, err := NewExport(filename)
	if err != nil {
		return nil, err
	}
	data := &Chrome{ExportTime: ex.Time}
	err = ex.Walk(func(f archive.File) error {
		name := f.Name()
		if filepath.Dir(name) != "Takeout/Chrome" {
			return nil
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		defer r.Close()
		switch base := filepath.Base(name); base {
		case "Autofill.json", "BrowserHistory.json", "Extensions.json",
			"SearchEngines.json", "SyncSettings.json":
			return jsonutil.Decode(r, data)
		case "Bookmarks.html":
			b, err := bookmark.ParseHTML(r)
			if err != nil {
				return err
			}
			data.Bookmarks = b
		case "Dictionary.csv": // TODO unknown structure
			if f.FileInfo().Size() != 0 {
				return errors.New("dictionary structure unknown")
			}
		default:
			return errors.New("unknown file")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
