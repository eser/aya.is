package http

import (
	"context"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/middlewares"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/healthcheck"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/openapi"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/profiling"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	mcpadapter "github.com/eser/aya.is/services/pkg/api/adapters/mcp"
	telegramadapter "github.com/eser/aya.is/services/pkg/api/adapters/telegram"
	"github.com/eser/aya.is/services/pkg/api/adapters/unsplash"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profile_envelopes"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/profile_questions"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/protection"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/story_interactions"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
	"github.com/eser/aya.is/services/pkg/api/business/uploads"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// TelegramProviders holds Telegram bot components (nil when Telegram is disabled).
type TelegramProviders struct {
	Bot           *telegramadapter.Bot
	Service       *telegrambiz.Service
	WebhookSecret string
}

func Run(
	ctx context.Context,
	baseURI string,
	config *httpfx.Config,
	logger *logfx.Logger,
	discloseErrors bool,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	profilePointsService *profile_points.Service,
	profileQuestionsService *profile_questions.Service,
	profileEnvelopesService *profile_envelopes.Service,
	storyService *stories.Service,
	storyInteractionService *story_interactions.Service,
	sessionService *sessions.Service,
	protectionService *protection.Service,
	uploadService *uploads.Service,
	unsplashClient *unsplash.Client,
	profileLinkProviders *ProfileLinkProviders,
	telegramProviders *TelegramProviders,
	aiModels *aifx.Registry,
	runtimeStatesService *runtime_states.Service,
	workerRegistry *workerfx.Registry,
) (func(), error) {
	httpfx.SetDiscloseErrors(discloseErrors)

	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHTTPService(config, routes, logger)

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(
		middlewares.TracingMiddleware(logger, "/health-check"),
	)
	routes.Use(CorsMiddlewareWithCustomDomains(authService.Config, profileService))
	routes.Use(middlewares.MetricsMiddleware(httpService.InnerMetrics)) //nolint:contextcheck

	// mcp adapter (must be registered before OPTIONS wildcard to avoid pattern conflict)
	mcpadapter.RegisterMCPRoutes(routes, profileService, storyService)

	// Global OPTIONS handler for preflight requests
	routes.Route("OPTIONS /{path...}", func(ctx *httpfx.Context) httpfx.Result {
		return ctx.Results.Ok()
	})

	// http modules
	healthcheck.RegisterHTTPRoutes(routes, config)
	openapi.RegisterHTTPRoutes(routes, config)
	profiling.RegisterHTTPRoutes(routes, config)

	// public http routes
	RegisterHTTPRoutesForUsers( //nolint:contextcheck
		baseURI,
		routes,
		logger,
		authService,
		userService,
		sessionService,
		profileService,
	)
	RegisterHTTPRoutesForSessions( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		sessionService,
		protectionService,
		profileEnvelopesService,
	)
	RegisterHTTPRoutesForProtection( //nolint:contextcheck
		routes,
		logger,
		protectionService,
	)
	RegisterHTTPRoutesForSite( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		uploadService,
		unsplashClient,
	)
	RegisterHTTPRoutesForProfiles( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		storyService,
		profilePointsService,
		aiModels,
	)
	RegisterHTTPRoutesForProfilePoints( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profilePointsService,
	)
	RegisterHTTPRoutesForStories( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		storyService,
		profilePointsService,
		aiModels,
	)
	RegisterHTTPRoutesForSearch( //nolint:contextcheck
		routes,
		profileService,
	)
	RegisterHTTPRoutesForProfileLinks( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profileLinkProviders,
		baseURI,
	)
	RegisterHTTPRoutesForProfileResources( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profileLinkProviders,
	)
	RegisterHTTPRoutesForProfileMemberships( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
	)
	RegisterHTTPRoutesForProfileQuestions( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profileQuestionsService,
	)

	var telegramServiceForEnvelopes *telegrambiz.Service
	if telegramProviders != nil {
		telegramServiceForEnvelopes = telegramProviders.Service
	}

	RegisterHTTPRoutesForProfileEnvelopes( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profileEnvelopesService,
		telegramServiceForEnvelopes,
	)
	RegisterHTTPRoutesForMailbox( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profileEnvelopesService,
	)
	RegisterHTTPRoutesForActivities( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		storyService,
	)
	RegisterHTTPRoutesForStoryInteractions( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		storyService,
		storyInteractionService,
	)
	RegisterHTTPRoutesForAdminProfiles( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		profilePointsService,
	)
	RegisterHTTPRoutesForAdminPoints( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profilePointsService,
	)
	RegisterHTTPRoutesForAdminWorkers( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		runtimeStatesService,
		workerRegistry,
	)

	if telegramProviders != nil {
		RegisterHTTPRoutesForTelegram( //nolint:contextcheck
			routes,
			logger,
			authService,
			userService,
			profileService,
			telegramProviders,
		)
	}

	// run
	return httpService.Start(ctx) //nolint:wrapcheck
}
