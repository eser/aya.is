package workers

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
)

const (
	maxSlugLength     = 80
	maxSummaryLength  = 500
	slugDatePrefixLen = 9 // "YYYYMMDD-"
)

// storyCreationRepo defines the minimal repository interface for creating and updating stories.
type storyCreationRepo interface {
	InsertStory(
		ctx context.Context,
		id string,
		authorProfileID string,
		slug string,
		kind string,
		storyPictureURI *string,
		properties map[string]any,
		isManaged bool,
	) (*stories.Story, error)
	InsertStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	InsertStoryPublication(
		ctx context.Context,
		id string,
		storyID string,
		profileID string,
		kind string,
		isFeatured bool,
		publishedAt *time.Time,
		properties map[string]any,
	) error
	UpdateStory(
		ctx context.Context,
		id string,
		slug string,
		storyPictureURI *string,
	) error
	UpsertStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	UpdateStoryPublicationDate(ctx context.Context, id string, publishedAt time.Time) error
	RemoveStoryPublication(ctx context.Context, id string) error
}

// YouTubeStoryProcessor creates and reconciles stories from synced YouTube video imports.
// Called by sync workers after completing a sync cycle.
type YouTubeStoryProcessor struct {
	config      *YouTubeSyncConfig
	logger      *logfx.Logger
	syncService *linksync.Service
	storyRepo   storyCreationRepo
	idGenerator func() string
}

// NewYouTubeStoryProcessor creates a new YouTube story processor.
func NewYouTubeStoryProcessor(
	config *YouTubeSyncConfig,
	logger *logfx.Logger,
	syncService *linksync.Service,
	storyRepo storyCreationRepo,
	idGenerator func() string,
) *YouTubeStoryProcessor {
	return &YouTubeStoryProcessor{
		config:      config,
		logger:      logger,
		syncService: syncService,
		storyRepo:   storyRepo,
		idGenerator: idGenerator,
	}
}

// ProcessStories creates stories from new imports and reconciles existing stories.
func (w *YouTubeStoryProcessor) ProcessStories(ctx context.Context) error {
	w.logger.DebugContext(ctx, "Starting YouTube story creation cycle")

	imports, err := w.syncService.ListImportsForStoryCreation(ctx, "youtube", w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(imports) == 0 {
		w.logger.DebugContext(ctx, "No YouTube imports need story creation")
	} else {
		w.logger.DebugContext(ctx, "Processing YouTube imports for story creation",
			slog.Int("count", len(imports)))

		created := 0

		for _, imp := range imports {
			err := w.createStoryFromImport(ctx, imp)
			if err != nil {
				w.logger.ErrorContext(ctx, "Failed to create story from import",
					slog.String("import_id", imp.ID),
					slog.String("remote_id", imp.RemoteID),
					slog.String("profile_id", imp.ProfileID),
					slog.Any("error", err))

				continue
			}

			created++
		}

		w.logger.DebugContext(ctx, "Completed YouTube story creation cycle",
			slog.Int("processed", len(imports)),
			slog.Int("created", created))
	}

	// Always reconcile existing stories with latest import data
	if err := w.reconcileExistingStories(ctx); err != nil {
		w.logger.ErrorContext(ctx, "Failed to reconcile existing stories",
			slog.Any("error", err))
	}

	return nil
}

// createStoryFromImport creates a story + story_tx + story_publication from a link import.
func (w *YouTubeStoryProcessor) createStoryFromImport( //nolint:cyclop,funlen
	ctx context.Context,
	imp *linksync.LinkImportForStoryCreation,
) error {
	// Extract video metadata
	videoMeta := extractVideoMetadata(imp.Properties)

	// Skip unlisted or private videos
	if isVideoNonPublic(w.logger, imp.RemoteID, imp.Properties) {
		w.logger.DebugContext(ctx, "Skipping non-public video",
			slog.String("remote_id", imp.RemoteID))

		return nil
	}

	// Detect locale
	locale := detectLocaleFromYouTubeVideo(videoMeta, imp.ProfileDefaultLocale)

	// Generate slug
	publishedAt := extractPublishedAt(imp.Properties)
	slug := generateSlugFromTitle(publishedAt, videoMeta.title)

	// Generate IDs
	storyID := w.idGenerator()
	publicationID := w.idGenerator()

	// Story picture from thumbnail
	thumbnailURI := extractThumbnailURI(videoMeta)

	var storyPictureURI *string
	if thumbnailURI != "" {
		storyPictureURI = &thumbnailURI
	}

	// Story properties
	properties := map[string]any{
		"managed_by": "youtube_sync_worker",
		"remote_id":  imp.RemoteID,
	}

	// Create story
	_, err := w.storyRepo.InsertStory(
		ctx,
		storyID,
		imp.ProfileID,
		slug,
		"content",
		storyPictureURI,
		properties,
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to insert story: %w", err)
	}

	// Build content: YouTube embed + description
	content := buildStoryContent(imp.RemoteID, videoMeta.description)

	// Truncate description for summary
	summary := truncateSummary(videoMeta.description, maxSummaryLength)

	// Create story translation
	err = w.storyRepo.InsertStoryTx(
		ctx,
		storyID,
		locale,
		videoMeta.title,
		summary,
		content,
	)
	if err != nil {
		return fmt.Errorf("failed to insert story translation: %w", err)
	}

	// Create story publication to author's profile
	err = w.storyRepo.InsertStoryPublication(
		ctx,
		publicationID,
		storyID,
		imp.ProfileID,
		"original",
		false,
		&publishedAt,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to insert story publication: %w", err)
	}

	w.logger.DebugContext(ctx, "Created story from YouTube import",
		slog.String("story_id", storyID),
		slog.String("remote_id", imp.RemoteID),
		slog.String("slug", slug),
		slog.String("locale", locale))

	return nil
}

// reconcileExistingStories updates existing stories to match the latest YouTube data.
func (w *YouTubeStoryProcessor) reconcileExistingStories(ctx context.Context) error {
	imports, err := w.syncService.ListImportsWithExistingStories(
		ctx,
		"youtube",
		w.config.FullSyncMaxStories,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSyncFailed, err)
	}

	if len(imports) == 0 {
		return nil
	}

	w.logger.WarnContext(ctx, "Reconciling existing YouTube stories",
		slog.Int("count", len(imports)))

	reconciled := 0
	nonPublicCount := 0

	for _, imp := range imports {
		if isVideoNonPublic(w.logger, imp.RemoteID, imp.Properties) {
			nonPublicCount++
		}

		err := w.reconcileStory(ctx, imp)
		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to reconcile story",
				slog.String("story_id", imp.StoryID),
				slog.String("remote_id", imp.RemoteID),
				slog.Any("error", err))

			continue
		}

		reconciled++
	}

	w.logger.WarnContext(ctx, "Completed YouTube story reconciliation",
		slog.Int("processed", len(imports)),
		slog.Int("reconciled", reconciled),
		slog.Int("non_public_detected", nonPublicCount))

	return nil
}

// reconcileStory updates a single story to match the latest YouTube data.
func (w *YouTubeStoryProcessor) reconcileStory(
	ctx context.Context,
	imp *linksync.LinkImportWithStory,
) error {
	videoMeta := extractVideoMetadata(imp.Properties)
	nonPublic := isVideoNonPublic(w.logger, imp.RemoteID, imp.Properties)

	// Handle publication based on privacy status
	if nonPublic && imp.PublicationID != nil {
		// Video became unlisted/private → remove author publication
		err := w.storyRepo.RemoveStoryPublication(ctx, *imp.PublicationID)
		if err != nil {
			return fmt.Errorf("failed to remove publication: %w", err)
		}

		w.logger.DebugContext(ctx, "Removed publication for non-public video",
			slog.String("story_id", imp.StoryID),
			slog.String("remote_id", imp.RemoteID))
	} else if !nonPublic && imp.PublicationID == nil {
		// Video is public but has no publication → add one
		publishedAt := extractPublishedAt(imp.Properties)
		publicationID := w.idGenerator()

		err := w.storyRepo.InsertStoryPublication(
			ctx,
			publicationID,
			imp.StoryID,
			imp.ProfileID,
			"original",
			false,
			&publishedAt,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to insert publication: %w", err)
		}

		w.logger.DebugContext(ctx, "Added publication for public video",
			slog.String("story_id", imp.StoryID),
			slog.String("remote_id", imp.RemoteID))
	}

	// Update publication date if publication exists and video is public
	if !nonPublic && imp.PublicationID != nil {
		publishedAt := extractPublishedAt(imp.Properties)

		err := w.storyRepo.UpdateStoryPublicationDate(ctx, *imp.PublicationID, publishedAt)
		if err != nil {
			w.logger.WarnContext(ctx, "Failed to update publication date",
				slog.String("story_id", imp.StoryID),
				slog.Any("error", err))
		}
	}

	// Update story fields to match YouTube data
	locale := detectLocaleFromYouTubeVideo(videoMeta, imp.ProfileDefaultLocale)
	publishedAt := extractPublishedAt(imp.Properties)
	slug := generateSlugFromTitle(publishedAt, videoMeta.title)

	thumbnailURI := extractThumbnailURI(videoMeta)

	var storyPictureURI *string
	if thumbnailURI != "" {
		storyPictureURI = &thumbnailURI
	}

	// Update story (slug + picture)
	err := w.storyRepo.UpdateStory(ctx, imp.StoryID, slug, storyPictureURI)
	if err != nil {
		return fmt.Errorf("failed to update story: %w", err)
	}

	// Update story translation (upsert handles locale changes)
	content := buildStoryContent(imp.RemoteID, videoMeta.description)
	summary := truncateSummary(videoMeta.description, maxSummaryLength)

	err = w.storyRepo.UpsertStoryTx(ctx, imp.StoryID, locale, videoMeta.title, summary, content)
	if err != nil {
		return fmt.Errorf("failed to upsert story translation: %w", err)
	}

	return nil
}

// isVideoNonPublic checks if a video's privacy status is unlisted or private.
// It checks both videoMetadata (from Videos API) and playlistItemMetadata (from Playlist Items API).
func isVideoNonPublic(logger *logfx.Logger, remoteID string, properties map[string]any) bool {
	// Check videoMetadata first (more authoritative)
	if videoMeta, ok := properties["videoMetadata"].(map[string]any); ok {
		if status, ok := videoMeta["status"].(map[string]any); ok {
			privacyStatus, _ := status["privacyStatus"].(string)

			logger.Warn("Privacy check: videoMetadata.status",
				slog.String("remote_id", remoteID),
				slog.String("privacyStatus", privacyStatus))

			if privacyStatus == "unlisted" || privacyStatus == "private" {
				return true
			}

			// If we have videoMetadata with a status, trust it
			if privacyStatus != "" {
				return false
			}
		} else {
			logger.Warn("Privacy check: videoMetadata exists but no 'status' field",
				slog.String("remote_id", remoteID))
		}
	} else {
		logger.Warn("Privacy check: no videoMetadata in properties",
			slog.String("remote_id", remoteID))
	}

	// Fallback: check playlistItemMetadata (available even when video metadata fetch fails)
	if playlistMeta, ok := properties["playlistItemMetadata"].(map[string]any); ok {
		if status, ok := playlistMeta["status"].(map[string]any); ok {
			privacyStatus, _ := status["privacyStatus"].(string)

			logger.Warn("Privacy check: playlistItemMetadata.status",
				slog.String("remote_id", remoteID),
				slog.String("privacyStatus", privacyStatus))

			if privacyStatus == "unlisted" || privacyStatus == "private" {
				return true
			}
		} else {
			logger.Warn("Privacy check: playlistItemMetadata exists but no 'status' field",
				slog.String("remote_id", remoteID))
		}
	} else {
		logger.Warn("Privacy check: no playlistItemMetadata in properties",
			slog.String("remote_id", remoteID))
	}

	logger.Warn("Privacy check: treating as public (no privacy status found)",
		slog.String("remote_id", remoteID))

	return false
}

// videoMetadata holds extracted metadata from YouTube video properties.
type videoMetadata struct {
	title                string
	description          string
	defaultLanguage      string
	defaultAudioLanguage string
	thumbnails           map[string]any
}

// extractVideoMetadata extracts video metadata from import properties.
func extractVideoMetadata(properties map[string]any) videoMetadata {
	meta := videoMetadata{} //nolint:exhaustruct

	videoMeta, ok := properties["videoMetadata"].(map[string]any)
	if !ok {
		// Try playlist item metadata as fallback
		playlistMeta, ok := properties["playlistItemMetadata"].(map[string]any)
		if !ok {
			return meta
		}

		snippet, _ := playlistMeta["snippet"].(map[string]any)
		if snippet != nil {
			meta.title, _ = snippet["title"].(string)
			meta.description, _ = snippet["description"].(string)
			meta.thumbnails, _ = snippet["thumbnails"].(map[string]any)
		}

		return meta
	}

	snippet, _ := videoMeta["snippet"].(map[string]any)
	if snippet == nil {
		return meta
	}

	meta.title, _ = snippet["title"].(string)
	meta.description, _ = snippet["description"].(string)
	meta.defaultLanguage, _ = snippet["defaultLanguage"].(string)
	meta.defaultAudioLanguage, _ = snippet["defaultAudioLanguage"].(string)
	meta.thumbnails, _ = snippet["thumbnails"].(map[string]any)

	return meta
}

// extractPublishedAt extracts the published timestamp from import properties.
func extractPublishedAt(properties map[string]any) time.Time {
	// Try videoMetadata first
	if videoMeta, ok := properties["videoMetadata"].(map[string]any); ok {
		if snippet, ok := videoMeta["snippet"].(map[string]any); ok {
			if publishedStr, ok := snippet["publishedAt"].(string); ok {
				if t, err := time.Parse(time.RFC3339, publishedStr); err == nil {
					return t
				}
			}
		}
	}

	// Try playlistItemMetadata
	if playlistMeta, ok := properties["playlistItemMetadata"].(map[string]any); ok {
		if snippet, ok := playlistMeta["snippet"].(map[string]any); ok {
			if publishedStr, ok := snippet["publishedAt"].(string); ok {
				if t, err := time.Parse(time.RFC3339, publishedStr); err == nil {
					return t
				}
			}
		}
	}

	return time.Now()
}

// extractThumbnailURI extracts the best available thumbnail URI.
func extractThumbnailURI(meta videoMetadata) string {
	if meta.thumbnails == nil {
		return ""
	}

	// Prefer higher quality thumbnails
	for _, key := range []string{"maxres", "high", "medium", "default"} {
		if thumb, ok := meta.thumbnails[key].(map[string]any); ok {
			if url, ok := thumb["url"].(string); ok && url != "" {
				return url
			}
		}
	}

	return ""
}

// supportedLocales is the set of locales supported by the platform.
var supportedLocales = map[string]bool{
	"ar": true, "de": true, "en": true, "es": true,
	"fr": true, "it": true, "ja": true, "ko": true,
	"nl": true, "pt-PT": true, "ru": true, "tr": true,
	"zh-CN": true,
}

// youtubeLocaleMap maps YouTube language codes to platform locale codes.
var youtubeLocaleMap = map[string]string{
	"pt":    "pt-PT",
	"pt-BR": "pt-PT",
	"zh":    "zh-CN",
	"zh-CN": "zh-CN",
	"zh-TW": "zh-CN",
}

// detectLocaleFromYouTubeVideo detects the locale from YouTube video metadata.
func detectLocaleFromYouTubeVideo(meta videoMetadata, fallback string) string {
	// Try defaultLanguage first, then defaultAudioLanguage
	for _, lang := range []string{meta.defaultLanguage, meta.defaultAudioLanguage} {
		if lang == "" {
			continue
		}

		// Check direct mapping
		if mapped, ok := youtubeLocaleMap[lang]; ok {
			return mapped
		}

		// Check if it's a supported locale directly
		langLower := strings.ToLower(lang)
		if supportedLocales[langLower] {
			return langLower
		}

		// Try just the language part (e.g., "en-US" -> "en")
		parts := strings.SplitN(lang, "-", 2)
		if len(parts) > 0 {
			baseLang := strings.ToLower(parts[0])
			if supportedLocales[baseLang] {
				return baseLang
			}

			if mapped, ok := youtubeLocaleMap[baseLang]; ok {
				return mapped
			}
		}
	}

	if fallback != "" {
		return fallback
	}

	return "en"
}

// slugSanitizeRegexp matches non-alphanumeric characters.
var slugSanitizeRegexp = regexp.MustCompile(`[^a-z0-9]+`)

// generateSlugFromTitle generates a slug from published date and title.
func generateSlugFromTitle(publishedAt time.Time, title string) string {
	datePrefix := publishedAt.Format("20060102") + "-"

	// Sanitize title: lowercase, replace non-alphanumeric with dashes
	slug := strings.ToLower(title)

	// Transliterate common unicode characters
	slug = transliterateBasic(slug)

	slug = slugSanitizeRegexp.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		slug = "video"
	}

	// Enforce max length
	maxContentLen := maxSlugLength - slugDatePrefixLen
	if len(slug) > maxContentLen {
		slug = slug[:maxContentLen]
		slug = strings.TrimRight(slug, "-")
	}

	return datePrefix + slug
}

// transliterateBasic performs basic transliteration for common characters.
func transliterateBasic(s string) string {
	var b strings.Builder

	b.Grow(len(s))

	for _, r := range s {
		if r < unicode.MaxASCII {
			b.WriteRune(r)

			continue
		}

		// Replace common characters; others become dashes during sanitization
		switch r {
		case 'ı':
			b.WriteRune('i')
		case 'ğ':
			b.WriteRune('g')
		case 'ü':
			b.WriteRune('u')
		case 'ş':
			b.WriteRune('s')
		case 'ö':
			b.WriteRune('o')
		case 'ç':
			b.WriteRune('c')
		case 'İ':
			b.WriteRune('i')
		case 'Ğ':
			b.WriteRune('g')
		case 'Ü':
			b.WriteRune('u')
		case 'Ş':
			b.WriteRune('s')
		case 'Ö':
			b.WriteRune('o')
		case 'Ç':
			b.WriteRune('c')
		case 'ä':
			b.WriteString("ae")
		case 'é', 'è', 'ê', 'ë':
			b.WriteRune('e')
		case 'á', 'à', 'â':
			b.WriteRune('a')
		case 'í', 'ì', 'î':
			b.WriteRune('i')
		case 'ó', 'ò', 'ô':
			b.WriteRune('o')
		case 'ú', 'ù', 'û':
			b.WriteRune('u')
		case 'ñ':
			b.WriteRune('n')
		default:
			b.WriteRune('-')
		}
	}

	return b.String()
}

// buildStoryContent builds the story content from a YouTube video.
func buildStoryContent(videoID string, description string) string {
	var b strings.Builder

	// YouTube embed
	b.WriteString(`<iframe width="560" height="315" src="https://www.youtube.com/embed/`)
	b.WriteString(videoID)
	b.WriteString(
		`" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
	)

	if description != "" {
		b.WriteString("\n\n")
		b.WriteString(description)
	}

	return b.String()
}

// truncateSummary truncates text to maxLen, breaking at word boundaries.
func truncateSummary(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	// Find last space before maxLen
	truncated := text[:maxLen]
	lastSpace := strings.LastIndex(truncated, " ")

	if lastSpace > maxLen/2 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}
