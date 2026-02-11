package http

import (
	"net/url"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// validateLocale extracts the locale path parameter and validates it against
// the list of supported locales. Returns the locale string and true if valid,
// or an empty string and false with a BadRequest result if invalid.
func validateLocale(ctx *httpfx.Context) (string, bool) {
	locale := ctx.Request.PathValue("locale")

	if !profiles.IsValidLocale(locale) {
		return "", false
	}

	return locale, true
}

// isAllowedCorsOrigin checks if a parsed URL's origin is in the CORS allowed origins list.
func isAllowedCorsOrigin(authService *auth.Service, u *url.URL) bool {
	origin := u.Scheme + "://" + u.Host

	for _, allowed := range authService.Config.GetCorsAllowedOrigins() {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}
