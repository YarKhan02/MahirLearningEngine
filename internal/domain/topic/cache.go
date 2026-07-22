package topic

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	topicCacheVersion = "v1"
	topicCacheName    = "topic"
	topicListTTL      = 5 * time.Minute
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func topicListKey(lessonID uuid.UUID) string {
	return "topic:" + topicCacheVersion + ":lesson:" + lessonID.String()
}

func (c *CachedRepository) invalidateByTopic(ctx context.Context, topicID uuid.UUID) {
	if lessonID, err := c.inner.GetTopicLessonID(ctx, topicID); err == nil {
		_ = c.cache.Delete(ctx, topicListKey(lessonID))
	}
}

func (c *CachedRepository) GetTopicsByLesson(ctx context.Context, lessonID uuid.UUID) ([]Topic, error) {
	key := topicListKey(lessonID)
	if v, ok := cache.GetJSON[[]Topic](ctx, c.cache, topicCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetTopicsByLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, topicListTTL)
	return list, nil
}

func (c *CachedRepository) InsertTopic(ctx context.Context, t Topic) error {
	if err := c.inner.InsertTopic(ctx, t); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, topicListKey(t.LessonID))
	return nil
}

func (c *CachedRepository) UpdateTopic(ctx context.Context, t UpdateTopic) error {
	if err := c.inner.UpdateTopic(ctx, t); err != nil {
		return err
	}
	c.invalidateByTopic(ctx, t.ID)
	return nil
}

func (c *CachedRepository) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	lessonID, lookupErr := c.inner.GetTopicLessonID(ctx, id)
	if err := c.inner.DeleteTopic(ctx, id); err != nil {
		return err
	}
	if lookupErr == nil {
		_ = c.cache.Delete(ctx, topicListKey(lessonID))
	}
	return nil
}

func (c *CachedRepository) ReorderTopic(ctx context.Context, topicID uuid.UUID, orderNo int) error {
	if err := c.inner.ReorderTopic(ctx, topicID, orderNo); err != nil {
		return err
	}
	c.invalidateByTopic(ctx, topicID)
	return nil
}

func (c *CachedRepository) LessonExists(ctx context.Context, lessonID uuid.UUID) (bool, error) {
	return c.inner.LessonExists(ctx, lessonID)
}

func (c *CachedRepository) GetTopicLessonID(ctx context.Context, topicID uuid.UUID) (uuid.UUID, error) {
	return c.inner.GetTopicLessonID(ctx, topicID)
}

func (c *CachedRepository) TopicExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return c.inner.TopicExists(ctx, id)
}

func (c *CachedRepository) UserHasLessonAccess(ctx context.Context, userID, lessonID uuid.UUID) (bool, error) {
	return c.inner.UserHasLessonAccess(ctx, userID, lessonID)
}
