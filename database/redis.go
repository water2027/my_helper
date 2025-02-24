package database

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func initRedisClient() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	client = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
}

func SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

func GetValue(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func DeleteValue(ctx context.Context, key string) error {
	return client.Del(ctx, key).Err()
}

