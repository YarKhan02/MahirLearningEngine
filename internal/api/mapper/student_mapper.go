package mapper

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
)

func ToRegisterStudent(req dto.RegisterStudentRequest) (*student.Student, error) {
	dob, err := time.Parse(constant.DateLayout, req.DOB)
	if err != nil {
		return nil, fmt.Errorf("invalid dob: %w", err)
	}

	return &student.Student{
		Email:       req.Email,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		DOB:         dob,
		Gender:      req.Gender,
	}, nil
}

func ToAdminStudentResponse(req student.StudentWithBatch) dto.AdminStudentResponse {
	resp := dto.AdminStudentResponse{
		ID:          req.ID.String(),
		FullName:    req.FullName,
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

func ToStudentCourseResponse(req student.StudentCourse) dto.StudentCourseResponse {
	return dto.StudentCourseResponse{
		ID:               req.ID.String(),
		Title:            req.Title,
		Level:            req.Level,
		Duration:         req.Duration,
		Description:      req.Description,
		TotalLessons:     req.TotalLessons,
		CompletedLessons: req.CompletedLessons,
	}
}

func ToStudentLessonResponse(req student.StudentLesson) dto.StudentLessonResponse {
	resp := dto.StudentLessonResponse{
		ID:          req.ID.String(),
		Title:       req.Title,
		Description: req.Description,
		OrderNo:     req.OrderNo,
		YoutubeURL:  req.YoutubeURL,
		Content:     req.Content,
		Completed:   req.Completed,
	}

	if req.CompletedAt != nil {
		resp.CompletedAt = req.CompletedAt.Format(time.RFC3339)
	}

	return resp
}
