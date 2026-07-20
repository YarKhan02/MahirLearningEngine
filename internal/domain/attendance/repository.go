package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetRoster(ctx context.Context, batchID uuid.UUID, date time.Time) ([]RosterEntry, error)
	Mark(ctx context.Context, req MarkAttendance) error
	GetStudentRecords(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]Record, error)
	CountStudentRecords(ctx context.Context, studentID uuid.UUID) (int, error)
	GetStudentSummary(ctx context.Context, studentID uuid.UUID) (Summary, error)
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}
