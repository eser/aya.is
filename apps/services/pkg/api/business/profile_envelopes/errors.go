package profile_envelopes

import "errors"

// Sentinel errors.
var (
	ErrEnvelopeNotFound    = errors.New("envelope not found")
	ErrInvalidStatus       = errors.New("invalid envelope status transition")
	ErrNotTargetProfile    = errors.New("envelope does not belong to this profile")
	ErrFailedToCreate      = errors.New("failed to create envelope")
	ErrFailedToUpdate      = errors.New("failed to update envelope")
	ErrAlreadyProcessed    = errors.New("envelope has already been processed")
	ErrInvalidEnvelopeKind = errors.New("invalid envelope kind")
)
