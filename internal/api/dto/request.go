package dto

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidURL          = errors.New("invalid origin url")
	ErrHTTPSRequired       = errors.New("https required")
	ErrFragmentNotAllowed  = errors.New("fragments not allowed")
	ErrUserInfoNotAllowed  = errors.New("userinfo not allowed")
	ErrLocalhostNotAllowed = errors.New("localhost not allowed")
	ErrPrivateIPNotAllowed = errors.New("private ip not allowed")
	ErrHostNotAllowed      = errors.New("host not allowed")
)

type RegisterRequest struct {
	Email    		string     	`json:"email"`
	Password 		string     	`json:"password"`
	ConfirmPassword string 		`json:"confirm_password"`
}

func (r RegisterRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.ConfirmPassword = strings.TrimSpace(r.ConfirmPassword)
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("invalid email")
	}
	if len(r.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if r.Password != r.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}

type LoginRequest struct {
	// Identifier is an email (admins) or a username (students).
	Identifier string     `json:"identifier"`
	Password   string     `json:"password"`
}

func (r LoginRequest) Validate() error {
	r.Identifier = strings.TrimSpace(r.Identifier)
	r.Password = strings.TrimSpace(r.Password)
	if r.Identifier == "" || r.Password == "" {
		return fmt.Errorf("identifier and password are required")
	}
	return nil
}

type InsertCourse struct {
	Title 		string 		`json:"title"`
	Level 		string 		`json:"level"`
	Duration 	int 		`json:"duration"`
	Description string 		`json:"description,omitempty"`
}

type InsertLesson struct {
	ID			uuid.UUID	`json:"id"`
	CourseID	uuid.UUID	`json:"courseId"`
	Title 		string 		`json:"title"`
	Description string 		`json:"description,omitempty"`
	OrderNo		int			`json:"orderNo"`
	YoutubeURL 	string		`json:"youtubeUrl,omitempty"`
	Content		string		`json:"content,omitempty"`
}

type UpdateLesson struct {
	ID			uuid.UUID	`json:"id"`
	CourseID	uuid.UUID	`json:"courseId"`
	Title 		*string 	`json:"title,omitempty"`
	Description *string 	`json:"description,omitempty"`
	YoutubeURL 	*string		`json:"youtubeUrl,omitempty"`
	Content		*string		`json:"content,omitempty"`
}

type UpdateLessonOrderRequest struct {
    OrderNo int `json:"orderNo" binding:"required"`
}

type CreateBatchRequest struct {
	BatchName	string		`json:"batchName"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Capacity 	int 		`json:"capacity"`
	Status		string 		`json:"status"`
	Price		int			`json:"price"`
}
type UpdateBatchRequest struct {
	BatchName	string		`json:"batchName"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Capacity 	int 		`json:"capacity"`
	Status		string 		`json:"status"`
	Price		int			`json:"price"`
}

type CreateAnnouncementRequest struct {
	BatchID		string	`json:"batchId" binding:"required"`
	Title		string	`json:"title" binding:"required"`
	Description	string	`json:"description" binding:"required"`
}

type CreateTimetableRequest struct {
	CourseID	string	`json:"courseId" binding:"required"`
	Weekdays	[]int	`json:"weekdays" binding:"required"`
	StartTime	string	`json:"startTime" binding:"required"`
	EndTime		string	`json:"endTime" binding:"required"`
}

type UpdateBatchCoursesRequest struct {
	AddCourseIDs	[]string	`json:"addCourseIds"`
	RemoveCourseIDs	[]string	`json:"removeCourseIds"`
}

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

type MarkAttendanceRequest struct {
	Date		string	`json:"date" binding:"required"`
	StudentID	string	`json:"studentId" binding:"required"`
	Status		string	`json:"status" binding:"required"`
}
