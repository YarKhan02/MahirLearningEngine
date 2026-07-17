package common

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrUsernameTaken = errors.New("username already taken")

type StudentProfile struct {
	FullName string
	Email    string
	Username string
}

type StudentProfileProvider interface {
	GetStudentProfile(ctx context.Context, userID uuid.UUID) (*StudentProfile, error)
}

type StudentAccountRegistrar interface {
	RegisterStudentAccount(ctx context.Context, username, password string) error
}
