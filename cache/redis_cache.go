package cache

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

// Rdb represents the REDIS connection
var Rdb *redis.Client
var Ctx = context.Background()

// InitRedis initializes the Redis client connection.
// It loads the Redis server address and password from environment variables,
// then establishes a connection to the Redis server by calling Ping.
// If the connection fails, it logs a fatal error and stops the application.
func InitRedis() {
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
