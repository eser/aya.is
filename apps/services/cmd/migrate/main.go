package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/eser/aya.is/services/pkg/ajan/connfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	"github.com/pressly/goose/v3"
)

var (
	ErrDatasourceNameRequired   = errors.New("datasource name is required")
	ErrCommandRequired          = errors.New("command is required")
	ErrAppContextNotInitialized = errors.New("app context is not initialized")
	ErrDatasourceNotInitialized = errors.New("datasource is not initialized")
	ErrDatasourceNotSQLDB       = errors.New("datasource is not an instance of *sql.DB")
	ErrFailedToRunGoose         = errors.New("failed to run goose")
)

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return ErrDatasourceNameRequired
	}

	if len(args) < 2 { //nolint:mnd
		return ErrCommandRequired
	}

	datasourceName := args[0]
	command := args[1]
	rest := args[2:]

	appContext := appcontext.New()

	err := appContext.Init(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAppContextNotInitialized, err)
	}

	datasource := appContext.Connections.GetNamed(datasourceName)
	if datasource == nil {
		return ErrDatasourceNotInitialized
	}

	sqlDB, err := connfx.GetTypedConnection[*sql.DB](appContext.Connections, datasourceName)
	if err != nil {
		return ErrDatasourceNotSQLDB
	}

	dialect := datasource.GetProtocol()

	err = goose.SetDialect(dialect)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRunGoose, err)
	}

	migrationsPath := fmt.Sprintf("./etc/data/%s/migrations", datasourceName)

	err = goose.RunContext(ctx, command, sqlDB, migrationsPath, rest...)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRunGoose, err)
	}

	return nil
}

func main() {
	ctx := context.Background()

	err := run(ctx, os.Args[1:])
	if err != nil {
		panic(err)
	}
}
