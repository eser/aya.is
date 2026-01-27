package http

import (
	"net/http"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func RegisterHTTPRoutesForAdminProfiles(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
) {
	// List all profiles (admin only)
	routes.
		Route(
			"GET /admin/profiles",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				// Check admin permission
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithPlainText(err.Error()))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithPlainText("Admin access required"),
					)
				}

				// Parse query parameters
				query := ctx.Request.URL.Query()
				locale := query.Get("locale")
				if locale == "" {
					locale = "en"
				}

				filterKind := query.Get("kind")

				limitStr := query.Get("limit")
				limit := 50
				if limitStr != "" {
					if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 &&
						parsed <= 100 {
						limit = parsed
					}
				}

				offsetStr := query.Get("offset")
				offset := 0
				if offsetStr != "" {
					if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
						offset = parsed
					}
				}

				// Call service method
				result, err := profileService.ListAllProfilesForAdmin(
					ctx.Request.Context(),
					locale,
					filterKind,
					limit,
					offset,
				)
				if err != nil {
					logger.Error(
						"failed to list profiles for admin",
						"error", err,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				return ctx.Results.JSON(result)
			},
		).
		HasSummary("List all profiles").
		HasDescription("List all profiles with pagination. Admin only.").
		HasResponse(http.StatusOK)

	// Get single profile by slug (admin only)
	routes.
		Route(
			"GET /admin/profiles/{slug}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithPlainText(err.Error()))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithPlainText("Admin access required"),
					)
				}

				slug := ctx.Request.PathValue("slug")
				if slug == "" {
					return ctx.Results.BadRequest(
						httpfx.WithPlainText("slug is required"),
					)
				}

				query := ctx.Request.URL.Query()
				locale := query.Get("locale")
				if locale == "" {
					locale = "en"
				}

				profile, err := profileService.GetAdminProfileBySlug(
					ctx.Request.Context(),
					locale,
					slug,
				)
				if err != nil {
					logger.Error(
						"failed to get admin profile",
						"error", err,
						"slug", slug,
					)

					return ctx.Results.NotFound(
						httpfx.WithPlainText("profile not found"),
					)
				}

				return ctx.Results.JSON(profile)
			},
		).
		HasSummary("Get profile by slug").
		HasDescription("Get a single profile by slug. Admin only.").
		HasResponse(http.StatusOK)
}
