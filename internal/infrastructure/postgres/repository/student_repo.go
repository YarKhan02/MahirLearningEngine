package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	
	"github.com/google/uuid"
)

//go:embed sql/student_create.sql
var studentCreateSQL string

//go:embed sql/student_batch_enroll.sql
var studentBatchEnrollSQL string

//go:embed sql/student_get_all.sql
var studentGetAllSQL string

//go:embed sql/student_get_by_id.sql
var studentGetByIDSQL string

//go:embed sql/student_update_status.sql
var studentUpdateStatusSQL string

//go:embed sql/student_batch_delete.sql
var studentBatchDeleteSQL string

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) RegisterStudent(ctx context.Context, s *student.Student, batchID uuid.UUID) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	s.ID = id
	s.Status = "pending"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("register student: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.ExecContext(
		ctx,
		studentCreateSQL,
		s.ID,
		s.Email,
		s.FullName,
		s.PhoneNumber,
		s.DOB,
		s.Gender,
		s.Status,
	)
	if err != nil {
		if strings.Contains(err.Error(), "students_email_key") {
			return student.ErrEmailAlreadyRegistered
		}
		return fmt.Errorf("register student: %w", err)
	}

	enrollID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, studentBatchEnrollSQL, enrollID, s.ID, batchID); err != nil {
		return fmt.Errorf("enroll student in batch: %w", err)
	}

	return tx.Commit()
}

func (r *StudentRepository) GetStudents(ctx context.Context, q string) ([]student.StudentWithBatch, error) {
	rows, err := r.db.QueryContext(ctx, studentGetAllSQL, q)
	if err != nil {
		return nil, fmt.Errorf("get students: %w", err)
	}
	defer rows.Close()

	var students []student.StudentWithBatch

	for rows.Next() {
		var s student.StudentWithBatch

		if err := rows.Scan(
			&s.ID,
			&s.Email,
			&s.FullName,
			&s.PhoneNumber,
			&s.DOB,
			&s.Gender,
			&s.Status,
			&s.BatchID,
			&s.BatchName,
			&s.HasAccount,
		); err != nil {
			return nil, fmt.Errorf("scan student: %w", err)
		}

		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate students: %w", err)
	}

	return students, nil
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, id uuid.UUID) (*student.Student, error) {
	var s student.Student

	err := r.db.QueryRowContext(ctx, studentGetByIDSQL, id).Scan(
		&s.ID,
		&s.Email,
		&s.FullName,
		&s.PhoneNumber,
		&s.DOB,
		&s.Gender,
		&s.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, student.ErrStudentNotFound
		}
		return nil, fmt.Errorf("get student: %w", err)
	}

	return &s, nil
}

func (r *StudentRepository) UpdateStudentStatus(ctx context.Context, id uuid.UUID, status string) error {
	if _, err := r.db.ExecContext(ctx, studentUpdateStatusSQL, id, status); err != nil {
		return fmt.Errorf("update student status: %w", err)
	}
	return nil
}

func (r *StudentRepository) UpdateStudentBatch(ctx context.Context, studentID uuid.UUID, batchID *uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update student batch: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx, studentBatchDeleteSQL, studentID); err != nil {
		return fmt.Errorf("remove student batch: %w", err)
	}

	if batchID != nil {
		enrollID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, studentBatchEnrollSQL, enrollID, studentID, *batchID); err != nil {
			return fmt.Errorf("enroll student in batch: %w", err)
		}
	}

	return tx.Commit()
}
