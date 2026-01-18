package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) GetSessionByID(
	ctx context.Context,
	id string,
) (*users.Session, error) {
	row, err := r.queries.GetSessionByID(ctx, GetSessionByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.Session{
		ID:                       row.ID,
		Status:                   row.Status,
		OauthRequestState:        row.OauthRequestState,
		OauthRequestCodeVerifier: row.OauthRequestCodeVerifier,
		OauthRedirectURI:         vars.ToStringPtr(row.OauthRedirectURI),
		LoggedInUserID:           vars.ToStringPtr(row.LoggedInUserID),
		LoggedInAt:               vars.ToTimePtr(row.LoggedInAt),
		ExpiresAt:                vars.ToTimePtr(row.ExpiresAt),
		CreatedAt:                row.CreatedAt,
		UpdatedAt:                vars.ToTimePtr(row.UpdatedAt),
	}

	return result, nil
}

func (r *Repository) CreateSession(
	ctx context.Context,
	session *users.Session,
) error {
	err := r.queries.CreateSession(ctx, CreateSessionParams{
		ID:                       session.ID,
		Status:                   session.Status,
		OauthRequestState:        session.OauthRequestState,
		OauthRequestCodeVerifier: session.OauthRequestCodeVerifier,
		OauthRedirectURI:         vars.ToSQLNullString(session.OauthRedirectURI),
		LoggedInUserID:           vars.ToSQLNullString(session.LoggedInUserID),
		LoggedInAt:               vars.ToSQLNullTime(session.LoggedInAt),
		ExpiresAt:                vars.ToSQLNullTime(session.ExpiresAt),
		CreatedAt:                session.CreatedAt,
		UpdatedAt:                vars.ToSQLNullTime(session.UpdatedAt),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateSessionLoggedInAt(
	ctx context.Context,
	id string,
	loggedInAt time.Time,
) error {
	err := r.queries.UpdateSessionLoggedInAt(ctx, UpdateSessionLoggedInAtParams{
		ID:         id,
		LoggedInAt: sql.NullTime{Time: loggedInAt, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

// Session Preferences

func (r *Repository) GetSessionPreferences(
	ctx context.Context,
	sessionID string,
) (map[string]string, error) {
	rows, err := r.queries.GetSessionPreferences(
		ctx,
		GetSessionPreferencesParams{SessionID: sessionID},
	)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.Key] = row.Value
	}

	return result, nil
}

func (r *Repository) GetSessionPreference(
	ctx context.Context,
	sessionID string,
	key string,
) (*SessionPreference, error) {
	row, err := r.queries.GetSessionPreference(ctx, GetSessionPreferenceParams{
		SessionID: sessionID,
		Key:       key,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return row, nil
}

func (r *Repository) SetSessionPreference(
	ctx context.Context,
	sessionID string,
	key string,
	value string,
) error {
	return r.queries.SetSessionPreference(ctx, SetSessionPreferenceParams{
		SessionID: sessionID,
		Key:       key,
		Value:     value,
	})
}

func (r *Repository) DeleteSessionPreference(
	ctx context.Context,
	sessionID string,
	key string,
) error {
	return r.queries.DeleteSessionPreference(ctx, DeleteSessionPreferenceParams{
		SessionID: sessionID,
		Key:       key,
	})
}

// Session Rate Limiting

func (r *Repository) CheckAndIncrementSessionRateLimit(
	ctx context.Context,
	ipHash string,
	limit int,
	windowSeconds int,
) (bool, error) {
	// First get current count
	row, err := r.queries.GetSessionRateLimit(ctx, GetSessionRateLimitParams{IpHash: ipHash})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	// Check if within limit
	if row != nil {
		// Check if window is still valid (1 hour)
		windowDuration := time.Duration(windowSeconds) * time.Second
		if time.Since(row.WindowStart) < windowDuration && int(row.Count) >= limit {
			return false, nil // Rate limited
		}
	}

	// Upsert the rate limit
	err = r.queries.UpsertSessionRateLimit(ctx, UpsertSessionRateLimitParams{IpHash: ipHash})
	if err != nil {
		return false, err
	}

	return true, nil
}

// sessions.Repository interface implementation
// These methods wrap the Session-prefixed methods to satisfy the interface.

// GetPreferences implements sessions.Repository.
func (r *Repository) GetPreferences(
	ctx context.Context,
	sessionID string,
) (sessions.SessionPreferences, error) {
	return r.GetSessionPreferences(ctx, sessionID)
}

// GetPreference implements sessions.Repository.
func (r *Repository) GetPreference(
	ctx context.Context,
	sessionID string,
	key string,
) (*sessions.SessionPreference, error) {
	row, err := r.GetSessionPreference(ctx, sessionID, key)
	if err != nil {
		return nil, err
	}

	if row == nil {
		return nil, nil //nolint:nilnil
	}

	return &sessions.SessionPreference{
		SessionID: row.SessionID,
		Key:       row.Key,
		Value:     row.Value,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// SetPreference implements sessions.Repository.
func (r *Repository) SetPreference(
	ctx context.Context,
	sessionID string,
	key string,
	value string,
) error {
	return r.SetSessionPreference(ctx, sessionID, key, value)
}

// DeletePreference implements sessions.Repository.
func (r *Repository) DeletePreference(
	ctx context.Context,
	sessionID string,
	key string,
) error {
	return r.DeleteSessionPreference(ctx, sessionID, key)
}

// CheckAndIncrementRateLimit implements sessions.Repository.
func (r *Repository) CheckAndIncrementRateLimit(
	ctx context.Context,
	ipHash string,
	limit int,
	windowSeconds int,
) (bool, error) {
	return r.CheckAndIncrementSessionRateLimit(ctx, ipHash, limit, windowSeconds)
}
