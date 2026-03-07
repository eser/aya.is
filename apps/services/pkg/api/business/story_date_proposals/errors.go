package story_date_proposals

import "errors"

var (
	ErrFailedToGetRecord     = errors.New("failed to get record")
	ErrFailedToInsertRecord  = errors.New("failed to insert record")
	ErrFailedToUpdateRecord  = errors.New("failed to update record")
	ErrFailedToRemoveRecord  = errors.New("failed to remove record")
	ErrFailedToListRecords   = errors.New("failed to list records")
	ErrProposalNotFound      = errors.New("date proposal not found")
	ErrStoryNotFound         = errors.New("story not found")
	ErrDateNotUndecided      = errors.New("activity date mode is not undecided")
	ErrDateAlreadyFinalized  = errors.New("activity date has already been finalized")
	ErrCannotRemoveFinalized = errors.New("cannot remove a finalized proposal")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInvalidVoteDirection  = errors.New("vote direction must be +1 or -1")
)
