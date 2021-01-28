package bookmark

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andrewarchi/archive/timefmt"
	"golang.org/x/net/html"
)

// Format reference:
// http://fileformats.archiveteam.org/wiki/Netscape_bookmarks
// https://docs.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/platform-apis/aa753582(v=vs.85)
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

// ParseNetscape parses a Netscape-format bookmark file.
func ParseNetscape(r io.Reader) ([]BookmarkEntry, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	if err := checkDoctype(doc, "netscape-bookmark-file-1"); err != nil {
		return nil, err
	}
	dl := doc.Find("body > dl")
	if dl.Length() != 1 {
		return nil, fmt.Errorf("root has %d lists", dl.Length())
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
				return fmt.Errorf("illegal doctype: %s", n.Data)
			}
			return nil
		}
	}
	return errors.New("doctype not found")
}

func parseFolder(dt *goquery.Selection) (*BookmarkFolder, error) {
	h3 := dt.ChildrenFiltered("h3").First()
	addDate, err := parseNumber(h3, "add_date")
	if err != nil {
		return nil, err
	}
	lastModified, err := parseNumber(h3, "last_modified")
	if err != nil {
		return nil, err
	}
	entries, err := parseFolderList(dt.ChildrenFiltered("dl").First())
	f := &BookmarkFolder{
		Title:        h3.Text(),
		AddDate:      timefmt.FromUnixMilli(addDate),
		LastModified: timefmt.FromUnixMilli(lastModified),
		Entries:      entries,
	}
	return f, nil
}

func parseFolderList(dl *goquery.Selection) ([]BookmarkEntry, error) {
	var err error
	children := dl.ChildrenFiltered("dt")
	entries := make([]BookmarkEntry, children.Length())
	children.EachWithBreak(func(_ int, dt *goquery.Selection) bool {
		var e BookmarkEntry
		a := dt.ChildrenFiltered("a").First()
		if a.Length() == 0 {
			e, err = parseFolder(dt)
			if err != nil {
				return false
			}
		} else {
			var addDate int64
			addDate, err = parseNumber(a, "add_date")
			if err != nil {
				return false
			}
			e = &Bookmark{
				Title:   a.Text(),
				URL:     a.AttrOr("href", ""),
				AddDate: timefmt.FromChromeMicro(addDate),
				IconURI: a.AttrOr("icon_uri", ""),
			}
		}
		entries = append(entries, e)
		return true
	})
	return entries, err
}

func parseNumber(e *goquery.Selection, attr string) (int64, error) {
	return strconv.ParseInt(e.AttrOr(attr, "0"), 10, 64)
}
