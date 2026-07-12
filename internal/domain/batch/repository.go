package batch

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateBatch(ctx context.Context, req *Batch) error 
	GetBatches(ctx context.Context) ([]Batch, error)
	GetBatchCourses(ctx context.Context, batchID uuid.UUID) ([]BatchCourse, error)
	UpdateBatchCourses(ctx context.Context, batchID uuid.UUID, add []uuid.UUID, remove []uuid.UUID, grantedBy *uuid.UUID) error
}
