package redis

import (
	"context"
	"time"

	redislib "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redislib.Client
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	opts, err := redislib.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redislib.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Close() error {
	if r == nil || r.client == nil {
		return nil
	}

	return r.client.Close()
}

func (r *RedisClient) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, 1, ttl).Result()
}

func (r *RedisClient) ReleaseLock(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
} 