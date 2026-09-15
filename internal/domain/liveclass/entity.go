package liveclass

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("live session not found")
	ErrForbidden   = errors.New("forbidden")
	ErrAlreadyLive = errors.New("a live session is already running for this batch")
	ErrNotLive     = errors.New("session is not live")
	ErrInvalid     = errors.New("invalid request")
)

const (
	StatusLive  = "live"
	StatusEnded = "ended"
)

// LiveSession is an ad-hoc online class started by an admin (teacher) for a
// batch taking a course.
type LiveSession struct {
	ID        uuid.UUID
	BatchID   uuid.UUID
	CourseID  uuid.UUID
	HostID    uuid.UUID
	Title     string
	ClassDate time.Time
	Status    string
	StartedAt time.Time
	EndedAt   *time.Time

	// Populated on reads for display.
	CourseTitle string
	BatchName   string
}

// BatchOption / CourseOption are lightweight pickers for the student browse UI.
type BatchOption struct {
	ID   uuid.UUID
	Name string
}

type CourseOption struct {
	ID    uuid.UUID
	Title string
}

// WhiteboardSnapshot is a saved board (scene JSON in R2 + optional PNG preview).
type WhiteboardSnapshot struct {
	ID         uuid.UUID
	SessionID  uuid.UUID
	StorageKey string
	ImageURL   string
	Version    int
	CreatedBy  uuid.UUID
	CreatedAt  time.Time
}

// SaveSnapshot carries the scene JSON and optional preview URL at save time.
type SaveSnapshot struct {
	Scene    json.RawMessage
	ImageURL string
}
