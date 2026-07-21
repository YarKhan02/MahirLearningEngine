package quiz

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func ToNewQuiz(req CreateQuizRequest, lessonID uuid.UUID) NewQuiz {
	questions := make([]NewQuestion, 0, len(req.Questions))
	for _, q := range req.Questions {
		opts := make([]NewOption, 0, len(q.Options))
		for _, o := range q.Options {
			opts = append(opts, NewOption{Text: o.Text, IsCorrect: o.IsCorrect})
		}
		questions = append(questions, NewQuestion{
			Prompt:        q.Prompt,
			Type:          q.Type,
			Marks:         q.Marks,
			AllowMultiple: q.AllowMultiple,
			Options:       opts,
		})
	}
	return NewQuiz{
		LessonID:    lessonID,
		Title:       req.Title,
		Description: req.Description,
		Questions:   questions,
	}
}

func ToNewSubmission(req SubmitQuizRequest, quizID uuid.UUID) (NewSubmission, error) {
	answers := make([]NewAnswer, 0, len(req.Answers))
	for _, a := range req.Answers {
		qid, err := uuid.Parse(a.QuestionID)
		if err != nil {
			return NewSubmission{}, fmt.Errorf("invalid questionId: %w", err)
		}
		opts := make([]uuid.UUID, 0, len(a.SelectedOptionIds))
		for _, o := range a.SelectedOptionIds {
			oid, err := uuid.Parse(o)
			if err != nil {
				return NewSubmission{}, fmt.Errorf("invalid optionId: %w", err)
			}
			opts = append(opts, oid)
		}
		answers = append(answers, NewAnswer{
			QuestionID:   qid,
			AnswerText:   a.AnswerText,
			SelectedOpts: opts,
		})
	}
	return NewSubmission{QuizID: quizID, Answers: answers}, nil
}

func ToGradeInput(req GradeQuizRequest, submissionID uuid.UUID) (GradeInput, error) {
	marks := make([]AnswerMark, 0, len(req.Marks))
	for _, m := range req.Marks {
		qid, err := uuid.Parse(m.QuestionID)
		if err != nil {
			return GradeInput{}, fmt.Errorf("invalid questionId: %w", err)
		}
		marks = append(marks, AnswerMark{QuestionID: qid, Marks: m.Marks})
	}
	return GradeInput{SubmissionID: submissionID, Remarks: req.Remarks, Marks: marks}, nil
}

func toOptionResponses(opts []Option, includeCorrect bool) []OptionResponse {
	out := make([]OptionResponse, 0, len(opts))
	for _, o := range opts {
		r := OptionResponse{ID: o.ID.String(), Text: o.Text, OrderNo: o.OrderNo}
		if includeCorrect {
			v := o.IsCorrect
			r.IsCorrect = &v
		}
		out = append(out, r)
	}
	return out
}

func toQuestionResponses(questions []Question, includeCorrect bool) []QuestionResponse {
	out := make([]QuestionResponse, 0, len(questions))
	for _, q := range questions {
		out = append(out, QuestionResponse{
			ID:            q.ID.String(),
			Prompt:        q.Prompt,
			Type:          q.Type,
			Marks:         q.Marks,
			AllowMultiple: q.AllowMultiple,
			OrderNo:       q.OrderNo,
			Options:       toOptionResponses(q.Options, includeCorrect),
		})
	}
	return out
}

func ToQuizResponse(q Quiz, includeCorrect bool) QuizResponse {
	return QuizResponse{
		ID:          q.ID.String(),
		Title:       q.Title,
		Description: q.Description,
		TotalMarks:  q.TotalMarks(),
		CreatedAt:   q.CreatedAt.Format(time.RFC3339),
		Questions:   toQuestionResponses(q.Questions, includeCorrect),
	}
}

func ToAdminQuizResponse(q Quiz, submissionCount, pendingCount int) AdminQuizResponse {
	return AdminQuizResponse{
		QuizResponse:    ToQuizResponse(q, true),
		SubmissionCount: submissionCount,
		PendingCount:    pendingCount,
	}
}

func ToStudentQuizResponse(sq StudentQuiz) StudentQuizResponse {
	resp := StudentQuizResponse{QuizResponse: ToQuizResponse(sq.Quiz, false)}
	if sq.Submission != nil {
		resp.Submission = toStudentSubmissionResponse(*sq.Submission, sq.Quiz.TotalMarks())
	}
	return resp
}

func toStudentSubmissionResponse(s Submission, totalMarks int) *StudentSubmissionResponse {
	answers := make([]StudentAnswerResponse, 0, len(s.Answers))
	for _, a := range s.Answers {
		opts := make([]string, 0, len(a.SelectedOpts))
		for _, o := range a.SelectedOpts {
			opts = append(opts, o.String())
		}
		answers = append(answers, StudentAnswerResponse{
			QuestionID:        a.QuestionID.String(),
			AnswerText:        a.AnswerText,
			SelectedOptionIds: opts,
			AwardedMarks:      a.AwardedMarks,
		})
	}
	return &StudentSubmissionResponse{
		ID:          s.ID.String(),
		Status:      s.Status,
		Score:       s.Score,
		TotalMarks:  totalMarks,
		Remarks:     s.Remarks,
		SubmittedAt: s.SubmittedAt.Format(time.RFC3339),
		Answers:     answers,
	}
}

func ToSubmissionRowResponse(r SubmissionRow) SubmissionRowResponse {
	return SubmissionRowResponse{
		ID:           r.ID.String(),
		StudentID:    r.StudentID.String(),
		StudentName:  r.StudentName,
		StudentEmail: r.StudentEmail,
		Status:       r.Status,
		Score:        r.Score,
		SubmittedAt:  r.SubmittedAt.Format(time.RFC3339),
	}
}

func ToSubmissionSummaryResponse(s SubmissionSummary) SubmissionSummaryResponse {
	return SubmissionSummaryResponse(s)
}

func ToSubmissionDetailResponse(q Quiz, sub Submission, studentName string) SubmissionDetailResponse {
	byQ := make(map[uuid.UUID]Answer, len(sub.Answers))
	for _, a := range sub.Answers {
		byQ[a.QuestionID] = a
	}

	questions := make([]GradedQuestionResponse, 0, len(q.Questions))
	for _, qq := range q.Questions {
		gq := GradedQuestionResponse{
			QuestionID: qq.ID.String(),
			Prompt:     qq.Prompt,
			Type:       qq.Type,
			Marks:      qq.Marks,
			Options:    toOptionResponses(qq.Options, true),
		}
		if a, ok := byQ[qq.ID]; ok {
			gq.AnswerText = a.AnswerText
			gq.AwardedMarks = a.AwardedMarks
			for _, o := range a.SelectedOpts {
				gq.SelectedOptionIds = append(gq.SelectedOptionIds, o.String())
			}
		}
		questions = append(questions, gq)
	}

	return SubmissionDetailResponse{
		ID:          sub.ID.String(),
		QuizTitle:   q.Title,
		StudentName: studentName,
		Status:      sub.Status,
		Score:       sub.Score,
		TotalMarks:  q.TotalMarks(),
		Remarks:     sub.Remarks,
		SubmittedAt: sub.SubmittedAt.Format(time.RFC3339),
		Questions:   questions,
	}
}
