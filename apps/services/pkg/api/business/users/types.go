package users

import (
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
)

// SessionStatus represents the status of a session.
type SessionStatus string

const (
	// SessionStatusActive indicates an active session.
	SessionStatusActive SessionStatus = "active"
	// SessionStatusLoggedOut indicates the user logged out.
	SessionStatusLoggedOut SessionStatus = "logged_out"
	// SessionStatusExpired indicates the session expired.
	SessionStatusExpired SessionStatus = "expired"
	// SessionStatusTerminated indicates the session was terminated by user (e.g., from another device).
	SessionStatusTerminated SessionStatus = "terminated"
)

// ValidSessionStatuses contains all valid session status values.
var ValidSessionStatuses = map[SessionStatus]bool{
	SessionStatusActive:     true,
	SessionStatusLoggedOut:  true,
	SessionStatusExpired:    true,
	SessionStatusTerminated: true,
}

// IsValid checks if a session status is valid.
func (s SessionStatus) IsValid() bool {
	return ValidSessionStatuses[s]
}

// String returns the string representation of the session status.
func (s SessionStatus) String() string {
	return string(s)
}

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(lib.IDsGenerateUnique())
}

type User struct {
	CreatedAt      time.Time `json:"created_at"`
	Email          *string   `json:"email"`
	Phone          *string   `json:"phone"`
	GithubHandle   *string   `json:"github_handle"`
	GithubRemoteID *string   `json:"github_remote_id"`
	BskyHandle     *string   `json:"bsky_handle"`
	// BskyRemoteID        *string    `json:"bsky_remote_id"`
	XHandle       *string `json:"x_handle"`
	AppleRemoteID *string `json:"apple_remote_id"`
	// XRemoteID           *string    `json:"x_remote_id"`
	IndividualProfileID *string    `json:"individual_profile_id"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
	ID                  string     `json:"id"`
	Kind                string     `json:"kind"`
	Name                string     `json:"name"`
}

type Session struct {
	CreatedAt                time.Time     `json:"created_at"`
	OauthRedirectURI         *string       `json:"oauth_redirect_uri"`
	LoggedInUserID           *string       `json:"logged_in_user_id"`
	LoggedInAt               *time.Time    `json:"logged_in_at"`
	LastActivityAt           *time.Time    `json:"last_activity_at"`
	ExpiresAt                *time.Time    `json:"expires_at"`
	UpdatedAt                *time.Time    `json:"updated_at"`
	UserAgent                *string       `json:"user_agent"`
	OAuthProvider            *string       `json:"-"`
	OAuthAccessToken         *string       `json:"-"`
	OAuthTokenScope          *string       `json:"-"`
	ID                       string        `json:"id"`
	Status                   SessionStatus `json:"status"`
	OauthRequestState        string        `json:"oauth_request_state"`
	OauthRequestCodeVerifier string        `json:"oauth_request_code_verifier"`
}
