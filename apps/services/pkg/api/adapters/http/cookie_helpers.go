package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

var ErrGetSessionCookie = errors.New("failed to get session cookie")

// SetSessionCookie sets an HttpOnly, Secure, SameSite=Lax cookie for same-site SSO.
// Lax is used instead of None because all services share the same eTLD+1 (aya.is),
// and SameSite=None cookies are blocked by Chrome's Tracking Protection in incognito mode.
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
		SameSite: http.SameSiteLaxMode,
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
		SameSite: http.SameSiteLaxMode,
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
// via the same domain as the session cookie (.aya.is), enabling same-site
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
		SameSite: http.SameSiteLaxMode,
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
		SameSite: http.SameSiteLaxMode,
	})
}
