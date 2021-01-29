package firefox

import (
	"io/ioutil"

	"github.com/andrewarchi/browser/timefmt"
)

type BookmarkBackup struct {
	GUID         string            `json:"guid"`
	Title        string            `json:"title"`
	Index        int               `json:"index"`
	DateAdded    timefmt.UnixMicro `json:"dateAdded"`
	LastModified timefmt.UnixMicro `json:"lastModified"`
	ID           int64             `json:"id"`
	TypeCode     int               `json:"typeCode"`
	Type         BookmarkType      `json:"type"`
	Root         string            `json:"root,omitempty"`
	Children     []BookmarkBackup  `json:"children,omitempty"`
	IconURI      string            `json:"iconuri,omitempty"`
	URI          string            `json:"uri,omitempty"`
}

type BookmarkTypeCode int
type BookmarkType string

const (
	PlaceCode          BookmarkTypeCode = 1
	PlaceContainerCode BookmarkTypeCode = 2

	Place          BookmarkType = "text/x-moz-place"
	PlaceContainer BookmarkType = "text/x-moz-place-container"
)

func ParseBookmarkBackup(filename string) (*BookmarkBackup, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var bookmarks BookmarkBackup
	if err := UnmarshalMozLz4Json(b, &bookmarks); err != nil {
		return nil, err
	}
	return &bookmarks, nil
}
