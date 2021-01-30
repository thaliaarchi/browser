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
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pierrec/lz4/v4"
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
// each file.
func Walk(filename string, walk WalkFunc) error {
	switch ext := Ext(filename); ext {
	case ".zip":
		return WalkZip(filename, walk)
	case ".tar":
		return WalkTar(filename, walk)
	case ".tar.gz", ".tgz":
		return WalkTarGz(filename, walk)
	case ".tar.lz4":
		return WalkTarLZ4(filename, walk)
	default:
		return fmt.Errorf("archive: unsupported extension: %q", ext)
	}
}

type zipFile struct {
	f *zip.File
}

func (zf zipFile) Name() string                 { return zf.f.Name }
func (zf zipFile) Open() (io.ReadCloser, error) { return zf.f.Open() }
func (zf zipFile) FileInfo() os.FileInfo        { return zf.f.FileInfo() }

// WalkZip traverses a ZIP archive and executes the given walk function
// on each file.
func WalkZip(filename string, walk WalkFunc) error {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if err := walk(zipFile{f}); err != nil {
			return fmt.Errorf("archive: walk %s:%s: %w", filename, f.Name, err)
		}
	}
	return nil
}

type tarFile struct {
	r *tar.Reader
	h *tar.Header
}

func (tf tarFile) Name() string                 { return tf.h.Name }
func (tf tarFile) Open() (io.ReadCloser, error) { return ioutil.NopCloser(tf.r), nil }
func (tf tarFile) FileInfo() os.FileInfo        { return tf.h.FileInfo() }

// WalkTar traverses a tar archive and executes the given walk function
// on each file.
func WalkTar(filename string, walk WalkFunc) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return walkTar(f, filename, walk)
}

// WalkTarGz traverses a gzip-compressed tar archive and executes the
// given walk function on each file.
func WalkTarGz(filename string, walk WalkFunc) error {
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
	return walkTar(gr, filename, walk)
}

// WalkTarLZ4 traverses an LZ4 compressed tar archive and executes the
// given walk function on each file.
func WalkTarLZ4(filename string, walk WalkFunc) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	lr := lz4.NewReader(f)
	return walkTar(lr, filename, walk)
}

func walkTar(r io.Reader, filename string, walk WalkFunc) error {
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

// Ext returns the extension of the filename. For .gz and .lz4 two-level
// extensions, the full extension is returned like .tar.gz.
func Ext(filename string) string {
	ext := filepath.Ext(filename)
	if ext == ".gz" || ext == ".lz4" {
		ext2 := filepath.Ext(filename[:len(filename)-len(ext)])
		return filename[len(filename)-len(ext)-len(ext2):] // avoid concat
	}
	return ext
}
