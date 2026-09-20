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

//go:embed sql/assignment_get.sql
var assignmentGetSQL string

//go:embed sql/assignment_update.sql
var assignmentUpdateSQL string

//go:embed sql/test_cases_delete_by_assignment.sql
var testCasesDeleteByAssignmentSQL string

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

//go:embed sql/submissions_by_batch_count.sql
var submissionsByBatchCountSQL string

//go:embed sql/submissions_by_batch_summary.sql
var submissionsByBatchSummarySQL string

//go:embed sql/submission_grade.sql
var submissionGradeSQL string

//go:embed sql/submissions_by_student.sql
var submissionsByStudentSQL string

//go:embed sql/submissions_by_student_count.sql
var submissionsByStudentCountSQL string

//go:embed sql/submissions_by_student_summary.sql
var submissionsByStudentSummarySQL string

//go:embed sql/student_by_user_get.sql
var assignmentStudentByUserSQL string

//go:embed sql/test_case_create.sql
var testCaseCreateSQL string

//go:embed sql/assignment_grading_meta.sql
var assignmentGradingMetaSQL string

//go:embed sql/test_cases_for_grading.sql
var testCasesForGradingSQL string

//go:embed sql/submission_autograde_save.sql
var submissionAutogradeSaveSQL string

//go:embed sql/submission_results_delete.sql
var submissionResultsDeleteSQL string

//go:embed sql/submission_result_insert.sql
var submissionResultInsertSQL string

//go:embed sql/submission_results_by_submission.sql
var submissionResultsBySubmissionSQL string

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

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create assignment: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.ExecContext(
		ctx,
		assignmentCreateSQL,
		a.ID,
		a.LessonID,
		a.Title,
		a.Description,
		a.StarterCode,
		a.Language,
		a.DueDate,
		a.TotalMarks,
	); err != nil {
		return fmt.Errorf("create assignment: %w", err)
	}

	for _, tc := range a.TestCases {
		tcID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, testCaseCreateSQL,
			tcID, a.ID, tc.Stdin, tc.ExpectedStdout, tc.Weight, tc.Ordinal,
		); err != nil {
			return fmt.Errorf("create test case: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create assignment: commit: %w", err)
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

func (r *AssignmentRepository) GetAssignmentWithTests(ctx context.Context, id uuid.UUID) (*assignment.Assignment, error) {
	var a assignment.Assignment
	err := r.db.QueryRowContext(ctx, assignmentGetSQL, id).Scan(
		&a.ID,
		&a.LessonID,
		&a.Title,
		&a.Description,
		&a.StarterCode,
		&a.Language,
		&a.DueDate,
		&a.TotalMarks,
		&a.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, assignment.ErrAssignmentNotFound
		}
		return nil, fmt.Errorf("get assignment: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, testCasesForGradingSQL, id)
	if err != nil {
		return nil, fmt.Errorf("get assignment test cases: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t assignment.TestCase
		if err := rows.Scan(&t.ID, &t.Stdin, &t.ExpectedStdout, &t.Weight, &t.Ordinal); err != nil {
			return nil, fmt.Errorf("scan test case: %w", err)
		}
		t.AssignmentID = id
		a.TestCases = append(a.TestCases, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate test cases: %w", err)
	}
	return &a, nil
}

func (r *AssignmentRepository) UpdateAssignment(ctx context.Context, a *assignment.Assignment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update assignment: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(
		ctx,
		assignmentUpdateSQL,
		a.ID,
		a.Title,
		a.Description,
		a.StarterCode,
		a.Language,
		a.DueDate,
		a.TotalMarks,
	)
	if err != nil {
		return fmt.Errorf("update assignment: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return assignment.ErrAssignmentNotFound
	}

	if _, err := tx.ExecContext(ctx, testCasesDeleteByAssignmentSQL, a.ID); err != nil {
		return fmt.Errorf("update assignment: clear test cases: %w", err)
	}
	for _, tc := range a.TestCases {
		tcID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, testCaseCreateSQL,
			tcID, a.ID, tc.Stdin, tc.ExpectedStdout, tc.Weight, tc.Ordinal,
		); err != nil {
			return fmt.Errorf("update assignment: create test case: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("update assignment: commit: %w", err)
	}
	return nil
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
			subAutoScore   sql.NullInt64
			subTestsTotal  sql.NullInt64
			subTestsPassed sql.NullInt64
			subStatus      sql.NullString
			subSubmittedAt sql.NullTime
		)

		if err := rows.Scan(
			&a.ID,
			&a.LessonID,
			&a.Title,
			&a.Description,
			&a.StarterCode,
			&a.Language,
			&a.DueDate,
			&a.TotalMarks,
			&a.CreatedAt,
			&subID,
			&subCode,
			&subRemarks,
			&subMarks,
			&subAutoScore,
			&subTestsTotal,
			&subTestsPassed,
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
			if subAutoScore.Valid {
				v := int(subAutoScore.Int64)
				sub.AutoScore = &v
			}
			if subTestsTotal.Valid {
				v := int(subTestsTotal.Int64)
				sub.TestsTotal = &v
			}
			if subTestsPassed.Valid {
				v := int(subTestsPassed.Int64)
				sub.TestsPassed = &v
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

func (r *AssignmentRepository) SubmitAssignment(ctx context.Context, studentID uuid.UUID, assignmentID uuid.UUID, code string) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	var submissionID uuid.UUID
	if err := r.db.QueryRowContext(ctx, submissionUpsertSQL, id, studentID, assignmentID, code).Scan(&submissionID); err != nil {
		return uuid.Nil, fmt.Errorf("submit assignment: %w", err)
	}
	return submissionID, nil
}

// GetGradingData returns the assignment's language, total marks, and hidden test
// cases (with expected outputs) for autograding.
func (r *AssignmentRepository) GetGradingData(ctx context.Context, assignmentID uuid.UUID) (string, int, []assignment.TestCase, error) {
	var language string
	var totalMarks int
	if err := r.db.QueryRowContext(ctx, assignmentGradingMetaSQL, assignmentID).Scan(&language, &totalMarks); err != nil {
		return "", 0, nil, fmt.Errorf("grading meta: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, testCasesForGradingSQL, assignmentID)
	if err != nil {
		return "", 0, nil, fmt.Errorf("grading test cases: %w", err)
	}
	defer rows.Close()

	var tests []assignment.TestCase
	for rows.Next() {
		var t assignment.TestCase
		if err := rows.Scan(&t.ID, &t.Stdin, &t.ExpectedStdout, &t.Weight, &t.Ordinal); err != nil {
			return "", 0, nil, fmt.Errorf("scan test case: %w", err)
		}
		t.AssignmentID = assignmentID
		tests = append(tests, t)
	}
	return language, totalMarks, tests, rows.Err()
}

// SaveAutoGrade stores the autoscore/counts and replaces the per-test results.
func (r *AssignmentRepository) SaveAutoGrade(ctx context.Context, submissionID uuid.UUID, autoScore, testsTotal, testsPassed int, results []assignment.TestResult) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("save autograde: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, submissionAutogradeSaveSQL, submissionID, autoScore, testsTotal, testsPassed); err != nil {
		return fmt.Errorf("save autograde: update: %w", err)
	}
	if _, err := tx.ExecContext(ctx, submissionResultsDeleteSQL, submissionID); err != nil {
		return fmt.Errorf("save autograde: clear results: %w", err)
	}
	for _, res := range results {
		rID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, submissionResultInsertSQL,
			rID, submissionID, res.TestCaseID, res.Passed, res.ActualStdout, res.Stderr, res.TimedOut, res.DurationMs, res.Ordinal,
		); err != nil {
			return fmt.Errorf("save autograde: insert result: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("save autograde: commit: %w", err)
	}
	return nil
}

func (r *AssignmentRepository) GetSubmissionResults(ctx context.Context, submissionID uuid.UUID) ([]assignment.TestResult, error) {
	rows, err := r.db.QueryContext(ctx, submissionResultsBySubmissionSQL, submissionID)
	if err != nil {
		return nil, fmt.Errorf("get submission results: %w", err)
	}
	defer rows.Close()

	var out []assignment.TestResult
	for rows.Next() {
		var res assignment.TestResult
		if err := rows.Scan(
			&res.TestCaseID, &res.Passed, &res.ActualStdout, &res.Stderr, &res.TimedOut, &res.DurationMs, &res.Ordinal,
		); err != nil {
			return nil, fmt.Errorf("scan result: %w", err)
		}
		out = append(out, res)
	}
	return out, rows.Err()
}

func (r *AssignmentRepository) GetBatchSubmissions(ctx context.Context, batchID uuid.UUID, q, status string, limit, offset int) ([]assignment.BatchSubmission, error) {
	rows, err := r.db.QueryContext(ctx, submissionsByBatchSQL, batchID, q, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get batch submissions: %w", err)
	}
	defer rows.Close()

	return scanBatchSubmissions(rows)
}

func (r *AssignmentRepository) CountBatchSubmissions(ctx context.Context, batchID uuid.UUID, q, status string) (int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, submissionsByBatchCountSQL, batchID, q, status).Scan(&total); err != nil {
		return 0, fmt.Errorf("count batch submissions: %w", err)
	}
	return total, nil
}

func (r *AssignmentRepository) GetBatchSubmissionSummary(ctx context.Context, batchID uuid.UUID, q string) (assignment.SubmissionSummary, error) {
	var s assignment.SubmissionSummary
	if err := r.db.QueryRowContext(ctx, submissionsByBatchSummarySQL, batchID, q).Scan(
		&s.Total, &s.Submitted, &s.Graded,
	); err != nil {
		return assignment.SubmissionSummary{}, fmt.Errorf("batch submission summary: %w", err)
	}
	return s, nil
}

func (r *AssignmentRepository) GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error {
	if _, err := r.db.ExecContext(ctx, submissionGradeSQL, submissionID, marks, remarks); err != nil {
		return fmt.Errorf("grade submission: %w", err)
	}
	return nil
}

func (r *AssignmentRepository) GetStudentSubmissions(ctx context.Context, studentID uuid.UUID, status string, limit, offset int) ([]assignment.BatchSubmission, error) {
	rows, err := r.db.QueryContext(ctx, submissionsByStudentSQL, studentID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get student submissions: %w", err)
	}
	defer rows.Close()

	return scanBatchSubmissions(rows)
}

func (r *AssignmentRepository) CountStudentSubmissions(ctx context.Context, studentID uuid.UUID, status string) (int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, submissionsByStudentCountSQL, studentID, status).Scan(&total); err != nil {
		return 0, fmt.Errorf("count student submissions: %w", err)
	}
	return total, nil
}

func (r *AssignmentRepository) GetStudentSubmissionSummary(ctx context.Context, studentID uuid.UUID) (assignment.SubmissionSummary, error) {
	var s assignment.SubmissionSummary
	if err := r.db.QueryRowContext(ctx, submissionsByStudentSummarySQL, studentID).Scan(
		&s.Total, &s.Submitted, &s.Graded,
	); err != nil {
		return assignment.SubmissionSummary{}, fmt.Errorf("student submission summary: %w", err)
	}
	return s, nil
}

func scanBatchSubmissions(rows *sql.Rows) ([]assignment.BatchSubmission, error) {
	var submissions []assignment.BatchSubmission

	for rows.Next() {
		var s assignment.BatchSubmission

		var (
			autoScore   sql.NullInt64
			testsTotal  sql.NullInt64
			testsPassed sql.NullInt64
		)
		if err := rows.Scan(
			&s.ID,
			&s.Code,
			&s.Remarks,
			&s.Marks,
			&autoScore,
			&testsTotal,
			&testsPassed,
			&s.Language,
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
			return nil, fmt.Errorf("scan submission: %w", err)
		}
		if autoScore.Valid {
			v := int(autoScore.Int64)
			s.AutoScore = &v
		}
		if testsTotal.Valid {
			v := int(testsTotal.Int64)
			s.TestsTotal = &v
		}
		if testsPassed.Valid {
			v := int(testsPassed.Int64)
			s.TestsPassed = &v
		}

		submissions = append(submissions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate submissions: %w", err)
	}

	return submissions, nil
}
