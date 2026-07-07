package role

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddRoleToUser(ctx context.Context, userID uuid.UUID, role string) error
	GetUserRole(ctx context.Context, userID uuid.UUID) (string, error)
}