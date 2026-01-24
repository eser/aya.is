package middlewares

import (
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

const (
	// HTTP status code threshold for error logging.
	httpErrorThreshold = 400
)

func TracingMiddleware(logger *logfx.Logger, skipLoggingPaths string) httpfx.Handler {
	skipLoggingPathsArray := strings.Split(skipLoggingPaths, ",")

	return func(ctx *httpfx.Context) httpfx.Result {
		// Process the request
		result := ctx.Next()

		if slices.Contains(skipLoggingPathsArray, ctx.Request.URL.Path) {
			return result
		}

		startTime := time.Now()

		attrs := []any{
			slog.String("scope_name", "http"),
			slog.String("http.method", ctx.Request.Method),
			slog.String("http.path", ctx.Request.URL.Path),
			slog.String("user_agent", ctx.Request.UserAgent()),
			slog.String("remote_addr", ctx.Request.RemoteAddr),
		}

		// Extract trace context from incoming request headers using W3C Trace Context
		requestCtx := logger.PropagatorExtract(ctx.Request.Context(), ctx.Request.Header)

		// Start span with the extracted context (or create new trace if none exists)
		newCtx, span := logger.StartSpan(requestCtx, "HTTP Request", attrs...)
		defer span.End()

		// Update the request context
		ctx.UpdateContext(newCtx)

		logger.InfoContext(newCtx, "HTTP request started", attrs...)

		// Calculate duration
		duration := time.Since(startTime)

		// Log request completion
		attrs = append(
			attrs,
			slog.Int("http.status_code", result.StatusCode()),
			slog.Duration("duration", duration),
		)

		// Inject trace context into response headers using W3C Trace Context format
		logger.PropagatorInject(newCtx, ctx.ResponseWriter.Header())

		span.SetAttributes(attrs...)

		if result.StatusCode() >= httpErrorThreshold {
			logger.WarnContext(
				newCtx,
				"HTTP request completed with error",
				attrs...,
			)
		} else {
			logger.DebugContext(
				newCtx,
				"HTTP request completed",
				attrs...,
			)
		}

		return result
	}
}
