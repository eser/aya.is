package http

import (
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/story_series"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForStorySeries( //nolint:funlen,cyclop,gocognit,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	storySeriesService *story_series.Service,
	storyService *stories.Service,
) {
	// List all series (public)
	routes.
		Route("GET /{locale}/series", func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			seriesList, err := storySeriesService.List(ctx.Request.Context(), localeParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(seriesList, nil)

			return ctx.Results.JSON(wrappedResponse)
		})

	// Get series by slug with its stories (public)
	routes.
		Route("GET /{locale}/series/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			series, err := storySeriesService.GetBySlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			if series == nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Series not found"),
				)
			}

			seriesStories, err := storyService.ListBySeriesID(
				ctx.Request.Context(),
				localeParam,
				series.ID,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(map[string]any{
				"series":  series,
				"stories": seriesStories,
			}, nil)

			return ctx.Results.JSON(wrappedResponse)
		})

	// Create series (admin only)
	routes.
		Route(
			"POST /{locale}/series",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != userKindAdmin {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				var requestBody struct {
					Slug             string  `json:"slug"`
					SeriesPictureURI *string `json:"series_picture_uri"`
					Title            string  `json:"title"`
					Description      string  `json:"description"`
				}

				parseErr := ctx.ParseJSONBody(&requestBody)
				if parseErr != nil {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
				}

				requestBody.Slug = strings.TrimSpace(requestBody.Slug)
				requestBody.Title = strings.TrimSpace(requestBody.Title)

				if requestBody.Slug == "" || requestBody.Title == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("slug and title are required"),
					)
				}

				series, createErr := storySeriesService.Create(
					ctx.Request.Context(),
					story_series.CreateParams{
						Slug:             requestBody.Slug,
						SeriesPictureURI: requestBody.SeriesPictureURI,
						LocaleCode:       localeParam,
						Title:            requestBody.Title,
						Description:      requestBody.Description,
					},
				)
				if createErr != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(createErr),
					)
				}

				wrappedResponse := cursors.WrapResponseWithCursor(series, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		)

	// Update series base fields (admin only)
	routes.
		Route(
			"PATCH /{locale}/series/{id}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != userKindAdmin {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				idParam := ctx.Request.PathValue("id")

				var requestBody struct {
					Slug             string  `json:"slug"`
					SeriesPictureURI *string `json:"series_picture_uri"`
				}

				parseErr := ctx.ParseJSONBody(&requestBody)
				if parseErr != nil {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
				}

				requestBody.Slug = strings.TrimSpace(requestBody.Slug)
				if requestBody.Slug == "" {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("slug is required"))
				}

				updateErr := storySeriesService.Update(
					ctx.Request.Context(),
					idParam,
					story_series.UpdateParams{
						Slug:             requestBody.Slug,
						SeriesPictureURI: requestBody.SeriesPictureURI,
					},
				)
				if updateErr != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(updateErr),
					)
				}

				return ctx.Results.Ok()
			},
		)

	// Upsert series translation (admin only)
	routes.
		Route(
			"PUT /{locale}/series/{id}/translations/{txLocale}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != userKindAdmin {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				idParam := ctx.Request.PathValue("id")
				txLocaleParam := ctx.Request.PathValue("txLocale")

				var requestBody struct {
					Title       string `json:"title"`
					Description string `json:"description"`
				}

				parseErr := ctx.ParseJSONBody(&requestBody)
				if parseErr != nil {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
				}

				requestBody.Title = strings.TrimSpace(requestBody.Title)
				if requestBody.Title == "" {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("title is required"))
				}

				upsertErr := storySeriesService.UpsertTranslation(
					ctx.Request.Context(),
					idParam,
					story_series.TranslationParams{
						LocaleCode:  txLocaleParam,
						Title:       requestBody.Title,
						Description: requestBody.Description,
					},
				)
				if upsertErr != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(upsertErr),
					)
				}

				return ctx.Results.Ok()
			},
		)

	// Delete series (admin only)
	routes.
		Route(
			"DELETE /{locale}/series/{id}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != userKindAdmin {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				idParam := ctx.Request.PathValue("id")

				deleteErr := storySeriesService.Delete(ctx.Request.Context(), idParam)
				if deleteErr != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(deleteErr),
					)
				}

				return ctx.Results.Ok()
			},
		)
}
