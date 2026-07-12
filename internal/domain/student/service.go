package student

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyRegistered = errors.New("this email is already registered")
	ErrStudentNotFound        = errors.New("student not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterStudent(ctx context.Context, req *Student, batchID uuid.UUID) error {
	return s.repo.RegisterStudent(ctx, req, batchID)
}

func (s *Service) GetStudents(ctx context.Context, q string) ([]StudentWithBatch, error) {
	return s.repo.GetStudents(ctx, q)
}

func (s *Service) GetStudentByID(ctx context.Context, id uuid.UUID) (*Student, error) {
	return s.repo.GetStudentByID(ctx, id)
}

func (s *Service) UpdateStudentStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateStudentStatus(ctx, id, status)
}

func (s *Service) UpdateStudentBatch(ctx context.Context, studentID uuid.UUID, batchID *uuid.UUID) error {
	return s.repo.UpdateStudentBatch(ctx, studentID, batchID)
}
