package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/google/uuid"
)

//go:embed sql/batch_create.sql
var batchCreateSQL string

//go:embed sql/batch_update.sql
var batchUpdateSQL string

//go:embed sql/batch_delete.sql
var batchDeleteSQL string

//go:embed sql/batch_get.sql
var batchGetAllSQL string

//go:embed sql/batch_courses_get.sql
var batchCoursesGetSQL string

//go:embed sql/batch_course_add.sql
var batchCourseAddSQL string

//go:embed sql/batch_course_remove.sql
var batchCourseRemoveSQL string

type BatchRepository struct {
	db *sql.DB
}

func NewBatchRepository(db *sql.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

func (r *BatchRepository) CreateBatch(ctx context.Context, batch *batch.Batch) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	batch.ID = id
	batch.Status = "upcoming"

	_, err = r.db.ExecContext(
		ctx,
		batchCreateSQL,
		batch.ID,
		batch.BatchName,
		batch.StartDate,
		batch.EndDate,
		batch.Capacity,
		batch.Days,
		batch.Status,
		batch.Price,
	)
	if err != nil {
		return fmt.Errorf("create batch: %w", err)
	}

	return nil
}

func (r *BatchRepository) UpdateBatch(ctx context.Context, b *batch.Batch) error {
	res, err := r.db.ExecContext(
		ctx,
		batchUpdateSQL,
		b.ID,
		b.BatchName,
		b.StartDate,
		b.EndDate,
		b.Capacity,
		b.Days,
		b.Status,
		b.Price,
	)
	if err != nil {
		return fmt.Errorf("update batch: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update batch: rows affected: %w", err)
	}

	if rows == 0 {
		return batch.ErrBatchNotFound
	}

	return nil
}

func (r *BatchRepository) DeleteBatch(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, batchDeleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete batch: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete batch: rows affected: %w", err)
	}

	if rows == 0 {
		return batch.ErrBatchNotFound
	}

	return nil
}

func (r *BatchRepository) GetBatches(ctx context.Context) ([]batch.Batch, error) {
	rows, err := r.db.QueryContext(ctx, batchGetAllSQL)
	if err != nil {
		return nil, fmt.Errorf("get batches: %w", err)
	}
	defer rows.Close()

	var batches []batch.Batch

	for rows.Next() {
		var b batch.Batch

		if err := rows.Scan(
			&b.ID,
			&b.BatchName,
			&b.StartDate,
			&b.EndDate,
			&b.Capacity,
			&b.Days,
			&b.Status,
			&b.Price,
		); err != nil {
			return nil, fmt.Errorf("scan batch: %w", err)
		}

		batches = append(batches, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate batches: %w", err)
	}

	return batches, nil
}

func (r *BatchRepository) GetBatchCourses(ctx context.Context, batchID uuid.UUID) ([]batch.BatchCourse, error) {
	rows, err := r.db.QueryContext(ctx, batchCoursesGetSQL, batchID)
	if err != nil {
		return nil, fmt.Errorf("get batch courses: %w", err)
	}
	defer rows.Close()

	var courses []batch.BatchCourse

	for rows.Next() {
		var bc batch.BatchCourse

		if err := rows.Scan(
			&bc.ID,
			&bc.CourseID,
			&bc.Title,
			&bc.Level,
			&bc.GrantedAt,
		); err != nil {
			return nil, fmt.Errorf("scan batch course: %w", err)
		}

		courses = append(courses, bc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate batch courses: %w", err)
	}

	return courses, nil
}

func (r *BatchRepository) UpdateBatchCourses(ctx context.Context, batchID uuid.UUID, add []uuid.UUID, remove []uuid.UUID, grantedBy *uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update batch courses: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	for _, courseID := range add {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, batchCourseAddSQL, id, batchID, courseID, grantedBy); err != nil {
			return fmt.Errorf("add course %s to batch: %w", courseID, err)
		}
	}

	for _, courseID := range remove {
		if _, err := tx.ExecContext(ctx, batchCourseRemoveSQL, batchID, courseID); err != nil {
			return fmt.Errorf("remove course %s from batch: %w", courseID, err)
		}
	}

	return tx.Commit()
}
