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

func RegisterHTTPRoutesForProfileQuestions( //nolint:funlen,gocognit,cyclop,maintidx
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

			err = json.NewDecoder(ctx.Request.Body).Decode(&body)
			if err != nil {
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
	answerHandler := questionAnswerHandler(
		logger, userService, profileQuestionsService, "create",
	)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_questions/{id}/answer",
		AuthMiddleware(authService, userService),
		answerHandler,
	).HasDescription("Answer a Q&A question (contributor+ access)")

	// Edit an answer (requires contributor+ access)
	editAnswerHandler := questionAnswerHandler(
		logger, userService, profileQuestionsService, answerModeEdit,
	)
	routes.Route(
		"PUT /{locale}/profiles/{slug}/_questions/{id}/answer",
		AuthMiddleware(authService, userService),
		editAnswerHandler,
	).HasDescription("Edit an answer on a Q&A question (contributor+ access)")

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

			err = json.NewDecoder(ctx.Request.Body).Decode(&body)
			if err != nil {
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

const answerModeEdit = "edit"

// questionAnswerHandler creates a handler for creating or editing answers on Q&A questions.
// The mode parameter must be "create" or "edit".
func questionAnswerHandler( //nolint:funlen
	logger *logfx.Logger,
	userService *users.Service,
	profileQuestionsService *profile_questions.Service,
	mode string,
) func(ctx *httpfx.Context) httpfx.Result {
	noProfileMsg := "you must have an individual profile to answer questions"
	if mode == answerModeEdit {
		noProfileMsg = "you must have an individual profile to edit answers"
	}

	return func(ctx *httpfx.Context) httpfx.Result {
		user, err := getUserFromContext(ctx, userService)
		if err != nil {
			return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
		}

		if user.IndividualProfileID == nil {
			return ctx.Results.BadRequest(httpfx.WithErrorMessage(noProfileMsg))
		}

		slugParam := ctx.Request.PathValue("slug")
		questionID := ctx.Request.PathValue("id")

		var body struct {
			AnswerURI     *string `json:"answer_uri"`
			AnswerKind    *string `json:"answer_kind"`
			AnswerContent string  `json:"answer_content"`
		}

		decErr := json.NewDecoder(ctx.Request.Body).Decode(&body)
		if decErr != nil {
			return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
		}

		if body.AnswerContent == "" {
			return ctx.Results.BadRequest(httpfx.WithErrorMessage("answer_content is required"))
		}

		params := profile_questions.AnswerQuestionParams{
			ProfileSlug:       slugParam,
			QuestionID:        questionID,
			UserID:            user.ID,
			AnswererProfileID: *user.IndividualProfileID,
			AnswerContent:     body.AnswerContent,
			AnswerURI:         body.AnswerURI,
			AnswerKind:        body.AnswerKind,
		}

		if mode == answerModeEdit {
			err = profileQuestionsService.EditAnswer(ctx.Request.Context(), params)
		} else {
			err = profileQuestionsService.AnswerQuestion(ctx.Request.Context(), params)
		}

		if err != nil {
			return handleQuestionAnswerError(ctx, logger, err, mode, slugParam, questionID)
		}

		return ctx.Results.JSON(map[string]any{
			"data": map[string]string{"status": "ok"},
		})
	}
}

// handleQuestionAnswerError maps business errors from answer create/edit to HTTP responses.
func handleQuestionAnswerError(
	ctx *httpfx.Context,
	logger *logfx.Logger,
	err error,
	mode, slugParam, questionID string,
) httpfx.Result {
	if errors.Is(err, profile_questions.ErrInsufficientPermission) {
		return ctx.Results.Error(
			http.StatusForbidden,
			httpfx.WithErrorMessage("insufficient permission"),
		)
	}

	if errors.Is(err, profile_questions.ErrQuestionNotFound) {
		return ctx.Results.NotFound(httpfx.WithErrorMessage("question not found"))
	}

	if mode == "create" && errors.Is(err, profile_questions.ErrQuestionAlreadyAnswered) {
		return ctx.Results.Error(
			http.StatusConflict,
			httpfx.WithErrorMessage("question has already been answered"),
		)
	}

	if mode == answerModeEdit && errors.Is(err, profile_questions.ErrQuestionNotAnswered) {
		return ctx.Results.Error(
			http.StatusConflict,
			httpfx.WithErrorMessage("question has not been answered yet"),
		)
	}

	operation := "Failed to answer question"
	if mode == answerModeEdit {
		operation = "Failed to edit answer"
	}

	logger.ErrorContext(ctx.Request.Context(), operation,
		slog.String("error", err.Error()),
		slog.String("slug", slugParam),
		slog.String("questionID", questionID))

	return ctx.Results.Error(
		http.StatusInternalServerError,
		httpfx.WithSanitizedError(err),
	)
}
