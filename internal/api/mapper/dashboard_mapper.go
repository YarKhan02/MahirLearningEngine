package mapper

import (
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/dashboard"
)

func ToAdminDashboardResponse(d *dashboard.AdminDashboard) dto.AdminDashboardResponse {
	resp := dto.AdminDashboardResponse{
		TotalStudents:      d.Counts.TotalStudents,
		ActiveStudents:     d.Counts.ActiveStudents,
		PendingStudents:    d.Counts.PendingStudents,
		PendingSubmissions: d.Counts.PendingSubmissions,
		RecentSubmissions:  make([]dto.DashboardSubmission, 0, len(d.RecentSubmissions)),
		UpcomingBatches:    make([]dto.DashboardUpcomingBatch, 0, len(d.UpcomingBatches)),
		RecentStudents:     make([]dto.DashboardStudent, 0, len(d.RecentStudents)),
	}

	for _, s := range d.RecentSubmissions {
		resp.RecentSubmissions = append(resp.RecentSubmissions, dto.DashboardSubmission{
			ID:              s.ID.String(),
			StudentName:     s.StudentName,
			AssignmentTitle: s.AssignmentTitle,
			CourseTitle:     s.CourseTitle,
			Status:          s.Status,
			SubmittedAt:     s.SubmittedAt.Format(time.RFC3339),
		})
	}

	for _, b := range d.UpcomingBatches {
		resp.UpcomingBatches = append(resp.UpcomingBatches, dto.DashboardUpcomingBatch{
			ID:        b.ID.String(),
			BatchName: b.BatchName,
			StartDate: b.StartDate.Format(constant.DateLayout),
			Days:      b.Days,
			Price:     b.Price,
			Capacity:  b.Capacity,
			Enrolled:  b.Enrolled,
		})
	}

	for _, s := range d.RecentStudents {
		st := dto.DashboardStudent{
			ID:        s.ID.String(),
			FullName:  s.FullName,
			Email:     s.Email,
			Status:    s.Status,
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
		}
		if s.BatchName != nil {
			st.BatchName = *s.BatchName
		}
		resp.RecentStudents = append(resp.RecentStudents, st)
	}

	return resp
}
