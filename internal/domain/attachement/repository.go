package attachement

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a Attachment) error
	GetPendingByKey(ctx context.Context, key string, uploadedBy uuid.UUID) (Attachment, error)
	ConfirmByKey(ctx context.Context, key string, uploadedBy uuid.UUID, sizeBytes int64, verifiedContentType string) (Attachment, error)
	ListByResource(ctx context.Context, resourceType, resourceID string) ([]Attachment, error)
	GetByID(ctx context.Context, id uuid.UUID) (Attachment, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	CourseExists(ctx context.Context, courseID string) (bool, error)
	UserHasCourseAccess(ctx context.Context, userID uuid.UUID, courseID string) (bool, error)
}
