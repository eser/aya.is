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
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
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
					httpfx.WithErrorMessage("Failed to list memberships"),
				)
			}

			// Ensure we always return a JSON array, never 204 No Content
			if memberships == nil {
				memberships = []*profiles.ProfileMembershipWithMember{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  memberships,
				"error": nil,
			})
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
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			query := ctx.Request.URL.Query().Get("q")

			if query == "" {
				return ctx.Results.JSON(map[string]any{
					"data":  []any{},
					"error": nil,
				})
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
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
					httpfx.WithErrorMessage("Failed to search users"),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  results,
				"error": nil,
			})
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
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			var input struct {
				MemberProfileID string `json:"member_profile_id"`
				Kind            string `json:"kind"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if input.MemberProfileID == "" || input.Kind == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("member_profile_id and kind are required"),
				)
			}

			err := profileService.AddMembership(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				user.Kind,
				user.IndividualProfileID,
				slugParam,
				input.MemberProfileID,
				input.Kind,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to add membership",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrCannotAssignHigherRole) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "ok"},
				"error": nil,
			})
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
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			membershipID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			var input struct {
				Kind string `json:"kind"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if input.Kind == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("kind is required"))
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
				if errors.Is(err, profiles.ErrCannotRemoveLastOwner) ||
					errors.Is(err, profiles.ErrCannotModifyOwnRole) ||
					errors.Is(err, profiles.ErrCannotAssignHigherRole) ||
					errors.Is(err, profiles.ErrCannotModifyHigherMember) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "ok"},
				"error": nil,
			})
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
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			membershipID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
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

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "ok"},
				"error": nil,
			})
		},
	).HasDescription("Delete a membership from a profile")

	// Follow a profile (self-service)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_follow",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil || user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("User profile not found"),
				)
			}

			err := profileService.FollowProfile(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				*user.IndividualProfileID,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to follow profile",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "ok"},
				"error": nil,
			})
		},
	).HasDescription("Follow a profile")

	// Unfollow a profile (self-service)
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_follow",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil || user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("User profile not found"),
				)
			}

			err := profileService.UnfollowProfile(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				*user.IndividualProfileID,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to unfollow profile",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "ok"},
				"error": nil,
			})
		},
	).HasDescription("Unfollow a profile")
}
