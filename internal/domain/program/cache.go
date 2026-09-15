package program

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	programCacheVersion = "v1"
	programCacheName    = "program"
	programListTTL      = 5 * time.Minute
)

func publishedKey() string {
	return "program:" + programCacheVersion + ":published"
}

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func (c *CachedRepository) invalidate(ctx context.Context) {
	_ = c.cache.Delete(ctx, publishedKey())
}


func (c *CachedRepository) ListPublished(ctx context.Context) ([]Program, error) {
	key := publishedKey()
	if v, ok := cache.GetJSON[[]Program](ctx, c.cache, programCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.ListPublished(ctx)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, programListTTL)
	return list, nil
}

func (c *CachedRepository) GetPublishedBySlug(ctx context.Context, slug string) (Program, error) {
	list, err := c.ListPublished(ctx)
	if err != nil {
		return Program{}, err
	}
	for _, p := range list {
		if p.Slug == slug {
			return p, nil
		}
	}
	return Program{}, ErrNotFound
}

func (c *CachedRepository) Create(ctx context.Context, p Program) error {
	if err := c.inner.Create(ctx, p); err != nil {
		return err
	}
	c.invalidate(ctx)
	return nil
}

func (c *CachedRepository) Update(ctx context.Context, id uuid.UUID, u Upsert) (Program, error) {
	p, err := c.inner.Update(ctx, id, u)
	if err != nil {
		return Program{}, err
	}
	c.invalidate(ctx)
	return p, nil
}

func (c *CachedRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.inner.Delete(ctx, id); err != nil {
		return err
	}
	c.invalidate(ctx)
	return nil
}

func (c *CachedRepository) SetPublished(ctx context.Context, id uuid.UUID, published bool) (Program, error) {
	p, err := c.inner.SetPublished(ctx, id, published)
	if err != nil {
		return Program{}, err
	}
	c.invalidate(ctx)
	return p, nil
}

func (c *CachedRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	if err := c.inner.Reorder(ctx, ids); err != nil {
		return err
	}
	c.invalidate(ctx)
	return nil
}

func (c *CachedRepository) ListAll(ctx context.Context) ([]Program, error) {
	return c.inner.ListAll(ctx)
}

func (c *CachedRepository) GetByID(ctx context.Context, id uuid.UUID) (Program, error) {
	return c.inner.GetByID(ctx, id)
}

func (c *CachedRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	return c.inner.SlugExists(ctx, slug)
}

func (c *CachedRepository) MaxOrderNo(ctx context.Context) (int, error) {
	return c.inner.MaxOrderNo(ctx)
}
