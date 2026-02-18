package siteimporter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
)

// Service provides site import orchestration.
type Service struct {
	logger      *logfx.Logger
	providers   map[string]SiteProvider
	syncService *linksync.Service
	idGenerator func() string
}

// NewService creates a new site importer service.
func NewService(
	logger *logfx.Logger,
	syncService *linksync.Service,
	idGenerator func() string,
) *Service {
	return &Service{
		logger:      logger,
		providers:   make(map[string]SiteProvider),
		syncService: syncService,
		idGenerator: idGenerator,
	}
}

// RegisterProvider registers a site provider by kind.
func (s *Service) RegisterProvider(provider SiteProvider) {
	s.providers[provider.Kind()] = provider
}

// CheckConnection validates a URL using the specified provider.
func (s *Service) CheckConnection(
	ctx context.Context,
	kind string,
	url string,
) (*CheckResult, error) {
	provider, ok := s.providers[kind]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, kind)
	}

	result, err := provider.Check(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("provider check failed: %w", err)
	}

	return result, nil
}

// SyncPublicLink fetches items from a public provider and upserts imports.
func (s *Service) SyncPublicLink(
	ctx context.Context,
	link *linksync.PublicManagedLink,
) (*linksync.SyncResult, error) {
	provider, ok := s.providers[link.Kind]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, link.Kind)
	}

	result := &linksync.SyncResult{LinkID: link.ID} //nolint:exhaustruct

	items, err := provider.FetchAll(ctx, link.RemoteID)
	if err != nil {
		result.Error = err

		return result, nil //nolint:nilerr // error is stored in result for caller to inspect
	}

	s.logger.DebugContext(ctx, "Fetched items from provider",
		slog.String("kind", link.Kind),
		slog.String("link_id", link.ID),
		slog.Int("count", len(items)))

	s.upsertItems(ctx, link, result, items)

	return result, nil
}

// upsertItems processes fetched items, upserting imports and marking deleted ones.
func (s *Service) upsertItems(
	ctx context.Context,
	link *linksync.PublicManagedLink,
	result *linksync.SyncResult,
	items []*ImportItem,
) {
	activeRemoteIDs := make([]string, 0, len(items))

	for _, item := range items {
		props := item.Properties
		if props == nil {
			props = make(map[string]any)
		}

		// Store standard metadata in properties
		props["title"] = item.Title
		props["description"] = item.Description
		props["published_at"] = item.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		props["link"] = item.Link
		props["thumbnail_url"] = item.ThumbnailURL
		props["story_kind"] = item.StoryKind

		err := s.syncService.UpsertImport(ctx, link.ID, item.RemoteID, props)
		if err != nil {
			s.logger.WarnContext(ctx, "Failed to upsert import",
				slog.String("link_id", link.ID),
				slog.String("remote_id", item.RemoteID),
				slog.Any("error", err))

			continue
		}

		activeRemoteIDs = append(activeRemoteIDs, item.RemoteID)
		result.ItemsAdded++
	}

	// Mark deleted items
	if len(activeRemoteIDs) > 0 {
		deletedCount, err := s.syncService.MarkDeletedImports(ctx, link.ID, activeRemoteIDs)
		if err != nil {
			s.logger.WarnContext(ctx, "Failed to mark deleted imports",
				slog.String("link_id", link.ID),
				slog.Any("error", err))
		} else {
			result.ItemsDeleted = int(deletedCount)
		}
	}
}
