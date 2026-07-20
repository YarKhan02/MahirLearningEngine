package attendance

import (
	"time"

	"github.com/google/uuid"
)

type RosterEntry struct {
	StudentID uuid.UUID
	FullName  string
	Email     string
	Status    *string
}

type Record struct {
	LessonDate time.Time
	Status     string
	BatchName  string
}

type Summary struct {
	Present     int
	Absent      int
	Total       int
	TodayStatus *string
}

type MarkAttendance struct {
	BatchID 	uuid.UUID
	Date 		time.Time
	StudentID 	uuid.UUID
	Status 		string
	CreatedBy 	uuid.UUID
}