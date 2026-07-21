package quiz

import (
	"context"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/cache"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

const (
	quizCacheVersion = "v1"
	quizCacheName    = "quiz"
	quizListTTL      = 5 * time.Minute
)

type CachedRepository struct {
	inner Repository
	cache *redis.RedisClient
}

func NewCachedRepository(inner Repository, c *redis.RedisClient) *CachedRepository {
	return &CachedRepository{inner: inner, cache: c}
}

func quizListKey(lessonID uuid.UUID) string {
	return "quiz:" + quizCacheVersion + ":lesson:" + lessonID.String()
}

func (c *CachedRepository) GetQuizzesByLesson(ctx context.Context, lessonID uuid.UUID) ([]Quiz, error) {
	key := quizListKey(lessonID)
	if v, ok := cache.GetJSON[[]Quiz](ctx, c.cache, quizCacheName, key); ok {
		return v, nil
	}
	list, err := c.inner.GetQuizzesByLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, c.cache, key, list, quizListTTL)
	return list, nil
}

func (c *CachedRepository) CreateQuiz(ctx context.Context, q NewQuiz) error {
	if err := c.inner.CreateQuiz(ctx, q); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, quizListKey(q.LessonID))
	return nil
}

func (c *CachedRepository) EditQuiz(ctx context.Context, quizID uuid.UUID, q NewQuiz) error {
	lessonID, lookupErr := c.inner.GetQuizLessonID(ctx, quizID)
	if err := c.inner.EditQuiz(ctx, quizID, q); err != nil {
		return err
	}
	if lookupErr == nil {
		_ = c.cache.Delete(ctx, quizListKey(lessonID))
	}
	return nil
}

func (c *CachedRepository) DeleteQuiz(ctx context.Context, quizID uuid.UUID) error {
	lessonID, lookupErr := c.inner.GetQuizLessonID(ctx, quizID)
	if err := c.inner.DeleteQuiz(ctx, quizID); err != nil {
		return err
	}
	if lookupErr == nil {
		_ = c.cache.Delete(ctx, quizListKey(lessonID))
	}
	return nil
}

// pass-throughs

func (c *CachedRepository) LessonExists(ctx context.Context, lessonID uuid.UUID) (bool, error) {
	return c.inner.LessonExists(ctx, lessonID)
}

func (c *CachedRepository) QuizExists(ctx context.Context, quizID uuid.UUID) (bool, error) {
	return c.inner.QuizExists(ctx, quizID)
}

func (c *CachedRepository) GetQuizLessonID(ctx context.Context, quizID uuid.UUID) (uuid.UUID, error) {
	return c.inner.GetQuizLessonID(ctx, quizID)
}

func (c *CachedRepository) GetQuizByID(ctx context.Context, quizID uuid.UUID) (Quiz, error) {
	return c.inner.GetQuizByID(ctx, quizID)
}

func (c *CachedRepository) GetLessonSubmissionCounts(ctx context.Context, lessonID uuid.UUID) (map[uuid.UUID]SubmissionSummary, error) {
	return c.inner.GetLessonSubmissionCounts(ctx, lessonID)
}

func (c *CachedRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return c.inner.GetStudentIDByUserID(ctx, userID)
}

func (c *CachedRepository) UserHasLessonAccess(ctx context.Context, userID, lessonID uuid.UUID) (bool, error) {
	return c.inner.UserHasLessonAccess(ctx, userID, lessonID)
}

func (c *CachedRepository) GetStudentSubmissionsForLesson(ctx context.Context, lessonID, studentID uuid.UUID) (map[uuid.UUID]Submission, error) {
	return c.inner.GetStudentSubmissionsForLesson(ctx, lessonID, studentID)
}

func (c *CachedRepository) HasSubmitted(ctx context.Context, quizID, studentID uuid.UUID) (bool, error) {
	return c.inner.HasSubmitted(ctx, quizID, studentID)
}

func (c *CachedRepository) CreateSubmission(ctx context.Context, sub Submission) error {
	return c.inner.CreateSubmission(ctx, sub)
}

func (c *CachedRepository) GetSubmissionsByQuiz(ctx context.Context, quizID uuid.UUID) ([]SubmissionRow, error) {
	return c.inner.GetSubmissionsByQuiz(ctx, quizID)
}

func (c *CachedRepository) GetSubmissionSummary(ctx context.Context, quizID uuid.UUID) (SubmissionSummary, error) {
	return c.inner.GetSubmissionSummary(ctx, quizID)
}

func (c *CachedRepository) GetSubmission(ctx context.Context, submissionID uuid.UUID) (Submission, error) {
	return c.inner.GetSubmission(ctx, submissionID)
}

func (c *CachedRepository) GetStudentName(ctx context.Context, studentID uuid.UUID) (string, error) {
	return c.inner.GetStudentName(ctx, studentID)
}

func (c *CachedRepository) GradeSubmission(ctx context.Context, in GradeInput) error {
	return c.inner.GradeSubmission(ctx, in)
}
