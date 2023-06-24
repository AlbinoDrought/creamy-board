package storage

import "io"

type Driver interface {
	Boot() error
	Read(path string) (io.ReadCloser, error)
	Write(path string, stream io.ReadSeeker) error
	Delete(path string) error
}
