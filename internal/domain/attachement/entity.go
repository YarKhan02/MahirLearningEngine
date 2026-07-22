package attachement

import (
	"time"

	"github.com/google/uuid"
)

type Presign struct {
	Filename     string
	ContentType  string
	SizeBytes    int64
	ResourceType string
	ResourceID   string
}

type PresignURL struct {
	URL       string
	Key       string
	ExpiresIn int
}

type Attachment struct {
	ID           uuid.UUID
	Key          string // r2_key
	Filename     string
	ContentType  string
	SizeBytes    *int64
	ResourceType string
	ResourceID   string
	UploadedBy          uuid.UUID
	Status              string
	VerifiedContentType string
	CreatedAt           time.Time
	ConfirmedAt         *time.Time
	DeletedAt           *time.Time

	DownloadURL string
}
