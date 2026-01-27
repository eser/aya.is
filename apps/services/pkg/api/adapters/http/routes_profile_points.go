package http

import (
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForProfilePoints(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	profilePointsService *profile_points.Service,
) {
	routes.
		Route(
			"GET /{locale}/profiles/{slug}/_points/transactions",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")

				if slugParam == "" {
					return ctx.Results.BadRequest(
						httpfx.WithPlainText("slug parameter is required"),
					)
				}

				// Get session ID from context (set by auth middleware)
				sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
				if !ok {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText("Session ID not found in context"),
					)
				}

				// Get session to get user ID
				session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
				if err != nil || session == nil || session.LoggedInUserID == nil {
					return ctx.Results.Unauthorized(httpfx.WithPlainText("Invalid session"))
				}

				// Get user to check if admin
				user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
				if userErr != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText("Failed to get user information"),
					)
				}

				// Admin users can edit any profile
				canEdit := false
				if user.Kind == "admin" {
					canEdit = true
				} else {
					// Check normal permissions
					var err error
					canEdit, err = profileService.CanUserEditProfile(
						ctx.Request.Context(),
						*session.LoggedInUserID,
						slugParam,
					)
					if err != nil {
						return ctx.Results.Error(
							http.StatusInternalServerError,
							httpfx.WithPlainText(err.Error()),
						)
					}
				}

				if !canEdit {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithPlainText(
							"You do not have permission to view this profile's transactions",
						),
					)
				}

				// Get profile to get the ID
				profile, err := profileService.GetBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				if profile == nil {
					return ctx.Results.NotFound(httpfx.WithPlainText("profile not found"))
				}

				cursor := cursors.NewCursorFromRequest(ctx.Request)

				transactions, err := profilePointsService.ListTransactions(
					ctx.Request.Context(),
					profile.ID,
					cursor,
				)
				if err != nil {
					logger.Error(
						"failed to list profile point transactions",
						"error", err,
						"locale", localeParam,
						"slug", slugParam,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				return ctx.Results.JSON(transactions)
			},
		).
		HasSummary("List profile point transactions").
		HasDescription("List point transactions for a profile. Requires edit permission.").
		HasResponse(http.StatusOK)
}
