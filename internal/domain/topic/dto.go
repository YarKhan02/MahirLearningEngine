package topic

type InsertTopicRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	YoutubeURL  string `json:"youtubeUrl,omitempty"`
}

type UpdateTopicRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Content     *string `json:"content,omitempty"`
	YoutubeURL  *string `json:"youtubeUrl,omitempty"`
}

type UpdateTopicOrderRequest struct {
	OrderNo int `json:"orderNo" binding:"required"`
}

type TopicResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	YoutubeURL  string `json:"youtubeUrl"`
	OrderNo     int    `json:"orderNo"`
}
