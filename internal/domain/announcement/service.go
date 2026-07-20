package announcement

import (
	"context"
	"errors"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/pagination"
	"github.com/google/uuid"
)

var (
	ErrEmptyTitle       = errors.New("title is required")
	ErrEmptyDescription = errors.New("description is required")
	ErrNotFound         = errors.New("announcement not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, a *Announcement) error {
	a.Title = strings.TrimSpace(a.Title)
	a.Description = strings.TrimSpace(a.Description)
	if a.Title == "" {
		return ErrEmptyTitle
	}
	if a.Description == "" {
		return ErrEmptyDescription
	}
	return s.repo.Create(ctx, a)
}

func (s *Service) GetAll(ctx context.Context) ([]Announcement, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetAllPaged(ctx context.Context, p pagination.Params) ([]Announcement, int, error) {
	all, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return pagination.Slice(all, p), len(all), nil
}

func (s *Service) GetForUserPaged(ctx context.Context, userID uuid.UUID, p pagination.Params) ([]Announcement, int, error) {
	all, err := s.repo.GetForUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return pagination.Slice(all, p), len(all), nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Announcement, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetForUser(ctx context.Context, userID uuid.UUID) ([]Announcement, error) {
	return s.repo.GetForUser(ctx, userID)
}
