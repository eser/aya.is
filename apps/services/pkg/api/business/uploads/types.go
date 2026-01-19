package uploads

import "time"

// Purpose represents the intended use of an uploaded file.
type Purpose string

const (
	PurposeContentImage   Purpose = "content-image"
	PurposeCoverImage     Purpose = "cover-image"
	PurposeProfilePicture Purpose = "profile-picture"
)

// PresignedURLRequest contains the request data for generating a presigned URL.
type PresignedURLRequest struct {
	Filename    string  `json:"filename"`
	ContentType string  `json:"content_type"`
	Purpose     Purpose `json:"purpose"`
}

// PresignedURLResponse contains the response data for a presigned URL request.
type PresignedURLResponse struct {
	UploadURL string    `json:"upload_url"`
	Key       string    `json:"key"`
	PublicURL string    `json:"public_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
