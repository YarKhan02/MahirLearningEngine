package program

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("program not found")
	ErrSlugTaken = errors.New("slug already in use")
	ErrInvalid   = errors.New("invalid program")
)

// Program is a standalone landing-page offering, fully admin-managed.
type Program struct {
	ID        uuid.UUID
	Slug      string
	Title     string
	Short     string
	Full      string
	Learn     []string
	Age       string
	Fee       string
	Level     string
	Bonus     string
	Accent    string
	ImageURL  string
	OrderNo   int
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Upsert carries the writable fields shared by create and update.
type Upsert struct {
	Title     string
	Short     string
	Full      string
	Learn     []string
	Age       string
	Fee       string
	Level     string
	Bonus     string
	Accent    string
	ImageURL  string
	Published bool
}
