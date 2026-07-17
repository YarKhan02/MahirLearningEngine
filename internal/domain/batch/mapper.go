package batch

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	
	"github.com/google/uuid"
)

func ToCreateBatch(req CreateBatchRequest) (*Batch, error) {
	startDate, err := time.Parse(constant.DateLayout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	endDate, err := time.Parse(constant.DateLayout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate: %w", err)
	}

	return &Batch{
		BatchName: req.BatchName,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
		Status:    req.Status,
		Price:     req.Price,
	}, nil
}

func ToUpdateBatch(id uuid.UUID, req UpdateBatchRequest) (*Batch, error) {
	startDate, err := time.Parse(constant.DateLayout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	endDate, err := time.Parse(constant.DateLayout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate: %w", err)
	}

	return &Batch{
		ID:        id,
		BatchName: req.BatchName,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
		Status:    req.Status,
		Price:     req.Price,
	}, nil
}

func ToBatchResponse(req Batch) BatchResponse {
	return BatchResponse{
		ID: req.ID.String(),
		BatchName: req.BatchName,
		StartDate: req.StartDate.Format(constant.DateLayout),
		EndDate: req.EndDate.Format(constant.DateLayout),
		Capacity: req.Capacity,
		Status: req.Status,
		Price: req.Price,
	}
}
func ToBatchCourseResponse(req BatchCourse) BatchCourseResponse {
	return BatchCourseResponse{
		ID:        req.ID.String(),
		CourseID:  req.CourseID.String(),
		Title:     req.Title,
		Level:     req.Level,
		GrantedAt: req.GrantedAt.Format(time.RFC3339),
	}
}

func ToPublicBatchResponse(req BatchWithCourses) PublicBatchResponse {
	courses := make([]BatchCourseResponse, 0, len(req.Courses))
	for _, c := range req.Courses {
		courses = append(courses, ToBatchCourseResponse(c))
	}

	return PublicBatchResponse{
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
