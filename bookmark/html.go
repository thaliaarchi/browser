// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package bookmark parses Netscape-style HTML bookmark files.
package bookmark

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
	"golang.org/x/net/html"
)

// Format reference:
// http://fileformats.archiveteam.org/wiki/Netscape_bookmarks
// http://fileformats.archiveteam.org/wiki/Chrome_bookmarks
// https://docs.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/platform-apis/aa753582(v=vs.85)
//
// Firefox:
// https://searchfox.org/mozilla-central/source/toolkit/components/places/BookmarkHTMLUtils.jsm
//
// Other services appear to use fields not used by
// Chrome bookmarks in Google Takeout.

type BookmarkEntry interface{} // BookmarkFolder or Bookmark

type BookmarkFolder struct {
	Title        string
	AddDate      time.Time
	LastModified time.Time
	Entries      []BookmarkEntry
}

type Bookmark struct {
	Title   string
	URL     string
	AddDate time.Time
	IconURI string
}

// ParseHTML parses a Netscape-style HTML bookmark file.
func ParseHTML(r io.Reader) ([]BookmarkEntry, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	if err := checkDoctype(doc, "netscape-bookmark-file-1"); err != nil {
		return nil, err
	}
	dl := doc.Find("body > dl")
	if dl.Length() != 1 {
		return nil, fmt.Errorf("bookmark: root has %d lists", dl.Length())
	}
	entries, err := parseFolderList(dl)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func checkDoctype(doc *goquery.Document, doctype string) error {
	for n := doc.Nodes[0].FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.DoctypeNode {
			if n.Data != doctype {
				return fmt.Errorf("bookmark: illegal doctype: %q", n.Data)
			}
			return nil
		}
	}
	return errors.New("bookmark: doctype not found")
}

func parseFolder(dt *goquery.Selection) (*BookmarkFolder, error) {
	h3 := dt.ChildrenFiltered("h3").First()
	add := h3.AttrOr("add_date", "")
	mod := h3.AttrOr("last_modified", "")
	addDate, err := timefmt.Parse(add, timefmt.Milli, timefmt.Unix)
	if err != nil {
		return nil, err
	}
	lastModified, err := timefmt.Parse(mod, timefmt.Milli, timefmt.Unix)
	if err != nil {
		return nil, err
	}
	entries, err := parseFolderList(dt.ChildrenFiltered("dl").First())
	f := &BookmarkFolder{
		Title:        h3.Text(),
		AddDate:      addDate,
		LastModified: lastModified,
		Entries:      entries,
	}
	return f, nil
}

func parseFolderList(dl *goquery.Selection) ([]BookmarkEntry, error) {
	var err error
	children := dl.ChildrenFiltered("dt")
	entries := make([]BookmarkEntry, 0, children.Length())
	children.EachWithBreak(func(_ int, dt *goquery.Selection) bool {
		var e BookmarkEntry
		a := dt.ChildrenFiltered("a").First()
		if a.Length() == 0 {
			e, err = parseFolder(dt)
			if err != nil {
				return false
			}
		} else {
			var addDate time.Time
			add := a.AttrOr("add_date", "0")
			addDate, err = timefmt.Parse(add, timefmt.Micro, timefmt.Windows)
			if err != nil {
				return false
			}
			e = &Bookmark{
				Title:   a.Text(),
				URL:     a.AttrOr("href", ""),
				AddDate: addDate,
				IconURI: a.AttrOr("icon_uri", ""),
			}
		}
		entries = append(entries, e)
		return true
	})
	return entries, err
}
