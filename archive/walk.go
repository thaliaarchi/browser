package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// File exposes a common interface for files in a zip or tar archive.
type File interface {
	Name() string
	Open() (io.ReadCloser, error)
	FileInfo() os.FileInfo
}

// WalkFunc is the type of function that is called for each file
// visited.
type WalkFunc func(File) error

type zipFile struct {
	f *zip.File
}

func (zf zipFile) Name() string                 { return zf.f.Name }
func (zf zipFile) Open() (io.ReadCloser, error) { return zf.f.Open() }
func (zf zipFile) FileInfo() os.FileInfo        { return zf.f.FileInfo() }

// WalkZip traverses a zip archive and executes the given walk function
// on each file.
func WalkZip(filename string, walk WalkFunc) error {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if err := walk(zipFile{f}); err != nil {
			return fmt.Errorf("walk %s:%s: %w", filename, f.Name, err)
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

// WalkTgz traverses a gzip compressed tar archive and executes the
// given walk function on each file. Zip archives are significantly
// faster to traverse and should be preferred.
func WalkTgz(filename string, walk WalkFunc) error {
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

// WalkTar traverses a tar archive and executes the given walk function
// on each file. Zip archives are significantly faster to traverse and
// should be preferred.
func WalkTar(filename string, walk WalkFunc) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return walkTar(f, filename, walk)
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
			return fmt.Errorf("walk %s:%s: %w", filename, header.Name, err)
		}
	}
	return nil
}
