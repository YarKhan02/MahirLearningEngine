package dashboard

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
)

const (
	dashboardCacheVersion = "v1"
	dashboardCacheName    = "dashboard"
	dashboardTTL = 60 * time.Second
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func dashboardKey() string {
	return "dashboard:" + dashboardCacheVersion + ":admin"
}

func (c *CachedRepository) GetAdminDashboard(ctx context.Context) (*AdminDashboard, error) {
	key := dashboardKey()
	if v, ok := cache.GetJSON[AdminDashboard](ctx, c.cache, dashboardCacheName, key); ok {
		return &v, nil
	}
	d, err := c.inner.GetAdminDashboard(ctx)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, *d, dashboardTTL)
	return d, nil
}
