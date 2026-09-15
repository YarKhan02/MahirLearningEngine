package program

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p Program) error
	Update(ctx context.Context, id uuid.UUID, u Upsert) (Program, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetPublished(ctx context.Context, id uuid.UUID, published bool) (Program, error)
	Reorder(ctx context.Context, ids []uuid.UUID) error

	// Admin reads (all programs, ordered).
	ListAll(ctx context.Context) ([]Program, error)
	GetByID(ctx context.Context, id uuid.UUID) (Program, error)

	// Public reads.
	ListPublished(ctx context.Context) ([]Program, error)
	GetPublishedBySlug(ctx context.Context, slug string) (Program, error)

	// Helpers.
	SlugExists(ctx context.Context, slug string) (bool, error)
	MaxOrderNo(ctx context.Context) (int, error)
}
