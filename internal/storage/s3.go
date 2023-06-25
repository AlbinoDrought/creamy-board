package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Driver struct {
	S3       s3iface.S3API
	Uploader *s3manager.Uploader
	Bucket   string
}

func (d *S3Driver) Boot(ctx context.Context) error {
	_, err := d.S3.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: &d.Bucket,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou {
				err = nil
			}
		}
	}
	return err
}

func (d *S3Driver) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := d.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &d.Bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (d *S3Driver) Write(ctx context.Context, path string, stream io.Reader) error {
	_, err := d.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: &d.Bucket,
		Key:    &path,
		Body:   stream,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *S3Driver) Delete(ctx context.Context, path string) error {
	_, err := d.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: &d.Bucket,
		Key:    &path,
	})
	if err != nil {
		return err
	}
	return nil
}

var _ Driver = &S3Driver{}
