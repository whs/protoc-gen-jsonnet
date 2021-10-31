package main

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
)

type FsLoader struct {
	Fs fs.FS
}

func (f *FsLoader) Abs(base, name string) string {
	return filepath.Join(filepath.Dir(base), name)
}

func (f *FsLoader) Get(path string) (io.Reader, error) {
	reader, err := f.Fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
