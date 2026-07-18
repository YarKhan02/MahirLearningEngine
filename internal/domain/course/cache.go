package course

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	courseCacheVersion = "v1"
	courseCacheName    = "course"
	courseListTTL      = 5 * time.Minute
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func courseListKey() string {
	return "course:" + courseCacheVersion + ":all"
}

// cached read

func (c *CachedRepository) GetCourse(ctx context.Context) ([]Course, error) {
	key := courseListKey()
	if v, ok := cache.GetJSON[[]Course](ctx, c.cache, courseCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetCourse(ctx)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, courseListTTL)
	return list, nil
}

// writes that invalidate the list

func (c *CachedRepository) InsertCourse(ctx context.Context, req Course) (*Course, error) {
	created, err := c.inner.InsertCourse(ctx, req)
	if err != nil {
		return nil, err
	}
	_ = c.cache.Delete(ctx, courseListKey())
	return created, nil
}

func (c *CachedRepository) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	if err := c.inner.DeleteCourse(ctx, id); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, courseListKey())
	return nil
}

// pass-throughs (uncached)

func (c *CachedRepository) GetLesson(ctx context.Context, id uuid.UUID) ([]Lesson, error) {
	return c.inner.GetLesson(ctx, id)
}
func (c *CachedRepository) CourseExists(ctx context.Context, id uuid.UUID) bool {
	return c.inner.CourseExists(ctx, id)
}
func (c *CachedRepository) InsertLesson(ctx context.Context, req Lesson) error {
	return c.inner.InsertLesson(ctx, req)
}
func (c *CachedRepository) LessonExists(ctx context.Context, id uuid.UUID) bool {
	return c.inner.LessonExists(ctx, id)
}
func (c *CachedRepository) UpdateLesson(ctx context.Context, req UpdateLesson) error {
	return c.inner.UpdateLesson(ctx, req)
}
func (c *CachedRepository) ReorderLesson(ctx context.Context, lessonID uuid.UUID, orderNo int) error {
	return c.inner.ReorderLesson(ctx, lessonID, orderNo)
}
