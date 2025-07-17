package utils

import (
	"os"
	"strconv"
	"sync"

	goredis "github.com/redis/go-redis/v9"
)

var (
	redisOnce   sync.Once
	redisClient *goredis.Client
)

// Redis returns a singleton *redis.Client configured from env vars or defaults.
// Default: REDIS_ADDR (localhost:6379), REDIS_DB (0), REDIS_PASSWORD (empty).
func Redis() *goredis.Client {
	redisOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}
		db := 0
		// optional REDIS_DB env.
		if v := os.Getenv("REDIS_DB"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil {
				db = parsed
			}
		}
		redisClient = goredis.NewClient(&goredis.Options{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
		})
	})
	return redisClient
}
