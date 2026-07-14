package dashboard

import (
	"time"

	"github.com/google/uuid"
)

type Counts struct {
	TotalStudents      int
	ActiveStudents     int
	PendingStudents    int
	PendingSubmissions int
}

type RecentSubmission struct {
	ID              uuid.UUID
	StudentName     string
	AssignmentTitle string
	CourseTitle     string
	Status          string
	SubmittedAt     time.Time
}

type UpcomingBatch struct {
	ID        uuid.UUID
	BatchName string
	StartDate time.Time
	Price     int
	Capacity  int
	Enrolled  int
}

type RecentStudent struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Status    string
	BatchName *string
	CreatedAt time.Time
}

// AdminDashboard is everything the admin landing page shows, in one payload.
type AdminDashboard struct {
	Counts            Counts
	RecentSubmissions []RecentSubmission
	UpcomingBatches   []UpcomingBatch
	RecentStudents    []RecentStudent
}
