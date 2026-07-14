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

func (s *Service) GetStudentRecords(ctx context.Context, studentID uuid.UUID) ([]Record, error) {
	return s.repo.GetStudentRecords(ctx, studentID)
}

// GetMyRecords returns the logged-in student's attendance history.
func (s *Service) GetMyRecords(ctx context.Context, userID uuid.UUID) ([]Record, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetStudentRecords(ctx, studentID)
}
