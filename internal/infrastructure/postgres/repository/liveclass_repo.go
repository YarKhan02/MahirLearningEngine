package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/liveclass"
	"github.com/google/uuid"
)

//go:embed sql/live_session_create.sql
var liveSessionCreateSQL string

//go:embed sql/live_session_end.sql
var liveSessionEndSQL string

//go:embed sql/live_session_get.sql
var liveSessionGetSQL string

//go:embed sql/live_session_live_by_batch.sql
var liveSessionLiveByBatchSQL string

//go:embed sql/live_session_student_can_join.sql
var liveSessionStudentCanJoinSQL string

//go:embed sql/live_session_for_user.sql
var liveSessionForUserSQL string

//go:embed sql/live_sessions_by_batch_course.sql
var liveSessionsByBatchCourseSQL string

//go:embed sql/live_student_batches.sql
var liveStudentBatchesSQL string

//go:embed sql/live_student_batch_courses.sql
var liveStudentBatchCoursesSQL string

//go:embed sql/live_student_in_batch.sql
var liveStudentInBatchSQL string

//go:embed sql/live_user_display_name.sql
var liveUserDisplayNameSQL string

//go:embed sql/whiteboard_snapshot_create.sql
var whiteboardSnapshotCreateSQL string

//go:embed sql/whiteboard_snapshot_next_version.sql
var whiteboardSnapshotNextVersionSQL string

//go:embed sql/whiteboard_snapshot_list.sql
var whiteboardSnapshotListSQL string

type LiveClassRepository struct {
	db *sql.DB
}

func NewLiveClassRepository(db *sql.DB) *LiveClassRepository {
	return &LiveClassRepository{db: db}
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func scanSession(s sessionScanner) (liveclass.LiveSession, error) {
	var (
		ls      liveclass.LiveSession
		endedAt sql.NullTime
	)
	if err := s.Scan(
		&ls.ID, &ls.BatchID, &ls.CourseID, &ls.HostID, &ls.Title, &ls.ClassDate,
		&ls.Status, &ls.StartedAt, &endedAt, &ls.CourseTitle, &ls.BatchName,
	); err != nil {
		return liveclass.LiveSession{}, err
	}
	if endedAt.Valid {
		t := endedAt.Time
		ls.EndedAt = &t
	}
	return ls, nil
}

func (r *LiveClassRepository) CreateSession(ctx context.Context, s liveclass.LiveSession) error {
	_, err := r.db.ExecContext(ctx, liveSessionCreateSQL, s.ID, s.BatchID, s.CourseID, s.HostID, s.Title)
	if err != nil {
		return fmt.Errorf("create live session: %w", err)
	}
	return nil
}

func (r *LiveClassRepository) SessionsByBatchCourse(ctx context.Context, batchID, courseID uuid.UUID) ([]liveclass.LiveSession, error) {
	rows, err := r.db.QueryContext(ctx, liveSessionsByBatchCourseSQL, batchID, courseID)
	if err != nil {
		return nil, fmt.Errorf("sessions by batch/course: %w", err)
	}
	defer rows.Close()

	out := make([]liveclass.LiveSession, 0)
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *LiveClassRepository) StudentBatches(ctx context.Context, userID uuid.UUID) ([]liveclass.BatchOption, error) {
	rows, err := r.db.QueryContext(ctx, liveStudentBatchesSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("student batches: %w", err)
	}
	defer rows.Close()

	out := make([]liveclass.BatchOption, 0)
	for rows.Next() {
		var b liveclass.BatchOption
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, fmt.Errorf("scan batch: %w", err)
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *LiveClassRepository) StudentBatchCourses(ctx context.Context, userID, batchID uuid.UUID) ([]liveclass.CourseOption, error) {
	rows, err := r.db.QueryContext(ctx, liveStudentBatchCoursesSQL, userID, batchID)
	if err != nil {
		return nil, fmt.Errorf("student batch courses: %w", err)
	}
	defer rows.Close()

	out := make([]liveclass.CourseOption, 0)
	for rows.Next() {
		var c liveclass.CourseOption
		if err := rows.Scan(&c.ID, &c.Title); err != nil {
			return nil, fmt.Errorf("scan course: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *LiveClassRepository) UserDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
	var name string
	if err := r.db.QueryRowContext(ctx, liveUserDisplayNameSQL, userID).Scan(&name); err != nil {
		return "", fmt.Errorf("user display name: %w", err)
	}
	return name, nil
}

func (r *LiveClassRepository) StudentInBatch(ctx context.Context, userID, batchID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, liveStudentInBatchSQL, userID, batchID).Scan(&ok); err != nil {
		return false, fmt.Errorf("student in batch: %w", err)
	}
	return ok, nil
}

func (r *LiveClassRepository) EndSession(ctx context.Context, id, hostID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, liveSessionEndSQL, id, hostID)
	if err != nil {
		return fmt.Errorf("end live session: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("end live session: rows affected: %w", err)
	}
	if n == 0 {
		return liveclass.ErrNotFound
	}
	return nil
}

func (r *LiveClassRepository) GetSession(ctx context.Context, id uuid.UUID) (liveclass.LiveSession, error) {
	row := r.db.QueryRowContext(ctx, liveSessionGetSQL, id)
	ls, err := scanSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return liveclass.LiveSession{}, liveclass.ErrNotFound
		}
		return liveclass.LiveSession{}, fmt.Errorf("get live session: %w", err)
	}
	return ls, nil
}

func (r *LiveClassRepository) GetLiveByBatch(ctx context.Context, batchID uuid.UUID) (liveclass.LiveSession, error) {
	row := r.db.QueryRowContext(ctx, liveSessionLiveByBatchSQL, batchID)
	ls, err := scanSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return liveclass.LiveSession{}, liveclass.ErrNotFound
		}
		return liveclass.LiveSession{}, fmt.Errorf("get live session by batch: %w", err)
	}
	return ls, nil
}

func (r *LiveClassRepository) GetLiveForUser(ctx context.Context, userID uuid.UUID) (liveclass.LiveSession, error) {
	row := r.db.QueryRowContext(ctx, liveSessionForUserSQL, userID)
	ls, err := scanSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return liveclass.LiveSession{}, liveclass.ErrNotFound
		}
		return liveclass.LiveSession{}, fmt.Errorf("get live session for user: %w", err)
	}
	return ls, nil
}

func (r *LiveClassRepository) StudentCanJoin(ctx context.Context, userID, sessionID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, liveSessionStudentCanJoinSQL, userID, sessionID).Scan(&ok); err != nil {
		return false, fmt.Errorf("student can join: %w", err)
	}
	return ok, nil
}

func (r *LiveClassRepository) CreateSnapshot(ctx context.Context, w liveclass.WhiteboardSnapshot) error {
	_, err := r.db.ExecContext(ctx, whiteboardSnapshotCreateSQL,
		w.ID, w.SessionID, w.StorageKey, w.ImageURL, w.Version, w.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}
	return nil
}

func (r *LiveClassRepository) NextSnapshotVersion(ctx context.Context, sessionID uuid.UUID) (int, error) {
	var v int
	if err := r.db.QueryRowContext(ctx, whiteboardSnapshotNextVersionSQL, sessionID).Scan(&v); err != nil {
		return 0, fmt.Errorf("next snapshot version: %w", err)
	}
	return v, nil
}

func (r *LiveClassRepository) ListSnapshots(ctx context.Context, sessionID uuid.UUID) ([]liveclass.WhiteboardSnapshot, error) {
	rows, err := r.db.QueryContext(ctx, whiteboardSnapshotListSQL, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}
	defer rows.Close()

	out := make([]liveclass.WhiteboardSnapshot, 0)
	for rows.Next() {
		var w liveclass.WhiteboardSnapshot
		if err := rows.Scan(
			&w.ID, &w.SessionID, &w.StorageKey, &w.ImageURL, &w.Version, &w.CreatedBy, &w.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan snapshot: %w", err)
		}
		out = append(out, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate snapshots: %w", err)
	}
	return out, nil
}
