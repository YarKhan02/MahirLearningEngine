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
	ID 		string		`json:"id"`
	Name 	string		`json:"name"`
	Email 	string		`json:"email"`
	Role 	string		`json:"role"`
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
	Days 		string 		`json:"days"`
	Status		string 		`json:"status"`
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
	Days		string					`json:"days"`
	Status		string					`json:"status"`
	Courses		[]BatchCourseResponse	`json:"courses"`
}

type AdminStudentResponse struct {
	ID			string	`json:"id"`
	FullName	string	`json:"fullName"`
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
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}
