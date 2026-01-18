package protection

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// Repository defines the interface for protection data operations.
type Repository interface {
	// POW Challenge operations
	CreatePOWChallenge(ctx context.Context, challenge *POWChallenge) error
	GetPOWChallengeByID(ctx context.Context, id string) (*POWChallenge, error)
	MarkPOWChallengeUsed(ctx context.Context, id string) error
	DeleteExpiredPOWChallenges(ctx context.Context) error
}

// Service handles protection-related business logic.
type Service struct {
	logger *logfx.Logger
	config *Config
	repo   Repository
	idGen  func() string
}

// NewService creates a new protection service.
func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	idGen func() string,
) *Service {
	return &Service{
		logger: logger,
		config: config,
		repo:   repo,
		idGen:  idGen,
	}
}

// CreatePOWChallenge creates a new PoW challenge for the given IP.
func (s *Service) CreatePOWChallenge(ctx context.Context, clientIP string) (*POWChallenge, error) {
	if !s.config.POWChallenge.Enabled {
		return nil, ErrPOWChallengeDisabled
	}

	// Generate random prefix (32 bytes = 64 hex chars)
	prefixBytes := make([]byte, 32)
	if _, err := rand.Read(prefixBytes); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateChallenge, err)
	}

	prefix := hex.EncodeToString(prefixBytes)
	now := time.Now()

	challenge := &POWChallenge{
		ID:         s.idGen(),
		Prefix:     prefix,
		Difficulty: s.config.POWChallenge.Difficulty,
		IPHash:     HashIP(clientIP),
		Used:       false,
		ExpiresAt:  now.Add(s.config.POWChallenge.Expiry),
		CreatedAt:  now,
	}

	err := s.repo.CreatePOWChallenge(ctx, challenge)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateChallenge, err)
	}

	return challenge, nil
}

// VerifyAndConsumePOWChallenge verifies a PoW challenge solution and marks it as used.
func (s *Service) VerifyAndConsumePOWChallenge(
	ctx context.Context,
	challengeID string,
	nonce string,
	clientIP string,
) error {
	if !s.config.POWChallenge.Enabled {
		// PoW disabled, skip verification
		return nil
	}

	challenge, err := s.repo.GetPOWChallengeByID(ctx, challengeID)
	if err != nil {
		return ErrPOWChallengeNotFound
	}

	// Check if challenge is already used
	if challenge.Used {
		return ErrPOWChallengeUsed
	}

	// Check if challenge has expired
	if time.Now().After(challenge.ExpiresAt) {
		return ErrPOWChallengeExpired
	}

	// Verify the solution
	if !VerifyPOWChallenge(challenge.Prefix, nonce, challenge.Difficulty) {
		return ErrPOWChallengeInvalid
	}

	// Mark as used
	if err := s.repo.MarkPOWChallengeUsed(ctx, challengeID); err != nil {
		s.logger.WarnContext(ctx, "Failed to mark PoW challenge as used",
			"challenge_id", challengeID,
			"error", err.Error())
		// Don't fail the operation - challenge was valid
	}

	return nil
}

// IsPOWChallengeEnabled returns whether PoW challenges are enabled.
func (s *Service) IsPOWChallengeEnabled() bool {
	return s.config.POWChallenge.Enabled
}

// CleanupExpiredChallenges removes expired challenges from the database.
func (s *Service) CleanupExpiredChallenges(ctx context.Context) error {
	return s.repo.DeleteExpiredPOWChallenges(ctx)
}
