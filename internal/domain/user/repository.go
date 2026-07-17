package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByEmailExists(ctx context.Context, email string) (bool, error)
	FindByIDExists(ctx context.Context, id uuid.UUID) (bool, error)
	FindByUsernameExists(ctx context.Context, username string) (bool, error)
	ResetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	// FindByLoginIdentifier matches a user by either email or username.
	FindByLoginIdentifier(ctx context.Context, identifier string) (*User, error)
	// Update(ctx context.Context, u *User) error
	UpdateFailedAttempts(ctx context.Context, id uuid.UUID, attempts int, lockedUntil *time.Time) error
	// Ban(ctx context.Context, id uuid.UUID) error
	// Unban(ctx context.Context, id uuid.UUID) error
}