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
	DueDate     *time.Time
	TotalMarks  int
	CreatedAt   time.Time
}

type Submission struct {
	ID           uuid.UUID
	StudentID    uuid.UUID
	AssignmentID uuid.UUID
	Code         string
	Remarks      *string
	Marks        *int
	Status       string
	SubmittedAt  time.Time
}

// StudentAssignment is an assignment paired with the student's submission (if any).
type StudentAssignment struct {
	Assignment
	Submission *Submission
}

// BatchSubmission is a submission enriched with student, assignment,
// lesson, and course context — the admin review view.
type BatchSubmission struct {
	ID              uuid.UUID
	Code            string
	Remarks         *string
	Marks           *int
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
