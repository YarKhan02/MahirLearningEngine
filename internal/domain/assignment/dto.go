package assignment

type CreateAssignmentRequest struct {
	Title		string				`json:"title" binding:"required"`
	Description	string				`json:"description"`
	StarterCode	string				`json:"starterCode"`
	Language	string				`json:"language"`
	DueDate		string				`json:"dueDate"`
	TotalMarks	int					`json:"totalMarks"`
	TestCases	[]TestCaseInput		`json:"testCases"`
}

type TestCaseInput struct {
	Stdin			string	`json:"stdin"`
	ExpectedStdout	string	`json:"expectedStdout"`
	Weight			int		`json:"weight"`
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
	Language	string	`json:"language,omitempty"`
	DueDate		string	`json:"dueDate,omitempty"`
	TotalMarks	int		`json:"totalMarks"`
	CreatedAt	string	`json:"createdAt"`
}

type StudentAssignmentResponse struct {
	AssignmentResponse
	Submission *SubmissionResponse `json:"submission,omitempty"`
}

type SubmissionResponse struct {
	Code		string					`json:"code"`
	Status		string					`json:"status"`
	Marks		*int					`json:"marks,omitempty"`
	Remarks		*string					`json:"remarks,omitempty"`
	AutoScore	*int					`json:"autoScore,omitempty"`
	TestsTotal	*int					`json:"testsTotal,omitempty"`
	TestsPassed	*int					`json:"testsPassed,omitempty"`
	Results		[]TestResultResponse	`json:"results,omitempty"`
	SubmittedAt	string					`json:"submittedAt"`
}

// TestResultResponse hides inputs/expected — students see only pass/fail.
type TestResultResponse struct {
	Ordinal		int		`json:"ordinal"`
	Passed		bool	`json:"passed"`
	TimedOut	bool	`json:"timedOut"`
}

type SubmissionSummaryResponse struct {
	Total     int `json:"total"`
	Submitted int `json:"submitted"`
	Graded    int `json:"graded"`
}

type BatchSubmissionResponse struct {
	ID				string	`json:"id"`
	Code			string	`json:"code"`
	Remarks			*string	`json:"remarks,omitempty"`
	Marks			*int	`json:"marks,omitempty"`
	AutoScore		*int	`json:"autoScore,omitempty"`
	TestsTotal		*int	`json:"testsTotal,omitempty"`
	TestsPassed		*int	`json:"testsPassed,omitempty"`
	Language		string	`json:"language,omitempty"`
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