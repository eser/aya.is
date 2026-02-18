package http

import (
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForActivities(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	storyService *stories.Service,
) {
	// List published activity stories
	routes.
		Route("GET /{locale}/activities", func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			records, err := storyService.ListActivities(ctx.Request.Context(), localeParam, nil)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List activity stories").
		HasDescription("List published activity stories sorted by activity_time_start.").
		HasResponse(http.StatusOK)

	// Get single activity by slug (reuses story GetBySlugForViewer)
	routes.
		Route("GET /{locale}/activities/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			// Optional viewer auth
			var viewerUserID *string

			sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
			if err == nil && sessionID != "" {
				session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
				if sessionErr == nil && session != nil && session.LoggedInUserID != nil {
					viewerUserID = session.LoggedInUserID
				}
			}

			record, err := storyService.GetBySlugForViewer(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				viewerUserID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			if record == nil || record.Kind != stories.KindActivity {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("activity not found"))
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get activity by slug").
		HasDescription("Get a single activity story by slug.").
		HasResponse(http.StatusOK)
}
