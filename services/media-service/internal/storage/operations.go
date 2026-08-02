package storage

import (
	"context"
	"io"
	"time"

	"github.com/diyorbeknematov/minitwitter/services/media-service/pkg/apperror"
	"github.com/minio/minio-go/v7"
)

func (m *MinIO) Upload(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	info, err := m.Client.PutObject(
		ctx,
		m.Bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", apperror.Wrap("storage", "Upload", "failed to upload object", err)
	}

	return info.Key, nil
}

func (m *MinIO) Delete(ctx context.Context, objectName string) error {
	err := m.Client.RemoveObject(ctx, m.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return apperror.Wrap("storage", "Delete", "failed to delete object", err)
	}

	return nil
}

func (m *MinIO) PresignedURL(ctx context.Context, objectName string, urlExpire time.Duration) (string, error) {
	url, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectName, urlExpire, nil)
	if err != nil {
		return "", apperror.Wrap("storage", "PresignedURL", "failed to generate presigned URL", err)
	}

	return url.String(), nil
}
