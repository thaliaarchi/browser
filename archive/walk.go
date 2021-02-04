// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package archive traverses ZIP and tar archives with a common method.
package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// File exposes a common interface for files in an archive.
type File interface {
	Name() string
	Open() (io.ReadCloser, error)
	FileInfo() os.FileInfo
}

// WalkFunc is the type of function that is called for each file
// visited.
type WalkFunc func(File) error

// Walk traverses an archive and executes the given walk function on
// each file. Supported archive and compression formats: ZIP, tar, gzip,
// XZ, and LZ4.
func Walk(filename string, walk WalkFunc) error {
	exts, err := splitExt(filepath.Base(filename))
	if err != nil {
		return err
	}
	if len(exts) == 1 && exts[0] == "zip" {
		return WalkZipFile(filename, walk)
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r := io.Reader(f)
	for _, ext := range exts {
		switch ext {
		case "zip":
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}
			return WalkZip(bytes.NewReader(b), int64(len(b)), filename, walk)
		case "tar":
			return WalkTar(r, filename, walk)
		case "gz":
			gr, err := gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer gr.Close()
			r = gr
		case "xz":
			xr, err := xz.NewReader(r)
			if err != nil {
				return err
			}
			r = xr
		case "lz4":
			r = lz4.NewReader(r)
		default:
			return fmt.Errorf("archive: unsupported extension: %q", ext)
		}
	}
	return fmt.Errorf("archive: no archive extension: %s", filename)
}

type zipFile struct {
	f *zip.File
}

func (zf zipFile) Name() string                 { return zf.f.Name }
func (zf zipFile) Open() (io.ReadCloser, error) { return zf.f.Open() }
func (zf zipFile) FileInfo() os.FileInfo        { return zf.f.FileInfo() }

func walkZip(zr *zip.Reader, filename string, walk WalkFunc) error {
	for _, f := range zr.File {
		if err := walk(zipFile{f}); err != nil {
			return fmt.Errorf("archive: walk %s:%s: %w", filename, f.Name, err)
		}
	}
	return nil
}

// WalkZipFile traverses a ZIP archive from a file and executes the
// given walk function on each file.
func WalkZipFile(filename string, walk WalkFunc) error {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()
	return walkZip(&zr.Reader, filename, walk)
}

// WalkZip traverses a ZIP archive from an io.ReaderAt and executes the
// given walk function on each file.
func WalkZip(r io.ReaderAt, size int64, filename string, walk WalkFunc) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}
	return walkZip(zr, filename, walk)
}

type tarFile struct {
	r *tar.Reader
	h *tar.Header
}

func (tf tarFile) Name() string                 { return tf.h.Name }
func (tf tarFile) Open() (io.ReadCloser, error) { return ioutil.NopCloser(tf.r), nil }
func (tf tarFile) FileInfo() os.FileInfo        { return tf.h.FileInfo() }

// WalkTar traverses a tar archive from an io.Reader and executes the
// given walk function on each file.
func WalkTar(r io.Reader, filename string, walk WalkFunc) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if err := walk(tarFile{tr, header}); err != nil {
			return fmt.Errorf("archive: walk %s:%s: %w", filename, header.Name, err)
		}
	}
	return nil
}

// WalkTarFile traverses a tar archive from a file and executes the
// given walk function on each file.
func WalkTarFile(filename string, walk WalkFunc) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return WalkTar(f, filename, walk)
}

// splitExt splits the filename into recognized extensions.
func splitExt(filename string) ([]string, error) {
	name := filename
	var exts []string
	for {
		switch ext := filepath.Ext(name); ext {
		case ".zip", ".tar":
			return append(exts, ext[1:]), nil
		case ".tgz", ".txz":
			return append(exts, ext[2:], "tar"), nil
		case ".gz", ".xz", ".lz4":
			exts = append(exts, ext[1:])
			name = name[:len(name)-len(ext)]
		case "":
			return nil, fmt.Errorf("archive: no archive extension: %q", filename)
		default:
			return nil, fmt.Errorf("archive: unrecognized extension %q: %q", ext, filename)
		}
	}
}
