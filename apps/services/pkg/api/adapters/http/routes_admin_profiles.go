package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func RegisterHTTPRoutesForAdminProfiles(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	profilePointsService *profile_points.Service,
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
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
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
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(map[string]any{
					"data":  result,
					"error": nil,
				})
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
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				slug := ctx.Request.PathValue("slug")
				if slug == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("slug is required"),
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
						httpfx.WithErrorMessage("profile not found"),
					)
				}

				return ctx.Results.JSON(map[string]any{
					"data":  profile,
					"error": nil,
				})
			},
		).
		HasSummary("Get profile by slug").
		HasDescription("Get a single profile by slug. Admin only.").
		HasResponse(http.StatusOK)

	// Add points to a profile (admin only)
	routes.
		Route(
			"POST /admin/profiles/{slug}/points",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				slug := ctx.Request.PathValue("slug")
				if slug == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("slug is required"),
					)
				}

				// Parse request body
				var body struct {
					Amount      uint64 `json:"amount"`
					Description string `json:"description"`
				}

				if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("invalid request body"),
					)
				}

				if body.Amount == 0 {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("amount must be greater than 0"),
					)
				}

				if body.Description == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("description is required"),
					)
				}

				// Get profile to get the ID
				profile, err := profileService.GetBySlug(
					ctx.Request.Context(),
					"en",
					slug,
				)
				if err != nil {
					logger.Error(
						"failed to get profile for awarding points",
						"error", err,
						"slug", slug,
					)

					return ctx.Results.NotFound(
						httpfx.WithErrorMessage("profile not found"),
					)
				}

				if profile == nil {
					return ctx.Results.NotFound(
						httpfx.WithErrorMessage("profile not found"),
					)
				}

				// Add points using GainPoints (direct admin award)
				triggeringEvent := "ADMIN_AWARD"
				tx, err := profilePointsService.GainPoints(
					ctx.Request.Context(),
					profile_points.GainParams{
						ActorID:         user.ID,
						TargetProfileID: profile.ID,
						Amount:          body.Amount,
						TriggeringEvent: &triggeringEvent,
						Description:     body.Description,
					},
				)
				if err != nil {
					logger.Error(
						"failed to add points",
						"error", err,
						"slug", slug,
						"amount", body.Amount,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(map[string]any{
					"data":  tx,
					"error": nil,
				})
			},
		).
		HasSummary("Add points to profile").
		HasDescription("Add points directly to a profile. Admin only.").
		HasResponse(http.StatusOK)
}
