package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

type telegramAdapter struct {
	repo *Repository
}

// NewTelegramRepository creates a new adapter that implements telegram.Repository.
func NewTelegramRepository(repo *Repository) telegrambiz.Repository {
	return &telegramAdapter{repo: repo}
}

func (a *telegramAdapter) CreateExternalCode(
	ctx context.Context,
	code *telegrambiz.ExternalCode,
) error {
	propsJSON, err := json.Marshal(code.Properties)
	if err != nil {
		return err
	}

	return a.repo.queries.CreateExternalCode(ctx, CreateExternalCodeParams{
		ID:             code.ID,
		Code:           code.Code,
		ExternalSystem: code.ExternalSystem,
		Properties:     pqtype.NullRawMessage{RawMessage: propsJSON, Valid: true},
		ExpiresAt:      time.Now().Add(telegrambiz.CodeExpiryMinutes * time.Minute),
	})
}

func (a *telegramAdapter) GetExternalCodeByCode(
	ctx context.Context,
	code string,
) (*telegrambiz.ExternalCode, error) {
	row, err := a.repo.queries.GetExternalCodeByCode(
		ctx,
		GetExternalCodeByCodeParams{
			Code: code,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, telegrambiz.ErrCodeNotFound
		}

		return nil, err
	}

	props := make(map[string]any)

	if row.Properties.Valid && len(row.Properties.RawMessage) > 0 {
		err = json.Unmarshal(row.Properties.RawMessage, &props)
		if err != nil {
			return nil, err
		}
	}

	result := &telegrambiz.ExternalCode{
		ID:             row.ID,
		Code:           row.Code,
		ExternalSystem: row.ExternalSystem,
		Properties:     props,
		CreatedAt:      row.CreatedAt,
		ExpiresAt:      row.ExpiresAt,
	}

	if row.ConsumedAt.Valid {
		t := row.ConsumedAt.Time
		result.ConsumedAt = &t
	}

	return result, nil
}

func (a *telegramAdapter) ConsumeExternalCode(ctx context.Context, code string) error {
	rowsAffected, err := a.repo.queries.ConsumeExternalCode(
		ctx,
		ConsumeExternalCodeParams{
			Code: code,
		},
	)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return telegrambiz.ErrCodeConsumed
	}

	return nil
}

func (a *telegramAdapter) CleanupExpiredCodes(ctx context.Context) error {
	_, err := a.repo.queries.CleanupExpiredExternalCodes(ctx)

	return err
}

func (a *telegramAdapter) GetProfileLinkByTelegramRemoteID(
	ctx context.Context,
	remoteID string,
) (*telegrambiz.ProfileLinkInfo, error) {
	row, err := a.repo.queries.GetProfileLinkByTelegramRemoteID(
		ctx,
		GetProfileLinkByTelegramRemoteIDParams{
			RemoteID: sql.NullString{String: remoteID, Valid: true},
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, telegrambiz.ErrNotLinked
		}

		return nil, err
	}

	return &telegrambiz.ProfileLinkInfo{
		ID:        row.ID,
		ProfileID: row.ProfileID,
		RemoteID:  row.RemoteID.String,
		PublicID:  row.PublicID.String,
	}, nil
}

func (a *telegramAdapter) GetProfileLinkByProfileIDAndTelegram(
	ctx context.Context,
	profileID string,
) (*telegrambiz.ProfileLinkInfo, error) {
	row, err := a.repo.queries.GetProfileLinkByProfileIDAndTelegram(
		ctx,
		GetProfileLinkByProfileIDAndTelegramParams{
			ProfileID: profileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, telegrambiz.ErrNotLinked
		}

		return nil, err
	}

	return &telegrambiz.ProfileLinkInfo{
		ID:        row.ID,
		ProfileID: row.ProfileID,
		RemoteID:  row.RemoteID.String,
		PublicID:  row.PublicID.String,
	}, nil
}

func (a *telegramAdapter) CreateTelegramProfileLink(
	ctx context.Context,
	params *telegrambiz.CreateProfileLinkParams,
) error {
	uriPtr := &params.URI
	if params.URI == "" {
		uriPtr = nil
	}

	addedByPtr := &params.AddedByProfileID
	if params.AddedByProfileID == "" {
		addedByPtr = nil
	}

	_, err := a.repo.queries.CreateProfileLink(ctx, CreateProfileLinkParams{
		ID:         params.ID,
		Kind:       "telegram",
		ProfileID:  params.ProfileID,
		LinkOrder:  int32(params.Order),
		IsManaged:  true,
		IsVerified: true,
		IsFeatured: false,
		Visibility: "public",
		RemoteID: sql.NullString{
			String: params.RemoteID,
			Valid:  params.RemoteID != "",
		},
		PublicID: sql.NullString{
			String: params.PublicID,
			Valid:  params.PublicID != "",
		},
		URI:                       vars.ToSQLNullString(uriPtr),
		AuthProvider:              sql.NullString{Valid: false},
		AuthAccessTokenScope:      sql.NullString{Valid: false},
		AuthAccessToken:           sql.NullString{Valid: false},
		AuthAccessTokenExpiresAt:  sql.NullTime{Valid: false},
		AuthRefreshToken:          sql.NullString{Valid: false},
		AuthRefreshTokenExpiresAt: sql.NullTime{Valid: false},
		AddedByProfileID:          vars.ToSQLNullString(addedByPtr),
	})

	return err
}

func (a *telegramAdapter) SoftDeleteTelegramProfileLink(
	ctx context.Context,
	remoteID string,
) error {
	rowsAffected, err := a.repo.queries.SoftDeleteTelegramProfileLink(
		ctx,
		SoftDeleteTelegramProfileLinkParams{
			RemoteID: sql.NullString{String: remoteID, Valid: true},
		},
	)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return telegrambiz.ErrNotLinked
	}

	return nil
}

func (a *telegramAdapter) GetMemberProfileTelegramLinks(
	ctx context.Context,
	memberProfileID string,
) ([]telegrambiz.RawGroupTelegramLink, error) {
	rows, err := a.repo.queries.GetMemberProfileTelegramLinks(
		ctx,
		GetMemberProfileTelegramLinksParams{
			MemberProfileID: sql.NullString{String: memberProfileID, Valid: true},
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]telegrambiz.RawGroupTelegramLink, 0, len(rows))

	for _, row := range rows {
		result = append(result, telegrambiz.RawGroupTelegramLink{
			ProfileSlug:    row.ProfileSlug,
			ProfileTitle:   row.ProfileTitle,
			MembershipKind: row.MembershipKind,
			LinkTitle:      row.LinkTitle,
			LinkURI:        row.URI.String,
			LinkPublicID:   row.LinkPublicID.String,
			LinkVisibility: row.LinkVisibility,
		})
	}

	return result, nil
}

func (a *telegramAdapter) GetMaxProfileLinkOrder(
	ctx context.Context,
	profileID string,
) (int, error) {
	return a.repo.GetMaxProfileLinkOrder(ctx, profileID)
}

func (a *telegramAdapter) GetProfileSlugByID(
	ctx context.Context,
	profileID string,
) (string, error) {
	slug, err := a.repo.queries.GetProfileSlugByIDForTelegram(
		ctx,
		GetProfileSlugByIDForTelegramParams{
			ID: profileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", telegrambiz.ErrNotLinked
		}

		return "", err
	}

	return slug, nil
}
