package takeout

import (
	"archive/tar"
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
	zr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer zr.Close()
	tr := tar.NewReader(zr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
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
				return err
			}
		case "Bookmarks.html":
			b, err := bookmark.ParseNetscape(tr)
			if err != nil {
				return err
			}
			data.Bookmarks = b
		case "Dictionary.csv": // TODO
		default:
			log.Printf("Unknown Chrome file: %s/%s", filename, header.Name)
		}
	}
}
