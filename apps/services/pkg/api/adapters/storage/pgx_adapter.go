package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrPgxExecFailed    = errors.New("pgx exec failed")
	ErrPgxQueryFailed   = errors.New("pgx query failed")
	ErrPgxAcquireFailed = errors.New("pgx acquire connection failed")
)

// PgxAdapter wraps a pgxpool.Pool to implement the DBTX interface.
type PgxAdapter struct {
	pool *pgxpool.Pool
}

// NewPgxAdapter creates a new adapter for pgx pool.
func NewPgxAdapter(pool *pgxpool.Pool) *PgxAdapter {
	return &PgxAdapter{pool: pool}
}

// ExecContext implements DBTX.
func (a *PgxAdapter) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	commandTag, err := a.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPgxExecFailed, err)
	}

	return pgxResult{commandTag}, nil
}

// PrepareContext implements DBTX.
func (a *PgxAdapter) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// For pgx, we don't need to prepare statements explicitly
	// The pool handles this internally
	// Return a stub that will use the pool
	return nil, sql.ErrTxDone // This signals that prepare is not needed
}

// QueryContext implements DBTX.
func (a *PgxAdapter) QueryContext(
	ctx context.Context,
	query string,
	args ...any,
) (*sql.Rows, error) {
	rows, err := a.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPgxQueryFailed, err)
	}

	// Convert pgx.Rows to *sql.Rows using stdlib
	conn, err := a.pool.Acquire(ctx)
	if err != nil {
		rows.Close()

		return nil, fmt.Errorf("%w: %w", ErrPgxAcquireFailed, err)
	}
	defer conn.Release()

	// We need to use stdlib to get sql.Rows compatibility
	// For now, we'll need to use the standard postgres driver for queries
	// that return rows. This is a limitation of mixing pgx with database/sql
	db := stdlib.OpenDBFromPool(a.pool)

	defer func() { _ = db.Close() }()

	return db.QueryContext(ctx, query, args...)
}

// QueryRowContext implements DBTX.
func (a *PgxAdapter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	// Similar to QueryContext, we need to use stdlib for Row compatibility
	db := stdlib.OpenDBFromPool(a.pool)

	defer func() { _ = db.Close() }()

	return db.QueryRowContext(ctx, query, args...)
}

// GetStdlibDB returns a *sql.DB that uses the pgx pool underneath
// This is needed for migrations and other tools that require *sql.DB.
func (a *PgxAdapter) GetStdlibDB() *sql.DB {
	return stdlib.OpenDBFromPool(a.pool)
}

// pgxResult wraps pgx CommandTag to implement sql.Result.
type pgxResult struct {
	tag pgconn.CommandTag
}

func (r pgxResult) LastInsertId() (int64, error) {
	// PostgreSQL doesn't support LastInsertId in the same way as MySQL
	// You need to use RETURNING clause instead
	return 0, sql.ErrNoRows
}

func (r pgxResult) RowsAffected() (int64, error) {
	return r.tag.RowsAffected(), nil
}
