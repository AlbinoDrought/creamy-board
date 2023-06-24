package storage

import (
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3Driver struct {
	S3     s3iface.S3API
	Bucket string
}

func (d *S3Driver) Boot() error {
	_, err := d.S3.CreateBucket(&s3.CreateBucketInput{
		Bucket: &d.Bucket,
	})
	return err
}

func (d *S3Driver) Read(path string) (io.ReadCloser, error) {
	output, err := d.S3.GetObject(&s3.GetObjectInput{
		Bucket: &d.Bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (d *S3Driver) Write(path string, stream io.ReadSeeker) error {
	_, err := d.S3.PutObject(&s3.PutObjectInput{
		Bucket: &d.Bucket,
		Key:    &path,
		Body:   stream,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *S3Driver) Delete(path string) error {
	_, err := d.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &d.Bucket,
		Key:    &path,
	})
	if err != nil {
		return err
	}
	return nil
}

var _ Driver = &S3Driver{}
