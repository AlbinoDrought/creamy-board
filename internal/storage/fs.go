package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type FSDriver struct {
	// Path to save all files in
	Path string
	// XOR read/write streams with this byte to avoid system thumbnail generation
	XOR byte
}

func safeFsPath(root string, user string) string {
	return filepath.Join(
		root,
		filepath.Clean(user),
	)
}

type xorReader struct {
	xor    byte
	source io.ReadCloser
}

func (r *xorReader) Read(p []byte) (n int, err error) {
	n, err = r.source.Read(p)
	for i := 0; i < n; i++ {
		p[i] ^= r.xor
	}
	return
}

func (r *xorReader) Close() error {
	return r.source.Close()
}

type xorWriter struct {
	xor    byte
	source io.Writer
}

func (w *xorWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	for i := range p {
		buf[i] = p[i] ^ w.xor
	}
	n, err = w.source.Write(buf)
	return
}

func (d *FSDriver) Boot(ctx context.Context) error {
	return os.MkdirAll(d.Path, os.ModePerm)
}

func (d *FSDriver) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	handle, err := os.Open(safeFsPath(d.Path, path))
	if err != nil {
		return nil, err
	}
	if d.XOR == 0 {
		return handle, nil
	}
	return &xorReader{
		source: handle,
		xor:    d.XOR,
	}, nil
}

func (d *FSDriver) Write(ctx context.Context, path string, stream io.Reader) error {
	os.MkdirAll(filepath.Dir(safeFsPath(d.Path, path)), os.ModePerm)

	handle, err := os.Create(safeFsPath(d.Path, path))
	if err != nil {
		return err
	}

	var writer io.Writer
	if d.XOR == 0 {
		writer = handle
	} else {
		writer = &xorWriter{
			source: handle,
			xor:    d.XOR,
		}
	}

	_, err = io.Copy(writer, stream)
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
