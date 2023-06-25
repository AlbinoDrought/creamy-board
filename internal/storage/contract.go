package storage

import (
	"context"
	"io"
)

type Driver interface {
	Boot(ctx context.Context) error
	Read(ctx context.Context, path string) (io.ReadCloser, error)
	Write(ctx context.Context, path string, stream io.Reader) error
	Delete(ctx context.Context, path string) error
}
