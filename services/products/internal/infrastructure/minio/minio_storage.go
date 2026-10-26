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
	// presignClient подписывает ссылки публичным адресом, доступным снаружи docker-сети
	presignClient *minio.Client
	bucketName string
}

func NewMinioStorage(
	endpoint, publicEndpoint, accessKey, secretKey, bucketName string, useSSL bool,
) (*MinioStorage, error) {
	client, err := minio.New(
		endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		},
	)

	if err != nil {
		return nil, err
	}

	presignClient := client

	if publicEndpoint != "" {
		// Region задан явно, чтобы при подписи не было запроса к publicEndpoint
		presignClient, err = minio.New(
			publicEndpoint,
			&minio.Options{
				Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
				Secure: useSSL,
				Region: "us-east-1",
			},
		)

		if err != nil {
			return nil, err
		}
	}

	return &MinioStorage{
		client: client,
		presignClient: presignClient,
		bucketName: bucketName,
	}, nil
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

	objectUrl, err := s.presignClient.PresignedGetObject(ctx, s.bucketName, objectName, expires, req)

	if err != nil {
		return "", err
	}

	return objectUrl.String(), nil
}

func (s *MinioStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
}