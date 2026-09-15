package program

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

var slugStrip = regexp.MustCompile(`[^a-z0-9]+`)

// slugify turns a title into a URL-safe slug.
func slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = slugStrip.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "program"
	}
	return s
}

// uniqueSlug appends -2, -3, … until the slug is free.
func (s *Service) uniqueSlug(ctx context.Context, base string) (string, error) {
	slug := base
	for i := 2; ; i++ {
		exists, err := s.repo.SlugExists(ctx, slug)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}

func (s *Service) Create(ctx context.Context, u Upsert) (Program, error) {
	log := logging.FromLogger(ctx)

	if strings.TrimSpace(u.Title) == "" {
		return Program{}, ErrInvalid
	}

	slug, err := s.uniqueSlug(ctx, slugify(u.Title))
	if err != nil {
		return Program{}, err
	}

	maxOrder, err := s.repo.MaxOrderNo(ctx)
	if err != nil {
		return Program{}, err
	}

	if u.Learn == nil {
		u.Learn = []string{}
	}

	p := Program{
		ID:        uuid.New(),
		Slug:      slug,
		Title:     u.Title,
		Short:     u.Short,
		Full:      u.Full,
		Learn:     u.Learn,
		Age:       u.Age,
		Fee:       u.Fee,
		Level:     u.Level,
		Bonus:     u.Bonus,
		Accent:    u.Accent,
		ImageURL:  u.ImageURL,
		OrderNo:   maxOrder + 1,
		Published: u.Published,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return Program{}, err
	}

	log.Info("program created",
		zap.String("event", "program_created"),
		zap.String("program_id", p.ID.String()),
		zap.String("slug", p.Slug),
		zap.Bool("published", p.Published),
	)
	return p, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, u Upsert) (Program, error) {
	if strings.TrimSpace(u.Title) == "" {
		return Program{}, ErrInvalid
	}
	if u.Learn == nil {
		u.Learn = []string{}
	}
	p, err := s.repo.Update(ctx, id, u)
	if err != nil {
		return Program{}, err
	}
	logging.FromLogger(ctx).Info("program updated",
		zap.String("event", "program_updated"),
		zap.String("program_id", p.ID.String()),
	)
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	logging.FromLogger(ctx).Info("program deleted",
		zap.String("event", "program_deleted"),
		zap.String("program_id", id.String()),
	)
	return nil
}

func (s *Service) SetPublished(ctx context.Context, id uuid.UUID, published bool) (Program, error) {
	return s.repo.SetPublished(ctx, id, published)
}

func (s *Service) Reorder(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.Reorder(ctx, ids)
}

func (s *Service) ListAll(ctx context.Context) ([]Program, error) {
	return s.repo.ListAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (Program, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListPublished(ctx context.Context) ([]Program, error) {
	return s.repo.ListPublished(ctx)
}

func (s *Service) GetPublishedBySlug(ctx context.Context, slug string) (Program, error) {
	return s.repo.GetPublishedBySlug(ctx, slug)
}
