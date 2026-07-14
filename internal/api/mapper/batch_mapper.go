package mapper

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/google/uuid"
)

func ToCreateBatch(req dto.CreateBatchRequest) (*batch.Batch, error) {
	startDate, err := time.Parse(constant.DateLayout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	endDate, err := time.Parse(constant.DateLayout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate: %w", err)
	}

	return &batch.Batch{
		BatchName: req.BatchName,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
		Status:    req.Status,
		Price:     req.Price,
	}, nil
}

func ToUpdateBatch(id uuid.UUID, req dto.UpdateBatchRequest) (*batch.Batch, error) {
	startDate, err := time.Parse(constant.DateLayout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	endDate, err := time.Parse(constant.DateLayout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate: %w", err)
	}

	return &batch.Batch{
		ID:        id,
		BatchName: req.BatchName,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
		Status:    req.Status,
		Price:     req.Price,
	}, nil
}

func ToBatchResponse(req batch.Batch) dto.BatchResponse {
	return dto.BatchResponse{
		ID: req.ID.String(),
		BatchName: req.BatchName,
		StartDate: req.StartDate.Format(constant.DateLayout),
		EndDate: req.EndDate.Format(constant.DateLayout),
		Capacity: req.Capacity,
		Status: req.Status,
		Price: req.Price,
	}
}
func ToBatchCourseResponse(req batch.BatchCourse) dto.BatchCourseResponse {
	return dto.BatchCourseResponse{
		ID:        req.ID.String(),
		CourseID:  req.CourseID.String(),
		Title:     req.Title,
		Level:     req.Level,
		GrantedAt: req.GrantedAt.Format(time.RFC3339),
	}
}

func ToPublicBatchResponse(req batch.BatchWithCourses) dto.PublicBatchResponse {
	courses := make([]dto.BatchCourseResponse, 0, len(req.Courses))
	for _, c := range req.Courses {
		courses = append(courses, ToBatchCourseResponse(c))
	}

	return dto.PublicBatchResponse{
		ID:        req.ID.String(),
		BatchName: req.BatchName,
		StartDate: req.StartDate.Format(constant.DateLayout),
		EndDate:   req.EndDate.Format(constant.DateLayout),
		Capacity:  req.Capacity,
		Status:    req.Status,
		Price:     req.Price,
		Courses:   courses,
	}
}
