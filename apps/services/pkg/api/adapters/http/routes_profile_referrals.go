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

// RegisterHTTPRoutesForProfileReferrals registers routes for managing profile membership referrals.
func RegisterHTTPRoutesForProfileReferrals(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
) {
	// List referrals for a profile (member+ only)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_referrals",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

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

			referrals, err := profileService.ListReferrals(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list referrals",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			if referrals == nil {
				referrals = []*profiles.ProfileMembershipReferral{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  referrals,
				"error": nil,
			})
		},
	).HasDescription("List profile membership referrals")

	// Create a new referral (member+ only)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_referrals",
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
				ReferredProfileSlug string   `json:"referred_profile_slug"`
				TeamIDs             []string `json:"team_ids"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			if input.ReferredProfileSlug == "" {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("referred_profile_slug is required"),
				)
			}

			referral, err := profileService.CreateReferral(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				input.ReferredProfileSlug,
				input.TeamIDs,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create referral",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrReferralAlreadyExists):
					statusCode = http.StatusConflict
				case errors.Is(err, profiles.ErrCannotReferSelf):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrCannotReferExistingMember):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrProfileNotFound):
					statusCode = http.StatusNotFound
				case errors.Is(err, profiles.ErrInvalidInput):
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  referral,
				"error": nil,
			})
		},
	).HasDescription("Create a membership referral")

	// Get votes for a referral (member+ only)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_referrals/{id}/votes",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			idParam := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			votes, err := profileService.GetReferralVotes(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				slugParam,
				idParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to get referral votes",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("referralId", idParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				} else if errors.Is(err, profiles.ErrReferralNotFound) {
					statusCode = http.StatusNotFound
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			if votes == nil {
				votes = []*profiles.ReferralVote{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  votes,
				"error": nil,
			})
		},
	).HasDescription("Get votes for a referral")

	// Cast or update a vote on a referral (member+ only)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_referrals/{id}/votes",
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
			idParam := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				Score   int16   `json:"score"`
				Comment *string `json:"comment"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&input); err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			vote, err := profileService.VoteOnReferral(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				idParam,
				input.Score,
				input.Comment,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to vote on referral",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("referralId", idParam))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrReferralNotFound):
					statusCode = http.StatusNotFound
				case errors.Is(err, profiles.ErrInvalidVoteScore):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrReferralNotVoting):
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  vote,
				"error": nil,
			})
		},
	).HasDescription("Cast or update a vote on a referral")
}
