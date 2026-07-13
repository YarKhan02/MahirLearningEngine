package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/dashboard"
)

//go:embed sql/dashboard_counts.sql
var dashboardCountsSQL string

//go:embed sql/dashboard_recent_submissions.sql
var dashboardRecentSubmissionsSQL string

//go:embed sql/dashboard_upcoming_batches.sql
var dashboardUpcomingBatchesSQL string

//go:embed sql/dashboard_recent_students.sql
var dashboardRecentStudentsSQL string

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetAdminDashboard(ctx context.Context) (*dashboard.AdminDashboard, error) {
	var d dashboard.AdminDashboard

	err := r.db.QueryRowContext(ctx, dashboardCountsSQL).Scan(
		&d.Counts.TotalStudents,
		&d.Counts.ActiveStudents,
		&d.Counts.PendingStudents,
		&d.Counts.PendingSubmissions,
	)
	if err != nil {
		return nil, fmt.Errorf("dashboard counts: %w", err)
	}

	if d.RecentSubmissions, err = r.recentSubmissions(ctx); err != nil {
		return nil, err
	}
	if d.UpcomingBatches, err = r.upcomingBatches(ctx); err != nil {
		return nil, err
	}
	if d.RecentStudents, err = r.recentStudents(ctx); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *DashboardRepository) recentSubmissions(ctx context.Context) ([]dashboard.RecentSubmission, error) {
	rows, err := r.db.QueryContext(ctx, dashboardRecentSubmissionsSQL)
	if err != nil {
		return nil, fmt.Errorf("recent submissions: %w", err)
	}
	defer rows.Close()

	var out []dashboard.RecentSubmission
	for rows.Next() {
		var s dashboard.RecentSubmission
		if err := rows.Scan(&s.ID, &s.StudentName, &s.AssignmentTitle, &s.CourseTitle, &s.Status, &s.SubmittedAt); err != nil {
			return nil, fmt.Errorf("scan recent submission: %w", err)
		}
		out = append(out, s)
	}

	return out, rows.Err()
}

func (r *DashboardRepository) upcomingBatches(ctx context.Context) ([]dashboard.UpcomingBatch, error) {
	rows, err := r.db.QueryContext(ctx, dashboardUpcomingBatchesSQL)
	if err != nil {
		return nil, fmt.Errorf("upcoming batches: %w", err)
	}
	defer rows.Close()

	var out []dashboard.UpcomingBatch
	for rows.Next() {
		var b dashboard.UpcomingBatch
		if err := rows.Scan(&b.ID, &b.BatchName, &b.StartDate, &b.Days, &b.Price, &b.Capacity, &b.Enrolled); err != nil {
			return nil, fmt.Errorf("scan upcoming batch: %w", err)
		}
		out = append(out, b)
	}

	return out, rows.Err()
}

func (r *DashboardRepository) recentStudents(ctx context.Context) ([]dashboard.RecentStudent, error) {
	rows, err := r.db.QueryContext(ctx, dashboardRecentStudentsSQL)
	if err != nil {
		return nil, fmt.Errorf("recent students: %w", err)
	}
	defer rows.Close()

	var out []dashboard.RecentStudent
	for rows.Next() {
		var s dashboard.RecentStudent
		if err := rows.Scan(&s.ID, &s.FullName, &s.Email, &s.Status, &s.BatchName, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan recent student: %w", err)
		}
		out = append(out, s)
	}

	return out, rows.Err()
}
