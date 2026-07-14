package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/timetable"
	"github.com/google/uuid"
)

//go:embed sql/timetable_create.sql
var timetableCreateSQL string

//go:embed sql/timetable_by_batch.sql
var timetableByBatchSQL string

//go:embed sql/timetable_for_user.sql
var timetableForUserSQL string

//go:embed sql/timetable_delete.sql
var timetableDeleteSQL string

type TimetableRepository struct {
	db *sql.DB
}

func NewTimetableRepository(db *sql.DB) *TimetableRepository {
	return &TimetableRepository{db: db}
}

func (r *TimetableRepository) Create(ctx context.Context, t *timetable.Timetable) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	t.ID = id

	_, err = r.db.ExecContext(
		ctx,
		timetableCreateSQL,
		t.ID,
		t.BatchID,
		t.CourseID,
		timetable.WeekdaysToMask(t.Weekdays),
		t.StartTime,
		t.EndTime,
	)
	if err != nil {
		return fmt.Errorf("create timetable: %w", err)
	}

	return nil
}

func (r *TimetableRepository) GetByBatch(ctx context.Context, batchID uuid.UUID) ([]timetable.Timetable, error) {
	rows, err := r.db.QueryContext(ctx, timetableByBatchSQL, batchID)
	if err != nil {
		return nil, fmt.Errorf("get batch timetable: %w", err)
	}
	defer rows.Close()

	return scanTimetables(rows)
}

func (r *TimetableRepository) GetRulesForUser(ctx context.Context, userID uuid.UUID) ([]timetable.Timetable, error) {
	rows, err := r.db.QueryContext(ctx, timetableForUserSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("get user timetable: %w", err)
	}
	defer rows.Close()

	return scanTimetables(rows)
}

func (r *TimetableRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, timetableDeleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete timetable: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete timetable: rows affected: %w", err)
	}
	if rows == 0 {
		return timetable.ErrTimetableNotFound
	}

	return nil
}

func scanTimetables(rows *sql.Rows) ([]timetable.Timetable, error) {
	var out []timetable.Timetable

	for rows.Next() {
		var t timetable.Timetable
		var mask int

		if err := rows.Scan(
			&t.ID,
			&t.BatchID,
			&t.CourseID,
			&mask,
			&t.StartTime,
			&t.EndTime,
			&t.CourseTitle,
			&t.BatchName,
			&t.BatchStart,
			&t.BatchEnd,
		); err != nil {
			return nil, fmt.Errorf("scan timetable: %w", err)
		}

		t.Weekdays = timetable.MaskToWeekdays(mask)
		out = append(out, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate timetable: %w", err)
	}

	return out, nil
}
