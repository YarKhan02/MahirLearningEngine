package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetRoster(ctx context.Context, batchID uuid.UUID, date time.Time) ([]RosterEntry, error)
	Mark(ctx context.Context, req MarkAttendance) error
	GetStudentRecords(ctx context.Context, studentID uuid.UUID) ([]Record, error)
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}
