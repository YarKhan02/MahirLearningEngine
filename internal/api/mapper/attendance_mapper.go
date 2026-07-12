package mapper

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
)

func ToRosterEntryResponse(req attendance.RosterEntry) dto.RosterEntryResponse {
	return dto.RosterEntryResponse{
		StudentID: req.StudentID.String(),
		FullName:  req.FullName,
		Email:     req.Email,
		Status:    req.Status,
	}
}

func ToAttendanceRecordResponse(req attendance.Record) dto.AttendanceRecordResponse {
	return dto.AttendanceRecordResponse{
		Date:      req.LessonDate.Format(dateLayout),
		Status:    req.Status,
		BatchName: req.BatchName,
	}
}
