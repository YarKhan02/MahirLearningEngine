package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/quiz"
	"github.com/google/uuid"
)

//go:embed sql/quiz_lesson_exists.sql
var quizLessonExistsSQL string

//go:embed sql/quiz_exists.sql
var quizExistsSQL string

//go:embed sql/quiz_lesson_id.sql
var quizLessonIDSQL string

//go:embed sql/quiz_insert.sql
var quizInsertSQL string

//go:embed sql/quiz_question_insert.sql
var quizQuestionInsertSQL string

//go:embed sql/quiz_option_insert.sql
var quizOptionInsertSQL string

//go:embed sql/quiz_delete.sql
var quizDeleteSQL string

//go:embed sql/quiz_update.sql
var quizUpdateSQL string

//go:embed sql/quiz_questions_delete_by_quiz.sql
var quizQuestionsDeleteByQuizSQL string

//go:embed sql/quizzes_by_lesson.sql
var quizzesByLessonSQL string

//go:embed sql/quiz_by_id.sql
var quizByIDSQL string

//go:embed sql/quiz_questions_by_lesson.sql
var quizQuestionsByLessonSQL string

//go:embed sql/quiz_options_by_lesson.sql
var quizOptionsByLessonSQL string

//go:embed sql/quiz_questions_by_quiz.sql
var quizQuestionsByQuizSQL string

//go:embed sql/quiz_options_by_quiz.sql
var quizOptionsByQuizSQL string

//go:embed sql/quiz_lesson_submission_counts.sql
var quizLessonSubmissionCountsSQL string

//go:embed sql/quiz_lesson_access.sql
var quizLessonAccessSQL string

//go:embed sql/student_by_user_get.sql
var quizStudentByUserSQL string

//go:embed sql/quiz_has_submitted.sql
var quizHasSubmittedSQL string

//go:embed sql/quiz_submission_insert.sql
var quizSubmissionInsertSQL string

//go:embed sql/quiz_answer_insert.sql
var quizAnswerInsertSQL string

//go:embed sql/quiz_answer_option_insert.sql
var quizAnswerOptionInsertSQL string

//go:embed sql/quiz_submissions_by_quiz.sql
var quizSubmissionsByQuizSQL string

//go:embed sql/quiz_submission_summary.sql
var quizSubmissionSummarySQL string

//go:embed sql/quiz_submission_by_id.sql
var quizSubmissionByIDSQL string

//go:embed sql/quiz_answers_by_submission.sql
var quizAnswersBySubmissionSQL string

//go:embed sql/quiz_answer_options_by_submission.sql
var quizAnswerOptionsBySubmissionSQL string

//go:embed sql/quiz_student_submissions_by_lesson.sql
var quizStudentSubmissionsByLessonSQL string

//go:embed sql/quiz_student_answers_by_lesson.sql
var quizStudentAnswersByLessonSQL string

//go:embed sql/quiz_student_answer_options_by_lesson.sql
var quizStudentAnswerOptionsByLessonSQL string

//go:embed sql/quiz_student_name.sql
var quizStudentNameSQL string

//go:embed sql/quiz_answer_grade.sql
var quizAnswerGradeSQL string

//go:embed sql/quiz_submission_grade.sql
var quizSubmissionGradeSQL string

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

/* ---------------- simple checks ---------------- */

func (r *QuizRepository) LessonExists(ctx context.Context, lessonID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, quizLessonExistsSQL, lessonID).Scan(&ok); err != nil {
		return false, fmt.Errorf("lesson exists: %w", err)
	}
	return ok, nil
}

func (r *QuizRepository) QuizExists(ctx context.Context, quizID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, quizExistsSQL, quizID).Scan(&ok); err != nil {
		return false, fmt.Errorf("quiz exists: %w", err)
	}
	return ok, nil
}

func (r *QuizRepository) GetQuizLessonID(ctx context.Context, quizID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.db.QueryRowContext(ctx, quizLessonIDSQL, quizID).Scan(&id); err != nil {
		return uuid.Nil, fmt.Errorf("quiz lesson id: %w", err)
	}
	return id, nil
}

func (r *QuizRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, quizStudentByUserSQL, userID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, quiz.ErrStudentNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("student by user: %w", err)
	}
	return id, nil
}

func (r *QuizRepository) UserHasLessonAccess(ctx context.Context, userID, lessonID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, quizLessonAccessSQL, userID, lessonID).Scan(&ok); err != nil {
		return false, fmt.Errorf("lesson access: %w", err)
	}
	return ok, nil
}

func (r *QuizRepository) GetStudentName(ctx context.Context, studentID uuid.UUID) (string, error) {
	var name string
	if err := r.db.QueryRowContext(ctx, quizStudentNameSQL, studentID).Scan(&name); err != nil {
		return "", fmt.Errorf("student name: %w", err)
	}
	return name, nil
}

func (r *QuizRepository) HasSubmitted(ctx context.Context, quizID, studentID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, quizHasSubmittedSQL, quizID, studentID).Scan(&ok); err != nil {
		return false, fmt.Errorf("has submitted: %w", err)
	}
	return ok, nil
}

func (r *QuizRepository) CreateQuiz(ctx context.Context, q quiz.NewQuiz) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	quizID := uuid.New()
	if _, err := tx.ExecContext(ctx, quizInsertSQL, quizID, q.LessonID, q.Title, nullStr(q.Description)); err != nil {
		return fmt.Errorf("insert quiz: %w", err)
	}
	if err := insertQuestions(ctx, tx, quizID, q.Questions); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *QuizRepository) EditQuiz(ctx context.Context, quizID uuid.UUID, q quiz.NewQuiz) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	res, err := tx.ExecContext(ctx, quizUpdateSQL, quizID, q.Title, nullStr(q.Description))
	if err != nil {
		return fmt.Errorf("update quiz: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update quiz rows: %w", err)
	}
	if n == 0 {
		return quiz.ErrQuizNotFound
	}

	if _, err := tx.ExecContext(ctx, quizQuestionsDeleteByQuizSQL, quizID); err != nil {
		return fmt.Errorf("clear questions: %w", err)
	}
	if err := insertQuestions(ctx, tx, quizID, q.Questions); err != nil {
		return err
	}

	return tx.Commit()
}

func insertQuestions(ctx context.Context, tx *sql.Tx, quizID uuid.UUID, questions []quiz.NewQuestion) error {
	for qi, question := range questions {
		questionID := uuid.New()
		if _, err := tx.ExecContext(ctx, quizQuestionInsertSQL,
			questionID, quizID, question.Prompt, question.Type, question.Marks, question.AllowMultiple, qi+1,
		); err != nil {
			return fmt.Errorf("insert question: %w", err)
		}
		for oi, opt := range question.Options {
			if _, err := tx.ExecContext(ctx, quizOptionInsertSQL,
				uuid.New(), questionID, opt.Text, opt.IsCorrect, oi+1,
			); err != nil {
				return fmt.Errorf("insert option: %w", err)
			}
		}
	}
	return nil
}

func (r *QuizRepository) DeleteQuiz(ctx context.Context, quizID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, quizDeleteSQL, quizID)
	if err != nil {
		return fmt.Errorf("delete quiz: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete quiz: rows: %w", err)
	}
	if n == 0 {
		return quiz.ErrQuizNotFound
	}
	return nil
}

type questionRow struct {
	quizID, id    uuid.UUID
	prompt, typ   string
	marks         int
	allowMultiple bool
	orderNo       int
}

type optionRow struct {
	questionID, id uuid.UUID
	text           string
	isCorrect      bool
	orderNo        int
}

func buildQuizzes(quizzes []quiz.Quiz, qrows []questionRow, orows []optionRow) []quiz.Quiz {
	optsByQ := make(map[uuid.UUID][]quiz.Option)
	for _, o := range orows {
		optsByQ[o.questionID] = append(optsByQ[o.questionID], quiz.Option{
			ID: o.id, Text: o.text, IsCorrect: o.isCorrect, OrderNo: o.orderNo,
		})
	}
	qByQuiz := make(map[uuid.UUID][]quiz.Question)
	for _, q := range qrows {
		qByQuiz[q.quizID] = append(qByQuiz[q.quizID], quiz.Question{
			ID: q.id, Prompt: q.prompt, Type: q.typ, Marks: q.marks,
			AllowMultiple: q.allowMultiple, OrderNo: q.orderNo, Options: optsByQ[q.id],
		})
	}
	for i := range quizzes {
		quizzes[i].Questions = qByQuiz[quizzes[i].ID]
	}
	return quizzes
}

func (r *QuizRepository) GetQuizzesByLesson(ctx context.Context, lessonID uuid.UUID) ([]quiz.Quiz, error) {
	quizzes, err := r.scanQuizRows(ctx, quizzesByLessonSQL, lessonID)
	if err != nil {
		return nil, err
	}
	if len(quizzes) == 0 {
		return quizzes, nil
	}

	qrows, err := r.scanQuestionRows(ctx, quizQuestionsByLessonSQL, true, lessonID)
	if err != nil {
		return nil, err
	}
	orows, err := r.scanOptionRows(ctx, quizOptionsByLessonSQL, lessonID)
	if err != nil {
		return nil, err
	}
	return buildQuizzes(quizzes, qrows, orows), nil
}

func (r *QuizRepository) GetQuizByID(ctx context.Context, quizID uuid.UUID) (quiz.Quiz, error) {
	var q quiz.Quiz
	err := r.db.QueryRowContext(ctx, quizByIDSQL, quizID).Scan(
		&q.ID, &q.LessonID, &q.Title, &q.Description, &q.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return quiz.Quiz{}, quiz.ErrQuizNotFound
	}
	if err != nil {
		return quiz.Quiz{}, fmt.Errorf("quiz by id: %w", err)
	}

	qrows, err := r.scanQuestionRows(ctx, quizQuestionsByQuizSQL, false, quizID)
	if err != nil {
		return quiz.Quiz{}, err
	}
	for i := range qrows {
		qrows[i].quizID = quizID
	}
	orows, err := r.scanOptionRows(ctx, quizOptionsByQuizSQL, quizID)
	if err != nil {
		return quiz.Quiz{}, err
	}
	return buildQuizzes([]quiz.Quiz{q}, qrows, orows)[0], nil
}

func (r *QuizRepository) scanQuizRows(ctx context.Context, query string, args ...any) ([]quiz.Quiz, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get quizzes: %w", err)
	}
	defer rows.Close()

	quizzes := make([]quiz.Quiz, 0)
	for rows.Next() {
		var q quiz.Quiz
		if err := rows.Scan(&q.ID, &q.LessonID, &q.Title, &q.Description, &q.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan quiz: %w", err)
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, rows.Err()
}

func (r *QuizRepository) scanQuestionRows(ctx context.Context, query string, withQuizID bool, args ...any) ([]questionRow, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get questions: %w", err)
	}
	defer rows.Close()

	var out []questionRow
	for rows.Next() {
		var q questionRow
		if withQuizID {
			err = rows.Scan(&q.quizID, &q.id, &q.prompt, &q.typ, &q.marks, &q.allowMultiple, &q.orderNo)
		} else {
			err = rows.Scan(&q.id, &q.prompt, &q.typ, &q.marks, &q.allowMultiple, &q.orderNo)
		}
		if err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func (r *QuizRepository) scanOptionRows(ctx context.Context, query string, args ...any) ([]optionRow, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get options: %w", err)
	}
	defer rows.Close()

	var out []optionRow
	for rows.Next() {
		var o optionRow
		if err := rows.Scan(&o.questionID, &o.id, &o.text, &o.isCorrect, &o.orderNo); err != nil {
			return nil, fmt.Errorf("scan option: %w", err)
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *QuizRepository) GetLessonSubmissionCounts(ctx context.Context, lessonID uuid.UUID) (map[uuid.UUID]quiz.SubmissionSummary, error) {
	rows, err := r.db.QueryContext(ctx, quizLessonSubmissionCountsSQL, lessonID)
	if err != nil {
		return nil, fmt.Errorf("submission counts: %w", err)
	}
	defer rows.Close()

	out := make(map[uuid.UUID]quiz.SubmissionSummary)
	for rows.Next() {
		var (
			quizID uuid.UUID
			s      quiz.SubmissionSummary
		)
		if err := rows.Scan(&quizID, &s.Total, &s.Pending, &s.Graded); err != nil {
			return nil, fmt.Errorf("scan counts: %w", err)
		}
		out[quizID] = s
	}
	return out, rows.Err()
}

func (r *QuizRepository) CreateSubmission(ctx context.Context, sub quiz.Submission) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	subID := uuid.New()
	if _, err := tx.ExecContext(ctx, quizSubmissionInsertSQL,
		subID, sub.QuizID, sub.StudentID, sub.Status, sub.Score, nullTime(sub.GradedAt),
	); err != nil {
		return fmt.Errorf("insert submission: %w", err)
	}

	for _, a := range sub.Answers {
		answerID := uuid.New()
		if _, err := tx.ExecContext(ctx, quizAnswerInsertSQL,
			answerID, subID, a.QuestionID, nullStr(a.AnswerText), nullIntPtr(a.AwardedMarks),
		); err != nil {
			return fmt.Errorf("insert answer: %w", err)
		}
		for _, optID := range a.SelectedOpts {
			if _, err := tx.ExecContext(ctx, quizAnswerOptionInsertSQL, answerID, optID); err != nil {
				return fmt.Errorf("insert answer option: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *QuizRepository) GetSubmissionsByQuiz(ctx context.Context, quizID uuid.UUID) ([]quiz.SubmissionRow, error) {
	rows, err := r.db.QueryContext(ctx, quizSubmissionsByQuizSQL, quizID)
	if err != nil {
		return nil, fmt.Errorf("submissions by quiz: %w", err)
	}
	defer rows.Close()

	var out []quiz.SubmissionRow
	for rows.Next() {
		var s quiz.SubmissionRow
		if err := rows.Scan(&s.ID, &s.StudentID, &s.StudentName, &s.StudentEmail, &s.Status, &s.SubmittedAt, &s.Score); err != nil {
			return nil, fmt.Errorf("scan submission row: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *QuizRepository) GetSubmissionSummary(ctx context.Context, quizID uuid.UUID) (quiz.SubmissionSummary, error) {
	var s quiz.SubmissionSummary
	if err := r.db.QueryRowContext(ctx, quizSubmissionSummarySQL, quizID).Scan(&s.Total, &s.Pending, &s.Graded); err != nil {
		return quiz.SubmissionSummary{}, fmt.Errorf("submission summary: %w", err)
	}
	return s, nil
}

func (r *QuizRepository) GetSubmission(ctx context.Context, submissionID uuid.UUID) (quiz.Submission, error) {
	var (
		s        quiz.Submission
		remarks  sql.NullString
		gradedAt sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, quizSubmissionByIDSQL, submissionID).Scan(
		&s.ID, &s.QuizID, &s.StudentID, &s.Status, &s.Score, &remarks, &s.SubmittedAt, &gradedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return quiz.Submission{}, quiz.ErrSubmissionNotFound
	}
	if err != nil {
		return quiz.Submission{}, fmt.Errorf("submission by id: %w", err)
	}
	if remarks.Valid {
		s.Remarks = &remarks.String
	}
	if gradedAt.Valid {
		s.GradedAt = &gradedAt.Time
	}

	answers, byAnswerID, err := r.scanAnswers(ctx, quizAnswersBySubmissionSQL, submissionID)
	if err != nil {
		return quiz.Submission{}, err
	}
	if err := r.attachAnswerOptions(ctx, quizAnswerOptionsBySubmissionSQL, byAnswerID, submissionID); err != nil {
		return quiz.Submission{}, err
	}
	for _, a := range answers {
		s.Answers = append(s.Answers, *a)
	}
	return s, nil
}

func (r *QuizRepository) GetStudentSubmissionsForLesson(ctx context.Context, lessonID, studentID uuid.UUID) (map[uuid.UUID]quiz.Submission, error) {
	rows, err := r.db.QueryContext(ctx, quizStudentSubmissionsByLessonSQL, lessonID, studentID)
	if err != nil {
		return nil, fmt.Errorf("student submissions: %w", err)
	}
	defer rows.Close()

	subsByID := make(map[uuid.UUID]*quiz.Submission)
	quizToSub := make(map[uuid.UUID]uuid.UUID)
	for rows.Next() {
		var (
			s        quiz.Submission
			remarks  sql.NullString
			gradedAt sql.NullTime
		)
		if err := rows.Scan(&s.ID, &s.QuizID, &s.StudentID, &s.Status, &s.Score, &remarks, &s.SubmittedAt, &gradedAt); err != nil {
			return nil, fmt.Errorf("scan student submission: %w", err)
		}
		if remarks.Valid {
			s.Remarks = &remarks.String
		}
		if gradedAt.Valid {
			s.GradedAt = &gradedAt.Time
		}
		cp := s
		subsByID[s.ID] = &cp
		quizToSub[s.QuizID] = s.ID
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	answers, byAnswerID, err := r.scanAnswersMulti(ctx, quizStudentAnswersByLessonSQL, lessonID, studentID)
	if err != nil {
		return nil, err
	}
	if err := r.attachAnswerOptions(ctx, quizStudentAnswerOptionsByLessonSQL, byAnswerID, lessonID, studentID); err != nil {
		return nil, err
	}
	for subID, list := range answers {
		if sub := subsByID[subID]; sub != nil {
			for _, a := range list {
				sub.Answers = append(sub.Answers, *a)
			}
		}
	}

	out := make(map[uuid.UUID]quiz.Submission)
	for quizID, subID := range quizToSub {
		if sub := subsByID[subID]; sub != nil {
			out[quizID] = *sub
		}
	}
	return out, nil
}

func (r *QuizRepository) scanAnswers(ctx context.Context, query string, args ...any) ([]*quiz.Answer, map[uuid.UUID]*quiz.Answer, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("get answers: %w", err)
	}
	defer rows.Close()

	var ordered []*quiz.Answer
	byID := make(map[uuid.UUID]*quiz.Answer)
	for rows.Next() {
		var (
			answerID uuid.UUID
			a        quiz.Answer
			marks    sql.NullInt64
		)
		if err := rows.Scan(&answerID, &a.QuestionID, &a.AnswerText, &marks); err != nil {
			return nil, nil, fmt.Errorf("scan answer: %w", err)
		}
		if marks.Valid {
			m := int(marks.Int64)
			a.AwardedMarks = &m
		}
		cp := a
		ordered = append(ordered, &cp)
		byID[answerID] = &cp
	}
	return ordered, byID, rows.Err()
}

func (r *QuizRepository) scanAnswersMulti(ctx context.Context, query string, args ...any) (map[uuid.UUID][]*quiz.Answer, map[uuid.UUID]*quiz.Answer, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("get answers: %w", err)
	}
	defer rows.Close()

	bySub := make(map[uuid.UUID][]*quiz.Answer)
	byID := make(map[uuid.UUID]*quiz.Answer)
	for rows.Next() {
		var (
			answerID, submissionID uuid.UUID
			a                      quiz.Answer
			marks                  sql.NullInt64
		)
		if err := rows.Scan(&answerID, &submissionID, &a.QuestionID, &a.AnswerText, &marks); err != nil {
			return nil, nil, fmt.Errorf("scan answer: %w", err)
		}
		if marks.Valid {
			m := int(marks.Int64)
			a.AwardedMarks = &m
		}
		cp := a
		bySub[submissionID] = append(bySub[submissionID], &cp)
		byID[answerID] = &cp
	}
	return bySub, byID, rows.Err()
}

func (r *QuizRepository) attachAnswerOptions(ctx context.Context, query string, byAnswerID map[uuid.UUID]*quiz.Answer, args ...any) error {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("get answer options: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var answerID, optionID uuid.UUID
		if err := rows.Scan(&answerID, &optionID); err != nil {
			return fmt.Errorf("scan answer option: %w", err)
		}
		if a := byAnswerID[answerID]; a != nil {
			a.SelectedOpts = append(a.SelectedOpts, optionID)
		}
	}
	return rows.Err()
}

func (r *QuizRepository) GradeSubmission(ctx context.Context, in quiz.GradeInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	for _, m := range in.Marks {
		if _, err := tx.ExecContext(ctx, quizAnswerGradeSQL, m.Marks, in.SubmissionID, m.QuestionID); err != nil {
			return fmt.Errorf("grade answer: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, quizSubmissionGradeSQL, nullStr(in.Remarks), in.SubmissionID); err != nil {
		return fmt.Errorf("grade submission: %w", err)
	}

	return tx.Commit()
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullIntPtr(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return *t
}
