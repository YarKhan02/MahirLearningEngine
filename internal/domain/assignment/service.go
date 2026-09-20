package assignment

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/codeexec"
	"github.com/google/uuid"
)

var (
	ErrStudentNotFound    = errors.New("student not found")
	ErrAccessDenied       = errors.New("you do not have access to this assignment")
	ErrAssignmentNotFound = errors.New("assignment not found")
)

// CodeGrader runs a submission's code against test inputs (implemented by codeexec.Runner).
type CodeGrader interface {
	Grade(ctx context.Context, source string, inputs []string) ([]codeexec.RunResult, error)
	Configured() bool
}

type Service struct {
	repo   Repository
	grader CodeGrader
}

func NewService(repo Repository, grader CodeGrader) *Service {
	return &Service{repo: repo, grader: grader}
}

func (s *Service) CreateAssignment(ctx context.Context, a *Assignment) error {
	return s.repo.CreateAssignment(ctx, a)
}

func (s *Service) GetLessonAssignments(ctx context.Context, lessonID uuid.UUID) ([]Assignment, error) {
	return s.repo.GetLessonAssignments(ctx, lessonID)
}

func (s *Service) GetAssignment(ctx context.Context, id uuid.UUID) (*Assignment, error) {
	return s.repo.GetAssignmentWithTests(ctx, id)
}

func (s *Service) UpdateAssignment(ctx context.Context, a *Assignment) error {
	return s.repo.UpdateAssignment(ctx, a)
}

func (s *Service) DeleteAssignment(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAssignment(ctx, id)
}

func (s *Service) GetStudentAssignments(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) ([]StudentAssignment, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hasAccess, err := s.repo.HasLessonAccess(ctx, studentID, lessonID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrAccessDenied
	}

	items, err := s.repo.GetStudentAssignments(ctx, lessonID, studentID)
	if err != nil {
		return nil, err
	}
	// Attach per-test results for graded code submissions.
	for i := range items {
		sub := items[i].Submission
		if sub != nil && sub.ID != uuid.Nil {
			if results, err := s.repo.GetSubmissionResults(ctx, sub.ID); err == nil {
				sub.Results = results
			}
		}
	}
	return items, nil
}

func (s *Service) SubmitAssignment(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID, code string) error {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	hasAccess, err := s.repo.HasAssignmentAccess(ctx, studentID, assignmentID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrAccessDenied
	}

	submissionID, err := s.repo.SubmitAssignment(ctx, studentID, assignmentID, code)
	if err != nil {
		return err
	}
	// Autograde synchronously so the student sees the score on the next fetch.
	s.autograde(ctx, submissionID, assignmentID, code)
	return nil
}

// autograde runs the assignment's hidden tests against the submitted code and
// stores the score + per-test results. Best-effort: any failure leaves the
// submission ungraded (teacher can still grade manually).
func (s *Service) autograde(ctx context.Context, submissionID, assignmentID uuid.UUID, code string) {
	if s.grader == nil || !s.grader.Configured() {
		return
	}
	language, totalMarks, tests, err := s.repo.GetGradingData(ctx, assignmentID)
	if err != nil || language != "python" || len(tests) == 0 {
		return
	}

	inputs := make([]string, len(tests))
	for i, t := range tests {
		inputs[i] = t.Stdin
	}

	gctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	runs, err := s.grader.Grade(gctx, code, inputs)
	if err != nil {
		return
	}

	results := make([]TestResult, 0, len(tests))
	passedCount, passedWeight, totalWeight := 0, 0, 0
	for i, t := range tests {
		totalWeight += t.Weight
		var run codeexec.RunResult
		if i < len(runs) {
			run = runs[i]
		}
		ok := !run.TimedOut && normalizeOutput(run.Stdout) == normalizeOutput(t.ExpectedStdout)
		if ok {
			passedCount++
			passedWeight += t.Weight
		}
		results = append(results, TestResult{
			TestCaseID:   t.ID,
			Passed:       ok,
			ActualStdout: run.Stdout,
			Stderr:       run.Stderr,
			TimedOut:     run.TimedOut,
			DurationMs:   run.DurationMs,
			Ordinal:      t.Ordinal,
		})
	}

	autoScore := 0
	if totalWeight > 0 {
		autoScore = int(math.Round(float64(passedWeight) / float64(totalWeight) * float64(totalMarks)))
	}
	_ = s.repo.SaveAutoGrade(ctx, submissionID, autoScore, len(tests), passedCount, results)
}

// normalizeOutput trims trailing whitespace per line and trailing blank lines,
// so cosmetic differences don't fail an otherwise-correct answer.
func normalizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

func (s *Service) GetBatchSubmissions(ctx context.Context, batchID uuid.UUID, q, status string, limit, offset int) ([]BatchSubmission, int, error) {
	total, err := s.repo.CountBatchSubmissions(ctx, batchID, q, status)
	if err != nil {
		return nil, 0, err
	}
	items, err := s.repo.GetBatchSubmissions(ctx, batchID, q, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) GetBatchSubmissionSummary(ctx context.Context, batchID uuid.UUID, q string) (SubmissionSummary, error) {
	return s.repo.GetBatchSubmissionSummary(ctx, batchID, q)
}

func (s *Service) GradeSubmission(ctx context.Context, submissionID uuid.UUID, marks int, remarks string) error {
	return s.repo.GradeSubmission(ctx, submissionID, marks, remarks)
}

func (s *Service) GetMySubmissions(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]BatchSubmission, int, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountStudentSubmissions(ctx, studentID, status)
	if err != nil {
		return nil, 0, err
	}
	items, err := s.repo.GetStudentSubmissions(ctx, studentID, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) GetMySubmissionSummary(ctx context.Context, userID uuid.UUID) (SubmissionSummary, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return SubmissionSummary{}, err
	}
	return s.repo.GetStudentSubmissionSummary(ctx, studentID)
}
