package timetable

import (
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	
	"github.com/google/uuid"
)

func ToCreateTimetable(batchID uuid.UUID, req CreateTimetableRequest) (*Timetable, error) {
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid courseId: %w", err)
	}

	return &Timetable{
		BatchID:   batchID,
		CourseID:  courseID,
		Weekdays:  req.Weekdays,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}, nil
}

func ToTimetableResponse(t Timetable) TimetableResponse {
	return TimetableResponse{
		ID:          t.ID.String(),
		BatchID:     t.BatchID.String(),
		CourseID:    t.CourseID.String(),
		CourseTitle: t.CourseTitle,
		Weekdays:    t.Weekdays,
		StartTime:   t.StartTime,
		EndTime:     t.EndTime,
	}
}

func ToClassSessionResponse(s ClassSession) ClassSessionResponse {
	return ClassSessionResponse{
		Date:        s.Date.Format(constant.DateLayout),
		Weekday:     s.Weekday,
		StartTime:   s.StartTime,
		EndTime:     s.EndTime,
		CourseID:    s.CourseID.String(),
		CourseTitle: s.CourseTitle,
		BatchID:     s.BatchID.String(),
		BatchName:   s.BatchName,
	}
}
