package repository

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/google/uuid"
)

//go:embed sql/course_insert.sql
var courseInsertSQL string

//go:embed sql/course_get.sql
var courseGetSQL string

type CoursRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CoursRepository {
	return &CoursRepository{
		db: db,
	}
}

func (r *CoursRepository) InsertCourse(ctx context.Context, req dto.InsertCourse) (*course.Course, error) {
	
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	var description any
	if req.Description == "" {
		description = nil
	} else {
		description = req.Description
	}

	var course course.Course

	err = r.db.QueryRowContext(ctx, courseInsertSQL, 
		id, 
		req.Title,
		req.Level, 
		req.Duration, 
		description,
	).Scan(
		&course.ID,
		&course.Title,
		&course.Level,
		&course.Duration,
		&course.Description,
		&course.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CoursRepository) GetCourse(ctx context.Context) ([]course.Course, error) {

	rows, err := r.db.QueryContext(ctx, courseGetSQL)
	if err != nil {
		return nil, err
	}

	courses := make([]course.Course, 0)

	for rows.Next() {
		var c course.Course

		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.Level,
			&c.Duration,
			&c.IsActive,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err 
	}

	return courses, nil
}