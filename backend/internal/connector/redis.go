package connector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"social-app/internal/config"
)

type RedisConnector struct {
	client *redis.Client
}

func NewRedisConnector(r config.Redis) *RedisConnector {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
		// PoolSize:     r.PoolSize,
		// MinIdleConns: r.MinIdle,
	})

	return &RedisConnector{
		client: client,
	}
}

func (rc *RedisConnector) Ping(ctx context.Context) error {
	_, err := rc.client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis ping failed: %v", err)
		return fmt.Errorf("failed to ping Redis: %w", err)
	}
	return nil
}

func (rc *RedisConnector) Close() error {
	if err := rc.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}

	return nil
}

func (rc *RedisConnector) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (rc *RedisConnector) Set(ctx context.Context, key, value string, exp time.Duration) error {
	err := rc.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (rc *RedisConnector) Delete(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func (rc *RedisConnector) DeleteIfExists(ctx context.Context, key string) (valKey string, isExists bool, err error) {
	val, err := rc.Get(ctx, key)
	if err != nil {
		return "", false, fmt.Errorf("failed to check if key %s exists: %w", key, err)
	}

	result, err := rc.client.Del(ctx, key).Result()
	if err != nil {
		return "", false, fmt.Errorf("failed to delete key %s if exists: %w", key, err)
	}
	return val, result > 0, nil
}

func (rc *RedisConnector) SetIfNotExists(ctx context.Context, key, value string) (bool, error) {
	result, err := rc.client.SetNX(ctx, key, value, 0).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key %s if not exists: %w", key, err)
	}
	return result, nil
}

func (rc *RedisConnector) IsExists(ctx context.Context, key string) (bool, error) {
	exists, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis error: %w", err)
	}
	return exists == 1, nil
}
