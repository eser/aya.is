package healthcheck

import (
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/processfx"
)

// WorkerHealthStatus represents the health status of a single worker.
type WorkerHealthStatus struct {
	State          string    `json:"state"`
	LastHeartbeat  time.Time `json:"last_heartbeat,omitempty"`
	RestartCount   int       `json:"restart_count"`
	TotalRestarts  int       `json:"total_restarts"`
	ItemsProcessed int64     `json:"items_processed"`
	Uptime         string    `json:"uptime,omitempty"`
	Error          string    `json:"error,omitempty"`
}

// HealthResponse represents the full health check response.
type HealthResponse struct {
	Status  string                        `json:"status"`
	Workers map[string]WorkerHealthStatus `json:"workers,omitempty"`
	Summary *HealthSummary                `json:"summary,omitempty"`
}

// HealthSummary provides aggregate worker health statistics.
type HealthSummary struct {
	Total      int `json:"total"`
	Healthy    int `json:"healthy"`
	Stuck      int `json:"stuck"`
	Restarting int `json:"restarting"`
	Failed     int `json:"failed"`
}

// supervisorRegistry holds the reference to the supervisor registry.
// This is set via SetSupervisorRegistry.
var supervisorRegistry *processfx.SupervisorRegistry //nolint:gochecknoglobals

// SetSupervisorRegistry sets the supervisor registry for health checks.
// This should be called during application initialization.
func SetSupervisorRegistry(registry *processfx.SupervisorRegistry) {
	supervisorRegistry = registry
}

func RegisterHTTPRoutes(routes *httpfx.Router, config *httpfx.Config) {
	if !config.HealthCheckEnabled {
		return
	}

	// Simple health check (backwards compatible).
	routes.
		Route("GET /health-check", func(ctx *httpfx.Context) httpfx.Result {
			return ctx.Results.Ok()
		}).
		HasSummary("Health Check").
		HasDescription("Simple health check endpoint").
		HasResponse(http.StatusNoContent)

	// Detailed health check with worker status.
	routes.
		Route("GET /health", func(ctx *httpfx.Context) httpfx.Result {
			response := HealthResponse{
				Status:  "healthy",
				Workers: make(map[string]WorkerHealthStatus),
			}

			// If supervisor registry is available, include worker status.
			if supervisorRegistry != nil {
				summary := supervisorRegistry.Summary()

				// Set overall status based on worker health.
				if !summary.IsHealthy {
					if summary.Failed > 0 {
						response.Status = "unhealthy"
					} else {
						response.Status = "degraded"
					}
				}

				// Build worker status map.
				for name, status := range summary.Supervisors {
					workerHealth := WorkerHealthStatus{
						State:          status.State.String(),
						RestartCount:   status.RestartCount,
						TotalRestarts:  status.TotalRestarts,
						ItemsProcessed: status.ItemsProcessed,
					}

					if !status.LastHeartbeat.IsZero() {
						workerHealth.LastHeartbeat = status.LastHeartbeat
					}

					if !status.StartedAt.IsZero() {
						workerHealth.Uptime = status.Uptime().Round(time.Second).String()
					}

					if status.LastError != nil {
						workerHealth.Error = status.LastError.Error()
					}

					response.Workers[name] = workerHealth
				}

				// Include summary.
				response.Summary = &HealthSummary{
					Total:      summary.Total,
					Healthy:    summary.Healthy,
					Stuck:      summary.Stuck,
					Restarting: summary.Restarting,
					Failed:     summary.Failed,
				}
			}

			// Return JSON response.
			// Note: If unhealthy, the "status" field in the response will indicate this.
			// Kubernetes probes can check the "status" field value.
			return ctx.Results.JSON(response)
		}).
		HasSummary("Detailed Health Check").
		HasDescription("Returns detailed health status including worker supervision state").
		HasResponse(http.StatusOK)
}
