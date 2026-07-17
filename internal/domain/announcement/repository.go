package announcement

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a *Announcement) error
	GetAll(ctx context.Context) ([]Announcement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Announcement, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetForUser(ctx context.Context, userID uuid.UUID) ([]Announcement, error)
}
