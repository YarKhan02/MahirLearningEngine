package attachement

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a Attachment) error
	ConfirmByKey(ctx context.Context, key string, uploadedBy uuid.UUID, sizeBytes int64) (Attachment, error)
	ListByResource(ctx context.Context, resourceType, resourceID string) ([]Attachment, error)
	GetByID(ctx context.Context, id uuid.UUID) (Attachment, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	UserHasCourseAccess(ctx context.Context, userID uuid.UUID, courseID string) (bool, error)
}
