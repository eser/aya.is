package ai

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
)

const (
	batchPollInterval = 10 * time.Second
	batchPollTimeout  = 10 * time.Minute
	summaryMaxTokens  = 256
)

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

// BulletinSummarizer generates AI summaries for bulletin stories.
type BulletinSummarizer struct {
	model  aifx.LanguageModel
	logger *logfx.Logger
}

// NewBulletinSummarizer creates a new AI summarizer adapter.
func NewBulletinSummarizer(
	model aifx.LanguageModel,
	logger *logfx.Logger,
) *BulletinSummarizer {
	return &BulletinSummarizer{
		model:  model,
		logger: logger,
	}
}

// SummarizeBatch generates AI summaries for the given stories.
// Uses the Anthropic Message Batches API if available, otherwise falls back to sequential calls.
func (s *BulletinSummarizer) SummarizeBatch(
	ctx context.Context,
	stories []*bulletinbiz.DigestStory,
	localeCode string,
) (map[string]string, error) {
	if len(stories) == 0 {
		return make(map[string]string), nil
	}

	batchModel, isBatchCapable := s.model.(aifx.BatchCapableModel)
	if isBatchCapable {
		return s.summarizeViaBatch(ctx, batchModel, stories, localeCode)
	}

	return s.summarizeSequentially(ctx, stories, localeCode)
}

// summarizeViaBatch uses the batch API for cost-efficient bulk summarization.
func (s *BulletinSummarizer) summarizeViaBatch(
	ctx context.Context,
	model aifx.BatchCapableModel,
	stories []*bulletinbiz.DigestStory,
	localeCode string,
) (map[string]string, error) {
	language := languageNameForLocale(localeCode)

	items := make([]aifx.BatchRequestItem, 0, len(stories))

	for _, story := range stories {
		customID := story.StoryID + "|" + localeCode

		items = append(items, aifx.BatchRequestItem{
			CustomID: customID,
			Options: aifx.GenerateTextOptions{
				System: "You are a content summarizer. Respond with a single paragraph summary only. No headers, no bullet points.",
				Messages: []aifx.Message{
					aifx.NewTextMessage(aifx.RoleUser, buildSummaryPrompt(story, language)),
				},
				MaxTokens: summaryMaxTokens,
			},
		})
	}

	s.logger.InfoContext(ctx, "Submitting batch summarization request",
		slog.Int("story_count", len(items)),
		slog.String("locale", localeCode))

	job, err := model.SubmitBatch(ctx, &aifx.BatchRequest{Items: items})
	if err != nil {
		return nil, fmt.Errorf("%w: submitting batch: %w", bulletinbiz.ErrSummarizationFailed, err)
	}

	// Poll until complete or timeout
	deadline := time.Now().Add(batchPollTimeout)

	for time.Now().Before(deadline) {
		job, err = model.GetBatchJob(ctx, job.ID)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: polling batch job: %w",
				bulletinbiz.ErrSummarizationFailed,
				err,
			)
		}

		if job.Status == aifx.BatchStatusCompleted || job.Status == aifx.BatchStatusFailed {
			break
		}

		s.logger.DebugContext(ctx, "Batch job still processing",
			slog.String("job_id", job.ID),
			slog.String("status", string(job.Status)))

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(batchPollInterval):
		}
	}

	if job.Status == aifx.BatchStatusFailed {
		return nil, fmt.Errorf(
			"%w: batch job failed: %s",
			bulletinbiz.ErrSummarizationFailed,
			job.Error,
		)
	}

	if job.Status != aifx.BatchStatusCompleted {
		return nil, fmt.Errorf(
			"%w: batch job timed out (status: %s)",
			bulletinbiz.ErrSummarizationFailed,
			job.Status,
		)
	}

	// Download results
	results, err := model.DownloadBatchResults(ctx, job)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: downloading batch results: %w",
			bulletinbiz.ErrSummarizationFailed,
			err,
		)
	}

	summaries := make(map[string]string, len(results))

	for _, result := range results {
		if result.Error != "" {
			s.logger.WarnContext(ctx, "Batch result error",
				slog.String("custom_id", result.CustomID),
				slog.String("error", result.Error))

			continue
		}

		// Parse CustomID to extract storyID
		parts := strings.SplitN(result.CustomID, "|", 2)
		if len(parts) < 1 {
			continue
		}

		storyID := parts[0]

		if result.Result != nil && result.Result.Text() != "" {
			summaries[storyID] = strings.TrimSpace(result.Result.Text())
		}
	}

	s.logger.InfoContext(ctx, "Batch summarization completed",
		slog.Int("total", len(items)),
		slog.Int("successful", len(summaries)))

	return summaries, nil
}

// summarizeSequentially falls back to individual API calls when batch is not available.
func (s *BulletinSummarizer) summarizeSequentially(
	ctx context.Context,
	stories []*bulletinbiz.DigestStory,
	localeCode string,
) (map[string]string, error) {
	language := languageNameForLocale(localeCode)
	summaries := make(map[string]string, len(stories))

	for _, story := range stories {
		result, err := s.model.GenerateText(ctx, &aifx.GenerateTextOptions{
			System: "You are a content summarizer. Respond with a single paragraph summary only. No headers, no bullet points.",
			Messages: []aifx.Message{
				aifx.NewTextMessage(aifx.RoleUser, buildSummaryPrompt(story, language)),
			},
			MaxTokens: summaryMaxTokens,
		})
		if err != nil {
			s.logger.WarnContext(ctx, "Sequential summarization failed for story",
				slog.String("story_id", story.StoryID),
				slog.String("error", err.Error()))

			continue
		}

		if result != nil && result.Text() != "" {
			summaries[story.StoryID] = strings.TrimSpace(result.Text())
		}
	}

	return summaries, nil
}

func buildSummaryPrompt(story *bulletinbiz.DigestStory, language string) string {
	return fmt.Sprintf(
		"Summarize this article in one paragraph in %s. Be concise and informative.\n\nTitle: %s\n\nSummary:\n%s",
		language,
		story.Title,
		story.Summary,
	)
}

func languageNameForLocale(locale string) string {
	if name, ok := localeToLanguageName[locale]; ok {
		return name
	}

	return locale
}
