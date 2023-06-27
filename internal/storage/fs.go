package storage

import (
	"context"
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

func (d *FSDriver) Boot(ctx context.Context) error {
	return os.MkdirAll(d.Path, os.ModePerm)
}

func (d *FSDriver) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(safeFsPath(d.Path, path))
}

func (d *FSDriver) Write(ctx context.Context, path string, stream io.Reader) error {
	os.MkdirAll(filepath.Dir(safeFsPath(d.Path, path)), os.ModePerm)

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

func (d *FSDriver) Delete(ctx context.Context, path string) error {
	return os.Remove(safeFsPath(d.Path, path))
}

var _ Driver = &FSDriver{}
