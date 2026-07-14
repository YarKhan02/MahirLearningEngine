package mapper

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/announcement"
	"github.com/google/uuid"
)

func ToCreateAnnouncement(req dto.CreateAnnouncementRequest) (*announcement.Announcement, error) {
	batchID, err := uuid.Parse(req.BatchID)
	if err != nil {
		return nil, fmt.Errorf("invalid batchId: %w", err)
	}

	return &announcement.Announcement{
		BatchID:     batchID,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

func ToAnnouncementResponse(a announcement.Announcement) dto.AnnouncementResponse {
	return dto.AnnouncementResponse{
		ID:          a.ID.String(),
		BatchID:     a.BatchID.String(),
		BatchName:   a.BatchName,
		Title:       a.Title,
		Description: a.Description,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
	}
}
