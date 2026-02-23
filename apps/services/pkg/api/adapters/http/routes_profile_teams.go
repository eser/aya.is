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

// RegisterHTTPRoutesForProfileTeams registers routes for managing profile teams.
func RegisterHTTPRoutesForProfileTeams(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
) {
	// List teams with member counts
	routes.Route(
		"GET /{locale}/profiles/{slug}/_teams",
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

			teams, err := profileService.ListTeams(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list teams",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			if teams == nil {
				teams = []*profiles.ProfileTeam{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  teams,
				"error": nil,
			})
		},
	).HasDescription("List profile teams with member counts")

	// Create team
	routes.Route(
		"POST /{locale}/profiles/{slug}/_teams",
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

			var input struct {
				Name        string  `json:"name"`
				Description *string `json:"description"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if input.Name == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("name is required"))
			}

			team, err := profileService.CreateTeam(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				input.Name,
				input.Description,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create team",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  team,
				"error": nil,
			})
		},
	).HasDescription("Create a new team for a profile")

	// Update team
	routes.Route(
		"PUT /{locale}/profiles/{slug}/_teams/{id}",
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
			teamID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				Name        string  `json:"name"`
				Description *string `json:"description"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if input.Name == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("name is required"))
			}

			err := profileService.UpdateTeam(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				teamID,
				input.Name,
				input.Description,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to update team",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("teamID", teamID))

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
	).HasDescription("Update a team's name and description")

	// Delete team
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_teams/{id}",
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
			teamID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			err := profileService.DeleteTeam(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				teamID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to delete team",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("teamID", teamID))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrCannotDeleteTeamWithMembers) {
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "ok"},
				"error": nil,
			})
		},
	).HasDescription("Delete a team from a profile")

	// Set membership teams
	routes.Route(
		"PUT /{locale}/profiles/{slug}/_memberships/{id}/teams",
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

			var input struct {
				TeamIDs []string `json:"team_ids"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			err := profileService.SetMembershipTeams(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				membershipID,
				input.TeamIDs,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to set membership teams",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("membershipID", membershipID))

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
	).HasDescription("Set teams for a membership")
}
