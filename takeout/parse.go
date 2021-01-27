package takeout

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewarchi/archive/bookmark"
)

func ParseChromeTgz(filename string) (*Chrome, error) {
	if !strings.HasSuffix(filename, "-001.tgz") {
		return nil, fmt.Errorf("archive must end with -001.tgz: %s", filename)
	}
	glob := strings.TrimSuffix(filename, "-001.tgz") + "-???.tgz"
	parts, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var data Chrome
	for _, part := range parts {
		if err := parseTgzPart(part, &data); err != nil {
			return nil, err
		}
	}
	return &data, nil
}

func parseTgzPart(filename string, data *Chrome) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg ||
			filepath.Dir(header.Name) != "Takeout/Chrome" {
			continue
		}
		switch base := filepath.Base(header.Name); base {
		case "Autofill.json", "BrowserHistory.json", "Extensions.json",
			"SearchEngines.json", "SyncSettings.json":
			d := json.NewDecoder(tr)
			d.DisallowUnknownFields()
			if err := d.Decode(data); err != nil {
				return fmt.Errorf("%s: %w", base, err)
			}
		case "Bookmarks.html":
			b, err := bookmark.ParseNetscape(tr)
			if err != nil {
				return fmt.Errorf("%s: %w", base, err)
			}
			data.Bookmarks = b
		case "Dictionary.csv": // TODO
		default:
			log.Printf("Unknown Chrome file: %s/%s", filename, header.Name)
		}
	}
	return nil
}

func ParseChromeZip(filename string) (*Chrome, error) {
	if !strings.HasSuffix(filename, "-001.zip") {
		return nil, fmt.Errorf("archive must end with -001.zip: %s", filename)
	}
	glob := strings.TrimSuffix(filename, "-001.zip") + "-???.zip"
	parts, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var data Chrome
	for _, part := range parts {
		if err := parseZipPart(part, &data); err != nil {
			return nil, err
		}
	}
	return &data, nil
}

func parseZipPart(filename string, data *Chrome) error {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if filepath.Dir(f.Name) != "Takeout/Chrome" {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		defer r.Close()
		switch base := filepath.Base(f.Name); base {
		case "Autofill.json", "BrowserHistory.json", "Extensions.json",
			"SearchEngines.json", "SyncSettings.json":
			d := json.NewDecoder(r)
			d.DisallowUnknownFields()
			if err := d.Decode(data); err != nil {
				return fmt.Errorf("%s: %w", base, err)
			}
		case "Bookmarks.html":
			b, err := bookmark.ParseNetscape(r)
			if err != nil {
				return fmt.Errorf("%s: %w", base, err)
			}
			data.Bookmarks = b
		case "Dictionary.csv": // TODO
		default:
			log.Printf("Unknown Chrome file: %s/%s", filename, f.Name)
		}
	}
	return nil
}
