package attendance

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidStatus   = errors.New("status must be present or absent")
	ErrStudentNotFound = errors.New("student not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetRoster(ctx context.Context, batchID uuid.UUID, date time.Time) ([]RosterEntry, error) {
	return s.repo.GetRoster(ctx, batchID, date)
}

func (s *Service) Mark(ctx context.Context, req MarkAttendance) error {
	if req.Status != "present" && req.Status != "absent" {
		return ErrInvalidStatus
	}
	return s.repo.Mark(ctx, req)
}

func (s *Service) GetStudentRecords(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]Record, int, error) {
	total, err := s.repo.CountStudentRecords(ctx, studentID)
	if err != nil {
		return nil, 0, err
	}
	records, err := s.repo.GetStudentRecords(ctx, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (s *Service) GetStudentSummary(ctx context.Context, studentID uuid.UUID) (Summary, error) {
	return s.repo.GetStudentSummary(ctx, studentID)
}

func (s *Service) GetMyRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Record, int, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return s.GetStudentRecords(ctx, studentID, limit, offset)
}

func (s *Service) GetMySummary(ctx context.Context, userID uuid.UUID) (Summary, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return Summary{}, err
	}
	return s.repo.GetStudentSummary(ctx, studentID)
}
