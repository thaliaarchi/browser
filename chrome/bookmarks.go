// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package chrome

import (
	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// Bookmarks contains Chrome bookmark information.
type Bookmarks struct {
	Checksum     jsonutil.Hex    `json:"checksum"`
	Roots        BookmarkRoots   `json:"roots"`
	SyncMetadata jsonutil.Base64 `json:"sync_metadata,omitempty"`
	Version      int             `json:"version"` // i.e. 1
}

// BookmarkRoots contains the root level bookmarks folders.
type BookmarkRoots struct {
	BookmarkBar BookmarkEntry `json:"bookmark_bar"` // "Bookmarks" folder
	Other       BookmarkEntry `json:"other"`        // "Other Bookmarks" folder
	Synced      BookmarkEntry `json:"synced"`       // "Mobile Bookmarks" folder
}

// BookmarkEntry is either a folder containing further entries or a URL.
type BookmarkEntry struct {
	Children     []BookmarkEntry      `json:"children"` // for folder type only
	DateAdded    timefmt.QuotedChrome `json:"date_added"`
	DateModified timefmt.QuotedChrome `json:"date_modified,omitempty"` // for folder type only
	GUID         string               `json:"guid"`                    // i.e. "01234567-89ab-cdef-0123-456789abcdef"
	ID           string               `json:"id"`                      // i.e. "567"
	Name         string               `json:"name"`
	Type         string               `json:"type"` // "folder" or "url"
	MetaInfo     *BookmarkMetaInfo    `json:"meta_info,omitempty"`
	URL          string               `json:"url,omitempty"` // for url type only
}

// BookmarkMetaInfo contains additional bookmark metadata.
type BookmarkMetaInfo struct {
	LastVisitedDesktop timefmt.QuotedChrome `json:"last_visited_desktop"`
}

// ParseBookmarks parses "Bookmarks" in a Chrome profile.
func ParseBookmarks(filename string) (*Bookmarks, error) {
	var bookmarks Bookmarks
	if err := jsonutil.DecodeFile(filename, &bookmarks); err != nil {
		return nil, err
	}
	return &bookmarks, nil
}
