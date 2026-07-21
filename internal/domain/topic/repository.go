package topic

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	LessonExists(ctx context.Context, lessonID uuid.UUID) (bool, error)
	InsertTopic(ctx context.Context, t Topic) error
	GetTopicsByLesson(ctx context.Context, lessonID uuid.UUID) ([]Topic, error)
	TopicExists(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateTopic(ctx context.Context, t UpdateTopic) error
	DeleteTopic(ctx context.Context, id uuid.UUID) error
	ReorderTopic(ctx context.Context, topicID uuid.UUID, orderNo int) error
	UserHasLessonAccess(ctx context.Context, userID, lessonID uuid.UUID) (bool, error)
}
