package http

import (
	"net/http"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForSearch(
	routes *httpfx.Router,
	profileService *profiles.Service,
) {
	routes.
		Route("GET /{locale}/search", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")

			// get query parameter
			query := ctx.Request.URL.Query().Get("q")
			if query == "" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("q parameter is required"))
			}

			// get limit parameter (default 20, max 100)
			limitStr := ctx.Request.URL.Query().Get("limit")
			limit := int32(20)
			if limitStr != "" {
				parsedLimit, err := strconv.Atoi(limitStr)
				if err == nil && parsedLimit > 0 && parsedLimit <= 100 {
					limit = int32(parsedLimit)
				}
			}

			// get profile parameter for scoped search (optional)
			var profileSlug *string
			profileParam := ctx.Request.URL.Query().Get("profile")
			if profileParam != "" {
				profileSlug = &profileParam
			}

			results, err := profileService.Search(
				ctx.Request.Context(),
				localeParam,
				query,
				profileSlug,
				limit,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(results, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Search across profiles, stories, and pages").
		HasDescription("Full-text search using PostgreSQL tsvector. Use 'profile' query param to scope search to a specific profile.").
		HasResponse(http.StatusOK)
}
