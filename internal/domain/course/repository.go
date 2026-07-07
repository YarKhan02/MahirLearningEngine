package course

import (
	"context"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
)

type Repository interface {
	InsertCourse(ctx context.Context, req dto.InsertCourse) (*Course, error)
	GetCourse(ctx context.Context) ([]Course, error)
}