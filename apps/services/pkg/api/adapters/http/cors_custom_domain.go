package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// CorsMiddlewareWithCustomDomains validates origins against:
// 1. Config-defined allowed origins (from auth.Config)
// 2. Database custom_domain field (cached at repository layer).
func CorsMiddlewareWithCustomDomains(
	authConfig *auth.Config,
	profileService *profiles.Service,
) httpfx.Handler {
	// Parse config values once at startup
	allowedOrigins := authConfig.GetCorsAllowedOrigins()
	allowedHeaders := strings.Join(authConfig.GetCorsAllowedHeaders(), ", ")
	allowedMethods := strings.Join(authConfig.GetCorsAllowedMethods(), ", ")

	return func(ctx *httpfx.Context) httpfx.Result {
		headers := ctx.ResponseWriter.Header()
		requestOrigin := ctx.Request.Header.Get("Origin")

		if requestOrigin == "" {
			return ctx.Next()
		}

		allowed := false

		// Check config-defined origins first (fast path)
		for _, origin := range allowedOrigins {
			if origin == requestOrigin {
				allowed = true

				break
			}
		}

		// If not in config list, check custom domains in database
		if !allowed {
			domain := extractDomainFromOrigin(requestOrigin, true) // strip www. for DB lookup
			if domain != "" {
				// GetByCustomDomain is cached at repository layer
				profile, _, _ := profileService.GetByCustomDomain(
					ctx.Request.Context(),
					"en", // locale doesn't matter for domain check
					domain,
				)
				if profile != nil {
					allowed = true
				}
			}
		}

		if allowed {
			headers.Set("Access-Control-Allow-Origin", requestOrigin)
			headers.Set("Access-Control-Allow-Credentials", "true")
			headers.Set("Access-Control-Allow-Headers", allowedHeaders)
			headers.Set("Access-Control-Allow-Methods", allowedMethods)
		}

		// Handle preflight
		if ctx.Request.Method == http.MethodOptions {
			return ctx.Results.Ok()
		}

		return ctx.Next()
	}
}

// extractDomainFromOrigin extracts domain from origin URL.
// e.g., "https://eser.dev:443" -> "eser.dev"
// If stripWWW is true, "www." prefix is removed.
func extractDomainFromOrigin(origin string, stripWWW bool) string {
	parsed, err := url.Parse(origin)
	if err != nil {
		return ""
	}

	host := parsed.Hostname() // strips port
	if stripWWW {
		return strings.TrimPrefix(host, "www.")
	}

	return host
}
