package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profile_questions"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func RegisterHTTPRoutesForProfileQuestions(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	profileQuestionsService *profile_questions.Service,
) {
	// List questions for a profile (public, optional auth for viewer vote state)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_questions",
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			// Optional: resolve viewer user ID from session cookie or Bearer token.
			// Cookie works on same-site (aya.is), Bearer token works on custom domains (eser.dev).
			var viewerUserID *string

			sessionID := GetSessionIDFromRequest(ctx.Request, authService)
			if sessionID != "" {
				session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
				if sessionErr == nil && session != nil && session.LoggedInUserID != nil {
					viewerUserID = session.LoggedInUserID
				}
			}

			// Show hidden questions to maintainer+ users
			includeHidden := false

			if viewerUserID != nil {
				hasAccess, accessErr := profileService.HasUserAccessToProfile(
					ctx.Request.Context(),
					*viewerUserID,
					slugParam,
					profiles.MembershipKindMaintainer,
				)
				if accessErr == nil && hasAccess {
					includeHidden = true
				}
			}

			result, err := profileQuestionsService.ListQuestions(
				ctx.Request.Context(),
				slugParam,
				localeParam,
				viewerUserID,
				includeHidden,
				nil,
			)
			if err != nil {
				if errors.Is(err, profile_questions.ErrQANotEnabled) {
					return ctx.Results.Error(
						http.StatusNotFound,
						httpfx.WithErrorMessage("Q&A is not enabled for this profile"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Failed to list questions",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(result)
		},
	).HasDescription("List Q&A questions for a profile")

	// Create a new question (requires authentication and an individual profile)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_questions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			// Require the user to have an individual profile
			if user.IndividualProfileID == nil {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("you must have an individual profile to ask questions"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			var body struct {
				Content     string `json:"content"`
				IsAnonymous bool   `json:"is_anonymous"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			question, err := profileQuestionsService.CreateQuestion(
				ctx.Request.Context(),
				profile_questions.CreateQuestionParams{
					ProfileSlug: slugParam,
					UserID:      user.ID,
					Content:     body.Content,
					IsAnonymous: body.IsAnonymous,
				},
			)
			if err != nil {
				if errors.Is(err, profile_questions.ErrQANotEnabled) {
					return ctx.Results.Error(
						http.StatusNotFound,
						httpfx.WithErrorMessage("Q&A is not enabled for this profile"),
					)
				}

				if errors.Is(err, profile_questions.ErrContentTooShort) ||
					errors.Is(err, profile_questions.ErrContentTooLong) {
					return ctx.Results.BadRequest(httpfx.WithSanitizedError(err))
				}

				logger.ErrorContext(ctx.Request.Context(), "Failed to create question",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data": question,
			})
		},
	).HasDescription("Create a new Q&A question on a profile")

	// Toggle vote on a question (requires authentication)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_questions/{id}/vote",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			slugParam := ctx.Request.PathValue("slug")
			questionID := ctx.Request.PathValue("id")

			voted, err := profileQuestionsService.ToggleVote(
				ctx.Request.Context(),
				profile_questions.VoteParams{
					ProfileSlug: slugParam,
					QuestionID:  questionID,
					UserID:      user.ID,
				},
			)
			if err != nil {
				if errors.Is(err, profile_questions.ErrQANotEnabled) {
					return ctx.Results.Error(
						http.StatusNotFound,
						httpfx.WithErrorMessage("Q&A is not enabled for this profile"),
					)
				}

				if errors.Is(err, profile_questions.ErrQuestionNotFound) {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("question not found"))
				}

				logger.ErrorContext(ctx.Request.Context(), "Failed to toggle vote",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("questionID", questionID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"voted": voted,
				},
			})
		},
	).HasDescription("Toggle vote on a Q&A question")

	// Answer a question (requires contributor+ access)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_questions/{id}/answer",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			slugParam := ctx.Request.PathValue("slug")
			questionID := ctx.Request.PathValue("id")

			var body struct {
				AnswerContent string  `json:"answer_content"`
				AnswerURI     *string `json:"answer_uri"`
				AnswerKind    *string `json:"answer_kind"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			if body.AnswerContent == "" {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("answer_content is required"))
			}

			err = profileQuestionsService.AnswerQuestion(
				ctx.Request.Context(),
				profile_questions.AnswerQuestionParams{
					ProfileSlug:   slugParam,
					QuestionID:    questionID,
					UserID:        user.ID,
					AnswerContent: body.AnswerContent,
					AnswerURI:     body.AnswerURI,
					AnswerKind:    body.AnswerKind,
				},
			)
			if err != nil {
				if errors.Is(err, profile_questions.ErrInsufficientPermission) {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("insufficient permission"),
					)
				}

				if errors.Is(err, profile_questions.ErrQuestionNotFound) {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("question not found"))
				}

				if errors.Is(err, profile_questions.ErrQuestionAlreadyAnswered) {
					return ctx.Results.Error(
						http.StatusConflict,
						httpfx.WithErrorMessage("question has already been answered"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Failed to answer question",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("questionID", questionID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Answer a Q&A question (contributor+ access)")

	// Hide/unhide a question (requires maintainer+ access)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_questions/{id}/hide",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			slugParam := ctx.Request.PathValue("slug")
			questionID := ctx.Request.PathValue("id")

			var body struct {
				IsHidden bool `json:"is_hidden"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			err = profileQuestionsService.HideQuestion(
				ctx.Request.Context(),
				profile_questions.HideQuestionParams{
					ProfileSlug: slugParam,
					QuestionID:  questionID,
					UserID:      user.ID,
					IsHidden:    body.IsHidden,
				},
			)
			if err != nil {
				if errors.Is(err, profile_questions.ErrInsufficientPermission) {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("insufficient permission"),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Failed to hide question",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("questionID", questionID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Hide or unhide a Q&A question (maintainer+ access)")
}
