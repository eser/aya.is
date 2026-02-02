package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForProfileMemberships registers the routes for managing profile memberships.
func RegisterHTTPRoutesForProfileMemberships(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
) {
	// List memberships for profile settings
	routes.Route(
		"GET /{locale}/profiles/{slug}/_memberships",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get user information"),
				)
			}

			memberships, err := profileService.ListMembershipsForSettings(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list memberships",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to list memberships"),
				)
			}

			return ctx.Results.Ok(httpfx.WithJSON(memberships))
		},
	).HasDescription("List profile memberships for settings page")

	// Search users for adding as members
	routes.Route(
		"GET /{locale}/profiles/{slug}/_memberships/search",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")
			query := ctx.Request.URL.Query().Get("q")

			if query == "" {
				return ctx.Results.Ok(httpfx.WithJSON([]any{}))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get user information"),
				)
			}

			results, err := profileService.SearchUsersForMembership(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				query,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to search users",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("query", query))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to search users"),
				)
			}

			return ctx.Results.Ok(httpfx.WithJSON(results))
		},
	).HasDescription("Search users for adding as profile members")

	// Add new membership
	routes.Route(
		"POST /{locale}/profiles/{slug}/_memberships",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get user information"),
				)
			}

			var input struct {
				MemberProfileID string `json:"member_profile_id"`
				Kind            string `json:"kind"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			if input.MemberProfileID == "" || input.Kind == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("member_profile_id and kind are required"),
				)
			}

			err := profileService.AddMembership(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				input.MemberProfileID,
				input.Kind,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to add membership",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			return ctx.Results.Ok(httpfx.WithJSON(map[string]string{"status": "ok"}))
		},
	).HasDescription("Add a new membership to a profile")

	// Update membership kind
	routes.Route(
		"PUT /{locale}/profiles/{slug}/_memberships/{id}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			membershipID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get user information"),
				)
			}

			var input struct {
				Kind string `json:"kind"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			if input.Kind == "" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("kind is required"))
			}

			err := profileService.UpdateMembership(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				user.IndividualProfileID,
				slugParam,
				membershipID,
				input.Kind,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to update membership",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("membershipID", membershipID))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrCannotRemoveLastOwner) {
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithPlainText(err.Error()))
			}

			return ctx.Results.Ok(httpfx.WithJSON(map[string]string{"status": "ok"}))
		},
	).HasDescription("Update a membership's access level")

	// Delete membership
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_memberships/{id}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			membershipID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get user information"),
				)
			}

			err := profileService.DeleteMembership(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				user.IndividualProfileID,
				slugParam,
				membershipID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to delete membership",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("membershipID", membershipID))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrCannotRemoveLastOwner) ||
					errors.Is(err, profiles.ErrCannotRemoveIndividualSelf) {
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithPlainText(err.Error()))
			}

			return ctx.Results.Ok(httpfx.WithJSON(map[string]string{"status": "ok"}))
		},
	).HasDescription("Delete a membership from a profile")
}
