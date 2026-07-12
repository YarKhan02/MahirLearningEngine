package assignment

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrAccessDenied    = errors.New("you do not have access to this assignment")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAssignment(ctx context.Context, a *Assignment) error {
	return s.repo.CreateAssignment(ctx, a)
}

func (s *Service) GetLessonAssignments(ctx context.Context, lessonID uuid.UUID) ([]Assignment, error) {
	return s.repo.GetLessonAssignments(ctx, lessonID)
}

func (s *Service) DeleteAssignment(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAssignment(ctx, id)
}

// GetStudentAssignments lists a lesson's assignments with the student's submissions.
func (s *Service) GetStudentAssignments(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) ([]StudentAssignment, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hasAccess, err := s.repo.HasLessonAccess(ctx, studentID, lessonID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrAccessDenied
	}

	return s.repo.GetStudentAssignments(ctx, lessonID, studentID)
}

// SubmitAssignment stores the student's code — no execution, storage only.
func (s *Service) SubmitAssignment(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID, code string) error {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	hasAccess, err := s.repo.HasAssignmentAccess(ctx, studentID, assignmentID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrAccessDenied
	}

	return s.repo.SubmitAssignment(ctx, studentID, assignmentID, code)
}

func (s *Service) GetBatchSubmissions(ctx context.Context, batchID uuid.UUID) ([]BatchSubmission, error) {
	return s.repo.GetBatchSubmissions(ctx, batchID)
}

func (s *Service) GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error {
	return s.repo.GradeSubmission(ctx, submissionID, marks, remarks)
}

// GetMySubmissions lists everything the logged-in student has submitted.
func (s *Service) GetMySubmissions(ctx context.Context, userID uuid.UUID) ([]BatchSubmission, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetStudentSubmissions(ctx, studentID)
}
