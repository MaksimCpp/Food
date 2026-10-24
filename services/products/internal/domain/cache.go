package domain

import (
	"context"
)

type Cache interface {
	Set(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, keys... string) error
	MGet(ctx context.Context, keys... string) ([]interface{}, error)
}