package story_interactions

import "errors"

var (
	ErrFailedToGetInteraction    = errors.New("failed to get interaction")
	ErrFailedToSetInteraction    = errors.New("failed to set interaction")
	ErrFailedToRemoveInteraction = errors.New("failed to remove interaction")
	ErrFailedToListInteractions  = errors.New("failed to list interactions")
	ErrFailedToCountInteractions = errors.New("failed to count interactions")
	ErrInteractionNotFound       = errors.New("interaction not found")
	ErrInvalidInteractionKind    = errors.New("invalid interaction kind")
)
