package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/google/uuid"
)

//go:embed sql/course_insert.sql
var courseInsertSQL string

//go:embed sql/course_get.sql
var courseGetSQL string

//go:embed sql/course_exists.sql
var courseExistsSQL string

//go:embed sql/course_insert_lesson.sql
var courseInsertLessonSQL string

//go:embed sql/course_get_lesson.sql
var courseGetLessonSQL string

//go:embed sql/course_lesson_exists.sql
var courseLessonExistsSQL string

//go:embed sql/course_lesson_update.sql
var courseLessonUpdateSQL string

//go:embed sql/course_lesson_order_no.sql
var courseLessonOrderNoSQL string

//go:embed sql/course_lesson_count.sql
var courseLessonCountSQL string

//go:embed sql/course_lesson_move_down.sql
var courseLessonMoveDownSQL string

//go:embed sql/course_lesson_move_up.sql
var courseLessonMoveUpSQL string

//go:embed sql/course_lesson_update_order.sql
var courseLessonUpdateOrderSQL string

var (
	ErrLessonNotFound  = errors.New("lesson not found")
	ErrInvalidOrderNo  = errors.New("order_no out of range")
)

type CoursRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CoursRepository {
	return &CoursRepository{
		db: db,
	}
}

func (r *CoursRepository) InsertCourse(ctx context.Context, req course.Course) (*course.Course, error) {
	
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	var description any
	if req.Description == "" {
		description = nil
	} else {
		description = req.Description
	}

	var course course.Course

	err = r.db.QueryRowContext(ctx, courseInsertSQL, 
		id, 
		req.Title,
		req.Level, 
		req.Duration, 
		description,
	).Scan(
		&course.ID,
		&course.Title,
		&course.Level,
		&course.Duration,
		&course.Description,
		&course.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CoursRepository) GetCourse(ctx context.Context) ([]course.Course, error) {

	rows, err := r.db.QueryContext(ctx, courseGetSQL)
	if err != nil {
		return nil, err
	}

	courses := make([]course.Course, 0)

	for rows.Next() {
		var c course.Course

		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.Level,
			&c.Duration,
			&c.IsActive,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err 
	}

	return courses, nil
}

func (r *CoursRepository) CourseExists(ctx context.Context, id uuid.UUID) bool {
	var exists bool

	err := r.db.QueryRowContext(ctx, courseExistsSQL, id).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (r *CoursRepository) InsertLesson(ctx context.Context, req course.Lesson) error {
	
	id, err := uuid.NewV7()
	if err != nil {
		return err 
	}

	req.ID = id

	_, err = r.db.ExecContext(ctx, courseInsertLessonSQL, 
		req.ID,
		req.CourseID,
		req.Title,
		req.Description,
		req.OrderNo,
		req.YoutubeURL,
		req.Content,
	)

	if err != nil {
		return err 
	} 

	return nil
}

func (r *CoursRepository) GetLesson(ctx context.Context, id uuid.UUID) ([]course.Lesson, error) {
	rows, err := r.db.QueryContext(ctx, courseGetLessonSQL)
	if err != nil {
		return nil, err
	}

	lessons := make([]course.Lesson, 0)

	for rows.Next() {
		var c course.Lesson

		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.OrderNo,
			&c.YoutubeURL,
			&c.Content,
		)
		if err != nil {
			return nil, err
		}

		lessons = append(lessons, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err 
	}

	return lessons, nil
}

func (r *CoursRepository) LessonExists(ctx context.Context, id uuid.UUID) bool {
	var exists bool

	err := r.db.QueryRowContext(ctx, courseLessonExistsSQL, id).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (r *CoursRepository) UpdateLesson(ctx context.Context, req course.UpdateLesson) error {
	query := courseLessonUpdateSQL + " "
	args := []any{}
    idx := 1

    if req.Title != nil {
        query += fmt.Sprintf("title = $%d,", idx)
        args = append(args, *req.Title)
        idx++
    }

    if req.Description != nil {
        query += fmt.Sprintf("description = $%d,", idx)
        args = append(args, *req.Description)
        idx++
    }

    if req.YoutubeURL != nil {
        query += fmt.Sprintf("youtube_url = $%d,", idx)
        args = append(args, *req.YoutubeURL)
        idx++
    }

    if req.Content != nil {
        query += fmt.Sprintf("content = $%d,", idx)
        args = append(args, *req.Content)
        idx++
    }

    if len(args) == 0 {
        return nil // nothing to update
    }

    // Remove trailing comma
    query = strings.TrimSuffix(query, ",")

    query += fmt.Sprintf(
        " WHERE id = $%d AND course_id = $%d",
        idx, idx+1,
    )

    args = append(args, req.ID, req.CourseID)

    _, err := r.db.ExecContext(ctx, query, args...)

	return err
}

func (r *CoursRepository) ReorderLesson(ctx context.Context, lessonID uuid.UUID, orderNo int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // no-op once committed; cleans up on any early return/panic

	var courseID uuid.UUID
	var oldNo int

	// Lock the row so a concurrent request against the same lesson can't
	// read a stale order_no while this transaction is in flight.
	err = tx.QueryRowContext(ctx, courseLessonOrderNoSQL, lessonID).Scan(&courseID, &oldNo)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrLessonNotFound
	} 
	if err != nil {
		return fmt.Errorf("lock lesson: %w", err)
	}

	// Already in the requested spot — nothing to do.
	if oldNo == orderNo {
		return tx.Commit()
	}

	var count int
	err = tx.QueryRowContext(ctx, courseLessonCountSQL, courseID).Scan(&count)
	if err != nil {
		return fmt.Errorf("count lessons: %w", err)
	}
	if orderNo < 1 || orderNo > count {
		return ErrInvalidOrderNo
	}

	// Shift the rows between old and new position to close/open the gap.
	if oldNo < orderNo {
		// Moving down the list: everything strictly after the old spot,
		// up to and including the target, shifts up by one.
		_, err = tx.ExecContext(ctx, courseLessonMoveDownSQL,
			courseID, 
			oldNo, 
			orderNo,
		)
	} else {
		// Moving up the list: everything from the target up to
		// (but not including) the old spot shifts down by one.
		_, err = tx.ExecContext(ctx, courseLessonMoveUpSQL,
			courseID, 
			orderNo, 
			oldNo,
		)
	}
	if err != nil {
		return fmt.Errorf("shift lessons: %w", err)
	}

	// Drop the moved lesson into its final slot.
	_, err = tx.ExecContext(ctx, courseLessonUpdateOrderSQL,
		orderNo, 
		lessonID,
	)
	if err != nil {
		return fmt.Errorf("set lesson position: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}