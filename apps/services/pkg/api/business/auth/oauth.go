package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

// Constants for OAuth operations.
const (
	OAuthStateSizeBytes = 32 // 256 bits of entropy
)

// Sentinel errors for OAuth operations.
var (
	ErrFailedToGenerateState = errors.New("failed to generate OAuth state")
)

// GenerateRandomState generates a cryptographically secure random state string.
func GenerateRandomState() (string, error) {
	stateBytes := make([]byte, OAuthStateSizeBytes)

	_, err := rand.Read(stateBytes)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGenerateState, err)
	}

	return hex.EncodeToString(stateBytes), nil
}
