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
	OrderNo		int			`json:"orderNo"`
}

type UpdateLessonOrderRequest struct {
    OrderNo int `json:"orderNo" binding:"required"`
}

type UpdateLessonRequest struct {
	ID			uuid.UUID	`json:"id"`
	CourseID	uuid.UUID	`json:"courseId"`
	Title 		*string 	`json:"title,omitempty"`
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
	OrderNo		int			`json:"orderNo"`
}