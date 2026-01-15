package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) GetUserByID(
	ctx context.Context,
	id string,
) (*users.User, error) {
	row, err := r.queries.GetUserByID(ctx, GetUserByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.User{
		ID:                  row.ID,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
		GithubRemoteID:      vars.ToStringPtr(row.GithubRemoteID),
		BskyHandle:          vars.ToStringPtr(row.BskyHandle),
		XHandle:             vars.ToStringPtr(row.XHandle),
		IndividualProfileID: vars.ToStringPtr(row.IndividualProfileID),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:           vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*users.User, error) {
	row, err := r.queries.GetUserByEmail(
		ctx,
		GetUserByEmailParams{Email: sql.NullString{String: email, Valid: true}},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.User{
		ID:                  row.ID,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
		GithubRemoteID:      vars.ToStringPtr(row.GithubRemoteID),
		BskyHandle:          vars.ToStringPtr(row.BskyHandle),
		XHandle:             vars.ToStringPtr(row.XHandle),
		IndividualProfileID: vars.ToStringPtr(row.IndividualProfileID),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:           vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) GetUserByGitHubRemoteID(
	ctx context.Context,
	githubRemoteID string,
) (*users.User, error) {
	row, err := r.queries.GetUserByGitHubRemoteID(
		ctx,
		GetUserByGitHubRemoteIDParams{
			GithubRemoteID: sql.NullString{String: githubRemoteID, Valid: true},
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.User{
		ID:                  row.ID,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
		GithubRemoteID:      vars.ToStringPtr(row.GithubRemoteID),
		BskyHandle:          vars.ToStringPtr(row.BskyHandle),
		XHandle:             vars.ToStringPtr(row.XHandle),
		IndividualProfileID: vars.ToStringPtr(row.IndividualProfileID),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:           vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) ListUsers(
	ctx context.Context,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*users.User], error) {
	var wrappedResponse cursors.Cursored[[]*users.User]

	rows, err := r.queries.ListUsers(
		ctx,
		ListUsersParams{
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*users.User, len(rows))
	for i, row := range rows {
		result[i] = &users.User{
			ID:                  row.ID,
			Kind:                row.Kind,
			Name:                row.Name,
			Email:               vars.ToStringPtr(row.Email),
			Phone:               vars.ToStringPtr(row.Phone),
			GithubHandle:        vars.ToStringPtr(row.GithubHandle),
			GithubRemoteID:      vars.ToStringPtr(row.GithubRemoteID),
			BskyHandle:          vars.ToStringPtr(row.BskyHandle),
			XHandle:             vars.ToStringPtr(row.XHandle),
			IndividualProfileID: vars.ToStringPtr(row.IndividualProfileID),
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
			DeletedAt:           vars.ToTimePtr(row.DeletedAt),
		}
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) CreateUser(
	ctx context.Context,
	user *users.User,
) error {
	err := r.queries.CreateUser(ctx, CreateUserParams{
		ID:                  user.ID,
		Kind:                user.Kind,
		Name:                user.Name,
		Email:               vars.ToSQLNullString(user.Email),
		Phone:               vars.ToSQLNullString(user.Phone),
		GithubHandle:        vars.ToSQLNullString(user.GithubHandle),
		GithubRemoteID:      vars.ToSQLNullString(user.GithubRemoteID),
		BskyHandle:          vars.ToSQLNullString(user.BskyHandle),
		BskyRemoteID:        sql.NullString{String: "", Valid: false},
		XHandle:             vars.ToSQLNullString(user.XHandle),
		XRemoteID:           sql.NullString{String: "", Valid: false},
		IndividualProfileID: vars.ToSQLNullString(user.IndividualProfileID),
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           vars.ToSQLNullTime(user.UpdatedAt),
		DeletedAt:           vars.ToSQLNullTime(user.DeletedAt),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUser(
	ctx context.Context,
	user *users.User,
) error {
	rowsAffected, err := r.queries.UpdateUser(ctx, UpdateUserParams{
		ID:                  user.ID,
		Kind:                user.Kind,
		Name:                user.Name,
		Email:               vars.ToSQLNullString(user.Email),
		Phone:               vars.ToSQLNullString(user.Phone),
		GithubHandle:        vars.ToSQLNullString(user.GithubHandle),
		GithubRemoteID:      vars.ToSQLNullString(user.GithubRemoteID),
		BskyHandle:          vars.ToSQLNullString(user.BskyHandle),
		BskyRemoteID:        sql.NullString{String: "", Valid: false},
		XHandle:             vars.ToSQLNullString(user.XHandle),
		XRemoteID:           sql.NullString{String: "", Valid: false},
		IndividualProfileID: vars.ToSQLNullString(user.IndividualProfileID),
	})
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) SetUserIndividualProfileID(
	ctx context.Context,
	userID string,
	profileID string,
) error {
	rowsAffected, err := r.queries.SetUserIndividualProfileID(ctx, SetUserIndividualProfileIDParams{
		ID:                  userID,
		IndividualProfileID: sql.NullString{String: profileID, Valid: true},
	})
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
