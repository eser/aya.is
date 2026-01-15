package http

import (
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForSite(
	routes *httpfx.Router,
	logger *logfx.Logger,
	profileService *profiles.Service,
) {
	routes.
		Route(
			"GET /{locale}/site/custom-domains/{domain}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				domainParam := ctx.Request.PathValue("domain")

				records, err := profileService.GetByCustomDomain(
					ctx.Request.Context(),
					localeParam,
					domainParam,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("Get profile by a custom domain").
		HasDescription("Get profile by a custom domain.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/site/spotlight", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")

			records, err := profileService.List(
				ctx.Request.Context(),
				localeParam,
				cursors.NewCursor(0, nil),
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			return ctx.Results.JSON(records)
		}).
		HasSummary("Gets spotlight metadata").
		HasDescription("Gets spotlight metadata.").
		HasResponse(http.StatusOK)
}
