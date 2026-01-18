package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/protection"
)

// CreatePOWChallenge creates a new PoW challenge.
func (r *Repository) CreatePOWChallenge(
	ctx context.Context,
	challenge *protection.POWChallenge,
) error {
	return r.queries.CreatePOWChallenge(ctx, CreatePOWChallengeParams{
		ID:         challenge.ID,
		Prefix:     challenge.Prefix,
		Difficulty: int16(challenge.Difficulty),
		IpHash:     challenge.IPHash,
		Used:       challenge.Used,
		ExpiresAt:  challenge.ExpiresAt,
		CreatedAt:  challenge.CreatedAt,
	})
}

// GetPOWChallengeByID gets a PoW challenge by ID.
func (r *Repository) GetPOWChallengeByID(
	ctx context.Context,
	id string,
) (*protection.POWChallenge, error) {
	row, err := r.queries.GetPOWChallengeByID(ctx, GetPOWChallengeByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &protection.POWChallenge{
		ID:         row.ID,
		Prefix:     row.Prefix,
		Difficulty: int(row.Difficulty),
		IPHash:     row.IpHash,
		Used:       row.Used,
		ExpiresAt:  row.ExpiresAt,
		CreatedAt:  row.CreatedAt,
	}, nil
}

// MarkPOWChallengeUsed marks a PoW challenge as used.
func (r *Repository) MarkPOWChallengeUsed(
	ctx context.Context,
	id string,
) error {
	return r.queries.MarkPOWChallengeUsed(ctx, MarkPOWChallengeUsedParams{ID: id})
}

// DeleteExpiredPOWChallenges deletes all expired PoW challenges.
func (r *Repository) DeleteExpiredPOWChallenges(
	ctx context.Context,
) error {
	return r.queries.DeleteExpiredPOWChallenges(ctx)
}
