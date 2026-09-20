package assignment

import (
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	
	"github.com/google/uuid"
)

func ToCreateAssignment(req CreateAssignmentRequest, lessonID uuid.UUID) (*Assignment, error) {
	a := &Assignment{
		LessonID:    lessonID,
		Title:       req.Title,
		Description: req.Description,
		StarterCode: req.StarterCode,
		Language:    req.Language,
		TotalMarks:  req.TotalMarks,
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse(constant.DateLayout, req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid dueDate: %w", err)
		}
		a.DueDate = &dueDate
	}

	for i, tc := range req.TestCases {
		w := tc.Weight
		if w <= 0 {
			w = 1
		}
		a.TestCases = append(a.TestCases, TestCase{
			Stdin:          tc.Stdin,
			ExpectedStdout: tc.ExpectedStdout,
			Weight:         w,
			Ordinal:        i,
		})
	}

	return a, nil
}

func ToUpdateAssignment(req UpdateAssignmentRequest, assignmentID uuid.UUID) (*Assignment, error) {
	a := &Assignment{
		ID:          assignmentID,
		Title:       req.Title,
		Description: req.Description,
		StarterCode: req.StarterCode,
		Language:    req.Language,
		TotalMarks:  req.TotalMarks,
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse(constant.DateLayout, req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid dueDate: %w", err)
		}
		a.DueDate = &dueDate
	}

	for i, tc := range req.TestCases {
		w := tc.Weight
		if w <= 0 {
			w = 1
		}
		a.TestCases = append(a.TestCases, TestCase{
			Stdin:          tc.Stdin,
			ExpectedStdout: tc.ExpectedStdout,
			Weight:         w,
			Ordinal:        i,
		})
	}

	return a, nil
}

func ToAssignmentWithTestsResponse(a Assignment) AssignmentWithTestsResponse {
	resp := AssignmentWithTestsResponse{
		AssignmentResponse: ToAssignmentResponse(a),
		TestCases:          make([]TestCaseResponse, 0, len(a.TestCases)),
	}
	for _, tc := range a.TestCases {
		resp.TestCases = append(resp.TestCases, TestCaseResponse{
			Stdin:          tc.Stdin,
			ExpectedStdout: tc.ExpectedStdout,
			Weight:         tc.Weight,
			Ordinal:        tc.Ordinal,
		})
	}
	return resp
}

func ToAssignmentResponse(req Assignment) AssignmentResponse {
	resp := AssignmentResponse{
		ID:          req.ID.String(),
		LessonID:    req.LessonID.String(),
		Title:       req.Title,
		Description: req.Description,
		StarterCode: req.StarterCode,
		Language:    req.Language,
		TotalMarks:  req.TotalMarks,
		CreatedAt:   req.CreatedAt.Format(time.RFC3339),
	}

	if req.DueDate != nil {
		resp.DueDate = req.DueDate.Format(constant.DateLayout)
	}

	return resp
}

func ToStudentAssignmentResponse(req StudentAssignment) StudentAssignmentResponse {
	resp := StudentAssignmentResponse{
		AssignmentResponse: ToAssignmentResponse(req.Assignment),
	}

	if req.Submission != nil {
		resp.Submission = &SubmissionResponse{
			Code:        req.Submission.Code,
			Status:      req.Submission.Status,
			Marks:       req.Submission.Marks,
			Remarks:     req.Submission.Remarks,
			AutoScore:   req.Submission.AutoScore,
			TestsTotal:  req.Submission.TestsTotal,
			TestsPassed: req.Submission.TestsPassed,
			Results:     toTestResultResponses(req.Submission.Results),
			SubmittedAt: req.Submission.SubmittedAt.Format(time.RFC3339),
		}
	}

	return resp
}

func toTestResultResponses(results []TestResult) []TestResultResponse {
	if len(results) == 0 {
		return nil
	}
	out := make([]TestResultResponse, 0, len(results))
	for _, r := range results {
		out = append(out, TestResultResponse{
			Ordinal:  r.Ordinal,
			Passed:   r.Passed,
			TimedOut: r.TimedOut,
		})
	}
	return out
}

func ToSubmissionSummaryResponse(s SubmissionSummary) SubmissionSummaryResponse {
	return SubmissionSummaryResponse(s)
}

func ToBatchSubmissionResponse(req BatchSubmission) BatchSubmissionResponse {
	return BatchSubmissionResponse{
		ID:              req.ID.String(),
		Code:            req.Code,
		Remarks:         req.Remarks,
		Marks:           req.Marks,
		AutoScore:       req.AutoScore,
		TestsTotal:      req.TestsTotal,
		TestsPassed:     req.TestsPassed,
		Language:        req.Language,
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
