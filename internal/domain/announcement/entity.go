package announcement

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID          uuid.UUID
	BatchID     uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time

	// Populated on reads for display.
	BatchName string
}
