package repository

import (
	"database/sql"
	"context"
	_ "embed"

	"github.com/google/uuid"
)

//go:embed sql/role_add_to_user.sql
var roleAddToUserSQL string

//go:embed sql/role_get_user.sql
var roleGetUserSQL string

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) AddRoleToUser(ctx context.Context, userID uuid.UUID, role string) error {
	_, err := r.db.ExecContext(ctx, roleAddToUserSQL, userID, role)
	if err != nil {
		return err
	}	
	return nil
}

func (r *RoleRepository) GetUserRole(ctx context.Context, userID uuid.UUID) (string, error) {
	var role string

	err := r.db.QueryRowContext(ctx, roleGetUserSQL, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return role, nil
}
