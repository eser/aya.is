package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/queue"
)

// PointsEventHandler handles queue items that award profile points.
type PointsEventHandler struct {
	logger               *logfx.Logger
	profilePointsService *profile_points.Service
}

// NewPointsEventHandler creates a new points event handler.
func NewPointsEventHandler(
	logger *logfx.Logger,
	profilePointsService *profile_points.Service,
) *PointsEventHandler {
	return &PointsEventHandler{
		logger:               logger,
		profilePointsService: profilePointsService,
	}
}

// HandleNewStory handles the NEW_STORY item and creates a pending award.
func (h *PointsEventHandler) HandleNewStory(ctx context.Context, item *queue.Item) error {
	// Parse the payload
	var payload struct {
		ProfileID string `json:"profile_id"`
		StoryID   string `json:"story_id"`
	}

	payloadBytes, err := json.Marshal(item.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if payload.ProfileID == "" {
		return profile_points.ErrMissingProfileID
	}

	h.logger.Info(
		"Processing NEW_STORY item for points",
		"profile_id", payload.ProfileID,
		"story_id", payload.StoryID,
		"item_id", item.ID,
	)

	// Create a pending award for the story publication
	metadata := map[string]any{
		"story_id": payload.StoryID,
		"item_id":  item.ID,
	}

	_, err = h.profilePointsService.AwardForEvent(
		ctx,
		profile_points.EventStoryPublished,
		payload.ProfileID,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to create pending award: %w", err)
	}

	h.logger.Info(
		"Created pending point award for story publication",
		"profile_id", payload.ProfileID,
		"story_id", payload.StoryID,
		"amount", profile_points.AwardCategories[profile_points.EventStoryPublished].Amount,
	)

	return nil
}

// RegisterHandlers registers all points-related queue handlers.
func (h *PointsEventHandler) RegisterHandlers(registry *queue.HandlerRegistry) {
	registry.Register(queue.ItemTypeNewStory, h.HandleNewStory)
}
