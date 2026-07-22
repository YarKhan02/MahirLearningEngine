package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attachement"
	"github.com/google/uuid"
)

//go:embed sql/attachment_create.sql
var attachmentCreateSQL string

//go:embed sql/attachment_confirm.sql
var attachmentConfirmSQL string

//go:embed sql/attachment_list_by_resource.sql
var attachmentListByResourceSQL string

//go:embed sql/attachment_get_by_id.sql
var attachmentGetByIDSQL string

//go:embed sql/attachment_soft_delete.sql
var attachmentSoftDeleteSQL string

//go:embed sql/attachment_course_access.sql
var attachmentCourseAccessSQL string

//go:embed sql/attachment_course_exists.sql
var attachmentCourseExistsSQL string

//go:embed sql/attachment_pending_by_key.sql
var attachmentPendingByKeySQL string

type AttachementRepository struct {
	db *sql.DB
}

func NewAttachementRepository(db *sql.DB) *AttachementRepository {
	return &AttachementRepository{db: db}
}

func (r *AttachementRepository) Create(ctx context.Context, a attachement.Attachment) error {
	var size any
	if a.SizeBytes != nil {
		size = *a.SizeBytes
	}

	_, err := r.db.ExecContext(ctx, attachmentCreateSQL,
		a.ID, a.Key, a.Filename, a.ContentType, size,
		a.ResourceType, a.ResourceID, a.UploadedBy,
	)
	if err != nil {
		return fmt.Errorf("create attachment: %w", err)
	}
	return nil
}

func (r *AttachementRepository) GetPendingByKey(ctx context.Context, key string, uploadedBy uuid.UUID) (attachement.Attachment, error) {
	row := r.db.QueryRowContext(ctx, attachmentPendingByKeySQL, key, uploadedBy)
	a, err := scanAttachment(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attachement.Attachment{}, attachement.ErrNotFound
		}
		return attachement.Attachment{}, fmt.Errorf("pending attachment: %w", err)
	}
	return a, nil
}

func (r *AttachementRepository) ConfirmByKey(ctx context.Context, key string, uploadedBy uuid.UUID, sizeBytes int64, verifiedContentType string) (attachement.Attachment, error) {
	row := r.db.QueryRowContext(ctx, attachmentConfirmSQL, key, uploadedBy, sizeBytes, verifiedContentType)
	a, err := scanAttachment(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attachement.Attachment{}, attachement.ErrNotFound
		}
		return attachement.Attachment{}, fmt.Errorf("confirm attachment: %w", err)
	}
	return a, nil
}

func (r *AttachementRepository) CourseExists(ctx context.Context, courseID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, attachmentCourseExistsSQL, courseID).Scan(&exists); err != nil {
		return false, fmt.Errorf("course exists: %w", err)
	}
	return exists, nil
}

func (r *AttachementRepository) ListByResource(ctx context.Context, resourceType, resourceID string) ([]attachement.Attachment, error) {
	rows, err := r.db.QueryContext(ctx, attachmentListByResourceSQL, resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("list attachments: %w", err)
	}
	defer rows.Close()

	var out []attachement.Attachment
	for rows.Next() {
		a, err := scanAttachment(rows)
		if err != nil {
			return nil, fmt.Errorf("scan attachment: %w", err)
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate attachments: %w", err)
	}
	return out, nil
}

func (r *AttachementRepository) GetByID(ctx context.Context, id uuid.UUID) (attachement.Attachment, error) {
	row := r.db.QueryRowContext(ctx, attachmentGetByIDSQL, id)
	a, err := scanAttachment(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attachement.Attachment{}, attachement.ErrNotFound
		}
		return attachement.Attachment{}, fmt.Errorf("get attachment: %w", err)
	}
	return a, nil
}

func (r *AttachementRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, attachmentSoftDeleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete attachment: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete attachment: rows affected: %w", err)
	}
	if n == 0 {
		return attachement.ErrNotFound
	}
	return nil
}

func (r *AttachementRepository) UserHasCourseAccess(ctx context.Context, userID uuid.UUID, courseID string) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, attachmentCourseAccessSQL, userID, courseID).Scan(&ok); err != nil {
		return false, fmt.Errorf("course access: %w", err)
	}
	return ok, nil
}

type attachmentScanner interface {
	Scan(dest ...any) error
}

// scanAttachment reads the shared attachment projection.
func scanAttachment(s attachmentScanner) (attachement.Attachment, error) {
	var (
		a           attachement.Attachment
		resourceID  uuid.UUID
		uploadedBy  uuid.NullUUID
		sizeBytes   sql.NullInt64
		confirmedAt sql.NullTime
	)

	if err := s.Scan(
		&a.ID,
		&a.Key,
		&a.Filename,
		&a.ContentType,
		&sizeBytes,
		&a.ResourceType,
		&resourceID,
		&uploadedBy,
		&a.Status,
		&a.CreatedAt,
		&confirmedAt,
	); err != nil {
		return attachement.Attachment{}, err
	}

	a.ResourceID = resourceID.String()
	if sizeBytes.Valid {
		v := sizeBytes.Int64
		a.SizeBytes = &v
	}
	if uploadedBy.Valid {
		a.UploadedBy = uploadedBy.UUID
	}
	if confirmedAt.Valid {
		t := confirmedAt.Time
		a.ConfirmedAt = &t
	}

	return a, nil
}
