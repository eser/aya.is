package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profile_points"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var ErrUnauthorized = errors.New("unauthorized")

func RegisterHTTPRoutesForAdminPoints(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profilePointsService *profile_points.Service,
) {
	// List pending awards
	routes.
		Route(
			"GET /admin/points/pending",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				// Check admin permission
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				// Parse status filter
				statusParam := ctx.Request.URL.Query().Get("status")
				var status *profile_points.PendingAwardStatus

				if statusParam != "" {
					s := profile_points.PendingAwardStatus(statusParam)
					status = &s
				}

				cursor := cursors.NewCursorFromRequest(ctx.Request)

				awards, err := profilePointsService.ListPendingAwards(
					ctx.Request.Context(),
					status,
					cursor,
				)
				if err != nil {
					logger.Error(
						"failed to list pending awards",
						"error", err,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(awards)
			},
		).
		HasSummary("List pending point awards").
		HasDescription("List all pending point awards. Admin only.").
		HasResponse(http.StatusOK)

	// Get pending award by ID
	routes.
		Route(
			"GET /admin/points/pending/{id}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				awardID := ctx.Request.PathValue("id")
				if awardID == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("award ID is required"),
					)
				}

				award, err := profilePointsService.GetPendingAward(
					ctx.Request.Context(),
					awardID,
				)
				if err != nil {
					return ctx.Results.NotFound(
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(award)
			},
		).
		HasSummary("Get pending point award").
		HasDescription("Get a single pending point award by ID. Admin only.").
		HasResponse(http.StatusOK)

	// Approve pending award
	routes.
		Route(
			"POST /admin/points/pending/{id}/approve",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				awardID := ctx.Request.PathValue("id")
				if awardID == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("award ID is required"),
					)
				}

				tx, err := profilePointsService.ApprovePendingAward(
					ctx.Request.Context(),
					awardID,
					user.ID,
				)
				if err != nil {
					logger.Error(
						"failed to approve pending award",
						"error", err,
						"award_id", awardID,
						"reviewer_id", user.ID,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(tx)
			},
		).
		HasSummary("Approve pending point award").
		HasDescription("Approve a pending point award and credit points. Admin only.").
		HasResponse(http.StatusOK)

	// Reject pending award
	routes.
		Route(
			"POST /admin/points/pending/{id}/reject",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				awardID := ctx.Request.PathValue("id")
				if awardID == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("award ID is required"),
					)
				}

				// Parse rejection reason from body
				var body struct {
					Reason string `json:"reason"`
				}

				if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("invalid request body"),
					)
				}

				err = profilePointsService.RejectPendingAward(
					ctx.Request.Context(),
					awardID,
					user.ID,
					body.Reason,
				)
				if err != nil {
					logger.Error(
						"failed to reject pending award",
						"error", err,
						"award_id", awardID,
						"reviewer_id", user.ID,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(map[string]string{
					"status": "rejected",
				})
			},
		).
		HasSummary("Reject pending point award").
		HasDescription("Reject a pending point award with optional reason. Admin only.").
		HasResponse(http.StatusOK)

	// Get pending awards stats
	routes.
		Route(
			"GET /admin/points/stats",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				stats, err := profilePointsService.GetPendingAwardsStats(
					ctx.Request.Context(),
				)
				if err != nil {
					logger.Error(
						"failed to get pending awards stats",
						"error", err,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				return ctx.Results.JSON(stats)
			},
		).
		HasSummary("Get pending awards statistics").
		HasDescription("Get statistics about pending point awards. Admin only.").
		HasResponse(http.StatusOK)

	// Bulk approve pending awards
	routes.
		Route(
			"POST /admin/points/pending/bulk-approve",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				var body struct {
					IDs []string `json:"ids"`
				}

				if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("invalid request body"),
					)
				}

				if len(body.IDs) == 0 {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("at least one ID is required"),
					)
				}

				approvedIDs, err := profilePointsService.BulkApprovePendingAwards(
					ctx.Request.Context(),
					body.IDs,
					user.ID,
				)
				if err != nil {
					logger.Error(
						"failed to bulk approve pending awards",
						"error", err,
						"reviewer_id", user.ID,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				// Build response with status for each ID
				approvedSet := make(map[string]struct{})
				for _, id := range approvedIDs {
					approvedSet[id] = struct{}{}
				}

				results := make([]map[string]any, 0, len(body.IDs))
				for _, id := range body.IDs {
					result := map[string]any{
						"id": id,
					}

					if _, wasApproved := approvedSet[id]; wasApproved {
						result["status"] = "approved"
					} else {
						result["status"] = "skipped"
					}

					results = append(results, result)
				}

				return ctx.Results.JSON(map[string]any{
					"results": results,
				})
			},
		).
		HasSummary("Bulk approve pending awards").
		HasDescription("Approve multiple pending point awards at once. Admin only.").
		HasResponse(http.StatusOK)

	// Bulk reject pending awards
	routes.
		Route(
			"POST /admin/points/pending/bulk-reject",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				user, err := getUserFromContext(ctx, userService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(err))
				}

				if user.Kind != "admin" {
					return ctx.Results.Error(
						http.StatusForbidden,
						httpfx.WithErrorMessage("Admin access required"),
					)
				}

				var body struct {
					IDs    []string `json:"ids"`
					Reason string   `json:"reason"`
				}

				if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("invalid request body"),
					)
				}

				if len(body.IDs) == 0 {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("at least one ID is required"),
					)
				}

				rejectedIDs, err := profilePointsService.BulkRejectPendingAwards(
					ctx.Request.Context(),
					body.IDs,
					user.ID,
					body.Reason,
				)
				if err != nil {
					logger.Error(
						"failed to bulk reject pending awards",
						"error", err,
						"reviewer_id", user.ID,
					)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(err),
					)
				}

				// Build response with status for each ID
				rejectedSet := make(map[string]struct{})
				for _, id := range rejectedIDs {
					rejectedSet[id] = struct{}{}
				}

				results := make([]map[string]any, 0, len(body.IDs))
				for _, id := range body.IDs {
					result := map[string]any{
						"id": id,
					}

					if _, wasRejected := rejectedSet[id]; wasRejected {
						result["status"] = "rejected"
					} else {
						result["status"] = "skipped"
					}

					results = append(results, result)
				}

				return ctx.Results.JSON(map[string]any{
					"results": results,
				})
			},
		).
		HasSummary("Bulk reject pending awards").
		HasDescription("Reject multiple pending point awards at once. Admin only.").
		HasResponse(http.StatusOK)
}

// getUserFromContext extracts the authenticated user from the request context.
func getUserFromContext(ctx *httpfx.Context, userService *users.Service) (*users.User, error) {
	sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
	if !ok {
		return nil, ErrUnauthorized
	}

	session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
	if err != nil || session == nil || session.LoggedInUserID == nil {
		return nil, ErrUnauthorized
	}

	user, err := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
	if err != nil || user == nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}
