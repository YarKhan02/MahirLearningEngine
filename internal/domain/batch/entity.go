package batch

import (
	"time"

	"github.com/google/uuid"
)

type Batch struct {
	ID 			uuid.UUID
	BatchName	string
	StartDate	time.Time
	EndDate		time.Time
	Capacity 	int
	Days 		string
	Status		string
	Price		int
}
type BatchCourse struct {
	ID			uuid.UUID
	CourseID	uuid.UUID
	Title		string
	Level		string
	GrantedAt	time.Time
}

type BatchWithCourses struct {
	Batch
	Courses []BatchCourse
}
