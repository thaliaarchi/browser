package chrome

import (
	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// Bookmarks contains Chrome bookmark information.
type Bookmarks struct {
	Checksum     string        `json:"checksum"` // hex
	Roots        BookmarkRoots `json:"roots"`
	SyncMetadata string        `json:"sync_metadata,omitempty"` // base64-encoded
	Version      int           `json:"version"`                 // i.e. 1
}

// BookmarkRoots contains the root level bookmarks folders.
type BookmarkRoots struct {
	BookmarkBar BookmarkEntry `json:"bookmark_bar"` // "Bookmarks" folder
	Other       BookmarkEntry `json:"other"`        // "Other Bookmarks" folder
	Synced      BookmarkEntry `json:"synced"`       // "Mobile Bookmarks" folder
}

// BookmarkEntry is either a folder containing further entries or a URL.
type BookmarkEntry struct {
	Children     []BookmarkEntry   `json:"children"`
	DateAdded    timefmt.Chrome    `json:"date_added"`
	DateModified timefmt.Chrome    `json:"date_modified,omitempty"` // for folder type only
	GUID         string            `json:"guid"`                    // i.e. "01234567-89ab-cdef-0123-456789abcdef"
	ID           string            `json:"id"`                      // i.e. "567"
	Name         string            `json:"name"`
	Type         string            `json:"type"` // "folder" or "url"
	MetaInfo     *BookmarkMetaInfo `json:"meta_info,omitempty"`
	URL          string            `json:"url,omitempty"` // for url type only
}

// BookmarkMetaInfo contains additional bookmark metadata.
type BookmarkMetaInfo struct {
	LastVisitedDesktop timefmt.Chrome `json:"last_visited_desktop"`
}

func ParseBookmarks(filename string) (*Bookmarks, error) {
	var bookmarks Bookmarks
	if err := jsonutil.Decode(filename, &bookmarks); err != nil {
		return nil, err
	}
	return &bookmarks, nil
}
