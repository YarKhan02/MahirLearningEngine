package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/program"
	"github.com/google/uuid"
)

//go:embed sql/program_create.sql
var programCreateSQL string

//go:embed sql/program_update.sql
var programUpdateSQL string

//go:embed sql/program_set_published.sql
var programSetPublishedSQL string

//go:embed sql/program_delete.sql
var programDeleteSQL string

//go:embed sql/program_list_all.sql
var programListAllSQL string

//go:embed sql/program_get_by_id.sql
var programGetByIDSQL string

//go:embed sql/program_list_published.sql
var programListPublishedSQL string

//go:embed sql/program_by_slug.sql
var programBySlugSQL string

//go:embed sql/program_slug_exists.sql
var programSlugExistsSQL string

//go:embed sql/program_max_order.sql
var programMaxOrderSQL string

type ProgramRepository struct {
	db *sql.DB
}

func NewProgramRepository(db *sql.DB) *ProgramRepository {
	return &ProgramRepository{db: db}
}

type programScanner interface {
	Scan(dest ...any) error
}

func scanProgram(s programScanner) (program.Program, error) {
	var p program.Program
	var learnRaw []byte
	if err := s.Scan(
		&p.ID, &p.Slug, &p.Title, &p.Short, &p.Full, &learnRaw,
		&p.Age, &p.Fee, &p.Level, &p.Bonus, &p.Accent, &p.ImageURL,
		&p.OrderNo, &p.Published, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return program.Program{}, err
	}
	if len(learnRaw) > 0 {
		if err := json.Unmarshal(learnRaw, &p.Learn); err != nil {
			return program.Program{}, fmt.Errorf("decode learn: %w", err)
		}
	}
	if p.Learn == nil {
		p.Learn = []string{}
	}
	return p, nil
}

func (r *ProgramRepository) Create(ctx context.Context, p program.Program) error {
	learn, err := json.Marshal(p.Learn)
	if err != nil {
		return fmt.Errorf("encode learn: %w", err)
	}
	_, err = r.db.ExecContext(ctx, programCreateSQL,
		p.ID, p.Slug, p.Title, p.Short, p.Full, learn,
		p.Age, p.Fee, p.Level, p.Bonus, p.Accent, p.ImageURL,
		p.OrderNo, p.Published,
	)
	if err != nil {
		return fmt.Errorf("create program: %w", err)
	}
	return nil
}

func (r *ProgramRepository) Update(ctx context.Context, id uuid.UUID, u program.Upsert) (program.Program, error) {
	learn, err := json.Marshal(u.Learn)
	if err != nil {
		return program.Program{}, fmt.Errorf("encode learn: %w", err)
	}
	row := r.db.QueryRowContext(ctx, programUpdateSQL,
		id, u.Title, u.Short, u.Full, learn, u.Age,
		u.Fee, u.Level, u.Bonus, u.Accent, u.ImageURL, u.Published,
	)
	p, err := scanProgram(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return program.Program{}, program.ErrNotFound
		}
		return program.Program{}, fmt.Errorf("update program: %w", err)
	}
	return p, nil
}

func (r *ProgramRepository) SetPublished(ctx context.Context, id uuid.UUID, published bool) (program.Program, error) {
	row := r.db.QueryRowContext(ctx, programSetPublishedSQL, id, published)
	p, err := scanProgram(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return program.Program{}, program.ErrNotFound
		}
		return program.Program{}, fmt.Errorf("set published: %w", err)
	}
	return p, nil
}

func (r *ProgramRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, programDeleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete program: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete program: rows affected: %w", err)
	}
	if n == 0 {
		return program.ErrNotFound
	}
	return nil
}

func (r *ProgramRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	for i, id := range ids {
		if _, err := tx.ExecContext(ctx,
			`UPDATE programs SET order_no = $1, updated_at = NOW() WHERE id = $2`, i, id,
		); err != nil {
			return fmt.Errorf("reorder program: %w", err)
		}
	}
	return tx.Commit()
}

func (r *ProgramRepository) ListAll(ctx context.Context) ([]program.Program, error) {
	return r.queryList(ctx, programListAllSQL)
}

func (r *ProgramRepository) ListPublished(ctx context.Context) ([]program.Program, error) {
	return r.queryList(ctx, programListPublishedSQL)
}

func (r *ProgramRepository) queryList(ctx context.Context, query string) ([]program.Program, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list programs: %w", err)
	}
	defer rows.Close()

	out := make([]program.Program, 0)
	for rows.Next() {
		p, err := scanProgram(rows)
		if err != nil {
			return nil, fmt.Errorf("scan program: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate programs: %w", err)
	}
	return out, nil
}

func (r *ProgramRepository) GetByID(ctx context.Context, id uuid.UUID) (program.Program, error) {
	row := r.db.QueryRowContext(ctx, programGetByIDSQL, id)
	p, err := scanProgram(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return program.Program{}, program.ErrNotFound
		}
		return program.Program{}, fmt.Errorf("get program: %w", err)
	}
	return p, nil
}

func (r *ProgramRepository) GetPublishedBySlug(ctx context.Context, slug string) (program.Program, error) {
	row := r.db.QueryRowContext(ctx, programBySlugSQL, slug)
	p, err := scanProgram(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return program.Program{}, program.ErrNotFound
		}
		return program.Program{}, fmt.Errorf("get program by slug: %w", err)
	}
	return p, nil
}

func (r *ProgramRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, programSlugExistsSQL, slug).Scan(&exists); err != nil {
		return false, fmt.Errorf("slug exists: %w", err)
	}
	return exists, nil
}

func (r *ProgramRepository) MaxOrderNo(ctx context.Context) (int, error) {
	var max int
	if err := r.db.QueryRowContext(ctx, programMaxOrderSQL).Scan(&max); err != nil {
		return 0, fmt.Errorf("max order: %w", err)
	}
	return max, nil
}
