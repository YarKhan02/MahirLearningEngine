package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
)

func GetJSON[T any](ctx context.Context, c *redis.RedisClient, name, key string) (T, bool) {
	var zero T

	raw, err := c.Get(ctx, key)
	if err != nil {
		metrics.RecordCache(name, "miss")
		return zero, false
	}

	var v T
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		metrics.RecordCache(name, "miss")
		return zero, false
	}

	metrics.RecordCache(name, "hit")
	return v, true
}

func SetJSON[T any](ctx context.Context, c *redis.RedisClient, key string, v T, ttl time.Duration) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	_ = c.Set(ctx, key, string(b), ttl)
}
