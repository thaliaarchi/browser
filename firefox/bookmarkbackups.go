package firefox

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// BookmarkBackup is a backup of Firefox bookmarks.
type BookmarkBackup struct {
	Date       time.Time // date of backup
	Count      int       // number of entries
	Hash       []byte    // hash of json contents
	Compressed bool      // true when mozlz4 compressed
	Bookmarks  *BookmarkBackupEntry
}

// BookmarkBackupEntry is an entry in a bookmark backup.
type BookmarkBackupEntry struct {
	GUID         string                `json:"guid"`
	Title        string                `json:"title"`
	Index        int                   `json:"index"`
	DateAdded    timefmt.UnixMicro     `json:"dateAdded"`
	LastModified timefmt.UnixMicro     `json:"lastModified"`
	ID           int64                 `json:"id"`
	TypeCode     int                   `json:"typeCode"` // place: 1, place-container: 2
	Type         string                `json:"type"`     // "text/x-moz-place", "text/x-moz-place-container"
	Root         string                `json:"root,omitempty"`
	Children     []BookmarkBackupEntry `json:"children,omitempty"`
	IconURI      string                `json:"iconuri,omitempty"`
	URI          string                `json:"uri,omitempty"`
}

// TODO handle all kinds of entry types from Firefox source.

// ParseBookmarkBackup parses a bookmarks file within bookmarkbackups in
// a Firefox profile.
func ParseBookmarkBackup(filename string) (*BookmarkBackup, error) {
	// JSON bookmark backup serialization:
	// https://searchfox.org/mozilla-central/source/toolkit/components/places/PlacesBackups.jsm#265

	backup, err := GetBookmarkBackupMetadata(filename)
	if err != nil {
		return nil, err
	}
	if backup.Compressed {
		if err := jsonutil.DecodeMozLz4File(filename, &backup.Bookmarks); err != nil {
			return nil, err
		}
	} else if err := jsonutil.DecodeFile(filename, &backup.Bookmarks); err != nil {
		return nil, err
	}
	return backup, nil
}

// bookmarkBackupPattern matches the backup filename:
//   0: file name
//   1: date in form 2006-01-02
//   2: bookmarks count
//   3: contents hash
//   4: file extension
//
// Pattern defined as PlacesBackups.filenamesRegex in
// https://searchfox.org/mozilla-central/source/toolkit/components/places/PlacesBackups.jsm#98
var bookmarkBackupPattern = regexp.MustCompile(`^bookmarks-([0-9-]+)(?:_([0-9]+)){0,1}(?:_([A-Za-z0-9=+-]{24})){0,1}\.(json(?:lz4)?)$`)

// TODO is the ordering -+ or +- ?
var filepathBase64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-+") // padding =

// GetBookmarkBackupMetadata reads the metadata from a bookmark backup
// filename. The returned BookmarkBackup has nil Bookmarks.
func GetBookmarkBackupMetadata(filename string) (*BookmarkBackup, error) {
	base := filepath.Base(filename)
	matches := bookmarkBackupPattern.FindStringSubmatch(base)
	if len(matches) != 5 {
		return nil, fmt.Errorf("firefox: filename is not a bookmark backup: %s", base)
	}
	var meta BookmarkBackup
	var err error
	meta.Date, err = time.ParseInLocation("2006-01-02", matches[1], time.Local)
	if err != nil {
		return nil, err
	}
	meta.Count = -1
	if matches[2] != "" {
		meta.Count, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, err
		}
	}
	if matches[3] != "" {
		meta.Hash, err = filepathBase64.DecodeString(matches[3])
		if err != nil {
			return nil, err
		}
		if len(meta.Hash) != 16 {
			return nil, fmt.Errorf("firefox: bookmark backup has is not 16 bytes: %s", matches[3])
		}
	}
	meta.Compressed = matches[4] == "jsonlz4"
	return &meta, err
}
