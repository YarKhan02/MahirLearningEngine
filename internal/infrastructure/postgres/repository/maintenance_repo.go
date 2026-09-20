package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/maintenance"
	"github.com/google/uuid"
)

//go:embed sql/maintenance_purge_tokens.sql
var maintenancePurgeTokensSQL string

//go:embed sql/maintenance_stale_pending_attachments.sql
var maintenanceStalePendingAttachmentsSQL string

//go:embed sql/maintenance_delete_attachment.sql
var maintenanceDeleteAttachmentSQL string

type MaintenanceRepository struct {
	db *sql.DB
}

func NewMaintenanceRepository(db *sql.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *MaintenanceRepository) PurgeExpiredTokens(ctx context.Context) (int64, error) {
	res, err := r.db.ExecContext(ctx, maintenancePurgeTokensSQL)
	if err != nil {
		return 0, fmt.Errorf("purge expired tokens: %w", err)
	}
	return res.RowsAffected()
}

func (r *MaintenanceRepository) ListStalePendingAttachments(ctx context.Context) ([]maintenance.StalePendingAttachment, error) {
	rows, err := r.db.QueryContext(ctx, maintenanceStalePendingAttachmentsSQL)
	if err != nil {
		return nil, fmt.Errorf("list stale pending attachments: %w", err)
	}
	defer rows.Close()

	var out []maintenance.StalePendingAttachment
	for rows.Next() {
		var a maintenance.StalePendingAttachment
		if err := rows.Scan(&a.ID, &a.R2Key); err != nil {
			return nil, fmt.Errorf("scan stale attachment: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *MaintenanceRepository) DeleteAttachmentByID(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, maintenanceDeleteAttachmentSQL, id); err != nil {
		return fmt.Errorf("delete attachment: %w", err)
	}
	return nil
}
