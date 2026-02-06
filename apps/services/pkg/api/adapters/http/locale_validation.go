package http

import (
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
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
