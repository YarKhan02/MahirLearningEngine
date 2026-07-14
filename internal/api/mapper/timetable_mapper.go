package mapper

import (
	"fmt"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/timetable"
	"github.com/google/uuid"
)

func ToCreateTimetable(batchID uuid.UUID, req dto.CreateTimetableRequest) (*timetable.Timetable, error) {
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid courseId: %w", err)
	}

	return &timetable.Timetable{
		BatchID:   batchID,
		CourseID:  courseID,
		Weekdays:  req.Weekdays,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}, nil
}

func ToTimetableResponse(t timetable.Timetable) dto.TimetableResponse {
	return dto.TimetableResponse{
		ID:          t.ID.String(),
		BatchID:     t.BatchID.String(),
		CourseID:    t.CourseID.String(),
		CourseTitle: t.CourseTitle,
		Weekdays:    t.Weekdays,
		StartTime:   t.StartTime,
		EndTime:     t.EndTime,
	}
}

func ToClassSessionResponse(s timetable.ClassSession) dto.ClassSessionResponse {
	return dto.ClassSessionResponse{
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
