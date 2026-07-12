package student

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	RegisterStudent(ctx context.Context, s *Student, batchID uuid.UUID) error
	GetStudents(ctx context.Context, q string) ([]StudentWithBatch, error)
	GetStudentByID(ctx context.Context, id uuid.UUID) (*Student, error)
	UpdateStudentStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateStudentBatch(ctx context.Context, studentID uuid.UUID, batchID *uuid.UUID) error
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetStudentCourses(ctx context.Context, userID uuid.UUID) ([]StudentCourse, error)
	HasCourseAccess(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (bool, error)
	GetStudentLessons(ctx context.Context, courseID uuid.UUID, studentID uuid.UUID) ([]StudentLesson, error)
	SetLessonProgress(ctx context.Context, studentID uuid.UUID, lessonID uuid.UUID, completed bool) error
}
