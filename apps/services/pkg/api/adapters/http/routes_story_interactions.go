package http

import (
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/story_interactions"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForStoryInteractions( //nolint:funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	storyService *stories.Service,
	interactionService *story_interactions.Service,
) {
	// Set interaction (authenticated)
	routes.Route(
		"POST /{locale}/stories/{slug}/interactions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			var requestBody struct {
				Kind string `json:"kind"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if requestBody.Kind == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("kind is required"))
			}

			// Resolve story
			story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			if story == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
			}

			// Get user's individual profile
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

			interactionKind := story_interactions.InteractionKind(requestBody.Kind)

			var interaction *story_interactions.StoryInteraction

			// Use SetRSVP for RSVP kinds (enforces mutual exclusivity)
			if story_interactions.IsRSVPKind(interactionKind) {
				interaction, err = interactionService.SetRSVP(
					ctx.Request.Context(),
					story.ID,
					*user.IndividualProfileID,
					interactionKind,
				)
			} else {
				interaction, err = interactionService.SetInteraction(
					ctx.Request.Context(),
					story.ID,
					*user.IndividualProfileID,
					requestBody.Kind,
				)
			}

			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(cursors.WrapResponseWithCursor(interaction, nil))
		},
	).
		HasSummary("Set story interaction").
		HasDescription("Set a story interaction (RSVP, like, etc.)").
		HasResponse(http.StatusOK)

	// Remove interaction (authenticated)
	routes.Route(
		"DELETE /{locale}/stories/{slug}/interactions/{kind}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			kindParam := ctx.Request.PathValue("kind")

			// Resolve story
			story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			if story == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
			}

			// Get user's individual profile
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

			err = interactionService.RemoveInteraction(
				ctx.Request.Context(),
				story.ID,
				*user.IndividualProfileID,
				kindParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.Ok()
		},
	).
		HasSummary("Remove story interaction").
		HasDescription("Remove a specific interaction for the current user.").
		HasResponse(http.StatusOK)

	// List interactions with profiles (public)
	routes.
		Route(
			"GET /{locale}/stories/{slug}/interactions",
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				slugParam := ctx.Request.PathValue("slug")
				filterKindParam := ctx.Request.URL.Query().Get("kind")

				// Resolve story
				story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				if story == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
				}

				var filterKind *string
				if filterKindParam != "" {
					filterKind = &filterKindParam
				}

				interactions, err := interactionService.ListInteractions(
					ctx.Request.Context(),
					localeParam,
					story.ID,
					filterKind,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(cursors.WrapResponseWithCursor(interactions, nil))
			},
		).
		HasSummary("List story interactions").
		HasDescription("List interactions for a story with profile info.").
		HasResponse(http.StatusOK)

	// Get my interactions (authenticated)
	routes.Route(
		"GET /{locale}/stories/{slug}/interactions/me",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Resolve story
			story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			if story == nil {
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
				return ctx.Results.JSON(
					cursors.WrapResponseWithCursor(
						[]*story_interactions.StoryInteraction{},
						nil,
					),
				)
			}

			interactions, err := interactionService.ListForProfile(
				ctx.Request.Context(),
				story.ID,
				*user.IndividualProfileID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(cursors.WrapResponseWithCursor(interactions, nil))
		},
	).
		HasSummary("Get my story interactions").
		HasDescription("Get current user's interactions for a story.").
		HasResponse(http.StatusOK)

	// Get interaction counts (public)
	routes.
		Route(
			"GET /{locale}/stories/{slug}/interactions/counts",
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}
				slugParam := ctx.Request.PathValue("slug")

				// Resolve story
				story, err := storyService.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				if story == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("story not found"))
				}

				counts, err := interactionService.CountInteractions(
					ctx.Request.Context(),
					story.ID,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(cursors.WrapResponseWithCursor(counts, nil))
			},
		).
		HasSummary("Get interaction counts").
		HasDescription("Get interaction counts grouped by kind for a story.").
		HasResponse(http.StatusOK)
}
