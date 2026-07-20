package attendance

type RosterEntryResponse struct {
	StudentID	string	`json:"studentId"`
	FullName	string	`json:"fullName"`
	Email		string	`json:"email"`
	Status		*string	`json:"status,omitempty"`
}

type MarkAttendanceRequest struct {
	Date		string	`json:"date" binding:"required"`
	StudentID	string	`json:"studentId" binding:"required"`
	Status		string	`json:"status" binding:"required"`
}

type AttendanceRecordResponse struct {
	Date		string	`json:"date"`
	Status		string	`json:"status"`
	BatchName	string	`json:"batchName"`
}

type AttendanceSummaryResponse struct {
	Present     int     `json:"present"`
	Absent      int     `json:"absent"`
	Total       int     `json:"total"`
	TodayStatus *string `json:"todayStatus,omitempty"`
}