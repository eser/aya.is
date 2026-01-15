package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eser/aya.is/services/pkg/lib/vars"
	"github.com/sqlc-dev/pqtype"
)

func (r *Repository) CacheGet(ctx context.Context, key string) (*[]byte, error) {
	row, err := r.queries.GetFromCache(ctx, GetFromCacheParams{Key: key})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := vars.ToRawMessage(row.Value)

	return &result, nil
}

func (r *Repository) CacheGetSince(
	ctx context.Context,
	key string,
	since time.Time,
) (*[]byte, error) {
	row, err := r.queries.GetFromCacheSince(ctx, GetFromCacheSinceParams{Key: key, Since: since})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := vars.ToRawMessage(row.Value)

	return &result, nil
}

func (r *Repository) CacheSet(ctx context.Context, key string, value []byte) error {
	_, err := r.queries.SetInCache(
		ctx,
		SetInCacheParams{Key: key, Value: pqtype.NullRawMessage{RawMessage: value, Valid: true}},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CacheRemove(ctx context.Context, key string) error {
	_, err := r.queries.RemoveFromCache(ctx, RemoveFromCacheParams{Key: key})

	return err
}

func (r *Repository) CacheRemoveAll(ctx context.Context) error {
	_, err := r.queries.RemoveAllFromCache(ctx)

	return err
}

func (r *Repository) CacheRemoveExpired(ctx context.Context, before time.Time) error {
	_, err := r.queries.RemoveExpiredFromCache(ctx, RemoveExpiredFromCacheParams{Before: before})

	return err
}
