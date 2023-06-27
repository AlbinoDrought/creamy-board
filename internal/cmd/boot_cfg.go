package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pgxpool4 "github.com/jackc/pgx/v4/pgxpool"
	pgx5 "github.com/jackc/pgx/v5"
	"go.albinodrought.com/creamy-board/internal/db"
	"go.albinodrought.com/creamy-board/internal/storage"
)

func bootDB4(ctx context.Context) (*pgxpool4.Pool, error) {
	return db.ConnectPool4(ctx, os.Getenv("CREAMY_DSN"))
}

func bootDB5(ctx context.Context) (*pgx5.Conn, error) {
	return db.Connect5(ctx, os.Getenv("CREAMY_DSN"))
}

var ErrUnknownStorageDriver = errors.New("unknown storage driver")
var ErrXorValueMustBeOneByte = errors.New("xor value must be one byte")

func bootStorage(ctx context.Context) (storage.Driver, error) {
	driver := os.Getenv("CREAMY_STORAGE_DRIVER")

	if driver == "" || driver == "fs" {
		path := os.Getenv("CREAMY_STORAGE_PATH")
		xorStr := os.Getenv("CREAMY_STORAGE_XOR")
		xor := byte(0)
		if xorStr != "" {
			xorBytes, err := hex.DecodeString(xorStr)
			if err != nil {
				return nil, err
			}
			if len(xorBytes) != 1 {
				return nil, ErrXorValueMustBeOneByte
			}
			xor = xorBytes[0]
		}
		return &storage.FSDriver{
			Path: path,
			XOR:  xor,
		}, nil
	}

	if driver == "minio" {
		s3Session, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("CREAMY_MINIO_KEY"),
				os.Getenv("CREAMY_MINIO_SECRET"),
				"",
			),
			Endpoint:         aws.String(os.Getenv("CREAMY_MINIO_ENDPOINT")),
			Region:           aws.String("us-west-2"),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		})
		if err != nil {
			return nil, err
		}
		s3Client := s3.New(s3Session)
		s3Uploader := s3manager.NewUploaderWithClient(s3Client)
		return &storage.S3Driver{
			S3:       s3Client,
			Uploader: s3Uploader,
			Bucket:   os.Getenv("CREAMY_MINIO_BUCKET"),
		}, nil
	}

	return nil, ErrUnknownStorageDriver
}
