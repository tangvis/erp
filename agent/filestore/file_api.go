package filestore

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	UploadBytes(ctx context.Context, filename string, data []byte) (string, error)
	Upload(ctx context.Context, filename string, body io.Reader) (string, error)
	DownloadFileByURL(ctx context.Context, location string) ([]byte, error)
	DeleteFileByURL(ctx context.Context, location string) error
	GetObject(ctx context.Context, bucket, key string) (*s3.GetObjectOutput, error)
	UploadWithFileName(ctx context.Context, filename string, body io.Reader) (string, error)
}
