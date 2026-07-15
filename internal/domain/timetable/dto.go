package timetable

type CreateTimetableRequest struct {
	CourseID	string	`json:"courseId" binding:"required"`
	Weekdays	[]int	`json:"weekdays" binding:"required"`
	StartTime	string	`json:"startTime" binding:"required"`
	EndTime		string	`json:"endTime" binding:"required"`
}

type TimetableResponse struct {
	ID			string	`json:"id"`
	BatchID		string	`json:"batchId"`
	CourseID	string	`json:"courseId"`
	CourseTitle	string	`json:"courseTitle"`
	Weekdays	[]int	`json:"weekdays"`
	StartTime	string	`json:"startTime"`
	EndTime		string	`json:"endTime"`
}

type ClassSessionResponse struct {
	Date		string	`json:"date"`
	Weekday		int		`json:"weekday"`
	StartTime	string	`json:"startTime"`
	EndTime		string	`json:"endTime"`
	CourseID	string	`json:"courseId"`
	CourseTitle	string	`json:"courseTitle"`
	BatchID		string	`json:"batchId"`
	BatchName	string	`json:"batchName"`
}