package story_interactions

import "context"

// Repository defines the storage operations for story interactions (port).
type Repository interface {
	// UpsertInteraction creates or refreshes an interaction.
	UpsertInteraction(
		ctx context.Context,
		id string,
		storyID string,
		profileID string,
		kind string,
	) (*StoryInteraction, error)

	// RemoveInteraction soft-deletes a specific interaction.
	RemoveInteraction(
		ctx context.Context,
		storyID string,
		profileID string,
		kind string,
	) (int64, error)

	// RemoveInteractionsByKinds soft-deletes all interactions matching given kinds.
	// Used for RSVP mutual exclusivity enforcement.
	RemoveInteractionsByKinds(
		ctx context.Context,
		storyID string,
		profileID string,
		kindsCSV string,
	) (int64, error)

	// GetInteraction returns a specific interaction.
	GetInteraction(
		ctx context.Context,
		storyID string,
		profileID string,
		kind string,
	) (*StoryInteraction, error)

	// ListInteractionsForProfile returns all interactions a profile has on a story.
	ListInteractionsForProfile(
		ctx context.Context,
		storyID string,
		profileID string,
	) ([]*StoryInteraction, error)

	// ListInteractions lists interactions on a story with profile info, optionally filtered by kind.
	ListInteractions(
		ctx context.Context,
		localeCode string,
		storyID string,
		filterKind *string,
	) ([]*InteractionWithProfile, error)

	// CountInteractionsByKind returns interaction counts grouped by kind.
	CountInteractionsByKind(
		ctx context.Context,
		storyID string,
	) ([]*InteractionCount, error)
}
