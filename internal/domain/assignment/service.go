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

func (s *Service) GetBatchSubmissions(ctx context.Context, batchID uuid.UUID, q, status string, limit, offset int) ([]BatchSubmission, int, error) {
	total, err := s.repo.CountBatchSubmissions(ctx, batchID, q, status)
	if err != nil {
		return nil, 0, err
	}
	items, err := s.repo.GetBatchSubmissions(ctx, batchID, q, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) GetBatchSubmissionSummary(ctx context.Context, batchID uuid.UUID, q string) (SubmissionSummary, error) {
	return s.repo.GetBatchSubmissionSummary(ctx, batchID, q)
}

func (s *Service) GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error {
	return s.repo.GradeSubmission(ctx, submissionID, marks, remarks)
}

func (s *Service) GetMySubmissions(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]BatchSubmission, int, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountStudentSubmissions(ctx, studentID, status)
	if err != nil {
		return nil, 0, err
	}
	items, err := s.repo.GetStudentSubmissions(ctx, studentID, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) GetMySubmissionSummary(ctx context.Context, userID uuid.UUID) (SubmissionSummary, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return SubmissionSummary{}, err
	}
	return s.repo.GetStudentSubmissionSummary(ctx, studentID)
}
