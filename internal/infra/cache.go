package infra

import (
	"github.com/redis/go-redis/v9"
)

func InitCache(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}
