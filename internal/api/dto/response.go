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