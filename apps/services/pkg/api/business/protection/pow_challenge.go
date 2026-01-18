package protection

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

var (
	ErrPOWChallengeNotFound    = errors.New("pow challenge not found")
	ErrPOWChallengeExpired     = errors.New("pow challenge expired")
	ErrPOWChallengeUsed        = errors.New("pow challenge already used")
	ErrPOWChallengeInvalid     = errors.New("pow challenge verification failed")
	ErrPOWChallengeDisabled    = errors.New("pow challenge is disabled")
	ErrFailedToCreateChallenge = errors.New("failed to create pow challenge")
)

// POWChallenge represents a proof-of-work challenge.
type POWChallenge struct {
	ID         string    `json:"id"`
	Prefix     string    `json:"prefix"`
	Difficulty int       `json:"difficulty"`
	IPHash     string    `json:"-"` // Not exposed to client
	Used       bool      `json:"-"` // Not exposed to client
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"-"` // Not exposed to client
}

// POWChallengeResponse is the response sent to the client.
type POWChallengeResponse struct {
	POWChallengeID string    `json:"pow_challenge_id"`
	Prefix         string    `json:"prefix"`
	Difficulty     int       `json:"difficulty"`
	ExpiresAt      time.Time `json:"expires_at"`
}

// ToResponse converts a POWChallenge to a POWChallengeResponse.
func (c *POWChallenge) ToResponse() POWChallengeResponse {
	return POWChallengeResponse{
		POWChallengeID: c.ID,
		Prefix:         c.Prefix,
		Difficulty:     c.Difficulty,
		ExpiresAt:      c.ExpiresAt,
	}
}

// VerifyPOWChallenge verifies that the nonce solves the challenge.
// It checks if SHA256(prefix + nonce) has at least 'difficulty' leading zero bits.
func VerifyPOWChallenge(prefix, nonce string, difficulty int) bool {
	data := prefix + nonce
	hash := sha256.Sum256([]byte(data))

	return hasLeadingZeroBits(hash[:], difficulty)
}

// hasLeadingZeroBits checks if the hash has at least n leading zero bits.
func hasLeadingZeroBits(hash []byte, n int) bool {
	fullBytes := n / 8
	remainingBits := n % 8

	// Check full zero bytes
	for i := range fullBytes {
		if i >= len(hash) {
			return false
		}

		if hash[i] != 0 {
			return false
		}
	}

	// Check remaining bits in the next byte
	if remainingBits > 0 && fullBytes < len(hash) {
		mask := byte(0xFF) << (8 - remainingBits)
		if hash[fullBytes]&mask != 0 {
			return false
		}
	}

	return true
}

// HashIP creates a SHA256 hash of an IP address for storage.
func HashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))

	return hex.EncodeToString(hash[:])
}
