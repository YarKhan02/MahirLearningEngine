package course

import "github.com/google/uuid"

type InsertCourse struct {
	Title 		string 		`json:"title"`
	Level 		string 		`json:"level"`
	Duration 	int 		`json:"duration"`
	Description string 		`json:"description,omitempty"`
}

type InsertLesson struct {
	ID			uuid.UUID	`json:"id"`
	CourseID	uuid.UUID	`json:"courseId"`
	Title 		string 		`json:"title"`
	Description string 		`json:"description,omitempty"`
	OrderNo		int			`json:"orderNo"`
	YoutubeURL 	string		`json:"youtubeUrl,omitempty"`
	Content		string		`json:"content,omitempty"`
}

type UpdateLessonOrderRequest struct {
    OrderNo int `json:"orderNo" binding:"required"`
}

type UpdateLessonRequest struct {
	ID			uuid.UUID	`json:"id"`
	CourseID	uuid.UUID	`json:"courseId"`
	Title 		*string 	`json:"title,omitempty"`
	Description *string 	`json:"description,omitempty"`
	YoutubeURL 	*string		`json:"youtubeUrl,omitempty"`
	Content		*string		`json:"content,omitempty"`
}

type CourseResponse struct {
	ID 			string		`json:"id"`
	Title 		string		`json:"title"`
	Level 		string		`json:"level"`
	Duration 	int			`json:"duration"`
	Description string		`json:"description"`
	Status		string		`json:"status"`
}

type LessonResponse struct {
	ID			string		`json:"id"`
	Title 		string		`json:"title"`
	Description string		`json:"description"`
	OrderNo		int			`json:"orderNo"`
	YoutubeURL 	string		`json:"youtubeUrl"`
	Content		string		`json:"content"`
}