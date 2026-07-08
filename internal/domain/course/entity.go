package course

import "github.com/google/uuid"

type Course struct {
	ID 			string
	Title 		string
	Level 		string
	Duration 	int
	Description string
	IsActive	bool
}

type Lesson struct {
	ID			uuid.UUID
	CourseID	uuid.UUID
	Title 		string
	Description string
	OrderNo		int
	YoutubeURL 	string
	Content		string
}

type UpdateLesson struct {
	ID			uuid.UUID
	CourseID	uuid.UUID
	Title 		*string
	Description *string
	YoutubeURL 	*string
	Content		*string
}