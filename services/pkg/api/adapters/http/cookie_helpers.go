package http

import (
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

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
		return "", err
	}

	return cookie.Value, nil
}
