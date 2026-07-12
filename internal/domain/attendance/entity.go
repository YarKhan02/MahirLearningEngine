package attendance

import (
	"time"

	"github.com/google/uuid"
)

// RosterEntry is a student in a batch with their status for one date
// (nil status = not marked yet).
type RosterEntry struct {
	StudentID uuid.UUID
	FullName  string
	Email     string
	Status    *string
}

// Record is one marked attendance day for a student.
type Record struct {
	LessonDate time.Time
	Status     string
	BatchName  string
}
