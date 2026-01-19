package s3client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ErrMissingConfig        = errors.New("missing required S3 configuration")
	ErrFailedToCreateClient = errors.New("failed to create S3 client")
	ErrFailedToPresignURL   = errors.New("failed to generate presigned URL")
	ErrFailedToRemoveObject = errors.New("failed to remove object")
)

// Config holds the S3 client configuration.
type Config struct {
	Endpoint        string // e.g., "https://<account>.r2.cloudflarestorage.com"
	Region          string // e.g., "auto" for R2
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicURL       string // e.g., "https://cdn.aya.is"
}

// Client is an S3-compatible storage client.
type Client struct {
	s3Client   *s3.Client
	presigner  *s3.PresignClient
	bucketName string
	publicURL  string
}

// New creates a new S3-compatible client.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" ||
		cfg.SecretAccessKey == "" || cfg.BucketName == "" {
		return nil, ErrMissingConfig
	}

	region := cfg.Region
	if region == "" {
		region = "auto"
	}

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateClient, err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})

	return &Client{
		s3Client:   s3Client,
		presigner:  s3.NewPresignClient(s3Client),
		bucketName: cfg.BucketName,
		publicURL:  cfg.PublicURL,
	}, nil
}

// GeneratePresignedUploadURL generates a pre-signed URL for uploading a file.
func (c *Client) GeneratePresignedUploadURL(
	ctx context.Context,
	key string,
	contentType string,
	expiresIn time.Duration,
) (string, error) {
	input := &s3.PutObjectInput{ //nolint:exhaustruct
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	presignedReq, err := c.presigner.PresignPutObject(ctx, input,
		s3.WithPresignExpires(expiresIn),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToPresignURL, err)
	}

	return presignedReq.URL, nil
}

// GeneratePresignedDownloadURL generates a pre-signed URL for downloading a file.
func (c *Client) GeneratePresignedDownloadURL(
	ctx context.Context,
	key string,
	expiresIn time.Duration,
) (string, error) {
	input := &s3.GetObjectInput{ //nolint:exhaustruct
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	presignedReq, err := c.presigner.PresignGetObject(ctx, input,
		s3.WithPresignExpires(expiresIn),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToPresignURL, err)
	}

	return presignedReq.URL, nil
}

// RemoveObject deletes an object from the bucket.
func (c *Client) RemoveObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{ //nolint:exhaustruct
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	_, err := c.s3Client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRemoveObject, err)
	}

	return nil
}

// GetPublicURL returns the public URL for a given key.
func (c *Client) GetPublicURL(key string) string {
	if c.publicURL == "" {
		return ""
	}

	return c.publicURL + "/" + key
}

// GetBucketName returns the configured bucket name.
func (c *Client) GetBucketName() string {
	return c.bucketName
}
