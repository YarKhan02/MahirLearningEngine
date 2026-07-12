package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/google/uuid"
)

//go:embed sql/assignment_create.sql
var assignmentCreateSQL string

//go:embed sql/assignments_by_lesson.sql
var assignmentsByLessonSQL string

//go:embed sql/assignment_delete.sql
var assignmentDeleteSQL string

//go:embed sql/lesson_access_check.sql
var lessonAccessCheckSQL string

//go:embed sql/assignment_access_check.sql
var assignmentAccessCheckSQL string

//go:embed sql/student_assignments_by_lesson.sql
var studentAssignmentsByLessonSQL string

//go:embed sql/submission_upsert.sql
var submissionUpsertSQL string

//go:embed sql/submissions_by_batch.sql
var submissionsByBatchSQL string

//go:embed sql/submission_grade.sql
var submissionGradeSQL string

//go:embed sql/submissions_by_student.sql
var submissionsByStudentSQL string

//go:embed sql/student_by_user_get.sql
var assignmentStudentByUserSQL string

type AssignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

func (r *AssignmentRepository) CreateAssignment(ctx context.Context, a *assignment.Assignment) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	a.ID = id

	_, err = r.db.ExecContext(
		ctx,
		assignmentCreateSQL,
		a.ID,
		a.LessonID,
		a.Title,
		a.Description,
		a.StarterCode,
		a.DueDate,
		a.TotalMarks,
	)
	if err != nil {
		return fmt.Errorf("create assignment: %w", err)
	}

	return nil
}

func (r *AssignmentRepository) GetLessonAssignments(ctx context.Context, lessonID uuid.UUID) ([]assignment.Assignment, error) {
	rows, err := r.db.QueryContext(ctx, assignmentsByLessonSQL, lessonID)
	if err != nil {
		return nil, fmt.Errorf("get assignments: %w", err)
	}
	defer rows.Close()

	var assignments []assignment.Assignment

	for rows.Next() {
		var a assignment.Assignment

		if err := rows.Scan(
			&a.ID,
			&a.LessonID,
			&a.Title,
			&a.Description,
			&a.StarterCode,
			&a.DueDate,
			&a.TotalMarks,
			&a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan assignment: %w", err)
		}

		assignments = append(assignments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assignments: %w", err)
	}

	return assignments, nil
}

func (r *AssignmentRepository) DeleteAssignment(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, assignmentDeleteSQL, id); err != nil {
		return fmt.Errorf("delete assignment: %w", err)
	}
	return nil
}

func (r *AssignmentRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx, assignmentStudentByUserSQL, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, assignment.ErrStudentNotFound
		}
		return uuid.Nil, fmt.Errorf("get student by user: %w", err)
	}

	return id, nil
}

func (r *AssignmentRepository) HasLessonAccess(ctx context.Context, studentID uuid.UUID, lessonID uuid.UUID) (bool, error) {
	var hasAccess bool

	err := r.db.QueryRowContext(ctx, lessonAccessCheckSQL, studentID, lessonID).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("check lesson access: %w", err)
	}

	return hasAccess, nil
}

func (r *AssignmentRepository) HasAssignmentAccess(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID) (bool, error) {
	var hasAccess bool

	err := r.db.QueryRowContext(ctx, assignmentAccessCheckSQL, studentID, assignmentID).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("check assignment access: %w", err)
	}

	return hasAccess, nil
}

func (r *AssignmentRepository) GetStudentAssignments(ctx context.Context, lessonID uuid.UUID, studentID uuid.UUID) ([]assignment.StudentAssignment, error) {
	rows, err := r.db.QueryContext(ctx, studentAssignmentsByLessonSQL, lessonID, studentID)
	if err != nil {
		return nil, fmt.Errorf("get student assignments: %w", err)
	}
	defer rows.Close()

	var assignments []assignment.StudentAssignment

	for rows.Next() {
		var a assignment.StudentAssignment
		var (
			subID          sql.Null[uuid.UUID]
			subCode        sql.NullString
			subRemarks     sql.NullString
			subMarks       sql.NullInt64
			subStatus      sql.NullString
			subSubmittedAt sql.NullTime
		)

		if err := rows.Scan(
			&a.ID,
			&a.LessonID,
			&a.Title,
			&a.Description,
			&a.StarterCode,
			&a.DueDate,
			&a.TotalMarks,
			&a.CreatedAt,
			&subID,
			&subCode,
			&subRemarks,
			&subMarks,
			&subStatus,
			&subSubmittedAt,
		); err != nil {
			return nil, fmt.Errorf("scan student assignment: %w", err)
		}

		if subID.Valid {
			sub := &assignment.Submission{
				ID:           subID.V,
				StudentID:    studentID,
				AssignmentID: a.ID,
				Code:         subCode.String,
				Status:       subStatus.String,
				SubmittedAt:  subSubmittedAt.Time,
			}
			if subRemarks.Valid {
				sub.Remarks = &subRemarks.String
			}
			if subMarks.Valid {
				m := int(subMarks.Int64)
				sub.Marks = &m
			}
			a.Submission = sub
		}

		assignments = append(assignments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student assignments: %w", err)
	}

	return assignments, nil
}

func (r *AssignmentRepository) SubmitAssignment(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID, code string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, submissionUpsertSQL, id, studentID, assignmentID, code); err != nil {
		return fmt.Errorf("submit assignment: %w", err)
	}

	return nil
}

func (r *AssignmentRepository) GetBatchSubmissions(ctx context.Context, batchID uuid.UUID) ([]assignment.BatchSubmission, error) {
	rows, err := r.db.QueryContext(ctx, submissionsByBatchSQL, batchID)
	if err != nil {
		return nil, fmt.Errorf("get batch submissions: %w", err)
	}
	defer rows.Close()

	var submissions []assignment.BatchSubmission

	for rows.Next() {
		var s assignment.BatchSubmission

		if err := rows.Scan(
			&s.ID,
			&s.Code,
			&s.Remarks,
			&s.Marks,
			&s.Status,
			&s.SubmittedAt,
			&s.StudentID,
			&s.StudentName,
			&s.StudentEmail,
			&s.AssignmentID,
			&s.AssignmentTitle,
			&s.TotalMarks,
			&s.LessonID,
			&s.LessonTitle,
			&s.CourseID,
			&s.CourseTitle,
		); err != nil {
			return nil, fmt.Errorf("scan batch submission: %w", err)
		}

		submissions = append(submissions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate batch submissions: %w", err)
	}

	return submissions, nil
}

func (r *AssignmentRepository) GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error {
	if _, err := r.db.ExecContext(ctx, submissionGradeSQL, submissionID, marks, remarks); err != nil {
		return fmt.Errorf("grade submission: %w", err)
	}
	return nil
}

func (r *AssignmentRepository) GetStudentSubmissions(ctx context.Context, studentID uuid.UUID) ([]assignment.BatchSubmission, error) {
	rows, err := r.db.QueryContext(ctx, submissionsByStudentSQL, studentID)
	if err != nil {
		return nil, fmt.Errorf("get student submissions: %w", err)
	}
	defer rows.Close()

	var submissions []assignment.BatchSubmission

	for rows.Next() {
		var s assignment.BatchSubmission

		if err := rows.Scan(
			&s.ID,
			&s.Code,
			&s.Remarks,
			&s.Marks,
			&s.Status,
			&s.SubmittedAt,
			&s.StudentID,
			&s.StudentName,
			&s.StudentEmail,
			&s.AssignmentID,
			&s.AssignmentTitle,
			&s.TotalMarks,
			&s.LessonID,
			&s.LessonTitle,
			&s.CourseID,
			&s.CourseTitle,
		); err != nil {
			return nil, fmt.Errorf("scan student submission: %w", err)
		}

		submissions = append(submissions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student submissions: %w", err)
	}

	return submissions, nil
}
