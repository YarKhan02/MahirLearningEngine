package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/google/uuid"
)

//go:embed sql/user_create.sql
var userCreateSQL string

//go:embed sql/user_find_by_email_exists.sql
var userFindByEmailExistsSQL string

//go:embed sql/user_find_by_email.sql
var userFindByEmailSQL string

//go:embed sql/user_find_by_id.sql
var userFindByIDSQL string

//go:embed sql/user_update_failed_attempts.sql
var userUpdateFailedAttemptsSQL string

type UserRepository struct {
	db	*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = id
	err = r.db.QueryRowContext(ctx, userCreateSQL,
		u.ID,
		u.Email,
		u.PasswordHash,
		u.IsVerified,
		u.IsBanned,
		u.FailedAttempts,
		u.LockedUntil,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
	return err
}

func (r *UserRepository) FindByEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, userFindByEmailExistsSQL, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User

	err := r.db.QueryRowContext(ctx, userFindByEmailSQL, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.IsVerified,
		&u.IsBanned,
		&u.FailedAttempts,
		&u.LockedUntil,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User

	err := r.db.QueryRowContext(ctx, userFindByIDSQL, id).Scan(
		&u.ID,
		&u.Email,
		&u.IsVerified,
		&u.IsBanned,
		&u.FailedAttempts,
		&u.LockedUntil,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateFailedAttempts(ctx context.Context, id uuid.UUID, attempts int, lockedUntil *time.Time) error {
	var lock sql.NullTime
	if lockedUntil != nil {
		lock = sql.NullTime{Time: *lockedUntil, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, userUpdateFailedAttemptsSQL, id, attempts, lock)
	return err
}