// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"archive/zip"
	"fmt"
	"io"
)

// OpenSingleFileZip opens a zip containing a single file for reading
// and returns the filename of the contained file.
func OpenSingleFileZip(filename string) (io.ReadCloser, string, error) {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return nil, "", err
	}
	if len(zr.File) != 1 {
		return nil, "", fmt.Errorf("archive: zip has %d files: %q", len(zr.File), filename)
	}
	f, err := zr.File[0].Open()
	if err != nil {
		return nil, "", err
	}
	return &singleFileZip{zr, f}, zr.File[0].Name, nil
}

type singleFileZip struct {
	zr *zip.ReadCloser
	f  io.ReadCloser
}

func (z *singleFileZip) Read(p []byte) (n int, err error) {
	return z.f.Read(p)
}

func (z *singleFileZip) Close() error {
	err1 := z.f.Close()
	err2 := z.zr.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
