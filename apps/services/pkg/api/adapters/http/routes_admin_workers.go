package http

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

type adminWorkerResponse struct {
	Name       string  `json:"name"`
	IsRunning  bool    `json:"is_running"`
	IsEnabled  bool    `json:"is_enabled"`
	LastRun    *string `json:"last_run"`
	NextRun    *string `json:"next_run"`
	LastError  *string `json:"last_error"`
	RunCount   int64   `json:"run_count"`
	ErrorCount int64   `json:"error_count"`
	Interval   string  `json:"interval"`
}

func RegisterHTTPRoutesForAdminWorkers(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	runtimeStates *runtime_states.Service,
	workerRegistry *workerfx.Registry,
) {
	// List all workers with status
	routes.
		Route(
			"GET /admin/workers",
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

				statuses := workerRegistry.List()
				response := make([]adminWorkerResponse, 0, len(statuses))

				for _, status := range statuses {
					item := adminWorkerResponse{
						Name:       status.Name,
						IsRunning:  status.IsRunning,
						IsEnabled:  true,
						RunCount:   status.RunCount,
						ErrorCount: status.ErrorCount,
						Interval:   status.Interval.String(),
					}

					// Last run
					if !status.LastRun.IsZero() {
						lastRun := status.LastRun.Format(time.RFC3339)
						item.LastRun = &lastRun
					}

					// Last error
					if status.LastError != nil {
						errMsg := status.LastError.Error()
						item.LastError = &errMsg
					}

					// Next run from runtime_state
					if status.StateKey != "" {
						nextRunKey := status.StateKey + ".next_run_at"
						nextRunAt, err := runtimeStates.GetTime(ctx.Request.Context(), nextRunKey)
						if err == nil {
							nextRun := nextRunAt.Format(time.RFC3339)
							item.NextRun = &nextRun
						}
					}

					// Check if disabled
					disabledKey := "worker." + status.Name + ".disabled"
					disabled, err := runtimeStates.Get(ctx.Request.Context(), disabledKey)
					if err == nil && disabled == "true" {
						item.IsEnabled = false
					}

					response = append(response, item)
				}

				slices.SortFunc(response, func(a, b adminWorkerResponse) int {
					return strings.Compare(a.Name, b.Name)
				})

				return ctx.Results.JSON(response)
			},
		).
		HasSummary("List all workers").
		HasDescription("List all background workers with their status. Admin only.").
		HasResponse(http.StatusOK)

	// Toggle worker enable/disable
	routes.
		Route(
			"POST /admin/workers/{name}/toggle",
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

				name := ctx.Request.PathValue("name")
				if name == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("worker name is required"),
					)
				}

				// Validate worker exists
				_, ok := workerRegistry.Get(name)
				if !ok {
					return ctx.Results.NotFound(
						httpfx.WithErrorMessage("worker not found"),
					)
				}

				disabledKey := "worker." + name + ".disabled"
				reqCtx := ctx.Request.Context()

				// Check current state
				disabled, err := runtimeStates.Get(reqCtx, disabledKey)
				isCurrentlyDisabled := err == nil && disabled == "true"

				var isEnabled bool

				if isCurrentlyDisabled {
					// Enable: remove the disabled key
					_ = runtimeStates.Remove(reqCtx, disabledKey)
					isEnabled = true
				} else {
					// Disable: set the disabled key
					_ = runtimeStates.Set(reqCtx, disabledKey, "true")
					isEnabled = false
				}

				return ctx.Results.JSON(map[string]any{
					"name":       name,
					"is_enabled": isEnabled,
				})
			},
		).
		HasSummary("Toggle worker enable/disable").
		HasDescription("Enable or disable a background worker. Admin only.").
		HasResponse(http.StatusOK)

	// Trigger worker immediately
	routes.
		Route(
			"POST /admin/workers/{name}/trigger",
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

				name := ctx.Request.PathValue("name")
				if name == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("worker name is required"),
					)
				}

				runner, ok := workerRegistry.Get(name)
				if !ok {
					return ctx.Results.NotFound(
						httpfx.WithErrorMessage("worker not found"),
					)
				}

				// Reset next_run_at to allow immediate execution
				status := runner.Status()
				if status.StateKey != "" {
					nextRunKey := status.StateKey + ".next_run_at"
					_ = runtimeStates.SetTime(
						ctx.Request.Context(),
						nextRunKey,
						time.Now().Add(-time.Minute),
					)
				}

				// Signal the runner to execute immediately
				runner.TriggerNow()

				logger.Info("Admin triggered worker",
					"worker", name,
					"triggered_by", user.ID,
				)

				return ctx.Results.JSON(map[string]any{
					"name":      name,
					"triggered": true,
				})
			},
		).
		HasSummary("Trigger worker immediately").
		HasDescription("Trigger a background worker to run immediately. Admin only.").
		HasResponse(http.StatusOK)
}
