package course

import (
	"context"
	"errors"
	
	"github.com/google/uuid"
)

var (
	ErrCourseNotFound       	= errors.New("course not found")
	ErrLessonNotFound       	= errors.New("lesson not found")
	ErrLessonInsertFailed       = errors.New("lesson insertion failed")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) InsertCourse(ctx context.Context, req Course) (*Course, error) {
	return s.repo.InsertCourse(ctx, req)
}

func (s *Service) GetCourse(ctx context.Context) ([]Course, error) {
	return s.repo.GetCourse(ctx)
}

func (s *Service) InsertLesson(ctx context.Context, req Lesson) error {

	isCourseID := s.repo.CourseExists(ctx, req.CourseID)
	if !isCourseID {
		return ErrCourseNotFound
	}

	err := s.repo.InsertLesson(ctx, req)
	if err != nil {
		return ErrLessonInsertFailed
	}

	return nil
}

func (s *Service) GetLesson(ctx context.Context, id uuid.UUID) ([]Lesson, error) {
	
	isCourseID := s.repo.CourseExists(ctx, id)
	if !isCourseID {
		return nil, ErrCourseNotFound
	}

	return s.repo.GetLesson(ctx, id)
}

func (s *Service) UpdateLesson(ctx context.Context, req UpdateLesson) error {
	
	isCourseID := s.repo.CourseExists(ctx, req.CourseID)
	if !isCourseID {
		return ErrCourseNotFound
	}

	isLessonID := s.repo.LessonExists(ctx, req.ID)
	if !isLessonID {
		return ErrLessonNotFound
	}

	return s.repo.UpdateLesson(ctx, req)
}

func (s *Service) ReorderLesson(ctx context.Context, lessonID uuid.UUID, orderNo int) error {

	isLessonID := s.repo.LessonExists(ctx, lessonID)
	if !isLessonID {
		return ErrLessonNotFound
	}

	return s.repo.ReorderLesson(ctx, lessonID, orderNo)
}