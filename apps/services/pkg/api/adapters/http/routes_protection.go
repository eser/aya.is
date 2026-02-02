package http

import (
	"errors"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/middlewares"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/protection"
)

func RegisterHTTPRoutesForProtection(
	routes *httpfx.Router,
	logger *logfx.Logger,
	protectionService *protection.Service,
) {
	// POST /{locale}/protection/pow-challenges - Create a new PoW challenge
	routes.Route(
		"POST /{locale}/protection/pow-challenges",
		func(ctx *httpfx.Context) httpfx.Result {
			// Get client IP from context (set by ResolveAddressMiddleware)
			clientIP := middlewares.GetClientAddrs(ctx.Request)

			// Create PoW challenge
			challenge, err := protectionService.CreatePOWChallenge(ctx.Request.Context(), clientIP)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create PoW challenge",
					"error", err.Error(),
					"client_ip_hash", protection.HashIP(clientIP))

				if errors.Is(err, protection.ErrPOWChallengeDisabled) {
					return ctx.Results.JSON(map[string]any{
						"data": map[string]any{
							"enabled": false,
						},
						"error": nil,
					})
				}

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to create PoW challenge"),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  challenge.ToResponse(),
				"error": nil,
			})
		}).
		HasSummary("Create PoW Challenge").
		HasDescription("Creates a new proof-of-work challenge for bot protection.").
		HasResponse(http.StatusOK)
}
