package profiles

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/auth"
)

// Constants for profile link state operations.
const (
	StateExpiryDuration = 10 * time.Minute
)

// Sentinel errors for profile link state operations.
var (
	ErrInvalidState = errors.New("invalid or expired profile link state")
	ErrStateExpired = errors.New("profile link state expired")
)

// EncodeProfileLinkState encodes a ProfileLinkState to a base64 URL-safe string.
func EncodeProfileLinkState(state *ProfileLinkState) (string, error) {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("%w: %w", auth.ErrFailedToGenerateState, err)
	}

	return base64.URLEncoding.EncodeToString(stateJSON), nil
}

// DecodeProfileLinkState decodes a base64 URL-safe string to a ProfileLinkState.
func DecodeProfileLinkState(encodedState string) (*ProfileLinkState, error) {
	stateJSON, err := base64.URLEncoding.DecodeString(encodedState)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid encoding", ErrInvalidState)
	}

	var state ProfileLinkState

	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid format", ErrInvalidState)
	}

	return &state, nil
}

// ValidateProfileLinkState validates a ProfileLinkState (checks expiry).
func ValidateProfileLinkState(state *ProfileLinkState) error {
	if time.Now().After(state.ExpiresAt) {
		return ErrStateExpired
	}

	return nil
}

// CreateProfileLinkState creates a new ProfileLinkState for profile link OAuth flows.
func CreateProfileLinkState(
	profileSlug, locale, redirectOrigin string,
) (*ProfileLinkState, string, error) {
	randomState, err := auth.GenerateRandomState()
	if err != nil {
		return nil, "", err //nolint:wrapcheck
	}

	state := &ProfileLinkState{
		State:          randomState,
		ProfileSlug:    profileSlug,
		Locale:         locale,
		RedirectOrigin: redirectOrigin,
		ExpiresAt:      time.Now().Add(StateExpiryDuration),
	}

	encodedState, err := EncodeProfileLinkState(state)
	if err != nil {
		return nil, "", err
	}

	return state, encodedState, nil
}
