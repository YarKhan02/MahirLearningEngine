package attendance

import (
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	"github.com/google/uuid"
)

func ToRosterEntryResponse(req RosterEntry) RosterEntryResponse {
	return RosterEntryResponse{
		StudentID: req.StudentID.String(),
		FullName:  req.FullName,
		Email:     req.Email,
		Status:    req.Status,
	}
}

func ToAttendanceRecordResponse(req Record) AttendanceRecordResponse {
	return AttendanceRecordResponse{
		Date:      req.LessonDate.Format(constant.DateLayout),
		Status:    req.Status,
		BatchName: req.BatchName,
	}
}

func ToMarkAttendance(batchID uuid.UUID, date time.Time, studentID uuid.UUID, status string, createdBy uuid.UUID) MarkAttendance {
	return MarkAttendance{
		BatchID: 	batchID,
		Date: 		date,
		StudentID: 	studentID,
		Status: 	status,
		CreatedBy: 	createdBy,
	}
}
