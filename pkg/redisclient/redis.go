// pkg/redisclient/redis.go
package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// InitRedisClient initializes a Redis client and tests the connection
func InitRedisClient(addr string) error {
	fmt.Println("🔌 Connecting to Redis at", addr)

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Ping to check connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	client = rdb
	fmt.Println("✅ Redis connected successfully")
	return nil
}

// CloseRedisClient safely closes the Redis connection
func CloseRedisClient() error {
	if client == nil {
		return nil
	}
	fmt.Println("🔒 Closing Redis connection...")
	return client.Close()
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	if client == nil {
		panic("Redis client is not initialized")
	}
	return client
}

// Get retrieves a value from Redis by key
func Get(key string) (string, error) {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get value from Redis: %w", err)
	}
	return val, nil
}

// Set sets a value in Redis with an expiration time
func Set(key string, value string, expiration time.Duration) error {
	err := client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}
	return nil
}

// Delete removes a key from Redis
func Delete(key string) error {
	err := client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %w", err)
	}
	return nil
}
