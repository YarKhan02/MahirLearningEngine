package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/announcement"
	"github.com/google/uuid"
)

//go:embed sql/announcement_create.sql
var announcementCreateSQL string

//go:embed sql/announcement_get_all.sql
var announcementGetAllSQL string

//go:embed sql/announcement_for_user.sql
var announcementForUserSQL string

//go:embed sql/announcement_delete.sql
var announcementDeleteSQL string

//go:embed sql/announcement_get_by_id.sql
var announcementGetByIDSQL string

type AnnouncementRepository struct {
	db *sql.DB
}

func NewAnnouncementRepository(db *sql.DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

func (r *AnnouncementRepository) Create(ctx context.Context, a *announcement.Announcement) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	a.ID = id

	_, err = r.db.ExecContext(ctx, announcementCreateSQL, a.ID, a.BatchID, a.Title, a.Description)
	if err != nil {
		return fmt.Errorf("create announcement: %w", err)
	}

	return nil
}

func (r *AnnouncementRepository) GetAll(ctx context.Context) ([]announcement.Announcement, error) {
	rows, err := r.db.QueryContext(ctx, announcementGetAllSQL)
	if err != nil {
		return nil, fmt.Errorf("get announcements: %w", err)
	}
	defer rows.Close()

	return scanAnnouncements(rows)
}

func (r *AnnouncementRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]announcement.Announcement, error) {
	rows, err := r.db.QueryContext(ctx, announcementForUserSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("get user announcements: %w", err)
	}
	defer rows.Close()

	return scanAnnouncements(rows)
}

func (r *AnnouncementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, announcementDeleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete announcement: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete announcement: rows affected: %w", err)
	}
	if rows == 0 {
		return announcement.ErrNotFound
	}

	return nil
}

func (r *AnnouncementRepository) GetByID(ctx context.Context, id uuid.UUID) (*announcement.Announcement, error) {
	row := r.db.QueryRowContext(ctx, announcementGetByIDSQL, id)

	var a announcement.Announcement
	if err := row.Scan(
		&a.ID,
		&a.BatchID,
		&a.Title,
		&a.Description,
		&a.CreatedAt,
		&a.BatchName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, announcement.ErrNotFound
		}
		return nil, fmt.Errorf("get announcement by id: %w", err)
	}

	return &a, nil
}

func scanAnnouncements(rows *sql.Rows) ([]announcement.Announcement, error) {
	var out []announcement.Announcement

	for rows.Next() {
		var a announcement.Announcement

		if err := rows.Scan(
			&a.ID,
			&a.BatchID,
			&a.Title,
			&a.Description,
			&a.CreatedAt,
			&a.BatchName,
		); err != nil {
			return nil, fmt.Errorf("scan announcement: %w", err)
		}

		out = append(out, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate announcements: %w", err)
	}

	return out, nil
}
