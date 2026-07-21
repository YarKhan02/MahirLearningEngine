package attachement

import (
	"time"

	"github.com/google/uuid"
)

func ToPresignRequestEntity(req PresignRequest) Presign {
	return Presign{
		Filename:     req.Filename,
		ContentType:  req.ContentType,
		SizeBytes:    req.SizeBytes,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
	}
}

func ToPresignResponseDTO(resp *PresignURL) *PresignResponse {
	if resp == nil {
		return nil
	}

	return &PresignResponse{
		UploadURL: resp.URL,
		Key:       resp.Key,
		ExpiresIn: resp.ExpiresIn,
	}
}

func ToAttachmentResponse(a Attachment) AttachmentResponse {
	return AttachmentResponse{
		ID:          a.ID.String(),
		FileName:    a.Filename,
		ContentType: a.ContentType,
		SizeBytes:   a.SizeBytes,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		DownloadURL: a.DownloadURL,
	}
}

func newPendingAttachment(p Presign, key string, uploadedBy uuid.UUID) Attachment {
	size := p.SizeBytes
	return Attachment{
		ID:           uuid.New(),
		Key:          key,
		Filename:     p.Filename,
		ContentType:  p.ContentType,
		SizeBytes:    &size,
		ResourceType: p.ResourceType,
		ResourceID:   p.ResourceID,
		UploadedBy:   uploadedBy,
		Status:       "pending",
	}
}
