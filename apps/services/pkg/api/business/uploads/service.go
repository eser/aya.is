package uploads

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/s3client"
)

const (
	DefaultPresignExpiration = 15 * time.Minute
)

var (
	ErrInvalidPurpose       = errors.New("invalid upload purpose")
	ErrInvalidContentType   = errors.New("invalid content type")
	ErrFailedToGenerateURL  = errors.New("failed to generate presigned URL")
	ErrFailedToRemoveObject = errors.New("failed to remove object")
)

// allowedContentTypes maps purposes to their allowed MIME types.
var allowedContentTypes = map[Purpose][]string{ //nolint:gochecknoglobals
	PurposeContentImage: {
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	},
	PurposeCoverImage: {
		"image/jpeg",
		"image/png",
		"image/webp",
	},
	PurposeProfilePicture: {
		"image/jpeg",
		"image/png",
		"image/webp",
	},
}

// Service handles upload operations.
type Service struct {
	logger   *logfx.Logger
	s3Client *s3client.Client
}

// NewService creates a new upload service.
func NewService(logger *logfx.Logger, s3Client *s3client.Client) *Service {
	return &Service{
		logger:   logger,
		s3Client: s3Client,
	}
}

// GenerateUploadURL generates a presigned URL for uploading a file.
func (s *Service) GenerateUploadURL(
	ctx context.Context,
	userID string,
	req PresignedURLRequest,
) (*PresignedURLResponse, error) {
	// Validate purpose
	allowedTypes, ok := allowedContentTypes[req.Purpose]
	if !ok {
		return nil, ErrInvalidPurpose
	}

	// Validate content type
	if !isAllowedContentType(req.ContentType, allowedTypes) {
		return nil, ErrInvalidContentType
	}

	// Generate unique key
	key := generateKey(userID, req.Purpose, req.Filename)

	// Generate presigned URL
	expiresAt := time.Now().Add(DefaultPresignExpiration)

	uploadURL, err := s.s3Client.GeneratePresignedUploadURL(
		ctx,
		key,
		req.ContentType,
		DefaultPresignExpiration,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGenerateURL, err)
	}

	// Get public URL
	publicURL := s.s3Client.GetPublicURL(key)

	return &PresignedURLResponse{
		UploadURL: uploadURL,
		Key:       key,
		PublicURL: publicURL,
		ExpiresAt: expiresAt,
	}, nil
}

// RemoveObject removes an uploaded file.
func (s *Service) RemoveObject(ctx context.Context, key string) error {
	err := s.s3Client.RemoveObject(ctx, key)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRemoveObject, err)
	}

	return nil
}

// generateKey creates a unique storage key for the upload.
// Format: {purpose}/{user-id}/{timestamp}-{random}.{ext}.
func generateKey(userID string, purpose Purpose, filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".bin"
	}

	timestamp := time.Now().Unix()
	random := lib.IDsGenerateUnique()[:8]

	return fmt.Sprintf("%s/%s/%d-%s%s", purpose, userID, timestamp, random, ext)
}

// isAllowedContentType checks if a content type is in the allowed list.
func isAllowedContentType(contentType string, allowed []string) bool {
	contentType = strings.ToLower(contentType)

	for _, allowedType := range allowed {
		if contentType == allowedType {
			return true
		}
	}

	return false
}
