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
