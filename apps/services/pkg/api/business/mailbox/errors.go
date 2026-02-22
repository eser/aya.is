package mailbox

import "errors"

// Sentinel errors.
var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrEnvelopeNotFound     = errors.New("envelope not found")
	ErrNotParticipant       = errors.New("profile is not a participant in this conversation")
	ErrInvalidEnvelopeKind  = errors.New("invalid envelope kind")
	ErrInvalidStatus        = errors.New("invalid envelope status transition")
	ErrNotTargetProfile     = errors.New("envelope does not belong to this profile")
	ErrAlreadyProcessed     = errors.New("envelope has already been processed")
	ErrFailedToCreate       = errors.New("failed to create resource")
	ErrFailedToUpdate       = errors.New("failed to update resource")
	ErrInvalidReaction      = errors.New("invalid reaction emoji")
	ErrSelfConversation     = errors.New("cannot create conversation with yourself")
)
