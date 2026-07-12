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
