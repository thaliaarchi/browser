package takeout

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/andrewarchi/browser/archive"
	"github.com/andrewarchi/browser/bookmark"
)

type Export struct {
	Time  time.Time // time of export from filename
	Ext   string    // zip or tgz
	Parts []string  // paths to multi-part archives
}

var exportPattern = regexp.MustCompile(`^takeout-\d{8}T\d{6}Z-\d{3}\.(?:tgz|zip)$`)

func Open(filename string) (*Export, error) {
	base := filepath.Base(filename)
	if !exportPattern.MatchString(base) {
		return nil, fmt.Errorf("path is not an export: %s", base)
	}
	timestamp := base[8:24]
	seq := base[25:28]
	ext := base[29:]
	if seq != "001" {
		return nil, fmt.Errorf("archive not first in sequence: %s", seq)
	}
	t, err := time.Parse("20060102T150405Z", timestamp)
	if err != nil {
		return nil, fmt.Errorf("archive timestamp: %w", err)
	}
	glob := filename[:len(filename)-len("-001.ext")] + "-???." + ext
	parts, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	return &Export{t, ext, parts}, nil
}

func (ex *Export) Walk(walk archive.WalkFunc) error {
	var walker func(string, archive.WalkFunc) error
	switch ex.Ext {
	case "zip":
		walker = archive.WalkZip
	case "tgz":
		walker = archive.WalkTgz
	default:
		return fmt.Errorf("illegal extension: %s", ex.Ext)
	}
	for _, part := range ex.Parts {
		if err := walker(part, walk); err != nil {
			return fmt.Errorf("takeout %s: %w", filepath.Base(part), err)
		}
	}
	return nil
}

func ParseChrome(filename string) (*Chrome, error) {
	ex, err := Open(filename)
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
			d := json.NewDecoder(r)
			d.DisallowUnknownFields()
			if err := d.Decode(data); err != nil {
				return err
			}
		case "Bookmarks.html":
			b, err := bookmark.ParseNetscape(r)
			if err != nil {
				return err
			}
			data.Bookmarks = b
		case "Dictionary.csv": // TODO unknown structure
			if f.FileInfo().Size() != 0 {
				return errors.New("non-empty Dictionary.csv: TODO support format")
			}
		default:
			return fmt.Errorf("unknown file: %s", name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
