package cache

import (
	"crypto/tls"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  pw,
		DB:        db,
		TLSConfig: &tls.Config{},
	})
}
