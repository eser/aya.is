package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/ai"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
)

const (
	lockIDStorySummaries      int64  = 100013
	storySummariesBatchPrefix string = "story-summaries.batch."
)

// trackedBatchJob is the JSON structure stored in runtime_state for each pending batch.
type trackedBatchJob struct {
	JobID       string    `json:"job_id"`
	SubmittedAt time.Time `json:"submitted_at"`
	ItemCount   int       `json:"item_count"`
}

// StorySummariesWorker periodically submits and processes AI summarization batches.
type StorySummariesWorker struct {
	config        *StorySummariesConfig
	logger        *logfx.Logger
	storyService  *stories.Service
	summarizer    *ai.StorySummarizer
	runtimeStates *runtime_states.Service
}

// NewStorySummariesWorker creates a new story summaries worker.
func NewStorySummariesWorker(
	config *StorySummariesConfig,
	logger *logfx.Logger,
	storyService *stories.Service,
	summarizer *ai.StorySummarizer,
	runtimeStates *runtime_states.Service,
) *StorySummariesWorker {
	return &StorySummariesWorker{
		config:        config,
		logger:        logger,
		storyService:  storyService,
		summarizer:    summarizer,
		runtimeStates: runtimeStates,
	}
}

// Name returns the worker name.
func (w *StorySummariesWorker) Name() string {
	return "story-summaries"
}

// Interval returns the check interval.
func (w *StorySummariesWorker) Interval() time.Duration {
	return w.config.CheckInterval
}

// Execute runs one cycle: process completed batches, then submit new ones.
func (w *StorySummariesWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return workerfx.ErrWorkerSkipped
	}

	// Try advisory lock to prevent concurrent execution
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDStorySummaries)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock for story-summaries",
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running story-summaries worker")

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDStorySummaries)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "Failed to release advisory lock for story-summaries",
				slog.String("error", releaseErr.Error()))
		}
	}()

	// Phase 1: Process completed batches
	hasPending, processErr := w.processCompletedBatches(ctx)
	if processErr != nil {
		w.logger.WarnContext(ctx, "Error processing completed batches",
			slog.String("error", processErr.Error()))
	}

	// Phase 2: Submit new batch (only if no pending batches)
	if !hasPending {
		submitErr := w.submitNewBatch(ctx)
		if submitErr != nil {
			return fmt.Errorf("submitting new batch: %w", submitErr)
		}
	}

	return nil
}

// processCompletedBatches checks all tracked batch jobs, persists results for completed ones,
// and cleans up finished/stale entries. Returns true if any batches are still pending.
func (w *StorySummariesWorker) processCompletedBatches( //nolint:funlen,cyclop
	ctx context.Context,
) (bool, error) {
	tracked, err := w.runtimeStates.ListByPrefix(ctx, storySummariesBatchPrefix)
	if err != nil {
		return false, fmt.Errorf("listing tracked batches: %w", err)
	}

	if len(tracked) == 0 {
		return false, nil
	}

	hasPending := false

	for _, state := range tracked {
		var batch trackedBatchJob

		unmarshalErr := json.Unmarshal([]byte(state.Value), &batch)
		if unmarshalErr != nil {
			w.logger.WarnContext(ctx, "Invalid batch tracking entry, removing",
				slog.String("key", state.Key),
				slog.String("error", unmarshalErr.Error()))

			_ = w.runtimeStates.Remove(ctx, state.Key)

			continue
		}

		// Check if batch is stale (exceeded max retry age)
		if time.Since(batch.SubmittedAt) > w.config.MaxRetryAge {
			w.logger.WarnContext(ctx, "Batch job exceeded max retry age, removing",
				slog.String("job_id", batch.JobID),
				slog.Duration("age", time.Since(batch.SubmittedAt)))

			_ = w.runtimeStates.Remove(ctx, state.Key)

			continue
		}

		// Poll batch status
		job, jobErr := w.summarizer.GetBatchJob(ctx, batch.JobID)
		if jobErr != nil {
			// Fatal configuration errors — stop retrying this batch
			if errors.Is(jobErr, aifx.ErrAuthFailed) ||
				errors.Is(jobErr, aifx.ErrInsufficientCredits) {
				w.logger.ErrorContext(ctx, "AI configuration error, removing batch tracking",
					slog.String("job_id", batch.JobID),
					slog.String("error", jobErr.Error()))

				_ = w.runtimeStates.Remove(ctx, state.Key)

				continue
			}

			w.logger.WarnContext(ctx, "Failed to get batch job status",
				slog.String("job_id", batch.JobID),
				slog.String("error", jobErr.Error()))

			hasPending = true

			continue
		}

		switch job.Status {
		case aifx.BatchStatusCompleted:
			w.processCompletedJob(ctx, job, state.Key)
		case aifx.BatchStatusFailed, aifx.BatchStatusCancelled:
			w.logger.WarnContext(ctx, "Batch job failed or cancelled",
				slog.String("job_id", batch.JobID),
				slog.String("status", string(job.Status)),
				slog.String("error", job.Error))

			_ = w.runtimeStates.Remove(ctx, state.Key)
		case aifx.BatchStatusPending, aifx.BatchStatusProcessing:
			w.logger.DebugContext(ctx, "Batch job still processing",
				slog.String("job_id", batch.JobID),
				slog.String("status", string(job.Status)),
				slog.Int("done", job.DoneCount),
				slog.Int("total", job.TotalCount))

			hasPending = true
		}
	}

	return hasPending, nil
}

// processCompletedJob downloads results from a completed batch and persists AI summaries.
func (w *StorySummariesWorker) processCompletedJob(
	ctx context.Context,
	job *aifx.BatchJob,
	stateKey string,
) {
	results, err := w.summarizer.DownloadBatchResults(ctx, job)
	if err != nil {
		w.logger.WarnContext(ctx, "Failed to download batch results",
			slog.String("job_id", job.ID),
			slog.String("error", err.Error()))

		return
	}

	successCount := 0

	for _, result := range results {
		if result.Error != "" {
			w.logger.WarnContext(ctx, "Batch result error",
				slog.String("custom_id", result.CustomID),
				slog.String("error", result.Error))

			continue
		}

		storyID, localeCode, parseErr := w.summarizer.ParseCustomID(result.CustomID)
		if parseErr != nil {
			w.logger.WarnContext(ctx, "Failed to parse custom_id",
				slog.String("custom_id", result.CustomID),
				slog.String("error", parseErr.Error()))

			continue
		}

		if result.Result == nil || result.Result.Text() == "" {
			continue
		}

		summary := strings.TrimSpace(result.Result.Text())

		persistErr := w.storyService.PersistSummaryAI(ctx, storyID, localeCode, summary)
		if persistErr != nil {
			w.logger.WarnContext(ctx, "Failed to persist AI summary",
				slog.String("story_id", storyID),
				slog.String("locale", localeCode),
				slog.String("error", persistErr.Error()))

			continue
		}

		successCount++
	}

	w.logger.InfoContext(ctx, "Processed batch summarization results",
		slog.String("job_id", job.ID),
		slog.Int("total_results", len(results)),
		slog.Int("persisted", successCount))

	// Clean up tracking entry
	_ = w.runtimeStates.Remove(ctx, stateKey)
}

// submitNewBatch fetches unsummarized stories and submits a new batch request.
func (w *StorySummariesWorker) submitNewBatch(ctx context.Context) error {
	unsummarized, err := w.storyService.GetUnsummarizedStories(ctx, w.config.MaxBatchSize)
	if err != nil {
		return fmt.Errorf("fetching unsummarized stories: %w", err)
	}

	if len(unsummarized) == 0 {
		w.logger.DebugContext(ctx, "No unsummarized stories found")

		return nil
	}

	w.logger.InfoContext(ctx, "Found unsummarized stories",
		slog.Int("count", len(unsummarized)))

	jobID, submitErr := w.summarizer.SubmitSummarizeBatch(ctx, unsummarized)
	if submitErr != nil {
		// Rate limited — silently skip this cycle, retry next interval
		if errors.Is(submitErr, aifx.ErrRateLimited) {
			w.logger.WarnContext(ctx, "AI rate limited, will retry next cycle")

			return nil
		}

		return fmt.Errorf("submitting summarize batch: %w", submitErr)
	}

	if jobID == "" {
		return nil
	}

	// Track the batch job in runtime_state
	tracking := trackedBatchJob{
		JobID:       jobID,
		SubmittedAt: time.Now().UTC(),
		ItemCount:   len(unsummarized),
	}

	trackingJSON, marshalErr := json.Marshal(tracking)
	if marshalErr != nil {
		return fmt.Errorf("marshaling batch tracking: %w", marshalErr)
	}

	stateKey := storySummariesBatchPrefix + jobID

	setErr := w.runtimeStates.Set(ctx, stateKey, string(trackingJSON))
	if setErr != nil {
		return fmt.Errorf("storing batch tracking: %w", setErr)
	}

	w.logger.InfoContext(ctx, "Submitted and tracked new summarization batch",
		slog.String("job_id", jobID),
		slog.Int("item_count", len(unsummarized)))

	return nil
}
