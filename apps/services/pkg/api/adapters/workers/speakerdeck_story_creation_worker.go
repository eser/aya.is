package workers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
)

// SpeakerDeckStoryProcessor creates stories from synced SpeakerDeck presentation imports.
type SpeakerDeckStoryProcessor struct {
	config      *SpeakerDeckSyncConfig
	logger      *logfx.Logger
	syncService *linksync.Service
	storyRepo   storyCreationRepo
	idGenerator func() string
}

// NewSpeakerDeckStoryProcessor creates a new SpeakerDeck story processor.
func NewSpeakerDeckStoryProcessor(
	config *SpeakerDeckSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	storyRepo storyCreationRepo,
	idGenerator func() string,
) *SpeakerDeckStoryProcessor {
	return &SpeakerDeckStoryProcessor{
		config:      config,
		logger:      logger,
		syncService: syncService,
		storyRepo:   storyRepo,
		idGenerator: idGenerator,
	}
}

// ProcessStories creates stories from new imports and reconciles existing stories.
func (w *SpeakerDeckStoryProcessor) ProcessStories(ctx context.Context) error {
	w.logger.DebugContext(ctx, "Starting SpeakerDeck story creation cycle")

	imports, err := w.syncService.ListImportsForStoryCreation(
		ctx,
		"speakerdeck",
		w.config.BatchSize,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(imports) == 0 {
		w.logger.DebugContext(ctx, "No SpeakerDeck imports need story creation")

		return nil
	}

	w.logger.DebugContext(ctx, "Processing SpeakerDeck imports for story creation",
		slog.Int("count", len(imports)))

	created := 0

	for _, imp := range imports {
		err := w.createStoryFromImport(ctx, imp)
		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to create story from SpeakerDeck import",
				slog.String("import_id", imp.ID),
				slog.String("remote_id", imp.RemoteID),
				slog.String("profile_id", imp.ProfileID),
				slog.Any("error", err))

			continue
		}

		created++
	}

	w.logger.DebugContext(ctx, "Completed SpeakerDeck story creation cycle",
		slog.Int("processed", len(imports)),
		slog.Int("created", created))

	return nil
}

// speakerDeckImportMeta holds extracted metadata from a SpeakerDeck import.
type speakerDeckImportMeta struct {
	title           string
	description     string
	thumbnailURL    string
	link            string
	pdfURL          string
	publishedAt     time.Time
	storyPictureURI *string
}

// extractImportMeta extracts metadata from an import's properties.
func extractImportMeta(imp *linksync.LinkImportForStoryCreation) *speakerDeckImportMeta {
	title, _ := imp.Properties["title"].(string)
	description, _ := imp.Properties["description"].(string)
	thumbnailURL, _ := imp.Properties["thumbnail_url"].(string)
	link, _ := imp.Properties["link"].(string)
	pdfURL, _ := imp.Properties["pdf_url"].(string)
	publishedAtStr, _ := imp.Properties["published_at"].(string)

	if title == "" {
		title = "Untitled Presentation"
	}

	publishedAt := time.Now()

	if publishedAtStr != "" {
		if parsed, err := time.Parse(time.RFC3339, publishedAtStr); err == nil {
			publishedAt = parsed
		}
	}

	var storyPictureURI *string
	if thumbnailURL != "" {
		storyPictureURI = &thumbnailURL
	}

	return &speakerDeckImportMeta{
		title:           title,
		description:     description,
		thumbnailURL:    thumbnailURL,
		link:            link,
		pdfURL:          pdfURL,
		publishedAt:     publishedAt,
		storyPictureURI: storyPictureURI,
	}
}

// createStoryFromImport creates a story from a SpeakerDeck import.
func (w *SpeakerDeckStoryProcessor) createStoryFromImport(
	ctx context.Context,
	imp *linksync.LinkImportForStoryCreation,
) error {
	meta := extractImportMeta(imp)

	locale := imp.ProfileDefaultLocale
	if locale == "" {
		locale = "en"
	}

	slug := generateSlugFromTitle(meta.publishedAt, meta.title)
	storyID := w.idGenerator()
	publicationID := w.idGenerator()

	properties := map[string]any{
		"managed_by": "speakerdeck_sync_worker",
		"remote_id":  imp.RemoteID,
	}

	_, err := w.storyRepo.InsertStory(
		ctx, storyID, imp.ProfileID, slug, "presentation",
		meta.storyPictureURI, properties, true, &imp.RemoteID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert story: %w", err)
	}

	content := buildSpeakerDeckStoryContent(meta.pdfURL, meta.link, meta.description)
	summary := truncateSummary(meta.description, maxSummaryLength)

	err = w.storyRepo.InsertStoryTx(ctx, storyID, locale, meta.title, summary, content)
	if err != nil {
		return fmt.Errorf("failed to insert story translation: %w", err)
	}

	err = w.storyRepo.InsertStoryPublication(
		ctx, publicationID, storyID, imp.ProfileID,
		"original", false, &meta.publishedAt, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to insert story publication: %w", err)
	}

	w.logger.DebugContext(ctx, "Created story from SpeakerDeck import",
		slog.String("story_id", storyID),
		slog.String("remote_id", imp.RemoteID),
		slog.String("slug", slug),
		slog.String("locale", locale))

	return nil
}

// buildSpeakerDeckStoryContent builds the story content from a SpeakerDeck presentation.
// Uses <PDF> MDX component when PDF URL is available, falls back to %[link] embed.
func buildSpeakerDeckStoryContent(
	pdfURL string,
	presentationLink string,
	description string,
) string {
	var content string

	if pdfURL != "" {
		content = `<PDF src="` + pdfURL + `" />`
	} else if presentationLink != "" {
		content = "%[" + presentationLink + "]"
	}

	if description != "" {
		if content != "" {
			content += "\n\n"
		}

		content += description
	}

	return content
}
