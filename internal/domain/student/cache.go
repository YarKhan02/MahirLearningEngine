package student

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	studentCacheVersion = "v1"
	studentCacheName    = "student"
	studentPortalTTL = 60 * time.Second
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func studentCoursesKey(userID uuid.UUID) string {
	return "student:" + studentCacheVersion + ":courses:" + userID.String()
}

func studentLessonsKey(studentID, courseID uuid.UUID) string {
	return "student:" + studentCacheVersion + ":lessons:" + studentID.String() + ":" + courseID.String()
}

// cached reads

func (c *CachedRepository) GetStudentCourses(ctx context.Context, userID uuid.UUID) ([]StudentCourse, error) {
	key := studentCoursesKey(userID)
	if v, ok := cache.GetJSON[[]StudentCourse](ctx, c.cache, studentCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetStudentCourses(ctx, userID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, studentPortalTTL)
	return list, nil
}

func (c *CachedRepository) GetStudentLessons(ctx context.Context, courseID, studentID uuid.UUID) ([]StudentLesson, error) {
	key := studentLessonsKey(studentID, courseID)
	if v, ok := cache.GetJSON[[]StudentLesson](ctx, c.cache, studentCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetStudentLessons(ctx, courseID, studentID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, studentPortalTTL)
	return list, nil
}

// pass-throughs

func (c *CachedRepository) RegisterStudent(ctx context.Context, s *Student, batchID uuid.UUID) error {
	return c.inner.RegisterStudent(ctx, s, batchID)
}
func (c *CachedRepository) GetStudents(ctx context.Context, q string) ([]StudentWithBatch, error) {
	return c.inner.GetStudents(ctx, q)
}
func (c *CachedRepository) GetStudentByID(ctx context.Context, id uuid.UUID) (*Student, error) {
	return c.inner.GetStudentByID(ctx, id)
}
func (c *CachedRepository) UpdateStudentStatus(ctx context.Context, id uuid.UUID, status string) error {
	return c.inner.UpdateStudentStatus(ctx, id, status)
}
func (c *CachedRepository) UpdateStudentBatch(ctx context.Context, studentID uuid.UUID, batchID *uuid.UUID) error {
	return c.inner.UpdateStudentBatch(ctx, studentID, batchID)
}
func (c *CachedRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return c.inner.GetStudentIDByUserID(ctx, userID)
}
func (c *CachedRepository) HasCourseAccess(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (bool, error) {
	return c.inner.HasCourseAccess(ctx, studentID, courseID)
}
func (c *CachedRepository) SetLessonProgress(ctx context.Context, studentID uuid.UUID, lessonID uuid.UUID, completed bool) error {
	return c.inner.SetLessonProgress(ctx, studentID, lessonID, completed)
}
