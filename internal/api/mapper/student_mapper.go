package mapper

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
)

func ToRegisterStudent(req dto.RegisterStudentRequest) (*student.Student, error) {
	dob, err := time.Parse(dateLayout, req.DOB)
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
		DOB:         req.DOB.Format(dateLayout),
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
