package infra

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	redisClient *redis.Client
}

type CacheInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
	Ping(ctx context.Context) error
}

func InitCache(addr, password string) CacheInterface {
	return &cache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		}),
	}
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *cache) Close() error {
	return c.redisClient.Close()
}

func (c *cache) Ping(ctx context.Context) error {
	return c.redisClient.Ping(ctx).Err()
}
