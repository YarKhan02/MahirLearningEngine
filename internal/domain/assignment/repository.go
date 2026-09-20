package assignment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAssignment(ctx context.Context, a *Assignment) error
	GetLessonAssignments(ctx context.Context, lessonID uuid.UUID) ([]Assignment, error)
	GetAssignmentWithTests(ctx context.Context, id uuid.UUID) (*Assignment, error)
	UpdateAssignment(ctx context.Context, a *Assignment) error
	DeleteAssignment(ctx context.Context, id uuid.UUID) error
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	HasLessonAccess(ctx context.Context, studentID uuid.UUID, lessonID uuid.UUID) (bool, error)
	HasAssignmentAccess(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID) (bool, error)
	GetStudentAssignments(ctx context.Context, lessonID uuid.UUID, studentID uuid.UUID) ([]StudentAssignment, error)
	SubmitAssignment(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID, code string) (uuid.UUID, error)

	// Autograde (code assignments).
	GetGradingData(ctx context.Context, assignmentID uuid.UUID) (language string, totalMarks int, tests []TestCase, err error)
	SaveAutoGrade(ctx context.Context, submissionID uuid.UUID, autoScore, testsTotal, testsPassed int, results []TestResult) error
	GetSubmissionResults(ctx context.Context, submissionID uuid.UUID) ([]TestResult, error)
	GetBatchSubmissions(ctx context.Context, batchID uuid.UUID, q, status string, limit, offset int) ([]BatchSubmission, error)
	CountBatchSubmissions(ctx context.Context, batchID uuid.UUID, q, status string) (int, error)
	GetBatchSubmissionSummary(ctx context.Context, batchID uuid.UUID, q string) (SubmissionSummary, error)
	GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error
	GetStudentSubmissions(ctx context.Context, studentID uuid.UUID, status string, limit, offset int) ([]BatchSubmission, error)
	CountStudentSubmissions(ctx context.Context, studentID uuid.UUID, status string) (int, error)
	GetStudentSubmissionSummary(ctx context.Context, studentID uuid.UUID) (SubmissionSummary, error)
}
