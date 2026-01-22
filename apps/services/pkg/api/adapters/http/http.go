package http

import (
	"context"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/middlewares"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/healthcheck"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/openapi"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/profiling"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	mcpadapter "github.com/eser/aya.is/services/pkg/api/adapters/mcp"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/protection"
	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/uploads"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func Run(
	ctx context.Context,
	baseURI string,
	config *httpfx.Config,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	storyService *stories.Service,
	sessionService *sessions.Service,
	protectionService *protection.Service,
	uploadService *uploads.Service,
	profileLinkProviders *ProfileLinkProviders,
) (func(), error) {
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHTTPService(config, routes, logger)

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.TracingMiddleware(logger)) //nolint:contextcheck
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
	)
	RegisterHTTPRoutesForSessions( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		sessionService,
		protectionService,
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
	)
	RegisterHTTPRoutesForProfiles( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		profileService,
		storyService,
	)
	RegisterHTTPRoutesForStories( //nolint:contextcheck
		routes,
		logger,
		authService,
		userService,
		storyService,
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

	// run
	return httpService.Start(ctx) //nolint:wrapcheck
}
