package course

import (
	"context"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) InsertCourse(ctx context.Context, req dto.InsertCourse) (*Course, error) {
	return s.repo.InsertCourse(ctx, req)
}

func (s *Service) GetCourse(ctx context.Context) ([]Course, error) {
	return s.repo.GetCourse(ctx)
}