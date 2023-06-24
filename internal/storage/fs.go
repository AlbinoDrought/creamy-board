package storage

import (
	"io"
	"os"
	"path/filepath"
)

type FSDriver struct {
	Path string
}

func safeFsPath(root string, user string) string {
	return filepath.Join(
		root,
		filepath.Clean(user),
	)
}

func (d *FSDriver) Boot() error {
	return os.MkdirAll(d.Path, os.ModePerm)
}

func (d *FSDriver) Read(path string) (io.ReadCloser, error) {
	return os.Open(safeFsPath(d.Path, path))
}

func (d *FSDriver) Write(path string, stream io.ReadSeeker) error {
	handle, err := os.Create(safeFsPath(d.Path, path))
	if err != nil {
		return err
	}
	_, err = io.Copy(handle, stream)
	if err != nil {
		handle.Close()
		return err
	}

	return handle.Close()
}

func (d *FSDriver) Delete(path string) error {
	return os.Remove(safeFsPath(d.Path, path))
}

var _ Driver = &FSDriver{}
