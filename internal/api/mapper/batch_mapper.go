package mapper

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
)

const dateLayout = "2006-01-02"

func ToCreateBatch(req dto.CreateBatchRequest) (*batch.Batch, error) {
	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate: %w", err)
	}

	return &batch.Batch{
		BatchName: req.BatchName,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
		Days:      req.Days,
		Status:    req.Status,
	}, nil
}

func ToBatchResponse(req batch.Batch) dto.BatchResponse {
	return dto.BatchResponse{
		ID: req.ID.String(),
		BatchName: req.BatchName,
		StartDate: req.StartDate.Format(dateLayout),
		EndDate: req.EndDate.Format(dateLayout),
		Capacity: req.Capacity,
		Days: req.Days,
		Status: req.Status,
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
		StartDate: req.StartDate.Format(dateLayout),
		EndDate:   req.EndDate.Format(dateLayout),
		Capacity:  req.Capacity,
		Days:      req.Days,
		Status:    req.Status,
		Courses:   courses,
	}
}
