package appcontext

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/configfx"
	"github.com/eser/aya.is/services/pkg/ajan/connfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/arcade"
	"github.com/eser/aya.is/services/pkg/api/adapters/auth_providers"
	"github.com/eser/aya.is/services/pkg/api/adapters/auth_tokens"
	"github.com/eser/aya.is/services/pkg/api/adapters/storage"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	_ "github.com/lib/pq"
)

var (
	ErrInitFailed       = errors.New("failed to initialize app context")
	ErrJWTSecretMissing = errors.New("JWT_SECRET environment variable is required")
)

type AppContext struct {
	// Adapters
	Config *AppConfig
	Logger *logfx.Logger

	HTTPClient *httpclient.Client

	Connections *connfx.Registry

	Arcade *arcade.Arcade

	Repository      *storage.Repository
	JWTTokenService *auth_tokens.JWTTokenService

	// Business
	AuthService    *auth.Service
	UserService    *users.Service
	ProfileService *profiles.Service
	StoryService   *stories.Service
}

func New() *AppContext {
	return &AppContext{} //nolint:exhaustruct
}

func (a *AppContext) Init(ctx context.Context) error { //nolint:funlen
	// ----------------------------------------------------
	// Adapter: Config
	// ----------------------------------------------------
	cl := configfx.NewConfigManager()

	a.Config = &AppConfig{} //nolint:exhaustruct

	err := cl.LoadDefaults(a.Config)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// Validate required secrets
	if a.Config.Auth.JwtSecret == "" {
		return fmt.Errorf("%w: %w", ErrInitFailed, ErrJWTSecretMissing)
	}

	// ----------------------------------------------------
	// Adapter: Logger
	// ----------------------------------------------------
	a.Logger = logfx.NewLogger(
		logfx.WithConfig(&a.Config.Log),
	)

	a.Logger.InfoContext(
		ctx,
		"[AppContext] Initialization in progress",
		slog.String("module", "appcontext"),
		slog.String("name", a.Config.AppName),
		slog.String("environment", a.Config.AppEnv),
		slog.String("site_uri", a.Config.SiteURI),
		slog.Any("features", a.Config.Features),
	)

	// ----------------------------------------------------
	// Adapter: HTTPClient
	// ----------------------------------------------------
	a.HTTPClient = httpclient.NewClient(
		httpclient.WithConfig(&a.Config.HTTPClient),
	)

	// ----------------------------------------------------
	// Adapter: Connections
	// ----------------------------------------------------
	a.Connections = connfx.NewRegistry(
		connfx.WithLogger(a.Logger),
		connfx.WithDefaultFactories(),
	)

	err = a.Connections.LoadFromConfig(ctx, &a.Config.Conn)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// // ----------------------------------------------------
	// // Adapter: Metrics
	// // ----------------------------------------------------
	// a.Metrics = metricsfx.NewMetricsProvider(
	// 	&a.Config.Metrics,
	// 	a.Connections,
	// )

	// err = a.Metrics.Init()
	// if err != nil {
	// 	return fmt.Errorf("%w: %w", ErrInitFailed, err)
	// }

	// ----------------------------------------------------
	// Adapter: Arcade
	// ----------------------------------------------------
	a.Arcade = arcade.New(
		a.Config.Externals.Arcade,
		a.HTTPClient,
	)

	// ----------------------------------------------------
	// Adapter: Repository
	// ----------------------------------------------------
	a.Repository, err = storage.NewRepositoryFromDefault(a.Logger, a.Connections)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// Run database migrations
	migrationsDir := a.Config.Data.MigrationsPath

	err = a.Repository.RunMigrations(ctx, migrationsDir)
	if err != nil {
		return fmt.Errorf("%w: failed to run migrations: %w", ErrInitFailed, err)
	}

	// Seed data if database is empty
	seedFilePath := a.Config.Data.SeedFilePath

	err = a.Repository.SeedData(ctx, seedFilePath)
	if err != nil {
		return fmt.Errorf("%w: failed to seed data: %w", ErrInitFailed, err)
	}

	// ----------------------------------------------------
	// Adapter: JWTTokenService
	// ----------------------------------------------------
	a.JWTTokenService = auth_tokens.NewJWTTokenService(&a.Config.Auth)

	// ----------------------------------------------------
	// Business Services
	// ----------------------------------------------------
	a.ProfileService = profiles.NewService(a.Logger, a.Repository)
	a.UserService = users.NewService(
		a.Logger,
		a.Repository,
	)
	a.StoryService = stories.NewService(a.Logger, a.Repository)

	a.AuthService = auth.NewService(
		a.Logger,
		a.JWTTokenService,
		&a.Config.Auth,
		a.UserService,
	)

	// ----------------------------------------------------
	// Auth Providers (adapters)
	// ----------------------------------------------------
	a.AuthService.RegisterProvider(
		"github",
		auth_providers.NewGitHubAuthProvider(
			&a.Config.Auth.GitHub,
			a.Logger,
			a.HTTPClient,
			a.JWTTokenService,
			a.UserService,
		),
	)

	return nil
}
