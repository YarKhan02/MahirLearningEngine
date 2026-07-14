package user

import (
	"context"
	"errors"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/crypto"
	
	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account temporarily locked")
	ErrAccountBanned      = errors.New("account banned")
)

const (
	maxFailedAttempts = 5
	lockDuration      = 15 * time.Minute
)

type RoleLoader interface {
	AddRoleToUser(ctx context.Context, userID uuid.UUID, role string) error
	GetUserRole(ctx context.Context, userID uuid.UUID) (string, error)
}

type Service struct {
	userRepo Repository
	roleRepo RoleLoader
}

func NewService(userRepo Repository, roleRepo RoleLoader) *Service {
	return &Service{
		userRepo: userRepo, 
		roleRepo: roleRepo,
	} 
}

func (s *Service) RegisterAdmin(ctx context.Context, req dto.RegisterRequest) (*User, error) {
	// Check if the email is already taken
	existingUser, err := s.userRepo.FindByEmailExists(ctx, req.Email)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if existingUser {
		return nil, ErrEmailTaken
	}

	// Hash the password
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create a new user entity
	user := &User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		IsVerified:   true,
		IsBanned:     false,
		Role:   	  "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the user to the repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Assign admin role to the user
	err = s.roleRepo.AddRoleToUser(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, identifier string, password string) (*User, error) {
	u, err := s.userRepo.FindByLoginIdentifier(ctx, identifier)
	if err != nil || u == nil {
		// constant-time even on not found to prevent timing attacks
		crypto.HashPassword(password) //nolint:errcheck
		return nil, ErrInvalidCredentials
	}

	if !u.IsVerified {
		return nil, ErrInvalidCredentials
	}

	if u.IsBanned {
		return nil, ErrAccountBanned
	}

	if u.IsLocked() {
		return nil, ErrAccountLocked
	}

	if !crypto.VerifyPassword(password, u.PasswordHash) {
		attempts := u.FailedAttempts + 1
		var lockedUntil *time.Time
		if attempts >= maxFailedAttempts {
			t := time.Now().Add(lockDuration)
			lockedUntil = &t
		}
		s.userRepo.UpdateFailedAttempts(ctx, u.ID, attempts, lockedUntil) //nolint:errcheck
		return nil, ErrInvalidCredentials
	}

	role, err := s.roleRepo.GetUserRole(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	u.Role = role

	// Reset failed attempts on success
	s.userRepo.UpdateFailedAttempts(ctx, u.ID, 0, nil) //nolint:errcheck

	return u, nil
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *Service) RegisterStudentAccount(ctx context.Context, username, password string) (*User, error) {
	usernameTaken, err := s.userRepo.FindByUsernameExists(ctx, username)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if usernameTaken {
		return nil, ErrUsernameTaken
	}

	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     username,
		PasswordHash: passwordHash,
		IsVerified:   true,
		IsBanned:     false,
		Role:         "student",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	if err := s.roleRepo.AddRoleToUser(ctx, user.ID, user.Role); err != nil {
		return nil, err
	}

	return user, nil
}
