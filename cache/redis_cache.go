package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient defines basic Redis operations.
type RedisClient interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

// Rdb is the global Redis client instance.
var Rdb RedisClient

// Ctx is the default context used for Redis operations.
var Ctx = context.Background()

// InitRedis initializes the Redis client connection.
// It loads the Redis server address and password from environment variables,
// then establishes a connection to the Redis server by calling Ping.
// If the connection fails, it logs a fatal error and stops the application.
func InitRedis(client RedisClient) {
	if client != nil {
		Rdb = client
		return
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}
