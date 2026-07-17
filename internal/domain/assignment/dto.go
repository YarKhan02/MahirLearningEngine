package assignment

type CreateAssignmentRequest struct {
	Title		string	`json:"title" binding:"required"`
	Description	string	`json:"description"`
	StarterCode	string	`json:"starterCode"`
	DueDate		string	`json:"dueDate"`
	TotalMarks	int		`json:"totalMarks"`
}

type SubmitAssignmentRequest struct {
	Code string `json:"code" binding:"required"`
}

type GradeSubmissionRequest struct {
	Marks	*int	`json:"marks" binding:"required"`
	Remarks	string	`json:"remarks"`
}

type AssignmentResponse struct {
	ID			string	`json:"id"`
	LessonID	string	`json:"lessonId"`
	Title		string	`json:"title"`
	Description	string	`json:"description"`
	StarterCode	string	`json:"starterCode"`
	DueDate		string	`json:"dueDate,omitempty"`
	TotalMarks	int		`json:"totalMarks"`
	CreatedAt	string	`json:"createdAt"`
}

type StudentAssignmentResponse struct {
	AssignmentResponse
	Submission *SubmissionResponse `json:"submission,omitempty"`
}

type SubmissionResponse struct {
	Code		string	`json:"code"`
	Status		string	`json:"status"`
	Marks		*int	`json:"marks,omitempty"`
	Remarks		*string	`json:"remarks,omitempty"`
	SubmittedAt	string	`json:"submittedAt"`
}

type BatchSubmissionResponse struct {
	ID				string	`json:"id"`
	Code			string	`json:"code"`
	Remarks			*string	`json:"remarks,omitempty"`
	Marks			*int	`json:"marks,omitempty"`
	Status			string	`json:"status"`
	SubmittedAt		string	`json:"submittedAt"`
	StudentID		string	`json:"studentId"`
	StudentName		string	`json:"studentName"`
	StudentEmail	string	`json:"studentEmail"`
	AssignmentID	string	`json:"assignmentId"`
	AssignmentTitle	string	`json:"assignmentTitle"`
	TotalMarks		int		`json:"totalMarks"`
	LessonID		string	`json:"lessonId"`
	LessonTitle		string	`json:"lessonTitle"`
	CourseID		string	`json:"courseId"`
	CourseTitle		string	`json:"courseTitle"`
}