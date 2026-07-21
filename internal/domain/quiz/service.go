package quiz

import (
	"context"
	"errors"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrLessonNotFound     = errors.New("lesson not found")
	ErrQuizNotFound       = errors.New("quiz not found")
	ErrStudentNotFound    = errors.New("student not found")
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrForbidden          = errors.New("forbidden")
	ErrAlreadySubmitted   = errors.New("already submitted")
	ErrInvalidQuiz        = errors.New("invalid quiz")
	ErrInvalidSubmission  = errors.New("invalid submission")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateQuiz(ctx context.Context, lessonID uuid.UUID, req CreateQuizRequest) error {
	if err := validateCreate(req); err != nil {
		return err
	}

	ok, err := s.repo.LessonExists(ctx, lessonID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrLessonNotFound
	}

	if err := s.repo.CreateQuiz(ctx, ToNewQuiz(req, lessonID)); err != nil {
		return err
	}

	logging.FromLogger(ctx).Info("quiz created",
		zap.String("event", "quiz_created"),
		zap.String("lesson_id", lessonID.String()),
		zap.String("title", req.Title),
		zap.Int("questions", len(req.Questions)),
	)
	return nil
}

func (s *Service) EditQuiz(ctx context.Context, quizID uuid.UUID, req CreateQuizRequest) error {
	if err := validateCreate(req); err != nil {
		return err
	}
	if err := s.repo.EditQuiz(ctx, quizID, ToNewQuiz(req, uuid.Nil)); err != nil {
		return err
	}
	logging.FromLogger(ctx).Info("quiz updated",
		zap.String("event", "quiz_updated"),
		zap.String("quiz_id", quizID.String()),
		zap.Int("questions", len(req.Questions)),
	)
	return nil
}

func (s *Service) ListForLesson(ctx context.Context, lessonID uuid.UUID) ([]Quiz, map[uuid.UUID]SubmissionSummary, error) {
	quizzes, err := s.repo.GetQuizzesByLesson(ctx, lessonID)
	if err != nil {
		return nil, nil, err
	}
	counts, err := s.repo.GetLessonSubmissionCounts(ctx, lessonID)
	if err != nil {
		return nil, nil, err
	}
	return quizzes, counts, nil
}

func (s *Service) DeleteQuiz(ctx context.Context, quizID uuid.UUID) error {
	if err := s.repo.DeleteQuiz(ctx, quizID); err != nil {
		return err
	}
	logging.FromLogger(ctx).Info("quiz deleted",
		zap.String("event", "quiz_deleted"),
		zap.String("quiz_id", quizID.String()),
	)
	return nil
}

func (s *Service) SubmissionsForQuiz(ctx context.Context, quizID uuid.UUID) ([]SubmissionRow, SubmissionSummary, error) {
	rows, err := s.repo.GetSubmissionsByQuiz(ctx, quizID)
	if err != nil {
		return nil, SubmissionSummary{}, err
	}
	summary, err := s.repo.GetSubmissionSummary(ctx, quizID)
	if err != nil {
		return nil, SubmissionSummary{}, err
	}
	return rows, summary, nil
}

func (s *Service) SubmissionDetail(ctx context.Context, submissionID uuid.UUID) (Quiz, Submission, string, error) {
	sub, err := s.repo.GetSubmission(ctx, submissionID)
	if err != nil {
		return Quiz{}, Submission{}, "", err
	}
	q, err := s.repo.GetQuizByID(ctx, sub.QuizID)
	if err != nil {
		return Quiz{}, Submission{}, "", err
	}
	name, err := s.repo.GetStudentName(ctx, sub.StudentID)
	if err != nil {
		return Quiz{}, Submission{}, "", err
	}
	return q, sub, name, nil
}

func (s *Service) Grade(ctx context.Context, in GradeInput) error {
	sub, err := s.repo.GetSubmission(ctx, in.SubmissionID)
	if err != nil {
		return err
	}
	if err := s.repo.GradeSubmission(ctx, in); err != nil {
		return err
	}
	logging.FromLogger(ctx).Info("quiz submission graded",
		zap.String("event", "quiz_graded"),
		zap.String("submission_id", in.SubmissionID.String()),
		zap.String("quiz_id", sub.QuizID.String()),
		zap.Int("graded_answers", len(in.Marks)),
	)
	return nil
}

func (s *Service) ListForStudent(ctx context.Context, userID, lessonID uuid.UUID) ([]StudentQuiz, error) {
	ok, err := s.repo.UserHasLessonAccess(ctx, userID, lessonID)
	if err != nil {
		return nil, err
	}
	if !ok {
		logging.FromLogger(ctx).Warn("student denied access to lesson quizzes",
			zap.String("event", "quiz_access_denied"),
			zap.String("user_id", userID.String()),
			zap.String("lesson_id", lessonID.String()),
		)
		return nil, ErrForbidden
	}

	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	quizzes, err := s.repo.GetQuizzesByLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	subs, err := s.repo.GetStudentSubmissionsForLesson(ctx, lessonID, studentID)
	if err != nil {
		return nil, err
	}

	out := make([]StudentQuiz, 0, len(quizzes))
	for _, q := range quizzes {
		sq := StudentQuiz{Quiz: q}
		if sub, ok := subs[q.ID]; ok {
			subCopy := sub
			sq.Submission = &subCopy
		}
		out = append(out, sq)
	}
	return out, nil
}

func (s *Service) Submit(ctx context.Context, userID, quizID uuid.UUID, req SubmitQuizRequest) error {
	log := logging.FromLogger(ctx)

	q, err := s.repo.GetQuizByID(ctx, quizID)
	if err != nil {
		return err
	}

	ok, err := s.repo.UserHasLessonAccess(ctx, userID, q.LessonID)
	if err != nil {
		return err
	}
	if !ok {
		log.Warn("student denied quiz submission",
			zap.String("event", "quiz_access_denied"),
			zap.String("user_id", userID.String()),
			zap.String("quiz_id", quizID.String()),
		)
		return ErrForbidden
	}

	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	submitted, err := s.repo.HasSubmitted(ctx, quizID, studentID)
	if err != nil {
		return err
	}
	if submitted {
		return ErrAlreadySubmitted
	}

	// Index the student's answers by question id.
	byQuestion := make(map[uuid.UUID]SubmitAnswerRequest)
	for _, a := range req.Answers {
		qid, err := uuid.Parse(a.QuestionID)
		if err != nil {
			return ErrInvalidSubmission
		}
		byQuestion[qid] = a
	}

	answers := make([]Answer, 0, len(q.Questions))
	autoScore := 0
	for _, question := range q.Questions {
		given := byQuestion[question.ID]
		ans := Answer{QuestionID: question.ID}

		if question.Type == TypeTyped {
			ans.AnswerText = given.AnswerText // awaits manual grading (AwardedMarks nil)
		} else {
			selected := parseSelected(given.SelectedOptionIds, question)
			ans.SelectedOpts = selected
			awarded := 0
			if isAllCorrect(selected, question) {
				awarded = question.Marks
			}
			autoScore += awarded
			ans.AwardedMarks = &awarded
		}
		answers = append(answers, ans)
	}

	status := StatusGraded
	var gradedAt *time.Time
	if q.HasTyped() {
		status = StatusSubmitted // teacher must grade the typed parts
	} else {
		now := time.Now()
		gradedAt = &now
	}

	sub := Submission{
		QuizID:    quizID,
		StudentID: studentID,
		Status:    status,
		Score:     autoScore,
		GradedAt:  gradedAt,
		Answers:   answers,
	}
	if err := s.repo.CreateSubmission(ctx, sub); err != nil {
		return err
	}

	log.Info("quiz submitted",
		zap.String("event", "quiz_submitted"),
		zap.String("quiz_id", quizID.String()),
		zap.String("student_id", studentID.String()),
		zap.String("status", status),
		zap.Int("auto_score", autoScore),
	)
	return nil
}

func validateCreate(req CreateQuizRequest) error {
	if req.Title == "" || len(req.Questions) == 0 {
		return ErrInvalidQuiz
	}
	for _, q := range req.Questions {
		if q.Prompt == "" || q.Marks < 0 {
			return ErrInvalidQuiz
		}
		switch q.Type {
		case TypeTyped:
			// no options required
		case TypeMCQ:
			if len(q.Options) < 2 {
				return ErrInvalidQuiz
			}
			correct := 0
			for _, o := range q.Options {
				if o.Text == "" {
					return ErrInvalidQuiz
				}
				if o.IsCorrect {
					correct++
				}
			}
			if correct == 0 {
				return ErrInvalidQuiz
			}
			if !q.AllowMultiple && correct != 1 {
				return ErrInvalidQuiz
			}
		default:
			return ErrInvalidQuiz
		}
	}
	return nil
}

func parseSelected(ids []string, question Question) []uuid.UUID {
	valid := make(map[uuid.UUID]bool, len(question.Options))
	for _, o := range question.Options {
		valid[o.ID] = true
	}
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		oid, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		if valid[oid] {
			out = append(out, oid)
		}
	}
	return out
}

func isAllCorrect(selected []uuid.UUID, question Question) bool {
	correct := make(map[uuid.UUID]bool)
	for _, o := range question.Options {
		if o.IsCorrect {
			correct[o.ID] = true
		}
	}
	sel := make(map[uuid.UUID]bool)
	for _, id := range selected {
		sel[id] = true
	}
	if len(sel) != len(correct) {
		return false
	}
	for id := range correct {
		if !sel[id] {
			return false
		}
	}
	return true
}
