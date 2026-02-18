package appcontext

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/configfx"
	"github.com/eser/aya.is/services/pkg/ajan/connfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/arcade"
	"github.com/eser/aya.is/services/pkg/api/adapters/auth_tokens"
	"github.com/eser/aya.is/services/pkg/api/adapters/github"
	"github.com/eser/aya.is/services/pkg/api/adapters/s3client"
	"github.com/eser/aya.is/services/pkg/api/adapters/speakerdeck"
	"github.com/eser/aya.is/services/pkg/api/adapters/storage"
	"github.com/eser/aya.is/services/pkg/api/adapters/unsplash"
	"github.com/eser/aya.is/services/pkg/api/adapters/workers"
	"github.com/eser/aya.is/services/pkg/api/adapters/youtube"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/linksync"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/profile_questions"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/protection"
	"github.com/eser/aya.is/services/pkg/api/business/resourcesync"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/siteimporter"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/uploads"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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
	AIModels    *aifx.Registry

	Arcade *arcade.Arcade

	Repository      *storage.Repository
	JWTTokenService *auth_tokens.JWTTokenService
	S3Client        *s3client.Client

	// External Services
	GitHubClient        *github.Client
	GitHubProvider      *github.Provider
	YouTubeProvider     *youtube.Provider
	SpeakerDeckProvider *speakerdeck.Provider
	UnsplashClient      *unsplash.Client

	// Business
	UploadService              *uploads.Service
	AuthService                *auth.Service
	UserService                *users.Service
	ProfileService             *profiles.Service
	ProfilePointsService       *profile_points.Service
	ProfileQuestionsService    *profile_questions.Service
	StoryService               *stories.Service
	SessionService             *sessions.Service
	ProtectionService          *protection.Service
	ProfileLinkSyncService     *linksync.Service
	SiteImporterService        *siteimporter.Service
	ProfileResourceSyncService *resourcesync.Service
	AuditService               *events.AuditService
	QueueService               *events.QueueService
	QueueRegistry              *events.HandlerRegistry
	RuntimeStateService        *runtime_states.Service
	WorkerRegistry             *workerfx.Registry
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

	a.Logger.DebugContext(
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

	// ----------------------------------------------------
	// Adapter: AI Models (optional - only if configured)
	// ----------------------------------------------------
	a.AIModels = aifx.NewRegistry(
		aifx.WithLogger(a.Logger),
		aifx.WithDefaultFactories(),
	)

	if len(a.Config.AI.Targets) > 0 {
		err = a.AIModels.LoadFromConfig(ctx, &a.Config.AI)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInitFailed, err)
		}
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

	goose.SetLogger(a.Logger)

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
	// Adapter: S3Client (optional - only if configured)
	// ----------------------------------------------------
	if a.Config.S3.IsConfigured() {
		a.S3Client, err = s3client.New(ctx, a.Config.S3)
		if err != nil {
			return fmt.Errorf("%w: failed to create S3 client: %w", ErrInitFailed, err)
		}

		a.Logger.DebugContext(
			ctx,
			"[AppContext] S3 client initialized",
			slog.String("module", "appcontext"),
			slog.String("endpoint", a.Config.S3.Endpoint),
			slog.String("bucket", a.Config.S3.BucketName),
		)
	}

	// ----------------------------------------------------
	// Business Services
	// ----------------------------------------------------

	// ID generator function using users.DefaultIDGenerator
	idGen := func() string {
		return string(users.DefaultIDGenerator())
	}

	// Event system services (audit + queue)
	a.AuditService = events.NewAuditService(
		a.Logger,
		a.Repository,
		idGen,
		func(ctx context.Context) *string {
			if sid, ok := ctx.Value(httpfx.ContextKey("session_id")).(string); ok {
				return &sid
			}

			return nil
		},
	)

	a.QueueService = events.NewQueueService(
		a.Logger,
		a.Repository,
		idGen,
	)

	a.QueueRegistry = events.NewHandlerRegistry()

	a.ProfileService = profiles.NewService(
		a.Logger,
		&a.Config.Profiles,
		a.Repository,
		a.AuditService,
	)
	a.UserService = users.NewService(
		a.Logger,
		a.Repository,
		a.AuditService,
	)
	a.StoryService = stories.NewService(a.Logger, &a.Config.Stories, a.Repository, a.AuditService)
	a.ProfilePointsService = profile_points.NewService(
		a.Logger,
		a.Repository,
		profile_points.DefaultIDGenerator,
		a.AuditService,
	)

	a.ProfileQuestionsService = profile_questions.NewService(
		a.Logger,
		a.Repository,
		a.ProfileService,
		a.AuditService,
		profile_questions.DefaultIDGenerator,
	)

	// UploadService (only if S3 client is available)
	if a.S3Client != nil {
		a.UploadService = uploads.NewService(a.Logger, a.S3Client)
	}

	a.AuthService = auth.NewService(
		a.Logger,
		a.JWTTokenService,
		&a.Config.Auth,
		a.UserService,
		func(ctx context.Context, domain string) bool {
			profile, _, _ := a.ProfileService.GetByCustomDomain(ctx, "en", domain)

			return profile != nil
		},
	)

	a.ProtectionService = protection.NewService(
		a.Logger,
		&a.Config.Protection,
		a.Repository,
		idGen,
	)

	a.SessionService = sessions.NewService(
		a.Logger,
		&a.Config.Sessions,
		a.Repository,
		a.UserService,
		idGen,
		a.AuditService,
	)

	a.ProfileLinkSyncService = linksync.NewService(
		a.Logger,
		a.Repository,
		idGen,
	)

	a.SiteImporterService = siteimporter.NewService(
		a.Logger,
		a.ProfileLinkSyncService,
		idGen,
	)

	a.ProfileResourceSyncService = resourcesync.NewService(
		a.Logger,
		storage.NewResourceSyncRepository(a.Repository),
	)

	// Register points event handler
	pointsEventHandler := workers.NewPointsEventHandler(
		a.Logger,
		a.ProfilePointsService,
	)
	pointsEventHandler.RegisterHandlers(a.QueueRegistry)

	a.RuntimeStateService = runtime_states.NewService(a.Logger, a.Repository)
	a.WorkerRegistry = workerfx.NewRegistry()

	// ----------------------------------------------------
	// External Services
	// ----------------------------------------------------

	// GitHub provider (used for both auth and profile links)
	a.GitHubClient = github.NewClient(
		&a.Config.Auth.GitHub,
		a.Logger,
		a.HTTPClient,
	)

	a.GitHubProvider = github.NewProvider(a.GitHubClient)

	// YouTube provider (for profile links)
	a.YouTubeProvider = youtube.NewProvider(
		&a.Config.Auth.YouTube,
		a.Logger,
		a.HTTPClient,
	)

	// SpeakerDeck provider (for profile links - no OAuth, RSS-based)
	a.SpeakerDeckProvider = speakerdeck.NewProvider(
		a.Logger,
		a.HTTPClient,
	)
	a.SiteImporterService.RegisterProvider(a.SpeakerDeckProvider)

	// Unsplash client (for background images, optional)
	if a.Config.Externals.Unsplash.IsConfigured() {
		a.UnsplashClient = unsplash.NewClient(
			&a.Config.Externals.Unsplash,
			a.Logger,
			a.HTTPClient,
		)

		a.Logger.DebugContext(
			ctx,
			"[AppContext] Unsplash client initialized",
			slog.String("module", "appcontext"),
		)
	}

	// ----------------------------------------------------
	// Auth Providers (adapters)
	// ----------------------------------------------------
	a.AuthService.RegisterProvider("github", a.GitHubProvider)

	return nil
}
