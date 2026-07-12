package assignment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAssignment(ctx context.Context, a *Assignment) error
	GetLessonAssignments(ctx context.Context, lessonID uuid.UUID) ([]Assignment, error)
	DeleteAssignment(ctx context.Context, id uuid.UUID) error
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	HasLessonAccess(ctx context.Context, studentID uuid.UUID, lessonID uuid.UUID) (bool, error)
	HasAssignmentAccess(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID) (bool, error)
	GetStudentAssignments(ctx context.Context, lessonID uuid.UUID, studentID uuid.UUID) ([]StudentAssignment, error)
	SubmitAssignment(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID, code string) error
	GetBatchSubmissions(ctx context.Context, batchID uuid.UUID) ([]BatchSubmission, error)
	GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error
	GetStudentSubmissions(ctx context.Context, studentID uuid.UUID) ([]BatchSubmission, error)
}
