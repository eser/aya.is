package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

func (r *Repository) GetState(
	ctx context.Context,
	key string,
) (*runtime_states.RuntimeState, error) {
	row, err := r.queries.GetRuntimeState(ctx, GetRuntimeStateParams{Key: key})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	return &runtime_states.RuntimeState{
		Key:       row.Key,
		Value:     row.Value,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *Repository) SetState(ctx context.Context, key string, value string) error {
	return r.queries.SetRuntimeState(ctx, SetRuntimeStateParams{
		Key:   key,
		Value: value,
	})
}

func (r *Repository) RemoveState(ctx context.Context, key string) error {
	_, err := r.queries.RemoveRuntimeState(ctx, RemoveRuntimeStateParams{Key: key})

	return err
}

func (r *Repository) TryAdvisoryLock(ctx context.Context, lockID int64) (bool, error) {
	acquired, err := r.queries.TryAdvisoryLock(ctx, TryAdvisoryLockParams{LockID: lockID})
	if err != nil {
		return false, err
	}

	return acquired, nil
}

func (r *Repository) ReleaseAdvisoryLock(ctx context.Context, lockID int64) error {
	_, err := r.queries.ReleaseAdvisoryLock(ctx, ReleaseAdvisoryLockParams{LockID: lockID})

	return err
}
