package dto

import "github.com/google/uuid"

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