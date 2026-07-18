package batch

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	batchCacheVersion = "v1"
	batchCacheName    = "batch"
	batchListTTL      = 5 * time.Minute
	batchCoursesTTL   = 5 * time.Minute
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func batchListKey() string {
	return "batch:" + batchCacheVersion + ":all"
}

func batchCoursesKey(id uuid.UUID) string {
	return "batch:" + batchCacheVersion + ":courses:" + id.String()
}

// cached reads

func (c *CachedRepository) GetBatches(ctx context.Context) ([]Batch, error) {
	key := batchListKey()
	if v, ok := cache.GetJSON[[]Batch](ctx, c.cache, batchCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetBatches(ctx)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, batchListTTL)
	return list, nil
}

func (c *CachedRepository) GetBatchCourses(ctx context.Context, batchID uuid.UUID) ([]BatchCourse, error) {
	key := batchCoursesKey(batchID)
	if v, ok := cache.GetJSON[[]BatchCourse](ctx, c.cache, batchCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetBatchCourses(ctx, batchID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, batchCoursesTTL)
	return list, nil
}

// writes (invalidate)

func (c *CachedRepository) CreateBatch(ctx context.Context, req *Batch) error {
	if err := c.inner.CreateBatch(ctx, req); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, batchListKey())
	return nil
}

func (c *CachedRepository) UpdateBatch(ctx context.Context, req *Batch) error {
	if err := c.inner.UpdateBatch(ctx, req); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, batchListKey())
	return nil
}

func (c *CachedRepository) DeleteBatch(ctx context.Context, id uuid.UUID) error {
	if err := c.inner.DeleteBatch(ctx, id); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, batchListKey(), batchCoursesKey(id))
	return nil
}

func (c *CachedRepository) UpdateBatchCourses(ctx context.Context, batchID uuid.UUID, add []uuid.UUID, remove []uuid.UUID, grantedBy *uuid.UUID) error {
	if err := c.inner.UpdateBatchCourses(ctx, batchID, add, remove, grantedBy); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, batchCoursesKey(batchID))
	return nil
}
