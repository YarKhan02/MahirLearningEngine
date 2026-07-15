package batch

type CreateBatchRequest struct {
	BatchName	string		`json:"batchName"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Capacity 	int 		`json:"capacity"`
	Status		string 		`json:"status"`
	Price		int			`json:"price"`
}
type UpdateBatchRequest struct {
	BatchName	string		`json:"batchName"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Capacity 	int 		`json:"capacity"`
	Status		string 		`json:"status"`
	Price		int			`json:"price"`
}

type UpdateBatchCoursesRequest struct {
	AddCourseIDs	[]string	`json:"addCourseIds"`
	RemoveCourseIDs	[]string	`json:"removeCourseIds"`
}

type BatchResponse struct {
	ID			string		`json:"id"`
	BatchName	string		`json:"batchName"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Capacity 	int 		`json:"capacity"`
	Status		string 		`json:"status"`
	Price		int			`json:"price"`
}

type BatchCourseResponse struct {
	ID			string	`json:"id"`
	CourseID	string	`json:"courseId"`
	Title		string	`json:"title"`
	Level		string	`json:"level"`
	GrantedAt	string	`json:"grantedAt"`
}

type PublicBatchResponse struct {
	ID			string					`json:"id"`
	BatchName	string					`json:"batchName"`
	StartDate	string					`json:"startDate"`
	EndDate		string					`json:"endDate"`
	Capacity	int						`json:"capacity"`
	Status		string					`json:"status"`
	Price		int						`json:"price"`
	Courses		[]BatchCourseResponse	`json:"courses"`
}