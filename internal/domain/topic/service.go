package topic

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrLessonNotFound = errors.New("lesson not found")
	ErrTopicNotFound  = errors.New("topic not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidOrderNo = errors.New("order_no out of range")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTopic(ctx context.Context, lessonID uuid.UUID, req InsertTopicRequest) error {
	ok, err := s.repo.LessonExists(ctx, lessonID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrLessonNotFound
	}
	return s.repo.InsertTopic(ctx, ToTopic(req, lessonID))
}

func (s *Service) GetTopics(ctx context.Context, lessonID uuid.UUID) ([]Topic, error) {
	return s.repo.GetTopicsByLesson(ctx, lessonID)
}

// GetTopicsForStudent enforces that the student can access the lesson's course.
func (s *Service) GetTopicsForStudent(ctx context.Context, userID, lessonID uuid.UUID) ([]Topic, error) {
	ok, err := s.repo.UserHasLessonAccess(ctx, userID, lessonID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}
	return s.repo.GetTopicsByLesson(ctx, lessonID)
}

func (s *Service) UpdateTopic(ctx context.Context, req UpdateTopic) error {
	ok, err := s.repo.TopicExists(ctx, req.ID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrTopicNotFound
	}
	return s.repo.UpdateTopic(ctx, req)
}

func (s *Service) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTopic(ctx, id)
}

func (s *Service) ReorderTopic(ctx context.Context, topicID uuid.UUID, orderNo int) error {
	ok, err := s.repo.TopicExists(ctx, topicID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrTopicNotFound
	}
	return s.repo.ReorderTopic(ctx, topicID, orderNo)
}
