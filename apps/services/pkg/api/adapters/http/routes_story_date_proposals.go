package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/story_date_proposals"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForStoryDateProposals( //nolint:funlen,gocognit,gocyclo,cyclop,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	storyService *stories.Service,
	dateProposalService *story_date_proposals.Service,
) {
	// List proposals (public, viewer vote included if authenticated)
	routes.
		Route(
			"GET /{locale}/stories/{slug}/date-proposals",
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				slugParam := ctx.Request.PathValue("slug")

				story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
				if err != nil || story == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
				}

				// Optionally extract viewer profile for vote display
				var viewerProfileID *string

				sessionID := GetSessionIDFromRequest(ctx.Request, authService)
				if sessionID != "" {
					session, sessionErr := userService.GetSessionByID(
						ctx.Request.Context(),
						sessionID,
					)
					if sessionErr == nil && session != nil && session.LoggedInUserID != nil {
						user, userErr := userService.GetByID(
							ctx.Request.Context(),
							*session.LoggedInUserID,
						)
						if userErr == nil && user != nil && user.IndividualProfileID != nil {
							viewerProfileID = user.IndividualProfileID
						}
					}
				}

				listResponse, err := dateProposalService.ListProposals(
					ctx.Request.Context(),
					localeParam,
					story.ID,
					viewerProfileID,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(cursors.WrapResponseWithCursor(listResponse, nil))
			},
		).
		HasSummary("List date proposals").
		HasDescription("List date proposals for an activity.").
		HasResponse(http.StatusOK)

	// Create proposal (authenticated)
	routes.Route(
		"POST /{locale}/stories/{slug}/date-proposals",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			var requestBody struct {
				DatetimeStart string  `json:"datetime_start"`
				DatetimeEnd   *string `json:"datetime_end"`
			}

			err := ctx.ParseJSONBody(&requestBody)
			if err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.DatetimeStart == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("datetime_start is required"))
			}

			datetimeStart, err := time.Parse(time.RFC3339, requestBody.DatetimeStart)
			if err != nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("datetime_start must be a valid RFC3339 timestamp"),
				)
			}

			var datetimeEnd *time.Time
			if requestBody.DatetimeEnd != nil && *requestBody.DatetimeEnd != "" {
				parsed, parseErr := time.Parse(time.RFC3339, *requestBody.DatetimeEnd)
				if parseErr != nil {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("datetime_end must be a valid RFC3339 timestamp"),
					)
				}

				datetimeEnd = &parsed
			}

			story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil || story == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil || session == nil || session.LoggedInUserID == nil {
				return ctx.Results.Error(
					http.StatusUnauthorized,
					httpfx.WithErrorMessage("Not authenticated"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil || user == nil || user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("User has no individual profile"),
				)
			}

			proposal, err := dateProposalService.CreateProposal(
				ctx.Request.Context(),
				story.ID,
				*user.IndividualProfileID,
				datetimeStart,
				datetimeEnd,
			)
			if err != nil {
				return mapDateProposalError(ctx, err)
			}

			return ctx.Results.JSON(cursors.WrapResponseWithCursor(proposal, nil))
		},
	).
		HasSummary("Create date proposal").
		HasDescription("Propose a date for an activity with undecided date.").
		HasResponse(http.StatusOK)

	// Delete proposal (authenticated)
	routes.Route(
		"DELETE /{locale}/stories/{slug}/date-proposals/{proposalId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")
			proposalID := ctx.Request.PathValue("proposalId")

			story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil || story == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil || session == nil || session.LoggedInUserID == nil {
				return ctx.Results.Error(
					http.StatusUnauthorized,
					httpfx.WithErrorMessage("Not authenticated"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil || user == nil || user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("User has no individual profile"),
				)
			}

			// Try as own proposal first
			err = dateProposalService.RemoveProposal(
				ctx.Request.Context(),
				proposalID,
				*user.IndividualProfileID,
			)
			if err != nil { //nolint:nestif
				// If unauthorized, check if user can edit the story (admin/maintainer+)
				if errors.Is(err, story_date_proposals.ErrUnauthorized) {
					canEdit, editErr := storyService.CanUserEditStory(
						ctx.Request.Context(),
						*session.LoggedInUserID,
						story.ID,
					)
					if editErr == nil && canEdit {
						adminErr := dateProposalService.RemoveProposalAsAdmin(
							ctx.Request.Context(),
							proposalID,
							*user.IndividualProfileID,
						)
						if adminErr != nil {
							return mapDateProposalError(ctx, adminErr)
						}

						return ctx.Results.Ok()
					}
				}

				return mapDateProposalError(ctx, err)
			}

			return ctx.Results.Ok()
		},
	).
		HasSummary("Remove date proposal").
		HasDescription("Remove a date proposal (proposer or maintainer+).").
		HasResponse(http.StatusOK)

	// Vote on proposal (authenticated)
	routes.Route(
		"POST /{locale}/stories/{slug}/date-proposals/{proposalId}/vote",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			proposalID := ctx.Request.PathValue("proposalId")

			var requestBody struct {
				Direction int `json:"direction"`
			}

			err := ctx.ParseJSONBody(&requestBody)
			if err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil || session == nil || session.LoggedInUserID == nil {
				return ctx.Results.Error(
					http.StatusUnauthorized,
					httpfx.WithErrorMessage("Not authenticated"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil || user == nil || user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("User has no individual profile"),
				)
			}

			voteResponse, err := dateProposalService.Vote(
				ctx.Request.Context(),
				proposalID,
				*user.IndividualProfileID,
				story_date_proposals.VoteDirection(requestBody.Direction),
			)
			if err != nil {
				return mapDateProposalError(ctx, err)
			}

			return ctx.Results.JSON(cursors.WrapResponseWithCursor(voteResponse, nil))
		},
	).
		HasSummary("Vote on date proposal").
		HasDescription("Vote agree (+1) or disagree (-1) on a date proposal.").
		HasResponse(http.StatusOK)

	// Finalize proposal (maintainer+ / author / admin)
	routes.Route(
		"POST /{locale}/stories/{slug}/date-proposals/{proposalId}/finalize",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")
			proposalID := ctx.Request.PathValue("proposalId")

			story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil || story == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found"),
				)
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil || session == nil || session.LoggedInUserID == nil {
				return ctx.Results.Error(
					http.StatusUnauthorized,
					httpfx.WithErrorMessage("Not authenticated"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil || user == nil || user.IndividualProfileID == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("User has no individual profile"),
				)
			}

			// Check permission: must be admin, story author, or maintainer+
			canEdit, editErr := storyService.CanUserEditStory(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				story.ID,
			)
			if editErr != nil || !canEdit {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithErrorMessage("You don't have permission to finalize the date"),
				)
			}

			err = dateProposalService.FinalizeProposal(
				ctx.Request.Context(),
				proposalID,
				story.ID,
				*user.IndividualProfileID,
			)
			if err != nil {
				return mapDateProposalError(ctx, err)
			}

			return ctx.Results.Ok()
		},
	).
		HasSummary("Finalize date proposal").
		HasDescription("Finalize a date proposal (maintainer+ / author / admin).").
		HasResponse(http.StatusOK)
}

func mapDateProposalError(ctx *httpfx.Context, err error) httpfx.Result {
	switch {
	case errors.Is(err, stories.ErrNotActivity):
		return ctx.Results.BadRequest(httpfx.WithErrorMessage("Story is not an activity"))
	case errors.Is(err, story_date_proposals.ErrDateNotUndecided):
		return ctx.Results.BadRequest(
			httpfx.WithErrorMessage("Activity date mode is not undecided"),
		)
	case errors.Is(err, story_date_proposals.ErrDateAlreadyFinalized):
		return ctx.Results.BadRequest(
			httpfx.WithErrorMessage("Activity date has already been finalized"),
		)
	case errors.Is(err, story_date_proposals.ErrCannotRemoveFinalized):
		return ctx.Results.BadRequest(httpfx.WithErrorMessage("Cannot remove a finalized proposal"))
	case errors.Is(err, story_date_proposals.ErrProposalNotFound):
		return ctx.Results.NotFound(httpfx.WithErrorMessage("Date proposal not found"))
	case errors.Is(err, story_date_proposals.ErrUnauthorized):
		return ctx.Results.Error(http.StatusForbidden, httpfx.WithErrorMessage("Unauthorized"))
	case errors.Is(err, story_date_proposals.ErrInvalidVoteDirection):
		return ctx.Results.BadRequest(httpfx.WithErrorMessage("Vote direction must be +1 or -1"))
	default:
		return ctx.Results.Error(http.StatusInternalServerError, httpfx.WithSanitizedError(err))
	}
}
