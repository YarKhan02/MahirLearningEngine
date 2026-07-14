package timetable

import (
	"time"

	"github.com/google/uuid"
)

// Timetable is a recurring schedule rule: a batch takes a course on the given
// weekdays at a fixed time. The concrete class dates are generated between the
// batch's start and end dates rather than stored.
type Timetable struct {
	ID        uuid.UUID
	BatchID   uuid.UUID
	CourseID  uuid.UUID
	Weekdays  []int // ISO weekdays: 1=Mon … 7=Sun
	StartTime string // "HH:MM"
	EndTime   string // "HH:MM"

	// Populated on reads for display.
	CourseTitle string
	BatchName   string
	BatchStart  time.Time
	BatchEnd    time.Time
}

// ClassSession is a single generated class occurrence — never persisted.
type ClassSession struct {
	Date        time.Time
	Weekday     int // ISO
	StartTime   string
	EndTime     string
	CourseID    uuid.UUID
	CourseTitle string
	BatchID     uuid.UUID
	BatchName   string
}
