package attachement

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/r2"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrForbidden = errors.New("forbidden")
	ErrFailed    = errors.New("failed to presign")
	ErrNotFound  = errors.New("object not found")
)

const ResourceTypeCourse = "course"

type Service struct {
	r2   *r2.Client
	repo Repository
}

func NewService(repo Repository, r2 *r2.Client) *Service {
	return &Service{
		repo: repo,
		r2:   r2,
	}
}

func (s *Service) PresignUpload(ctx context.Context, userID uuid.UUID, req Presign) (*PresignURL, error) {
	log := logging.FromLogger(ctx)

	if req.ResourceType != ResourceTypeCourse {
		return nil, ErrFailed
	}
	if _, err := uuid.Parse(req.ResourceID); err != nil {
		return nil, ErrFailed
	}

	ext := filepath.Ext(req.Filename)
	key := fmt.Sprintf("%s/%s/%s%s", req.ResourceType, req.ResourceID, uuid.New().String(), ext)

	presigned, err := s.r2.Presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.r2.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(req.ContentType),
	})
	if err != nil {
		log.Error("presign upload failed",
			zap.String("event", "material_presign_failed"),
			zap.String("resource_id", req.ResourceID),
			zap.String("uploaded_by", userID.String()),
			zap.Error(err),
		)
		return nil, ErrFailed
	}

	a := newPendingAttachment(req, key, userID)
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}

	log.Info("material upload presigned",
		zap.String("event", "material_presigned"),
		zap.String("attachment_id", a.ID.String()),
		zap.String("resource_id", req.ResourceID),
		zap.String("uploaded_by", userID.String()),
		zap.String("content_type", req.ContentType),
		zap.Int64("size_bytes", req.SizeBytes),
		zap.String("r2_key", key),
	)

	return &PresignURL{URL: presigned.URL, Key: key, ExpiresIn: 300}, nil
}

func (s *Service) ConfirmUpload(ctx context.Context, userID uuid.UUID, key string) (*Attachment, error) {
	log := logging.FromLogger(ctx)

	head, err := s.r2.S3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.r2.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Warn("confirm upload: object not found in bucket",
			zap.String("event", "material_confirm_object_missing"),
			zap.String("uploaded_by", userID.String()),
			zap.String("r2_key", key),
		)
		return nil, ErrNotFound
	}

	var size int64
	if head.ContentLength != nil {
		size = *head.ContentLength
	}

	if size > MaxUploadSize {
		if delErr := s.r2.DeleteObject(ctx, key); delErr != nil {
			log.Warn("failed to delete oversized object",
				zap.String("event", "material_r2_delete_failed"),
				zap.String("r2_key", key),
				zap.Error(delErr),
			)
		}
		log.Warn("confirm upload: rejected oversized file",
			zap.String("event", "material_rejected_oversize"),
			zap.String("uploaded_by", userID.String()),
			zap.String("r2_key", key),
			zap.Int64("size_bytes", size),
		)
		return nil, ErrFailed
	}

	a, err := s.repo.ConfirmByKey(ctx, key, userID, size)
	if err != nil {
		return nil, err
	}

	log.Info("material uploaded",
		zap.String("event", "material_uploaded"),
		zap.String("attachment_id", a.ID.String()),
		zap.String("resource_id", a.ResourceID),
		zap.String("uploaded_by", userID.String()),
		zap.String("content_type", a.ContentType),
		zap.Int64("size_bytes", size),
	)

	return &a, nil
}

func (s *Service) ListCourseMaterials(ctx context.Context, courseID string) ([]Attachment, error) {
	items, err := s.repo.ListByResource(ctx, ResourceTypeCourse, courseID)
	if err != nil {
		return nil, err
	}

	for i := range items {
		url, err := s.r2.PresignGet(ctx, items[i].Key, time.Duration(DownloadURLTTL)*time.Second)
		if err != nil {
			logging.FromLogger(ctx).Error("failed to presign download url",
				zap.String("event", "material_download_presign_failed"),
				zap.String("attachment_id", items[i].ID.String()),
				zap.String("r2_key", items[i].Key),
				zap.Error(err),
			)
			return nil, err
		}
		items[i].DownloadURL = url
	}
	return items, nil
}

func (s *Service) ListCourseMaterialsForStudent(ctx context.Context, userID uuid.UUID, courseID string) ([]Attachment, error) {
	ok, err := s.repo.UserHasCourseAccess(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		logging.FromLogger(ctx).Warn("student denied access to course materials",
			zap.String("event", "material_access_denied"),
			zap.String("user_id", userID.String()),
			zap.String("course_id", courseID),
		)
		return nil, ErrForbidden
	}
	return s.ListCourseMaterials(ctx, courseID)
}

func (s *Service) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	log := logging.FromLogger(ctx)

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if delErr := s.r2.DeleteObject(ctx, a.Key); delErr != nil {
		log.Warn("failed to delete object from bucket",
			zap.String("event", "material_r2_delete_failed"),
			zap.String("attachment_id", a.ID.String()),
			zap.String("r2_key", a.Key),
			zap.Error(delErr),
		)
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	log.Info("material deleted",
		zap.String("event", "material_deleted"),
		zap.String("attachment_id", a.ID.String()),
		zap.String("resource_id", a.ResourceID),
		zap.String("r2_key", a.Key),
	)

	return nil
}
