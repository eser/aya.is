package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/privacy"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

var (
	ErrUnknownPrivacyDataset                    = errors.New("unknown privacy dataset")
	ErrPrivacyQuery                             = errors.New("privacy query failed")
	ErrPrivacyTransaction                       = errors.New("privacy transaction failed")
	_                        privacy.Repository = (*Repository)(nil)
)

func (r *Repository) GetAccount(ctx context.Context, userID string) (*privacy.Account, error) {
	row, err := r.queries.GetPrivacyStateByUserID(
		ctx,
		GetPrivacyStateByUserIDParams{UserID: userID},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return &privacy.Account{
		HiddenAt:    vars.ToTimePtr(row.HiddenAt),
		UserID:      row.UserID,
		ProfileID:   row.ProfileID,
		ProfileSlug: row.ProfileSlug,
		Redacted:    row.UserRedactedAt.Valid || row.ProfileRedactedAt.Valid,
	}, nil
}

func (r *Repository) GetContact(ctx context.Context, userID string) (*privacy.Contact, error) {
	row, err := r.queries.GetPrivacyContactByUserID(
		ctx,
		GetPrivacyContactByUserIDParams{UserID: userID},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	contact := &privacy.Contact{Email: "", Locale: row.DefaultLocale, Slug: row.Slug}
	if row.Email.Valid {
		contact.Email = row.Email.String
	}

	return contact, nil
}

// invalidateProfileCaches drops cached lookups that would keep serving a profile (or its old username)
// after its visibility changed.
func (r *Repository) invalidateProfileCaches(ctx context.Context, userID string, slugs ...string) {
	for _, slug := range slugs {
		_ = r.cache.Invalidate(ctx, "profile_id_by_slug:"+slug)
		_ = r.cache.Invalidate(ctx, "profile_slug_exists:"+slug)
	}

	_ = r.cache.Invalidate(ctx, "user_brief_info:"+userID)
}

func (r *Repository) HideProfile(
	ctx context.Context,
	userID string,
	profileID string,
) (bool, error) {
	account, err := r.GetAccount(ctx, userID)
	if err != nil {
		return false, err
	}

	dbTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrPrivacyTransaction, err)
	}

	defer func() { _ = dbTx.Rollback() }()

	queriesTx := r.queries.WithTx(dbTx)

	affected, err := queriesTx.HideProfileByID(ctx, HideProfileByIDParams{ID: profileID})
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	// Ending the sessions is what makes "hidden until you log in again" true.
	_, err = queriesTx.TerminateActiveSessionsByUserID(ctx, TerminateActiveSessionsByUserIDParams{
		UserID: vars.ToSQLNullString(&userID),
	})
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	err = dbTx.Commit()
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrPrivacyTransaction, err)
	}

	if account != nil {
		r.invalidateProfileCaches(ctx, userID, account.ProfileSlug)
	}

	return affected > 0, nil
}

func (r *Repository) UnhideProfile(ctx context.Context, userID string) (bool, error) {
	account, err := r.GetAccount(ctx, userID)
	if err != nil {
		return false, err
	}

	affected, err := r.queries.UnhideProfileByUserID(
		ctx,
		UnhideProfileByUserIDParams{UserID: userID},
	)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	if affected > 0 && account != nil {
		r.invalidateProfileCaches(ctx, userID, account.ProfileSlug)
	}

	return affected > 0, nil
}

func (r *Repository) ListSoleOwnedProfiles(
	ctx context.Context,
	profileID string,
) ([]privacy.OwnedProfile, error) {
	rows, err := r.queries.ListSoleOwnedProfiles(ctx, ListSoleOwnedProfilesParams{
		ProfileID: vars.ToSQLNullString(&profileID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	owned := make([]privacy.OwnedProfile, 0, len(rows))
	for _, row := range rows {
		owned = append(owned, privacy.OwnedProfile{
			ID:    row.ID,
			Slug:  row.Slug,
			Kind:  row.Kind,
			Title: row.Title,
		})
	}

	return owned, nil
}

func toDataExport(row *ProfileDataExport) *privacy.DataExport {
	return &privacy.DataExport{
		RequestedAt:  row.RequestedAt,
		StartedAt:    vars.ToTimePtr(row.StartedAt),
		CompletedAt:  vars.ToTimePtr(row.CompletedAt),
		ExpiresAt:    vars.ToTimePtr(row.ExpiresAt),
		NotifiedAt:   vars.ToTimePtr(row.NotifiedAt),
		StorageKey:   vars.ToStringPtr(row.StorageKey),
		SizeBytes:    toInt64Ptr(row.SizeBytes),
		ErrorMessage: vars.ToStringPtr(row.ErrorMessage),
		ID:           row.ID,
		UserID:       row.UserID,
		ProfileID:    row.ProfileID,
		Status:       privacy.ExportStatus(row.Status),
	}
}

func toInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}

	return &value.Int64
}

func (r *Repository) CreateExport(ctx context.Context, export *privacy.DataExport) error {
	err := r.queries.CreateDataExport(ctx, CreateDataExportParams{
		ID:        export.ID,
		UserID:    export.UserID,
		ProfileID: export.ProfileID,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return nil
}

func (r *Repository) GetLatestExport(
	ctx context.Context,
	userID string,
) (*privacy.DataExport, error) {
	row, err := r.queries.GetLatestDataExportByUserID(
		ctx,
		GetLatestDataExportByUserIDParams{UserID: userID},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return toDataExport(row), nil
}

func (r *Repository) GetExport(ctx context.Context, id string) (*privacy.DataExport, error) {
	row, err := r.queries.GetDataExportByID(ctx, GetDataExportByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return toDataExport(row), nil
}

func (r *Repository) ClaimExport(ctx context.Context, id string) (bool, error) {
	affected, err := r.queries.ClaimDataExport(ctx, ClaimDataExportParams{ID: id})
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return affected > 0, nil
}

func (r *Repository) MarkExportReady(
	ctx context.Context,
	exportID string,
	storageKey string,
	sizeBytes int64,
	expiresAt time.Time,
) error {
	err := r.queries.MarkDataExportReady(ctx, MarkDataExportReadyParams{
		ID:         exportID,
		ExpiresAt:  sql.NullTime{Time: expiresAt, Valid: true},
		StorageKey: sql.NullString{String: storageKey, Valid: true},
		SizeBytes:  sql.NullInt64{Int64: sizeBytes, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return nil
}

func (r *Repository) MarkExportFailed(ctx context.Context, id string, message string) error {
	err := r.queries.MarkDataExportFailed(ctx, MarkDataExportFailedParams{
		ID:           id,
		ErrorMessage: sql.NullString{String: message, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return nil
}

func (r *Repository) MarkExportNotified(ctx context.Context, id string) error {
	err := r.queries.MarkDataExportNotified(ctx, MarkDataExportNotifiedParams{ID: id})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return nil
}

func (r *Repository) ListExpiredExports(
	ctx context.Context,
	limit int,
) ([]*privacy.DataExport, error) {
	rows, err := r.queries.ListExpiredDataExports(
		ctx,
		ListExpiredDataExportsParams{BatchSize: int32(limit)},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	exports := make([]*privacy.DataExport, 0, len(rows))
	for _, row := range rows {
		exports = append(exports, toDataExport(row))
	}

	return exports, nil
}

func (r *Repository) MarkExportExpired(ctx context.Context, exportID string) error {
	err := r.queries.MarkDataExportExpired(ctx, MarkDataExportExpiredParams{ID: exportID})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return nil
}

func (r *Repository) ListExportsByUser(
	ctx context.Context,
	userID string,
) ([]*privacy.DataExport, error) {
	rows, err := r.queries.ListDataExportsByUserID(
		ctx,
		ListDataExportsByUserIDParams{UserID: userID},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	exports := make([]*privacy.DataExport, 0, len(rows))
	for _, row := range rows {
		exports = append(exports, toDataExport(row))
	}

	return exports, nil
}

// existingTables returns the tables of the public schema.
func (r *Repository) existingTables(ctx context.Context) (map[string]bool, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT table_name FROM information_schema.tables
		 WHERE table_schema = current_schema() AND table_type = 'BASE TABLE'`)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	defer func() { _ = rows.Close() }()

	tables := map[string]bool{}

	for rows.Next() {
		var name string

		scanErr := rows.Scan(&name)
		if scanErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, scanErr)
		}

		tables[name] = true
	}

	return tables, rows.Err() //nolint:wrapcheck
}

// ExportDatasetNames lists the datasets whose tables exist in this database.
func (r *Repository) ExportDatasetNames(ctx context.Context) ([]string, error) {
	existing, err := r.existingTables(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(privacyDatasets))

	for _, dataset := range privacyDatasets {
		if allTablesExist(dataset.query, existing) {
			names = append(names, dataset.name)
		}
	}

	return names, nil
}

func (r *Repository) StreamDataset(
	ctx context.Context,
	name string,
	userID string,
	profileID string,
	handleRow func(row json.RawMessage) error,
) error {
	dataset, ok := privacyDatasetsByName()[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownPrivacyDataset, name)
	}

	rows, err := r.db.QueryContext(ctx, dataset.query, userID, profileID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var doc []byte

		scanErr := rows.Scan(&doc)
		if scanErr != nil {
			return fmt.Errorf("%w: %w", ErrPrivacyQuery, scanErr)
		}

		handlerErr := handleRow(json.RawMessage(doc))
		if handlerErr != nil {
			return handlerErr
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	return nil
}
