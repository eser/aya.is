package uploads

import (
	"context"
	"time"
)

// StorageClient defines the interface for object storage operations.
// This interface allows the business layer to remain decoupled from specific
// storage implementations (e.g., S3, R2, local filesystem).
type StorageClient interface {
	// GeneratePresignedUploadURL generates a pre-signed URL for uploading a file.
	GeneratePresignedUploadURL(
		ctx context.Context,
		key string,
		contentType string,
		expiresIn time.Duration,
	) (string, error)

	// GetPublicURL returns the public URL for a given key.
	GetPublicURL(key string) string

	// RemoveObject deletes an object from storage.
	RemoveObject(ctx context.Context, key string) error
}

// Purpose represents the intended use of an uploaded file.
type Purpose string

const (
	PurposeContentImage   Purpose = "content-image"
	PurposeStoryPicture   Purpose = "story-picture"
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
