package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/discussions"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func RegisterHTTPRoutesForDiscussions( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	discussionsService *discussions.Service,
) {
	// List top-level comments for a story discussion.
	routes.Route(
		"GET /{locale}/stories/{slug}/_discussions",
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			viewerUserID := resolveViewerUserID(ctx, authService, userService)
			includeHidden := resolveIncludeHidden(ctx, viewerUserID, slugParam, profileService)

			params := buildListParams(ctx, localeParam, viewerUserID, includeHidden)

			thread, err := discussionsService.GetOrCreateThreadByStorySlug(
				ctx.Request.Context(),
				slugParam,
			)
			if err != nil {
				return handleDiscussionError(ctx, logger, err, "story", slugParam)
			}

			params.ThreadID = thread.ID

			comments, err := discussionsService.ListComments(ctx.Request.Context(), params)
			if err != nil {
				return handleDiscussionError(ctx, logger, err, "story", slugParam)
			}

			return ctx.Results.JSON(discussions.ListResponse{
				Thread:   thread,
				Comments: comments,
			})
		},
	).HasDescription("List discussion comments for a story")

	// List top-level comments for a profile discussion.
	routes.Route(
		"GET /{locale}/profiles/{slug}/_discussions",
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			viewerUserID := resolveViewerUserID(ctx, authService, userService)
			includeHidden := resolveIncludeHidden(ctx, viewerUserID, slugParam, profileService)

			params := buildListParams(ctx, localeParam, viewerUserID, includeHidden)

			thread, err := discussionsService.GetOrCreateThreadByProfileSlug(
				ctx.Request.Context(),
				slugParam,
			)
			if err != nil {
				return handleDiscussionError(ctx, logger, err, "profile", slugParam)
			}

			params.ThreadID = thread.ID

			comments, err := discussionsService.ListComments(ctx.Request.Context(), params)
			if err != nil {
				return handleDiscussionError(ctx, logger, err, "profile", slugParam)
			}

			return ctx.Results.JSON(discussions.ListResponse{
				Thread:   thread,
				Comments: comments,
			})
		},
	).HasDescription("List discussion comments for a profile")

	// List replies to a comment.
	routes.Route(
		"GET /{locale}/discussions/comments/{commentId}/replies",
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			commentID := ctx.Request.PathValue("commentId")

			// Get the parent comment to find the thread ID.
			parent, err := discussionsService.GetComment(
				ctx.Request.Context(),
				commentID,
				localeParam,
			)
			if err != nil || parent == nil {
				return ctx.Results.NotFound(httpfx.WithErrorMessage("comment not found"))
			}

			viewerUserID := resolveViewerUserID(ctx, authService, userService)

			params := buildListParams(ctx, localeParam, viewerUserID, false)
			params.ParentID = &commentID
			params.ThreadID = parent.ThreadID

			comments, err := discussionsService.ListComments(ctx.Request.Context(), params)
			if err != nil {
				return handleDiscussionError(ctx, logger, err, "replies", commentID)
			}

			return ctx.Results.JSON(map[string]any{
				"comments": comments,
			})
		},
	).HasDescription("List replies to a discussion comment")

	// Create a comment on a story discussion.
	routes.Route(
		"POST /{locale}/stories/{slug}/_discussions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			var body struct {
				Content  string  `json:"content"`
				ParentID *string `json:"parent_id"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			comment, err := discussionsService.CreateComment(
				ctx.Request.Context(),
				discussions.CreateCommentParams{
					StorySlug: &slugParam,
					Locale:    localeParam,
					UserID:    user.ID,
					ParentID:  body.ParentID,
					Content:   body.Content,
				},
			)
			if err != nil {
				return handleCreateCommentError(ctx, logger, err, "story", slugParam)
			}

			return ctx.Results.JSON(map[string]any{
				"data": comment,
			})
		},
	).HasDescription("Create a comment on a story discussion")

	// Create a comment on a profile discussion.
	routes.Route(
		"POST /{locale}/profiles/{slug}/_discussions",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			slugParam := ctx.Request.PathValue("slug")

			var body struct {
				Content  string  `json:"content"`
				ParentID *string `json:"parent_id"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			comment, err := discussionsService.CreateComment(
				ctx.Request.Context(),
				discussions.CreateCommentParams{
					ProfileSlug: &slugParam,
					Locale:      localeParam,
					UserID:      user.ID,
					ParentID:    body.ParentID,
					Content:     body.Content,
				},
			)
			if err != nil {
				return handleCreateCommentError(ctx, logger, err, "profile", slugParam)
			}

			return ctx.Results.JSON(map[string]any{
				"data": comment,
			})
		},
	).HasDescription("Create a comment on a profile discussion")

	// Edit a comment (author only).
	routes.Route(
		"PUT /{locale}/discussions/comments/{commentId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			commentID := ctx.Request.PathValue("commentId")

			var body struct {
				Content string `json:"content"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			err = discussionsService.EditComment(
				ctx.Request.Context(),
				discussions.EditCommentParams{
					CommentID: commentID,
					UserID:    user.ID,
					Content:   body.Content,
				},
			)
			if err != nil {
				return handleMutationError(ctx, logger, err, "edit", commentID)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Edit a discussion comment (author only)")

	// Delete a comment (author or contributor+).
	routes.Route(
		"DELETE /{locale}/discussions/comments/{commentId}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			commentID := ctx.Request.PathValue("commentId")

			query := ctx.Request.URL.Query()
			profileSlug := query.Get("profile_slug")

			err = discussionsService.DeleteComment(
				ctx.Request.Context(),
				discussions.DeleteCommentParams{
					CommentID:   commentID,
					UserID:      user.ID,
					ProfileSlug: profileSlug,
				},
			)
			if err != nil {
				return handleMutationError(ctx, logger, err, "delete", commentID)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Delete a discussion comment (author or contributor+)")

	// Vote on a comment.
	routes.Route(
		"POST /{locale}/discussions/comments/{commentId}/vote",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			commentID := ctx.Request.PathValue("commentId")

			var body struct {
				Direction int `json:"direction"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			result, err := discussionsService.Vote(
				ctx.Request.Context(),
				discussions.VoteParams{
					CommentID: commentID,
					UserID:    user.ID,
					Direction: discussions.VoteDirection(body.Direction),
				},
			)
			if err != nil {
				return handleMutationError(ctx, logger, err, "vote", commentID)
			}

			return ctx.Results.JSON(map[string]any{
				"data": result,
			})
		},
	).HasDescription("Vote on a discussion comment")

	// Hide/unhide a comment (contributor+ only).
	routes.Route(
		"POST /{locale}/discussions/comments/{commentId}/hide",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			commentID := ctx.Request.PathValue("commentId")

			var body struct {
				IsHidden    bool   `json:"is_hidden"`
				ProfileSlug string `json:"profile_slug"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			err = discussionsService.HideComment(
				ctx.Request.Context(),
				discussions.HideCommentParams{
					CommentID:   commentID,
					UserID:      user.ID,
					ProfileSlug: body.ProfileSlug,
					IsHidden:    body.IsHidden,
				},
			)
			if err != nil {
				return handleMutationError(ctx, logger, err, "hide", commentID)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Hide or unhide a discussion comment (contributor+)")

	// Pin/unpin a comment (contributor+ only).
	routes.Route(
		"POST /{locale}/discussions/comments/{commentId}/pin",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			commentID := ctx.Request.PathValue("commentId")

			var body struct {
				IsPinned    bool   `json:"is_pinned"`
				ProfileSlug string `json:"profile_slug"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			err = discussionsService.PinComment(
				ctx.Request.Context(),
				discussions.PinCommentParams{
					CommentID:   commentID,
					UserID:      user.ID,
					ProfileSlug: body.ProfileSlug,
					IsPinned:    body.IsPinned,
				},
			)
			if err != nil {
				return handleMutationError(ctx, logger, err, "pin", commentID)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Pin or unpin a discussion comment (contributor+)")

	// Lock/unlock a thread (contributor+ only).
	routes.Route(
		"POST /{locale}/discussions/threads/{threadId}/lock",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			user, err := getUserFromContext(ctx, userService)
			if err != nil {
				return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
			}

			threadID := ctx.Request.PathValue("threadId")

			var body struct {
				IsLocked    bool   `json:"is_locked"`
				ProfileSlug string `json:"profile_slug"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
			}

			err = discussionsService.LockThread(
				ctx.Request.Context(),
				discussions.LockThreadParams{
					ThreadID:    threadID,
					UserID:      user.ID,
					ProfileSlug: body.ProfileSlug,
					IsLocked:    body.IsLocked,
				},
			)
			if err != nil {
				return handleMutationError(ctx, logger, err, "lock", threadID)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]string{"status": "ok"},
			})
		},
	).HasDescription("Lock or unlock a discussion thread (contributor+)")
}

// resolveViewerUserID extracts the viewer user ID from the session cookie or Bearer token.
func resolveViewerUserID(
	ctx *httpfx.Context,
	authService *auth.Service,
	userService *users.Service,
) *string {
	sessionID := GetSessionIDFromRequest(ctx.Request, authService)
	if sessionID == "" {
		return nil
	}

	session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
	if sessionErr != nil || session == nil || session.LoggedInUserID == nil {
		return nil
	}

	return session.LoggedInUserID
}

// resolveIncludeHidden determines if hidden comments should be included (contributor+ access).
func resolveIncludeHidden(
	ctx *httpfx.Context,
	viewerUserID *string,
	profileSlug string,
	profileService *profiles.Service,
) bool {
	if viewerUserID == nil {
		return false
	}

	hasAccess, accessErr := profileService.HasUserAccessToProfile(
		ctx.Request.Context(),
		*viewerUserID,
		profileSlug,
		profiles.MembershipKindContributor,
	)

	return accessErr == nil && hasAccess
}

// buildListParams creates ListCommentsParams from query string.
func buildListParams(
	ctx *httpfx.Context,
	locale string,
	viewerUserID *string,
	includeHidden bool,
) discussions.ListCommentsParams {
	query := ctx.Request.URL.Query()

	limit := discussions.DefaultPageLimit

	if l := query.Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	offset := 0

	if o := query.Get("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	sort := discussions.SortMode(query.Get("sort"))

	return discussions.ListCommentsParams{
		Locale:        locale,
		ViewerUserID:  viewerUserID,
		IncludeHidden: includeHidden,
		Sort:          sort,
		Limit:         limit,
		Offset:        offset,
	}
}

// handleDiscussionError handles errors for list-type discussion endpoints.
func handleDiscussionError(
	ctx *httpfx.Context,
	logger *logfx.Logger,
	err error,
	entityType string,
	entitySlug string,
) httpfx.Result {
	if errors.Is(err, discussions.ErrDiscussionsNotEnabled) {
		return ctx.Results.Error(
			http.StatusNotFound,
			httpfx.WithErrorMessage("discussions are not enabled"),
		)
	}

	if errors.Is(err, discussions.ErrThreadNotFound) {
		return ctx.Results.NotFound(httpfx.WithErrorMessage("thread not found"))
	}

	logger.ErrorContext(ctx.Request.Context(), "Failed to list discussion comments",
		slog.String("error", err.Error()),
		slog.String("entity_type", entityType),
		slog.String("entity_slug", entitySlug))

	return ctx.Results.Error(
		http.StatusInternalServerError,
		httpfx.WithSanitizedError(err),
	)
}

// handleCreateCommentError handles errors for comment creation.
func handleCreateCommentError(
	ctx *httpfx.Context,
	logger *logfx.Logger,
	err error,
	entityType string,
	entitySlug string,
) httpfx.Result {
	if errors.Is(err, discussions.ErrDiscussionsNotEnabled) {
		return ctx.Results.Error(
			http.StatusNotFound,
			httpfx.WithErrorMessage("discussions are not enabled"),
		)
	}

	if errors.Is(err, discussions.ErrThreadLocked) {
		return ctx.Results.Error(
			http.StatusForbidden,
			httpfx.WithErrorMessage("this thread is locked"),
		)
	}

	if errors.Is(err, discussions.ErrContentTooShort) ||
		errors.Is(err, discussions.ErrContentTooLong) {
		return ctx.Results.BadRequest(httpfx.WithSanitizedError(err))
	}

	if errors.Is(err, discussions.ErrMaxNestingDepth) {
		return ctx.Results.BadRequest(httpfx.WithErrorMessage("maximum nesting depth reached"))
	}

	if errors.Is(err, discussions.ErrCommentNotFound) {
		return ctx.Results.NotFound(httpfx.WithErrorMessage("parent comment not found"))
	}

	logger.ErrorContext(ctx.Request.Context(), "Failed to create discussion comment",
		slog.String("error", err.Error()),
		slog.String("entity_type", entityType),
		slog.String("entity_slug", entitySlug))

	return ctx.Results.Error(
		http.StatusInternalServerError,
		httpfx.WithSanitizedError(err),
	)
}

// handleMutationError handles errors for comment/thread mutation endpoints.
func handleMutationError(
	ctx *httpfx.Context,
	logger *logfx.Logger,
	err error,
	action string,
	entityID string,
) httpfx.Result {
	if errors.Is(err, discussions.ErrCommentNotFound) ||
		errors.Is(err, discussions.ErrThreadNotFound) {
		return ctx.Results.NotFound(httpfx.WithErrorMessage("not found"))
	}

	if errors.Is(err, discussions.ErrInsufficientPermission) {
		return ctx.Results.Error(
			http.StatusForbidden,
			httpfx.WithErrorMessage("insufficient permission"),
		)
	}

	if errors.Is(err, discussions.ErrInvalidVoteDirection) {
		return ctx.Results.BadRequest(httpfx.WithErrorMessage("vote direction must be +1 or -1"))
	}

	if errors.Is(err, discussions.ErrContentTooShort) ||
		errors.Is(err, discussions.ErrContentTooLong) {
		return ctx.Results.BadRequest(httpfx.WithSanitizedError(err))
	}

	logger.ErrorContext(ctx.Request.Context(), "Failed to "+action+" discussion entity",
		slog.String("error", err.Error()),
		slog.String("action", action),
		slog.String("entity_id", entityID))

	return ctx.Results.Error(
		http.StatusInternalServerError,
		httpfx.WithSanitizedError(err),
	)
}
