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

//go:embed sql/student_courses_get.sql
var studentCoursesGetSQL string

//go:embed sql/student_lessons_get.sql
var studentLessonsGetSQL string

//go:embed sql/student_by_user_get.sql
var studentByUserGetSQL string

//go:embed sql/student_course_access_check.sql
var studentCourseAccessCheckSQL string

//go:embed sql/lesson_progress_upsert.sql
var lessonProgressUpsertSQL string

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

func (r *StudentRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx, studentByUserGetSQL, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, student.ErrStudentNotFound
		}
		return uuid.Nil, fmt.Errorf("get student by user: %w", err)
	}

	return id, nil
}

func (r *StudentRepository) GetStudentCourses(ctx context.Context, userID uuid.UUID) ([]student.StudentCourse, error) {
	rows, err := r.db.QueryContext(ctx, studentCoursesGetSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("get student courses: %w", err)
	}
	defer rows.Close()

	var courses []student.StudentCourse

	for rows.Next() {
		var c student.StudentCourse

		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Level,
			&c.Duration,
			&c.Description,
			&c.TotalLessons,
			&c.CompletedLessons,
		); err != nil {
			return nil, fmt.Errorf("scan student course: %w", err)
		}

		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student courses: %w", err)
	}

	return courses, nil
}

func (r *StudentRepository) HasCourseAccess(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (bool, error) {
	var hasAccess bool

	err := r.db.QueryRowContext(ctx, studentCourseAccessCheckSQL, studentID, courseID).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("check course access: %w", err)
	}

	return hasAccess, nil
}

func (r *StudentRepository) GetStudentLessons(ctx context.Context, courseID uuid.UUID, studentID uuid.UUID) ([]student.StudentLesson, error) {
	rows, err := r.db.QueryContext(ctx, studentLessonsGetSQL, courseID, studentID)
	if err != nil {
		return nil, fmt.Errorf("get student lessons: %w", err)
	}
	defer rows.Close()

	var lessons []student.StudentLesson

	for rows.Next() {
		var l student.StudentLesson

		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.OrderNo,
			&l.YoutubeURL,
			&l.Content,
			&l.Completed,
			&l.CompletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan student lesson: %w", err)
		}

		lessons = append(lessons, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student lessons: %w", err)
	}

	return lessons, nil
}

func (r *StudentRepository) SetLessonProgress(ctx context.Context, studentID uuid.UUID, lessonID uuid.UUID, completed bool) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, lessonProgressUpsertSQL, id, studentID, lessonID, completed); err != nil {
		return fmt.Errorf("set lesson progress: %w", err)
	}

	return nil
}
