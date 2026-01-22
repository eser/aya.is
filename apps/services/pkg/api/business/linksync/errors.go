package linksync

import "errors"

// Sentinel errors.
var (
	ErrFailedToGetLinks        = errors.New("failed to get managed links")
	ErrFailedToGetLatestImport = errors.New("failed to get latest import")
	ErrFailedToUpsertImport    = errors.New("failed to upsert import")
	ErrFailedToMarkDeleted     = errors.New("failed to mark imports as deleted")
	ErrFailedToUpdateTokens    = errors.New("failed to update tokens")
	ErrFailedToFetchStories    = errors.New("failed to fetch remote stories")
	ErrFailedToRefreshToken    = errors.New("failed to refresh access token")
	ErrNoRefreshToken          = errors.New("no refresh token available")
	ErrTokenExpired            = errors.New("access token expired and refresh failed")
	ErrImportNotFound          = errors.New("import not found")
)
