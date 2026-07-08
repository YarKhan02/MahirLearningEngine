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
	Email    string     `json:"email"`
	Password string     `json:"password"`
}

func (r LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	if r.Email == "" || r.Password == "" {
		return fmt.Errorf("email and password are required")
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