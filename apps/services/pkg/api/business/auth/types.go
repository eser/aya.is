package auth

import (
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// Config types.
type GitHubAuthProviderConfig struct {
	ClientID     string `conf:"client_id"`
	ClientSecret string `conf:"client_secret"`
	Scope        string `conf:"scope"         default:"read:user user:email"`
}

type YouTubeOAuthConfig struct {
	ClientID     string `conf:"client_id"`
	ClientSecret string `conf:"client_secret"`
	Scope        string `conf:"scope"         default:"https://www.googleapis.com/auth/youtube.readonly"`
}

type Config struct {
	GitHub    GitHubAuthProviderConfig `conf:"github"`
	YouTube   YouTubeOAuthConfig       `conf:"youtube"`
	JwtSecret string                   `conf:"jwt_secret"`                 // Required - no default for security
	TokenTTL  time.Duration            `conf:"token_ttl"  default:"8760h"` // 365 days in hours (Go doesn't support "d")

	// Cookie settings for cross-domain SSO
	CookieDomain string `conf:"cookie_domain" default:".aya.is"`
	CookieName   string `conf:"cookie_name"   default:"aya_session"`
	SecureCookie bool   `conf:"secure_cookie" default:"true"`

	// CORS settings (comma-separated)
	CorsAllowedOrigins string `conf:"cors_allowed_origins" default:"https://aya.is,https://www.aya.is,http://localhost:3000,http://localhost:5173,http://localhost:4173"`
	CorsAllowedHeaders string `conf:"cors_allowed_headers" default:"Accept,Authorization,Content-Type,Origin,X-Requested-With,Traceparent,Tracestate"`
	CorsAllowedMethods string `conf:"cors_allowed_methods" default:"GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS"`
}

// GetCorsAllowedOrigins parses comma-separated origins into a slice.
func (c *Config) GetCorsAllowedOrigins() []string {
	return splitAndTrim(c.CorsAllowedOrigins)
}

// GetCorsAllowedHeaders parses comma-separated headers into a slice.
func (c *Config) GetCorsAllowedHeaders() []string {
	return splitAndTrim(c.CorsAllowedHeaders)
}

// GetCorsAllowedMethods parses comma-separated methods into a slice.
func (c *Config) GetCorsAllowedMethods() []string {
	return splitAndTrim(c.CorsAllowedMethods)
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// Auth types

type OAuthState struct {
	State       string
	RedirectURI string
}

type AuthResult struct {
	ExpiresAt   time.Time
	User        *users.User
	SessionID   string
	JWT         string
	RedirectURI string
}

type JWTClaims struct {
	UserID    string
	SessionID string
	ExpiresAt int64
}

// OAuthCallbackResult contains all information from an OAuth callback.
// Different use cases (login, profile linking) use different fields.
type OAuthCallbackResult struct {
	// Identity
	RemoteID string // Provider's user/channel ID
	Username string // Handle/login (e.g., GitHub username, YouTube channel handle)
	Name     string // Display name
	Email    string // Email (may be empty for non-login providers like YouTube)
	URI      string // Profile/channel URL

	// Tokens (used for profile linking to store for future API calls)
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiresAt *time.Time
	Scope                string
}
