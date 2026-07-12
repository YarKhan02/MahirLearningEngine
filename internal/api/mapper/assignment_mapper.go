package mapper

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/google/uuid"
)

func ToCreateAssignment(req dto.CreateAssignmentRequest, lessonID uuid.UUID) (*assignment.Assignment, error) {
	a := &assignment.Assignment{
		LessonID:    lessonID,
		Title:       req.Title,
		Description: req.Description,
		StarterCode: req.StarterCode,
		TotalMarks:  req.TotalMarks,
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse(dateLayout, req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid dueDate: %w", err)
		}
		a.DueDate = &dueDate
	}

	return a, nil
}

func ToAssignmentResponse(req assignment.Assignment) dto.AssignmentResponse {
	resp := dto.AssignmentResponse{
		ID:          req.ID.String(),
		LessonID:    req.LessonID.String(),
		Title:       req.Title,
		Description: req.Description,
		StarterCode: req.StarterCode,
		TotalMarks:  req.TotalMarks,
		CreatedAt:   req.CreatedAt.Format(time.RFC3339),
	}

	if req.DueDate != nil {
		resp.DueDate = req.DueDate.Format(dateLayout)
	}

	return resp
}

func ToStudentAssignmentResponse(req assignment.StudentAssignment) dto.StudentAssignmentResponse {
	resp := dto.StudentAssignmentResponse{
		AssignmentResponse: ToAssignmentResponse(req.Assignment),
	}

	if req.Submission != nil {
		resp.Submission = &dto.SubmissionResponse{
			Code:        req.Submission.Code,
			Status:      req.Submission.Status,
			Marks:       req.Submission.Marks,
			Remarks:     req.Submission.Remarks,
			SubmittedAt: req.Submission.SubmittedAt.Format(time.RFC3339),
		}
	}

	return resp
}

func ToBatchSubmissionResponse(req assignment.BatchSubmission) dto.BatchSubmissionResponse {
	return dto.BatchSubmissionResponse{
		ID:              req.ID.String(),
		Code:            req.Code,
		Remarks:         req.Remarks,
		Marks:           req.Marks,
		Status:          req.Status,
		SubmittedAt:     req.SubmittedAt.Format(time.RFC3339),
		StudentID:       req.StudentID.String(),
		StudentName:     req.StudentName,
		StudentEmail:    req.StudentEmail,
		AssignmentID:    req.AssignmentID.String(),
		AssignmentTitle: req.AssignmentTitle,
		TotalMarks:      req.TotalMarks,
		LessonID:        req.LessonID.String(),
		LessonTitle:     req.LessonTitle,
		CourseID:        req.CourseID.String(),
		CourseTitle:     req.CourseTitle,
	}
}
