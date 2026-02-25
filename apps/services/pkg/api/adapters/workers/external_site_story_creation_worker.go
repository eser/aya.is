package workers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
)

// ExternalSiteStoryProcessor creates stories from synced external site imports.
type ExternalSiteStoryProcessor struct {
	config      *ExternalSiteSyncConfig
	logger      *logfx.Logger
	syncService *linksync.Service
	storyRepo   storyCreationRepo
	idGenerator func() string
}

// NewExternalSiteStoryProcessor creates a new external site story processor.
func NewExternalSiteStoryProcessor(
	config *ExternalSiteSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	storyRepo storyCreationRepo,
	idGenerator func() string,
) *ExternalSiteStoryProcessor {
	return &ExternalSiteStoryProcessor{
		config:      config,
		logger:      logger,
		syncService: syncService,
		storyRepo:   storyRepo,
		idGenerator: idGenerator,
	}
}

// ProcessStories creates stories from new imports and reconciles existing stories.
// It drains all pending imports in batches rather than processing a single batch,
// because external site imports are cheap DB operations (unlike API-based providers).
func (w *ExternalSiteStoryProcessor) ProcessStories(ctx context.Context) error {
	w.logger.DebugContext(ctx, "Starting external site story creation cycle")

	totalCreated := 0

	for {
		imports, err := w.syncService.ListImportsForStoryCreation(
			ctx,
			"external-site",
			w.config.BatchSize,
		)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrSyncFailed, err)
		}

		if len(imports) == 0 {
			break
		}

		w.logger.DebugContext(ctx, "Processing external site imports for story creation",
			slog.Int("batch_count", len(imports)))

		for _, imp := range imports {
			createErr := w.createStoryFromImport(ctx, imp)
			if createErr != nil {
				w.logger.ErrorContext(ctx, "Failed to create story from external site import",
					slog.String("import_id", imp.ID),
					slog.String("remote_id", imp.RemoteID),
					slog.String("profile_id", imp.ProfileID),
					slog.Any("error", createErr))

				continue
			}

			totalCreated++
		}

		// If we got fewer than BatchSize, all pending imports have been processed
		if len(imports) < w.config.BatchSize {
			break
		}
	}

	if totalCreated == 0 {
		w.logger.DebugContext(ctx, "No external site imports need story creation")
	} else {
		w.logger.DebugContext(ctx, "Completed external site story creation cycle",
			slog.Int("created", totalCreated))
	}

	// Always reconcile existing stories with latest import data
	err := w.reconcileExistingStories(ctx)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to reconcile existing external site stories",
			slog.Any("error", err))
	}

	return nil
}

// externalSiteImportMeta holds extracted metadata from an external site import.
type externalSiteImportMeta struct {
	title       string
	description string
	content     string
	slug        string
	language    string
	link        string
	siteURL     string
	sourcePath  string
	storyKind   string
	publishedAt time.Time
}

// extractExternalSiteImportMeta extracts metadata from an external site import's properties.
func extractExternalSiteImportMeta(
	imp *linksync.LinkImportForStoryCreation,
) *externalSiteImportMeta {
	title, _ := imp.Properties["title"].(string)
	description, _ := imp.Properties["description"].(string)
	content, _ := imp.Properties["content"].(string)
	slug, _ := imp.Properties["slug"].(string)
	language, _ := imp.Properties["language"].(string)
	link, _ := imp.Properties["link"].(string)
	siteURL, _ := imp.Properties["site_url"].(string)
	sourcePath, _ := imp.Properties["source_path"].(string)
	storyKind, _ := imp.Properties["story_kind"].(string)
	publishedAtStr, _ := imp.Properties["published_at"].(string)

	if title == "" {
		title = "Untitled Post"
	}

	if storyKind == "" {
		storyKind = "article"
	}

	publishedAt := time.Now()

	if publishedAtStr != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, publishedAtStr); parseErr == nil &&
			!parsed.IsZero() {
			publishedAt = parsed
		}
	}

	return &externalSiteImportMeta{
		title:       title,
		description: description,
		content:     content,
		slug:        slug,
		language:    language,
		link:        link,
		siteURL:     siteURL,
		sourcePath:  sourcePath,
		storyKind:   storyKind,
		publishedAt: publishedAt,
	}
}

// buildPublishedURL constructs the published URL for an external site post.
// For Hugo/Jekyll/Zola sites, the pattern is typically {site_url}/{section}/{slug}/.
// The section is extracted from the source path (e.g., "content/posts/.../index.md" → "posts").
func buildPublishedURL(siteURL, sourcePath, slug string) string {
	if siteURL == "" || slug == "" {
		return ""
	}

	siteURL = strings.TrimRight(siteURL, "/")

	// Extract section from source path (e.g., "content/posts/..." → "posts")
	section := ""

	parts := strings.Split(sourcePath, "/")
	if len(parts) >= 2 { //nolint:mnd
		// Skip "content" or "_posts" prefix, take the next segment as section
		start := 0

		if parts[0] == "content" || parts[0] == "_posts" {
			start = 1
		}

		if start < len(parts)-1 {
			section = parts[start]
		}
	}

	if section != "" {
		return siteURL + "/" + section + "/" + slug + "/"
	}

	return siteURL + "/" + slug + "/"
}

// reconcileExistingStories updates existing stories to match the latest import data.
func (w *ExternalSiteStoryProcessor) reconcileExistingStories(ctx context.Context) error {
	imports, err := w.syncService.ListImportsWithExistingStories(
		ctx,
		"external-site",
		w.config.BatchSize,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(imports) == 0 {
		return nil
	}

	w.logger.DebugContext(ctx, "Reconciling existing external site stories",
		slog.Int("count", len(imports)))

	reconciled := 0

	for _, imp := range imports {
		reconcileErr := w.reconcileStory(ctx, imp)
		if reconcileErr != nil {
			w.logger.ErrorContext(ctx, "Failed to reconcile external site story",
				slog.String("story_id", imp.StoryID),
				slog.String("remote_id", imp.RemoteID),
				slog.Any("error", reconcileErr))

			continue
		}

		reconciled++
	}

	w.logger.DebugContext(ctx, "Completed external site story reconciliation",
		slog.Int("processed", len(imports)),
		slog.Int("reconciled", reconciled))

	return nil
}

// reconcileStory updates a single story's content to match the latest import data.
func (w *ExternalSiteStoryProcessor) reconcileStory(
	ctx context.Context,
	imp *linksync.LinkImportWithStory,
) error {
	title, _ := imp.Properties["title"].(string)
	description, _ := imp.Properties["description"].(string)
	content, _ := imp.Properties["content"].(string)
	language, _ := imp.Properties["language"].(string)

	if title == "" {
		title = "Untitled Post"
	}

	locale := language
	if locale == "" {
		locale = imp.ProfileDefaultLocale
	}

	if locale == "" {
		locale = "en"
	}

	summary := truncateSummary(description, maxSummaryLength)

	err := w.storyRepo.UpsertStoryTx(ctx, imp.StoryID, locale, title, summary, content, true)
	if err != nil {
		return fmt.Errorf("failed to upsert story translation: %w", err)
	}

	return nil
}

// createStoryFromImport creates a story from an external site import.
func (w *ExternalSiteStoryProcessor) createStoryFromImport(
	ctx context.Context,
	imp *linksync.LinkImportForStoryCreation,
) error {
	meta := extractExternalSiteImportMeta(imp)

	locale := meta.language
	if locale == "" {
		locale = imp.ProfileDefaultLocale
	}

	if locale == "" {
		locale = "en"
	}

	// Use the frontmatter slug if available, otherwise generate from title
	slug := meta.slug
	if slug == "" {
		slug = generateSlugFromTitle(meta.publishedAt, meta.title)
	} else {
		// Prefix with date for uniqueness
		datePrefix := meta.publishedAt.Format("20060102") + "-"
		slug = datePrefix + slug
	}

	storyID := w.idGenerator()
	publicationID := w.idGenerator()

	sourceURL := buildPublishedURL(meta.siteURL, meta.sourcePath, meta.slug)
	if sourceURL == "" {
		sourceURL = meta.link // fallback to GitHub blob URL
	}

	properties := map[string]any{
		"managed_by": "external_site_sync_worker",
		"remote_id":  imp.RemoteID,
		"source_url": sourceURL,
	}

	_, err := w.storyRepo.InsertStory(
		ctx, storyID, imp.ProfileID, slug, meta.storyKind,
		nil, properties, true, &imp.RemoteID,
		"public", false,
	)
	if err != nil {
		return fmt.Errorf("failed to insert story: %w", err)
	}

	summary := truncateSummary(meta.description, maxSummaryLength)

	err = w.storyRepo.InsertStoryTx(ctx, storyID, locale, meta.title, summary, meta.content, true)
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

	w.logger.DebugContext(ctx, "Created story from external site import",
		slog.String("story_id", storyID),
		slog.String("remote_id", imp.RemoteID),
		slog.String("slug", slug),
		slog.String("locale", locale),
		slog.String("kind", meta.storyKind))

	return nil
}
