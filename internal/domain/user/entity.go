package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID				uuid.UUID
	Email 			string
	PasswordHash 	string
	IsVerified 		bool
	IsBanned 		bool
	FailedAttempts 	int
	LockedUntil 	*time.Time
	Role 			string
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
}

func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}

func (u *User) HasGlobalRole(role string) bool {
    return u.Role == role
}