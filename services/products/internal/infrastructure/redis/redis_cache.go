package rediscache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisCache(address string, ttl time.Duration) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
		DB: 0,
	})

	return &RedisCache{
		rdb: rdb,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, data []byte) error {
	return c.rdb.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

func (c *RedisCache) Delete(ctx context.Context, keys... string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c RedisCache) MGet(ctx context.Context, keys... string) ([]interface{}, error) {
	return c.rdb.MGet(ctx, keys...).Result()
}