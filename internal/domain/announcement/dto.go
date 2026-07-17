package announcement

type CreateAnnouncementRequest struct {
	BatchID		string	`json:"batchId" binding:"required"`
	Title		string	`json:"title" binding:"required"`
	Description	string	`json:"description" binding:"required"`
}

type AnnouncementResponse struct {
	ID			string	`json:"id"`
	BatchID		string	`json:"batchId"`
	BatchName	string	`json:"batchName"`
	Title		string	`json:"title"`
	Description	string	`json:"description"`
	CreatedAt	string	`json:"createdAt"`
}