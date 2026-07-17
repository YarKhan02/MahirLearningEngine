package announcement

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func ToCreateAnnouncement(req CreateAnnouncementRequest) (*Announcement, error) {
	batchID, err := uuid.Parse(req.BatchID)
	if err != nil {
		return nil, fmt.Errorf("invalid batchId: %w", err)
	}

	return &Announcement{
		BatchID:     batchID,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

func ToAnnouncementResponse(a Announcement) AnnouncementResponse {
	return AnnouncementResponse{
		ID:          a.ID.String(),
		BatchID:     a.BatchID.String(),
		BatchName:   a.BatchName,
		Title:       a.Title,
		Description: a.Description,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
	}
}
