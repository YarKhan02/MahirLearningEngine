package attachement

type PresignRequest struct {
	Filename     string `json:"filename"`
	ContentType  string `json:"content_type"`
	SizeBytes    int64  `json:"size_bytes"`
	ResourceType string `json:"resource_type"` // e.g. "course"
	ResourceID   string `json:"resource_id"`   // e.g. the course id
}

type ConfirmRequest struct {
	Key string `json:"key"`
}

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
	ExpiresIn int    `json:"expires_in"`
}

type AttachmentResponse struct {
	ID          string `json:"id"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	SizeBytes   *int64 `json:"sizeBytes,omitempty"`
	CreatedAt   string `json:"createdAt"`
	DownloadURL string `json:"downloadUrl"`
}
