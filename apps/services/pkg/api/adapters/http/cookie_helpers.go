package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

var ErrGetSessionCookie = errors.New("failed to get session cookie")

// SetSessionCookie sets an HttpOnly, Secure, SameSite=None cookie for cross-domain SSO.
func SetSessionCookie(
	w http.ResponseWriter,
	sessionID string,
	expiresAt time.Time,
	config *auth.Config,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieName,
		Value:    sessionID,
		Path:     "/",
		Domain:   config.CookieDomain,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(w http.ResponseWriter, config *auth.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieName,
		Value:    "",
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   -1,
		Secure:   config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

// GetSessionIDFromCookie extracts the session ID from the request cookie.
func GetSessionIDFromCookie(r *http.Request, config *auth.Config) (string, error) {
	cookie, err := r.Cookie(config.CookieName)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGetSessionCookie, err)
	}

	return cookie.Value, nil
}

// SetThemeCookie sets a non-HttpOnly cookie for the theme preference.
// This cookie is readable by JavaScript and shared across all subdomains
// via the same domain as the session cookie (.aya.is), enabling cross-domain
// theme persistence without extra API calls.
func SetThemeCookie(w http.ResponseWriter, theme string, config *auth.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:     "site_theme",
		Value:    theme,
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		Secure:   config.SecureCookie,
		HttpOnly: false, // Readable by JavaScript for FOUC prevention
		SameSite: http.SameSiteNoneMode,
	})
}

// ClearThemeCookie removes the theme cookie.
func ClearThemeCookie(w http.ResponseWriter, config *auth.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:     "site_theme",
		Value:    "",
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   -1,
		Secure:   config.SecureCookie,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	})
}
