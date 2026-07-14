package course

import (
	"context"
	
	"github.com/google/uuid"
)

type Repository interface {
	InsertCourse(ctx context.Context, req Course) (*Course, error)
	GetCourse(ctx context.Context) ([]Course, error)
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	CourseExists(ctx context.Context, id uuid.UUID) bool
	InsertLesson(ctx context.Context, req Lesson) error
	GetLesson(ctx context.Context, id uuid.UUID) ([]Lesson, error)
	LessonExists(ctx context.Context, id uuid.UUID) bool
	UpdateLesson(ctx context.Context, req UpdateLesson) error
	ReorderLesson(ctx context.Context, lessonID uuid.UUID, orderNo int) error
}