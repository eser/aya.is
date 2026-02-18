package resourcesync

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// Service handles resource synchronization business logic.
type Service struct {
	logger *logfx.Logger
	repo   Repository
}

// NewService creates a new resource sync service.
func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

// GetGitHubResourcesForSync returns GitHub resources that need syncing.
func (s *Service) GetGitHubResourcesForSync(
	ctx context.Context,
	batchSize int,
) ([]*GitHubResourceForSync, error) {
	resources, err := s.repo.ListGitHubResourcesForSync(ctx, batchSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetResources, err)
	}

	return resources, nil
}

// UpdateResourceProperties updates a resource's stored properties.
func (s *Service) UpdateResourceProperties(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	err := s.repo.UpdateProfileResourceProperties(ctx, id, properties)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateResource, err)
	}

	return nil
}

// UpdateMembershipProperties updates a membership's properties.
func (s *Service) UpdateMembershipProperties(
	ctx context.Context,
	id string,
	properties map[string]any,
) error {
	err := s.repo.UpdateProfileMembershipProperties(ctx, id, properties)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUpdateMembership, err)
	}

	return nil
}

// GetMembershipByProfiles finds a membership between two profiles.
func (s *Service) GetMembershipByProfiles(
	ctx context.Context,
	profileID string,
	memberProfileID string,
) (string, error) {
	return s.repo.GetMembershipByProfiles(ctx, profileID, memberProfileID)
}

// GetProfileLinkByRemoteID finds a profile (by its linked GitHub remote ID).
func (s *Service) GetProfileLinkByRemoteID(
	ctx context.Context,
	kind string,
	remoteID string,
) (string, error) {
	return s.repo.GetProfileLinkByRemoteID(ctx, kind, remoteID)
}

// GetProfileLinksByRemoteIDs batch-loads profile_links by kind and multiple remote_ids.
func (s *Service) GetProfileLinksByRemoteIDs(
	ctx context.Context,
	kind string,
	remoteIDs []string,
) (map[string]string, error) {
	return s.repo.GetProfileLinksByRemoteIDs(ctx, kind, remoteIDs)
}

// GetMembershipsByProfilePairs batch-loads memberships for multiple profile pairs.
func (s *Service) GetMembershipsByProfilePairs(
	ctx context.Context,
	profileIDs []string,
	memberProfileIDs []string,
) (map[string]string, error) {
	return s.repo.GetMembershipsByProfilePairs(ctx, profileIDs, memberProfileIDs)
}
