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

var ExtByType = map[string]string{
	"image/png":       ".png",
	"image/jpeg":      ".jpg",
	"image/webp":      ".webp",
	"image/gif":       ".gif",
	"application/pdf": ".pdf",
	"application/vnd.ms-powerpoint": ".ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
}

var OfficeTypes = map[string]bool{
	"application/vnd.ms-powerpoint": true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
}

const HeaderSniffBytes = 32 * 1024

const MaxFileNameLen = 200

// Inline content images (embedded inside topic rich-text). Images only, capped
// smaller than course materials since they're received in-memory by the server.
var InlineImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/webp": true,
	"image/gif":  true,
}

const MaxInlineImageSize = 5 * 1024 * 1024 // 5MB

const InlineURLTTL = 60 * 60 // seconds; TTL of the presigned URL the public route redirects to
