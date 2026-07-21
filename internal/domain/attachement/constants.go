package attachement

var AllowedTypes = map[string]bool{
	"image/png":       true,
	"image/jpeg":      true,
	"image/webp":      true,
	"image/gif":       true,
	"application/pdf": true,
	"application/vnd.ms-powerpoint": true, // .ppt
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
}

const MaxUploadSize = 25 * 1024 * 1024 // 25MB

const DownloadURLTTL = 60 * 60 // 1 hour, in seconds
