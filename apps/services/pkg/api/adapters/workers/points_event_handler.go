package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
)

// PointsEventHandler handles events that award profile points.
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

// HandleNewStory handles the NEW_STORY event and creates a pending award.
func (h *PointsEventHandler) HandleNewStory(ctx context.Context, event *events.Event) error {
	// Parse the payload
	var payload struct {
		ProfileID string `json:"profile_id"`
		StoryID   string `json:"story_id"`
	}

	payloadBytes, err := json.Marshal(event.Payload)
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
		"Processing NEW_STORY event for points",
		"profile_id", payload.ProfileID,
		"story_id", payload.StoryID,
		"event_id", event.ID,
	)

	// Create a pending award for the story publication
	metadata := map[string]any{
		"story_id": payload.StoryID,
		"event_id": event.ID,
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

// RegisterHandlers registers all points-related event handlers.
func (h *PointsEventHandler) RegisterHandlers(registry *events.HandlerRegistry) {
	registry.Register(events.EventTypeNewStory, h.HandleNewStory)
}
