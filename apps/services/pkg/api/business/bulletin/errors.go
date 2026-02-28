package bulletin

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("bulletin subscription not found")
	ErrNoStoriesForDigest   = errors.New("not enough stories for digest")
	ErrChannelNotAvailable  = errors.New("bulletin channel not available")
)
