package users

import (
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
)

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
	XHandle *string `json:"x_handle"`
	// XRemoteID           *string    `json:"x_remote_id"`
	IndividualProfileID *string    `json:"individual_profile_id"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
	ID                  string     `json:"id"`
	Kind                string     `json:"kind"`
	Name                string     `json:"name"`
}

type Session struct {
	CreatedAt                time.Time  `json:"created_at"`
	OauthRedirectURI         *string    `json:"oauth_redirect_uri"`
	LoggedInUserID           *string    `json:"logged_in_user_id"`
	LoggedInAt               *time.Time `json:"logged_in_at"`
	ExpiresAt                *time.Time `json:"expires_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
	ID                       string     `json:"id"`
	Status                   string     `json:"status"`
	OauthRequestState        string     `json:"oauth_request_state"`
	OauthRequestCodeVerifier string     `json:"oauth_request_code_verifier"`
}
