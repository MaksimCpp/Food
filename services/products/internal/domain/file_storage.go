package domain

import (
	"context"
	"time"
)

type FileStorage interface {
	Upload(
		ctx context.Context, objectName string, fileBytes []byte, objectSize int64, contentType string,
	) (string, error)
	GetImage(
		ctx context.Context, objectName string, expires time.Duration,
	) (string, error)
	Delete(ctx context.Context, objectName string) error
}