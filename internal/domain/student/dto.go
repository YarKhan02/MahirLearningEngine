package student

type RegisterStudentRequest struct {
	FullName	string	`json:"fullName" binding:"required"`
	Username	string	`json:"username" binding:"required"`
	Email		string	`json:"email" binding:"required,email"`
	PhoneNumber	string	`json:"phoneNumber" binding:"required"`
	DOB			string	`json:"dob" binding:"required"`
	Gender		string	`json:"gender" binding:"required"`
	BatchID		string	`json:"batchId" binding:"required"`
}

type UpdateStudentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateStudentBatchRequest struct {
	// Empty string removes the student from their current batch.
	BatchID string `json:"batchId"`
}

type SetLessonProgressRequest struct {
	Completed *bool `json:"completed" binding:"required"`
}

type AdminStudentResponse struct {
	ID			string	`json:"id"`
	FullName	string	`json:"fullName"`
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	PhoneNumber	string	`json:"phoneNumber"`
	DOB			string	`json:"dob"`
	Gender		string	`json:"gender"`
	Status		string	`json:"status"`
	BatchID		string	`json:"batchId,omitempty"`
	BatchName	string	`json:"batchName,omitempty"`
	HasAccount	bool	`json:"hasAccount"`
}

type StudentAccountResponse struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

type StudentCourseResponse struct {
	ID					string	`json:"id"`
	Title				string	`json:"title"`
	Level				string	`json:"level"`
	Duration			int		`json:"duration"`
	Description			string	`json:"description"`
	TotalLessons		int		`json:"totalLessons"`
	CompletedLessons	int		`json:"completedLessons"`
}

type StudentLessonResponse struct {
	ID			string	`json:"id"`
	Title		string	`json:"title"`
	OrderNo		int		`json:"orderNo"`
	Completed	bool	`json:"completed"`
	CompletedAt	string	`json:"completedAt,omitempty"`
}
