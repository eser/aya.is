package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
