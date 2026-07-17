package announcement

import (
	"context"
	"encoding/json"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/google/uuid"
)

const (
	cacheVersion 	= "v1"
	allTTL			= 2 * time.Minute
	userTTL			= 2 * time.Minute
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

func (c *CachedRepository) getList(ctx context.Context, key string) ([]Announcement, bool) {
	
	raw, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, false
	}
	var list []Announcement
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil, false
	}
	return list, true
}

func (c *CachedRepository) setList(ctx context.Context, key string, list []Announcement, ttl time.Duration) {
	
	b, err := json.Marshal(list)
	if err != nil {
		return
	}
	_ = c.cache.Set(ctx, key, string(b), ttl)
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
	if list, ok := c.getList(ctx, key); ok {
		return list, nil
	}
	list, err := c.inner.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	c.setList(ctx, key, list, allTTL)
	return list, nil
}

func (c *CachedRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]Announcement, error) {
	
	key := keyUser(userID)
	if list, ok := c.getList(ctx, key); ok {
		return list, nil
	}
	list, err := c.inner.GetForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	c.setList(ctx, key, list, userTTL)
	return list, nil
}