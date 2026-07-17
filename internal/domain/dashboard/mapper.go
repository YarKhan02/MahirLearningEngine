package dashboard

import (
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
)

func ToAdminDashboardResponse(d *AdminDashboard) AdminDashboardResponse {
	resp := AdminDashboardResponse{
		TotalStudents:      d.Counts.TotalStudents,
		ActiveStudents:     d.Counts.ActiveStudents,
		PendingStudents:    d.Counts.PendingStudents,
		PendingSubmissions: d.Counts.PendingSubmissions,
		RecentSubmissions:  make([]DashboardSubmission, 0, len(d.RecentSubmissions)),
		UpcomingBatches:    make([]DashboardUpcomingBatch, 0, len(d.UpcomingBatches)),
		RecentStudents:     make([]DashboardStudent, 0, len(d.RecentStudents)),
	}

	for _, s := range d.RecentSubmissions {
		resp.RecentSubmissions = append(resp.RecentSubmissions, DashboardSubmission{
			ID:              s.ID.String(),
			StudentName:     s.StudentName,
			AssignmentTitle: s.AssignmentTitle,
			CourseTitle:     s.CourseTitle,
			Status:          s.Status,
			SubmittedAt:     s.SubmittedAt.Format(time.RFC3339),
		})
	}

	for _, b := range d.UpcomingBatches {
		resp.UpcomingBatches = append(resp.UpcomingBatches, DashboardUpcomingBatch{
			ID:        b.ID.String(),
			BatchName: b.BatchName,
			StartDate: b.StartDate.Format(constant.DateLayout),
			Price:     b.Price,
			Capacity:  b.Capacity,
			Enrolled:  b.Enrolled,
		})
	}

	for _, s := range d.RecentStudents {
		st := DashboardStudent{
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
