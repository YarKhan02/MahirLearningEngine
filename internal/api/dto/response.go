package dto

import (
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	ID          uuid.UUID           `json:"id"`
	Email       string              `json:"email"`
	Role 		string            	`json:"role"`
}

type TokenResponse struct {
	AccessToken string 	`json:"access_token"`
	ExpiresIn 	int64 	`json:"expires_in"`
}

type AuthUser struct {
	ID 			string		`json:"id"`
	Name 		string		`json:"name"`
	Username	string		`json:"username,omitempty"`
	Email 		string		`json:"email"`
	Role 		string		`json:"role"`
}

type LoginResponse struct {
	AccessToken string 		`json:"token"`
	User 		AuthUser 	`json:"user"`
}

type CourseResponse struct {
	ID 			string		`json:"id"`
	Title 		string		`json:"title"`
	Level 		string		`json:"level"`
	Duration 	int			`json:"duration"`
	Description string		`json:"description"`
	Status		string		`json:"status"`
}

type LessonResponse struct {
	ID			string		`json:"id"`
	Title 		string		`json:"title"`
	Description string		`json:"description"`
	OrderNo		int			`json:"orderNo"`
	YoutubeURL 	string		`json:"youtubeUrl"`
	Content		string		`json:"content"`
}

type BatchResponse struct {
	ID			string		`json:"id"`
	BatchName	string		`json:"batchName"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Capacity 	int 		`json:"capacity"`
	Status		string 		`json:"status"`
	Price		int			`json:"price"`
}
type BatchCourseResponse struct {
	ID			string	`json:"id"`
	CourseID	string	`json:"courseId"`
	Title		string	`json:"title"`
	Level		string	`json:"level"`
	GrantedAt	string	`json:"grantedAt"`
}

type PublicBatchResponse struct {
	ID			string					`json:"id"`
	BatchName	string					`json:"batchName"`
	StartDate	string					`json:"startDate"`
	EndDate		string					`json:"endDate"`
	Capacity	int						`json:"capacity"`
	Status		string					`json:"status"`
	Price		int						`json:"price"`
	Courses		[]BatchCourseResponse	`json:"courses"`
}

type TimetableResponse struct {
	ID			string	`json:"id"`
	BatchID		string	`json:"batchId"`
	CourseID	string	`json:"courseId"`
	CourseTitle	string	`json:"courseTitle"`
	Weekdays	[]int	`json:"weekdays"`
	StartTime	string	`json:"startTime"`
	EndTime		string	`json:"endTime"`
}

type ClassSessionResponse struct {
	Date		string	`json:"date"`
	Weekday		int		`json:"weekday"`
	StartTime	string	`json:"startTime"`
	EndTime		string	`json:"endTime"`
	CourseID	string	`json:"courseId"`
	CourseTitle	string	`json:"courseTitle"`
	BatchID		string	`json:"batchId"`
	BatchName	string	`json:"batchName"`
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
	Description	string	`json:"description"`
	OrderNo		int		`json:"orderNo"`
	YoutubeURL	string	`json:"youtubeUrl"`
	Content		string	`json:"content"`
	Completed	bool	`json:"completed"`
	CompletedAt	string	`json:"completedAt,omitempty"`
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

type SubmissionResponse struct {
	Code		string	`json:"code"`
	Status		string	`json:"status"`
	Marks		*int	`json:"marks,omitempty"`
	Remarks		*string	`json:"remarks,omitempty"`
	SubmittedAt	string	`json:"submittedAt"`
}

type StudentAssignmentResponse struct {
	AssignmentResponse
	Submission *SubmissionResponse `json:"submission,omitempty"`
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

type RosterEntryResponse struct {
	StudentID	string	`json:"studentId"`
	FullName	string	`json:"fullName"`
	Email		string	`json:"email"`
	Status		*string	`json:"status,omitempty"`
}

type AttendanceRecordResponse struct {
	Date		string	`json:"date"`
	Status		string	`json:"status"`
	BatchName	string	`json:"batchName"`
}

type AdminDashboardResponse struct {
	TotalStudents		int							`json:"totalStudents"`
	ActiveStudents		int							`json:"activeStudents"`
	PendingStudents		int							`json:"pendingStudents"`
	PendingSubmissions	int							`json:"pendingSubmissions"`
	RecentSubmissions	[]DashboardSubmission		`json:"recentSubmissions"`
	UpcomingBatches		[]DashboardUpcomingBatch	`json:"upcomingBatches"`
	RecentStudents		[]DashboardStudent			`json:"recentStudents"`
}

type DashboardSubmission struct {
	ID				string	`json:"id"`
	StudentName		string	`json:"studentName"`
	AssignmentTitle	string	`json:"assignmentTitle"`
	CourseTitle		string	`json:"courseTitle"`
	Status			string	`json:"status"`
	SubmittedAt		string	`json:"submittedAt"`
}

type DashboardUpcomingBatch struct {
	ID			string	`json:"id"`
	BatchName	string	`json:"batchName"`
	StartDate	string	`json:"startDate"`
	Price		int		`json:"price"`
	Capacity	int		`json:"capacity"`
	Enrolled	int		`json:"enrolled"`
}

type DashboardStudent struct {
	ID			string	`json:"id"`
	FullName	string	`json:"fullName"`
	Email		string	`json:"email"`
	Status		string	`json:"status"`
	BatchName	string	`json:"batchName,omitempty"`
	CreatedAt	string	`json:"createdAt"`
}
