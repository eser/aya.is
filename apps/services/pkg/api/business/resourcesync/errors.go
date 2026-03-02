package resourcesync

import "errors"

// Sentinel errors.
var (
	ErrFailedToGetResources     = errors.New("failed to get resources for sync")
	ErrFailedToUpdateResource   = errors.New("failed to update resource properties")
	ErrFailedToUpdateMembership = errors.New("failed to update membership properties")
	ErrFailedToGetMembership    = errors.New("failed to get membership")
	ErrFailedToGetProfileLink   = errors.New("failed to get profile link")
	ErrFailedToGetProfileLinks  = errors.New("failed to get profile links")
	ErrFailedToGetMemberships   = errors.New("failed to get memberships")
)
