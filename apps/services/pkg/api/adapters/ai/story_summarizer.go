package ai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
)

var errInvalidCustomID = errors.New("invalid custom_id format")

const summaryMaxTokens = 256

//
//nolint:lll
const summarySystemPrompt = "You are a content summarizer for a publishing platform. Generate a single concise paragraph (2-4 sentences) summarizing the article. Write in the specified language. No headers, bullet points, or markdown formatting."

// localeToLanguageName maps locale codes to full language names for AI prompts.
// Using full names prevents LLMs from defaulting to English when given short codes.
var localeToLanguageName = map[string]string{ //nolint:gochecknoglobals
	"ar":    "Arabic",
	"de":    "German",
	"en":    "English",
	"es":    "Spanish",
	"fr":    "French",
	"it":    "Italian",
	"ja":    "Japanese",
	"ko":    "Korean",
	"nl":    "Dutch",
	"pt-PT": "Portuguese (Portugal)",
	"ru":    "Russian",
	"tr":    "Turkish",
	"zh-CN": "Chinese (Simplified)",
}

// StorySummarizer generates AI summaries for stories via the batch API.
type StorySummarizer struct {
	model  aifx.BatchCapableModel
	logger *logfx.Logger
}

// NewStorySummarizer creates a new AI summarizer adapter.
func NewStorySummarizer(
	model aifx.BatchCapableModel,
	logger *logfx.Logger,
) *StorySummarizer {
	return &StorySummarizer{
		model:  model,
		logger: logger,
	}
}

// SubmitSummarizeBatch submits a batch of stories for AI summarization.
// Returns the batch job ID for tracking. Does NOT poll — the caller is
// responsible for checking completion via GetBatchJob.
func (s *StorySummarizer) SubmitSummarizeBatch(
	ctx context.Context,
	storiesToSummarize []*stories.UnsummarizedStory,
) (string, error) {
	if len(storiesToSummarize) == 0 {
		return "", nil
	}

	items := make([]aifx.BatchRequestItem, 0, len(storiesToSummarize))

	for _, story := range storiesToSummarize {
		customID := story.StoryID + "_" + story.LocaleCode
		language := languageNameForLocale(story.LocaleCode)

		items = append(items, aifx.BatchRequestItem{
			CustomID: customID,
			Options: aifx.GenerateTextOptions{ //nolint:exhaustruct
				System: summarySystemPrompt,
				Messages: []aifx.Message{
					aifx.NewTextMessage(aifx.RoleUser, buildSummaryPrompt(story, language)),
				},
				MaxTokens: summaryMaxTokens,
			},
		})
	}

	s.logger.InfoContext(ctx, "Submitting batch summarization request",
		slog.Int("story_count", len(items)))

	job, err := s.model.SubmitBatch(ctx, &aifx.BatchRequest{Items: items})
	if err != nil {
		return "", fmt.Errorf("submitting batch: %w", err)
	}

	s.logger.InfoContext(ctx, "Batch summarization submitted",
		slog.String("job_id", job.ID),
		slog.Int("item_count", len(items)))

	return job.ID, nil
}

// GetBatchJob retrieves the current status of a batch job.
func (s *StorySummarizer) GetBatchJob(
	ctx context.Context,
	jobID string,
) (*aifx.BatchJob, error) {
	return s.model.GetBatchJob(ctx, jobID) //nolint:wrapcheck
}

// DownloadBatchResults downloads the results of a completed batch job.
func (s *StorySummarizer) DownloadBatchResults(
	ctx context.Context,
	job *aifx.BatchJob,
) ([]*aifx.BatchResult, error) {
	return s.model.DownloadBatchResults(ctx, job) //nolint:wrapcheck
}

// ParseCustomID extracts storyID and localeCode from a batch result's CustomID.
// Format: "{storyID}_{localeCode}".
func (s *StorySummarizer) ParseCustomID(customID string) (string, string, error) {
	parts := strings.SplitN(customID, "_", 2) //nolint:mnd
	if len(parts) != 2 {                      //nolint:mnd
		return "", "", fmt.Errorf("%w: %q", errInvalidCustomID, customID)
	}

	return parts[0], parts[1], nil
}

func buildSummaryPrompt(story *stories.UnsummarizedStory, language string) string {
	return fmt.Sprintf(
		"Summarize this article in one paragraph in %s. Be concise and informative.\n\nTitle: %s\n\nContent:\n%s",
		language,
		story.Title,
		story.Content,
	)
}

func languageNameForLocale(locale string) string {
	if name, ok := localeToLanguageName[locale]; ok {
		return name
	}

	return locale
}
