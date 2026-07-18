package announcement

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/google/uuid"
)

const (
	cacheVersion 			= "v1"
	announcementCacheName 	= "announcement"
	allTTL					= 2 * time.Minute
	userTTL					= 2 * time.Minute
)

type CachedRepository struct {
	inner	Repository
	cache	*redis.RedisClient
}

func NewCachedRepository(inner Repository, cache *redis.RedisClient) *CachedRepository {
	return &CachedRepository{
		inner: inner,
		cache: cache,
	}
}

func keyAll() string {
	return "announcements:" + cacheVersion + ":all"
}

func keyUser(id uuid.UUID) string {
	return "announcements:" + cacheVersion + ":user" + id.String()
}

func keyBatch(id uuid.UUID) string {
	return "announcements:" + cacheVersion + ":batch" + id.String()
}

// repository interface methods

func (c *CachedRepository) Create(ctx context.Context, a *Announcement) error {
	
	if err := c.inner.Create(ctx, a); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, keyAll())
	_ = c.cache.Delete(ctx, keyBatch(a.BatchID))
	return nil
}

func (c *CachedRepository) Delete(ctx context.Context, id uuid.UUID) error {
	
	existing, err := c.inner.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := c.inner.Delete(ctx, id); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, keyAll())
	_ = c.cache.Delete(ctx, keyBatch(existing.BatchID))
	return nil
}

func (c *CachedRepository) GetByID(ctx context.Context, id uuid.UUID) (*Announcement, error) {
	return c.inner.GetByID(ctx, id)
}

func (c *CachedRepository) GetAll(ctx context.Context) ([]Announcement, error) {
	
	key := keyAll()
	if v, ok := cache.GetJSON[[]Announcement](ctx, c.cache, announcementCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, allTTL)
	return list, nil
}

func (c *CachedRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]Announcement, error) {
	
	key := keyUser(userID)
	if v, ok := cache.GetJSON[[]Announcement](ctx, c.cache, announcementCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, userTTL)
	return list, nil
}