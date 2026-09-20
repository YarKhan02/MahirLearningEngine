package maintenance

import (
	"context"

	"github.com/google/uuid"
)

// Report summarizes what a cleanup run did (small — safe to return in the HTTP body).
type Report struct {
	ExpiredTokensDeleted    int64 `json:"expiredTokensDeleted"`
	PendingAttachmentsSwept int   `json:"pendingAttachmentsSwept"`
	R2DeleteFailures        int   `json:"r2DeleteFailures"`
}

// StalePendingAttachment is an abandoned two-phase upload (row + R2 object) to reclaim.
type StalePendingAttachment struct {
	ID    uuid.UUID
	R2Key string
}

type Repository interface {
	Ping(ctx context.Context) error
	PurgeExpiredTokens(ctx context.Context) (int64, error)
	ListStalePendingAttachments(ctx context.Context) ([]StalePendingAttachment, error)
	DeleteAttachmentByID(ctx context.Context, id uuid.UUID) error
}

// ObjectStore is the subset of the R2 client the cleanup needs.
type ObjectStore interface {
	DeleteObject(ctx context.Context, key string) error
}
