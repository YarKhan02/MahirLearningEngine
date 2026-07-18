package timetable

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	timetableCacheVersion = "v1"
	timetableCacheName    = "timetable"
	timetableUserTTL      = 2 * time.Minute
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func timetableUserKey(userID uuid.UUID) string {
	return "timetable:" + timetableCacheVersion + ":user:" + userID.String()
}

// cached read

func (c *CachedRepository) GetRulesForUser(ctx context.Context, userID uuid.UUID) ([]Timetable, error) {
	key := timetableUserKey(userID)
	if v, ok := cache.GetJSON[[]Timetable](ctx, c.cache, timetableCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetRulesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, timetableUserTTL)
	return list, nil
}

// pass-throughs

func (c *CachedRepository) Create(ctx context.Context, t *Timetable) error {
	return c.inner.Create(ctx, t)
}
func (c *CachedRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return c.inner.Delete(ctx, id)
}
func (c *CachedRepository) GetByBatch(ctx context.Context, batchID uuid.UUID) ([]Timetable, error) {
	return c.inner.GetByBatch(ctx, batchID)
}
