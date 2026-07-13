package mapper

import (
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/google/uuid"
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
		Date:      req.LessonDate.Format(constant.DateLayout),
		Status:    req.Status,
		BatchName: req.BatchName,
	}
}

func ToMarkAttendance(batchID uuid.UUID, date time.Time, studentID uuid.UUID, status string, createdBy uuid.UUID) attendance.MarkAttendance {
	return attendance.MarkAttendance{
		BatchID: 	batchID,
		Date: 		date,
		StudentID: 	studentID,
		Status: 	status,
		CreatedBy: 	createdBy,
	}
}
