package stories

import (
	"context"
	"errors"
)

// Sentinel errors for AI summarization.
var ErrSummarizationNotAvailable = errors.New("AI summarization not available")

// ContentSummarizer defines the port for AI-powered content summarization.
// Implementations live in the adapter layer (e.g., AI adapter using aifx).
type ContentSummarizer interface {
	SubmitSummarizeBatch(
		ctx context.Context,
		stories []*UnsummarizedStory,
	) (jobID string, err error)
}
