package quiz

import (
	"time"

	"github.com/google/uuid"
)

const (
	TypeMCQ   = "mcq"
	TypeTyped = "typed"

	StatusSubmitted = "submitted"
	StatusGraded    = "graded"
)

type Option struct {
	ID        uuid.UUID
	Text      string
	IsCorrect bool
	OrderNo   int
}

type Question struct {
	ID            uuid.UUID
	Prompt        string
	Type          string // TypeMCQ | TypeTyped
	Marks         int
	AllowMultiple bool
	OrderNo       int
	Options       []Option
}

type Quiz struct {
	ID          uuid.UUID
	LessonID    uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	Questions   []Question
}

func (q Quiz) TotalMarks() int {
	total := 0
	for _, qq := range q.Questions {
		total += qq.Marks
	}
	return total
}

func (q Quiz) HasTyped() bool {
	for _, qq := range q.Questions {
		if qq.Type == TypeTyped {
			return true
		}
	}
	return false
}

type NewQuiz struct {
	LessonID    uuid.UUID
	Title       string
	Description string
	Questions   []NewQuestion
}

type NewQuestion struct {
	Prompt        string
	Type          string
	Marks         int
	AllowMultiple bool
	Options       []NewOption
}

type NewOption struct {
	Text      string
	IsCorrect bool
}

type Answer struct {
	QuestionID   uuid.UUID
	AnswerText   string
	SelectedOpts []uuid.UUID
	AwardedMarks *int
}

type Submission struct {
	ID          uuid.UUID
	QuizID      uuid.UUID
	StudentID   uuid.UUID
	Status      string
	Score       int
	Remarks     *string
	SubmittedAt time.Time
	GradedAt    *time.Time
	Answers     []Answer
}

type StudentQuiz struct {
	Quiz
	Submission *Submission
}

type SubmissionSummary struct {
	Total   int
	Pending int // status = submitted
	Graded  int
}

type SubmissionRow struct {
	ID           uuid.UUID
	StudentID    uuid.UUID
	StudentName  string
	StudentEmail string
	Status       string
	Score        int
	SubmittedAt  time.Time
}

type NewSubmission struct {
	QuizID  uuid.UUID
	Answers []NewAnswer
}

type NewAnswer struct {
	QuestionID   uuid.UUID
	AnswerText   string
	SelectedOpts []uuid.UUID
}

type GradeInput struct {
	SubmissionID uuid.UUID
	Remarks      string
	Marks        []AnswerMark
}

type AnswerMark struct {
	QuestionID uuid.UUID
	Marks      int
}
