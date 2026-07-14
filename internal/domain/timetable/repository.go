package timetable

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *Timetable) error
	GetByBatch(ctx context.Context, batchID uuid.UUID) ([]Timetable, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// GetRulesForUser returns the schedule rules for every batch the user (a
	// student) is enrolled in, with batch dates and course titles attached.
	GetRulesForUser(ctx context.Context, userID uuid.UUID) ([]Timetable, error)
}
