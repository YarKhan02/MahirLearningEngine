package student

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
)

func ToRegisterStudent(req RegisterStudentRequest) (*Student, error) {
	dob, err := time.Parse(constant.DateLayout, req.DOB)
	if err != nil {
		return nil, fmt.Errorf("invalid dob: %w", err)
	}

	return &Student{
		Email:       req.Email,
		Username:    req.Username,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		DOB:         dob,
		Gender:      req.Gender,
	}, nil
}

func ToAdminStudentResponse(req StudentWithBatch) AdminStudentResponse {
	resp := AdminStudentResponse{
		ID:          req.ID.String(),
		FullName:    req.FullName,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		DOB:         req.DOB.Format(constant.DateLayout),
		Gender:      req.Gender,
		Status:      req.Status,
		HasAccount:  req.HasAccount,
	}

	if req.BatchID != nil {
		resp.BatchID = req.BatchID.String()
	}
	if req.BatchName != nil {
		resp.BatchName = *req.BatchName
	}

	return resp
}

func ToStudentCourseResponse(req StudentCourse) StudentCourseResponse {
	return StudentCourseResponse{
		ID:               req.ID.String(),
		Title:            req.Title,
		Level:            req.Level,
		Duration:         req.Duration,
		Description:      req.Description,
		TotalLessons:     req.TotalLessons,
		CompletedLessons: req.CompletedLessons,
	}
}

func ToStudentLessonResponse(req StudentLesson) StudentLessonResponse {
	resp := StudentLessonResponse{
		ID:        req.ID.String(),
		Title:     req.Title,
		OrderNo:   req.OrderNo,
		Completed: req.Completed,
	}

	if req.CompletedAt != nil {
		resp.CompletedAt = req.CompletedAt.Format(time.RFC3339)
	}

	return resp
}
