package topic

import "github.com/google/uuid"

// Topic is an ordered unit of material inside a lesson: a title, optional
// description, rich-text content, and an optional YouTube video.
type Topic struct {
	ID          uuid.UUID
	LessonID    uuid.UUID
	Title       string
	Description string
	Content     string
	YoutubeURL  string
	OrderNo     int
}

// UpdateTopic is a partial update keyed by topic id; nil fields are left as-is.
type UpdateTopic struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Content     *string
	YoutubeURL  *string
}
