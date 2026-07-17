package timetable

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoWeekdays    		= errors.New("select at least one weekday")
	ErrInvalidTime   		= errors.New("start and end time are required")
	ErrTimeOrder     		= errors.New("end time must be after start time")
	ErrTimetableNotFound 	= errors.New("timetable not found")
)

// upcomingWindowDays is how far ahead a student sees classes — the next week.
const upcomingWindowDays = 7

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, t *Timetable) error {
	if len(t.Weekdays) == 0 {
		return ErrNoWeekdays
	}
	if t.StartTime == "" || t.EndTime == "" {
		return ErrInvalidTime
	}
	// "HH:MM" is zero-padded, so lexical comparison matches chronological order.
	if t.EndTime <= t.StartTime {
		return ErrTimeOrder
	}
	return s.repo.Create(ctx, t)
}

func (s *Service) GetByBatch(ctx context.Context, batchID uuid.UUID) ([]Timetable, error) {
	return s.repo.GetByBatch(ctx, batchID)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetUpcomingForUser returns the class occurrences in the next week for the
// batches the given student is enrolled in.
func (s *Service) GetUpcomingForUser(ctx context.Context, userID uuid.UUID) ([]ClassSession, error) {
	rules, err := s.repo.GetRulesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return GenerateSessions(rules, now, now.AddDate(0, 0, upcomingWindowDays-1)), nil
}
