package dashboard

type AdminDashboardResponse struct {
	TotalStudents		int							`json:"totalStudents"`
	ActiveStudents		int							`json:"activeStudents"`
	PendingStudents		int							`json:"pendingStudents"`
	PendingSubmissions	int							`json:"pendingSubmissions"`
	RecentSubmissions	[]DashboardSubmission		`json:"recentSubmissions"`
	UpcomingBatches		[]DashboardUpcomingBatch	`json:"upcomingBatches"`
	RecentStudents		[]DashboardStudent			`json:"recentStudents"`
}

type DashboardSubmission struct {
	ID				string	`json:"id"`
	StudentName		string	`json:"studentName"`
	AssignmentTitle	string	`json:"assignmentTitle"`
	CourseTitle		string	`json:"courseTitle"`
	Status			string	`json:"status"`
	SubmittedAt		string	`json:"submittedAt"`
}

type DashboardUpcomingBatch struct {
	ID			string	`json:"id"`
	BatchName	string	`json:"batchName"`
	StartDate	string	`json:"startDate"`
	Price		int		`json:"price"`
	Capacity	int		`json:"capacity"`
	Enrolled	int		`json:"enrolled"`
}

type DashboardStudent struct {
	ID			string	`json:"id"`
	FullName	string	`json:"fullName"`
	Email		string	`json:"email"`
	Status		string	`json:"status"`
	BatchName	string	`json:"batchName,omitempty"`
	CreatedAt	string	`json:"createdAt"`
}