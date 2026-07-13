package batch

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrBatchNotFound = errors.New("batch not found")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBatch(ctx context.Context, req *Batch) error {
	return s.repo.CreateBatch(ctx, req)
}

func (s *Service) UpdateBatch(ctx context.Context, req *Batch) error {
	return s.repo.UpdateBatch(ctx, req)
}

func (s *Service) DeleteBatch(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteBatch(ctx, id)
}

func (s *Service) GetBatches(ctx context.Context) ([]Batch, error) {
	return s.repo.GetBatches(ctx)
}
func (s *Service) GetBatchCourses(ctx context.Context, batchID uuid.UUID) ([]BatchCourse, error) {
	return s.repo.GetBatchCourses(ctx, batchID)
}

func (s *Service) UpdateBatchCourses(ctx context.Context, batchID uuid.UUID, add []uuid.UUID, remove []uuid.UUID, grantedBy *uuid.UUID) error {
	if len(add) == 0 && len(remove) == 0 {
		return nil
	}
	return s.repo.UpdateBatchCourses(ctx, batchID, add, remove, grantedBy)
}

func (s *Service) GetOpenBatchesWithCourses(ctx context.Context) ([]BatchWithCourses, error) {
	batches, err := s.repo.GetBatches(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]BatchWithCourses, 0, len(batches))
	for _, b := range batches {
		if b.Status == "completed" {
			continue
		}

		courses, err := s.repo.GetBatchCourses(ctx, b.ID)
		if err != nil {
			return nil, err
		}

		out = append(out, BatchWithCourses{Batch: b, Courses: courses})
	}

	return out, nil
}
