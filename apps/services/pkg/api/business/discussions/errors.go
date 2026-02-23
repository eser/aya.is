package discussions

import "errors"

// Sentinel errors.
var (
	ErrDiscussionsNotEnabled  = errors.New("discussions are not enabled for this profile")
	ErrThreadNotFound         = errors.New("discussion thread not found")
	ErrCommentNotFound        = errors.New("comment not found")
	ErrContentTooShort        = errors.New("comment content is too short")
	ErrContentTooLong         = errors.New("comment content is too long")
	ErrMaxNestingDepth        = errors.New("maximum nesting depth reached")
	ErrThreadLocked           = errors.New("this discussion thread is locked")
	ErrInsufficientPermission = errors.New("insufficient permission for this action")
	ErrInvalidVoteDirection   = errors.New("vote direction must be +1 or -1")

	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToInsertRecord = errors.New("failed to insert record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
	ErrFailedToListRecords  = errors.New("failed to list records")
)
