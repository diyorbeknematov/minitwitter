package storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Upload(
		ctx context.Context,
		objectName string,
		reader io.Reader,
		size int64,
		contentType string,
	) (string, error)

	Delete(ctx context.Context, objectName string) error

	PresignedURL(ctx context.Context, objectName string) (string, error)
}
