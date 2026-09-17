package assignment

import (
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	ID          uuid.UUID
	LessonID    uuid.UUID
	Title       string
	Description string
	StarterCode string
	Language    string // "" = plain assignment; "python" = code assignment (runner + autograde)
	DueDate     *time.Time
	TotalMarks  int
	CreatedAt   time.Time
	TestCases   []TestCase // used on create only
}

// TestCase is a hidden input/expected-output pair for autograding.
type TestCase struct {
	ID             uuid.UUID
	AssignmentID   uuid.UUID
	Stdin          string
	ExpectedStdout string
	Weight         int
	Ordinal        int
}

// TestResult is one test's outcome for a submission (shown without inputs).
type TestResult struct {
	TestCaseID   uuid.UUID
	Passed       bool
	ActualStdout string
	Stderr       string
	TimedOut     bool
	DurationMs   int
	Ordinal      int
}

type Submission struct {
	ID           uuid.UUID
	StudentID    uuid.UUID
	AssignmentID uuid.UUID
	Code         string
	Remarks      *string
	Marks        *int
	AutoScore    *int
	TestsTotal   *int
	TestsPassed  *int
	Status       string
	SubmittedAt  time.Time
	Results      []TestResult
}

type StudentAssignment struct {
	Assignment
	Submission *Submission
}

type SubmissionSummary struct {
	Total     int
	Submitted int
	Graded    int
}

type BatchSubmission struct {
	ID              uuid.UUID
	Code            string
	Remarks         *string
	Marks           *int
	AutoScore       *int
	TestsTotal      *int
	TestsPassed     *int
	Language        string
	Status          string
	SubmittedAt     time.Time
	StudentID       uuid.UUID
	StudentName     string
	StudentEmail    string
	AssignmentID    uuid.UUID
	AssignmentTitle string
	TotalMarks      int
	LessonID        uuid.UUID
	LessonTitle     string
	CourseID        uuid.UUID
	CourseTitle     string
}
