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

//go:embed sql/user_find_by_id_exists.sql
var userFindByIDExistsSQL string

//go:embed sql/user_find_by_email.sql
var userFindByEmailSQL string

//go:embed sql/user_find_by_login.sql
var userFindByLoginSQL string

//go:embed sql/user_find_by_username_exists.sql
var userFindByUsernameExistsSQL string

//go:embed sql/user_find_by_id.sql
var userFindByIDSQL string

//go:embed sql/user_update_failed_attempts.sql
var userUpdateFailedAttemptsSQL string

//go:embed sql/user_reset_password.sql
var userResetPasswordSQL string

type UserRepository struct {
	db *sql.DB
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

	// Admins log in by email (no username); students log in by username and
	// store NO email (siblings may share one). Empty values are stored as NULL
	// so the UNIQUE indexes ignore them.
	email := sql.NullString{String: u.Email, Valid: u.Email != ""}
	username := sql.NullString{String: u.Username, Valid: u.Username != ""}

	err = r.db.QueryRowContext(ctx, userCreateSQL,
		u.ID,
		email,
		username,
		u.PasswordHash,
		u.IsVerified,
		u.IsBanned,
		u.FailedAttempts,
		u.LockedUntil,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
	return err
}

func (r *UserRepository) FindByLoginIdentifier(ctx context.Context, identifier string) (*user.User, error) {
	var u user.User
	var email, username sql.NullString

	err := r.db.QueryRowContext(ctx, userFindByLoginSQL, identifier).Scan(
		&u.ID,
		&email,
		&username,
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

	u.Email = email.String
	u.Username = username.String
	return &u, nil
}

func (r *UserRepository) FindByUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, userFindByUsernameExistsSQL, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) FindByEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, userFindByEmailExistsSQL, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) FindByIDExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, userFindByIDExistsSQL, id).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	var emailVal sql.NullString

	err := r.db.QueryRowContext(ctx, userFindByEmailSQL, email).Scan(
		&u.ID,
		&emailVal,
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

	u.Email = emailVal.String
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User

	var email sql.NullString

	err := r.db.QueryRowContext(ctx, userFindByIDSQL, id).Scan(
		&u.ID,
		&email,
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

	u.Email = email.String
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

func (r *UserRepository) ResetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	res, err := r.db.ExecContext(ctx, userResetPasswordSQL, passwordHash, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
