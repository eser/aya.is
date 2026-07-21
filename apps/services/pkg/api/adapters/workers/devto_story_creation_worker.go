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

// DevtoStoryProcessor creates stories from synced Dev.to imports.
type DevtoStoryProcessor struct {
  config      *DevtoSyncConfig
  logger      *logfx.Logger
  syncService *linksync.Service
  storyRepo   storyCreationRepo
  idGenerator func() string
}

// NewDevtoStoryProcessor creates a new Dev.to story processor.
func NewDevtoStoryProcessor(
  config *DevtoSyncConfig,
  logger *logfx.Logger,
  syncService *linksync.Service,
  storyRepo storyCreationRepo,
  idGenerator func() string,
) *DevtoStoryProcessor {
  return &DevtoStoryProcessor{
    config:      config,
    logger:      logger,
    syncService: syncService,
    storyRepo:   storyRepo,
    idGenerator: idGenerator,
  }
}

// ProcessStories creates stories from new imports.
func (w *DevtoStoryProcessor) ProcessStories(ctx context.Context) error {
  w.logger.DebugContext(ctx, "Starting Dev.to story creation cycle")

  imports, err := w.syncService.ListImportsForStoryCreation(
    ctx,
    "devto",
    w.config.BatchSize,
  )
  if err != nil {
    return fmt.Errorf("%w: %w", ErrSyncFailed, err)
  }

  if len(imports) == 0 {
    w.logger.DebugContext(ctx, "No Dev.to imports need story creation")

    return nil
  }

  w.logger.DebugContext(ctx, "Processing Dev.to imports for story creation",
    slog.Int("count", len(imports)))

  created := 0

  for _, imp := range imports {
    err := w.createStoryFromImport(ctx, imp)
    if err != nil {
      w.logger.ErrorContext(ctx, "Failed to create story from Dev.to import",
        slog.String("import_id", imp.ID),
        slog.String("remote_id", imp.RemoteID),
        slog.String("profile_id", imp.ProfileID),
        slog.Any("error", err))

      continue
    }

    created++
  }

  w.logger.DebugContext(ctx, "Completed Dev.to story creation cycle",
    slog.Int("processed", len(imports)),
    slog.Int("created", created))

  return nil
}

type devtoImportMeta struct {
  publishedAt     time.Time
  storyPictureURI *string
  title           string
  description     string
  content         string
  link            string
  canonicalURL    string
  devtoURL        string
}

func extractDevtoImportMeta(imp *linksync.LinkImportForStoryCreation) *devtoImportMeta {
  title, _ := imp.Properties["title"].(string)
  description, _ := imp.Properties["description"].(string)
  content, _ := imp.Properties["content"].(string)
  link, _ := imp.Properties["link"].(string)
  canonicalURL, _ := imp.Properties["canonical_url"].(string)
  devtoURL, _ := imp.Properties["devto_url"].(string)
  thumbnailURL, _ := imp.Properties["thumbnail_url"].(string)
  publishedAtStr, _ := imp.Properties["published_at"].(string)

  if title == "" {
    title = "Untitled Article"
  }

  publishedAt := time.Now()
  if publishedAtStr != "" {
    parsed, parseErr := time.Parse(time.RFC3339, publishedAtStr)
    if parseErr == nil {
      publishedAt = parsed
    }
  }

  var storyPictureURI *string
  if thumbnailURL != "" {
    storyPictureURI = &thumbnailURL
  }

  return &devtoImportMeta{
    title:           title,
    description:     description,
    content:         content,
    link:            link,
    canonicalURL:    canonicalURL,
    devtoURL:        devtoURL,
    publishedAt:     publishedAt,
    storyPictureURI: storyPictureURI,
  }
}

func (w *DevtoStoryProcessor) createStoryFromImport(
  ctx context.Context,
  imp *linksync.LinkImportForStoryCreation,
) error {
  meta := extractDevtoImportMeta(imp)

  locale := imp.ProfileDefaultLocale
  if locale == "" {
    locale = "en"
  }

  slug := generateSlugFromTitle(meta.publishedAt, meta.title)
  storyID := w.idGenerator()
  publicationID := w.idGenerator()

  sourceURL := meta.canonicalURL
  if sourceURL == "" {
    sourceURL = meta.devtoURL
  }
  if sourceURL == "" {
    sourceURL = meta.link
  }

  properties := map[string]any{
    "managed_by": "devto_sync_worker",
    "remote_id":  imp.RemoteID,
  }

  if sourceURL != "" {
    properties["source_url"] = sourceURL
  }

  _, err := w.storyRepo.InsertStory(
    ctx, storyID, imp.ProfileID, slug, "article",
    meta.storyPictureURI, properties, true, &imp.RemoteID,
    "public", false,
  )
  if err != nil {
    return fmt.Errorf("failed to insert story: %w", err)
  }

  content := buildDevtoStoryContent(meta.content, sourceURL, meta.description)
  summary := truncateSummary(meta.description)

  err = w.storyRepo.InsertStoryTx(ctx, storyID, locale, meta.title, summary, content, true)
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

  w.logger.DebugContext(ctx, "Created story from Dev.to import",
    slog.String("story_id", storyID),
    slog.String("remote_id", imp.RemoteID),
    slog.String("slug", slug),
    slog.String("locale", locale))

  return nil
}

func buildDevtoStoryContent(markdown string, link string, description string) string {
  content := strings.TrimSpace(markdown)
  if content != "" {
    return content
  }

  if link != "" {
    content = "%[" + link + "]"
  }

  if description != "" {
    if content != "" {
      content += "\n\n"
    }

    content += description
  }

  return content
}
