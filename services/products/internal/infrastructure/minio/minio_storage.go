package miniostorage

import (
	"bytes"
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) *MinioStorage {
	client, err := minio.New(
		endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		},
	)

	if err != nil {
		return nil
	}

	return &MinioStorage{
		client: client,
		bucketName: bucketName,
	}
}

func (s *MinioStorage) Upload(
	ctx context.Context, objectName string, fileBytes []byte, objectSize int64, contentType string,
) (string, error) {
	isExists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return "", err
	}

	if !isExists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})

		if err != nil {
			return "", err
		}
	}

	reader := bytes.NewReader(fileBytes)

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (s *MinioStorage) GetImage(
	ctx context.Context, objectName string, expires time.Duration,
) (string, error) {
	req := make(url.Values)

	objectUrl, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expires, req)

	if err != nil {
		return "", err
	}

	return objectUrl.String(), nil
}

func (s *MinioStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
}