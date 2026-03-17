package aifx

import "time"

// BatchStatus represents the current state of a batch job.
type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "pending"
	BatchStatusProcessing BatchStatus = "processing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusFailed     BatchStatus = "failed"
	BatchStatusCancelled  BatchStatus = "cancelled"
)

// BatchRequest contains a set of generation requests to process asynchronously.
type BatchRequest struct {
	Items []BatchRequestItem
}

// BatchRequestItem is a single item in a batch request.
type BatchRequestItem struct {
	CustomID string              // caller-defined ID for matching results
	Options  GenerateTextOptions // same options as real-time generation
}

// BatchJob represents the status of a submitted batch job.
type BatchJob struct {
	CreatedAt   time.Time
	CompletedAt *time.Time
	Storage     *BatchStorage // provider-specific storage reference
	ID          string
	Status      BatchStatus
	Error       string
	TotalCount  int
	DoneCount   int
	FailedCount int
}

// BatchResult holds the result for a single item in a batch.
type BatchResult struct {
	CustomID string
	Result   *GenerateTextResult // nil if error
	Error    string
}

// BatchStorage holds provider-specific storage references for batch data.
type BatchStorage struct {
	Properties map[string]any // provider-specific extras (bucket, project, etc.)
	Type       string         // "openai_file", "cloud_storage", "inline"
	InputRef   string         // file ID, GCS URI, etc.
	OutputRef  string         // file ID, GCS URI, etc.
}

// ListBatchOptions configures batch listing requests.
type ListBatchOptions struct {
	After string // cursor for pagination
	Limit int
}
