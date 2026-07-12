package student

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID			uuid.UUID
	Email		string
	FullName	string
	PhoneNumber	string
	DOB			time.Time
	Gender		string
	Status		string
}

// StudentWithBatch is an admin list row — student joined with their batch and account existence.
type StudentWithBatch struct {
	Student
	BatchID    *uuid.UUID
	BatchName  *string
	HasAccount bool
}

// StudentCourse is a course the student can access, with their progress.
type StudentCourse struct {
	ID               uuid.UUID
	Title            string
	Level            string
	Duration         int
	Description      string
	TotalLessons     int
	CompletedLessons int
}

// StudentLesson is a lesson with the student's completion state.
type StudentLesson struct {
	ID          uuid.UUID
	Title       string
	Description string
	OrderNo     int
	YoutubeURL  string
	Content     string
	Completed   bool
	CompletedAt *time.Time
}
