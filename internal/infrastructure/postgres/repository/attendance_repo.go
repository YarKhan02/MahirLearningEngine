package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/google/uuid"
)

//go:embed sql/attendance_roster.sql
var attendanceRosterSQL string

//go:embed sql/attendance_session_upsert.sql
var attendanceSessionUpsertSQL string

//go:embed sql/attendance_mark_upsert.sql
var attendanceMarkUpsertSQL string

//go:embed sql/attendance_by_student.sql
var attendanceByStudentSQL string

//go:embed sql/student_by_user_get.sql
var attendanceStudentByUserSQL string

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) GetRoster(ctx context.Context, batchID uuid.UUID, date time.Time) ([]attendance.RosterEntry, error) {
	rows, err := r.db.QueryContext(ctx, attendanceRosterSQL, batchID, date)
	if err != nil {
		return nil, fmt.Errorf("get roster: %w", err)
	}
	defer rows.Close()

	var roster []attendance.RosterEntry

	for rows.Next() {
		var e attendance.RosterEntry

		if err := rows.Scan(
			&e.StudentID,
			&e.FullName,
			&e.Email,
			&e.Status,
		); err != nil {
			return nil, fmt.Errorf("scan roster entry: %w", err)
		}

		roster = append(roster, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate roster: %w", err)
	}

	return roster, nil
}

func (r *AttendanceRepository) Mark(ctx context.Context, req attendance.MarkAttendance) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("mark attendance: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	sessionID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	// Get-or-create the session for this batch + date.
	if err := tx.QueryRowContext(
		ctx,
		attendanceSessionUpsertSQL,
		sessionID,
		req.BatchID,
		req.Date,
		req.CreatedBy,
	).Scan(&sessionID); err != nil {
		return fmt.Errorf("upsert attendance session: %w", err)
	}

	markID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, attendanceMarkUpsertSQL, markID, sessionID, req.StudentID, req.Status); err != nil {
		return fmt.Errorf("upsert attendance: %w", err)
	}

	return tx.Commit()
}

func (r *AttendanceRepository) GetStudentRecords(ctx context.Context, studentID uuid.UUID) ([]attendance.Record, error) {
	rows, err := r.db.QueryContext(ctx, attendanceByStudentSQL, studentID)
	if err != nil {
		return nil, fmt.Errorf("get attendance records: %w", err)
	}
	defer rows.Close()

	var records []attendance.Record

	for rows.Next() {
		var rec attendance.Record

		if err := rows.Scan(
			&rec.LessonDate,
			&rec.Status,
			&rec.BatchName,
		); err != nil {
			return nil, fmt.Errorf("scan attendance record: %w", err)
		}

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate attendance records: %w", err)
	}

	return records, nil
}

func (r *AttendanceRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx, attendanceStudentByUserSQL, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, attendance.ErrStudentNotFound
		}
		return uuid.Nil, fmt.Errorf("get student by user: %w", err)
	}

	return id, nil
}
