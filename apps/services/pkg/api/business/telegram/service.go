package telegram

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

const (
	// TokenLength is the number of random bytes (produces 2x hex chars).
	TokenLength = 32
	// TokenExpiryMinutes is the TTL for a link token.
	TokenExpiryMinutes = 10
)

// Sentinel errors.
var (
	ErrTokenNotFound             = errors.New("link token not found or expired")
	ErrTokenConsumed             = errors.New("link token already consumed")
	ErrAlreadyLinked             = errors.New("telegram account already linked to a profile")
	ErrProfileAlreadyHasTelegram = errors.New("profile already has a telegram link")
	ErrFailedToCreateToken       = errors.New("failed to create link token")
	ErrFailedToLink              = errors.New("failed to link telegram account")
	ErrFailedToUnlink            = errors.New("failed to unlink telegram account")
	ErrNotLinked                 = errors.New("telegram account is not linked")
)

// Repository defines storage operations for the Telegram service.
type Repository interface {
	CreateLinkToken(ctx context.Context, token *TelegramLinkToken) error
	GetLinkTokenByToken(ctx context.Context, token string) (*TelegramLinkToken, error)
	ConsumeLinkToken(ctx context.Context, token string) error
	CleanupExpiredTokens(ctx context.Context) error

	GetProfileLinkByTelegramRemoteID(ctx context.Context, remoteID string) (*ProfileLinkInfo, error)
	GetProfileLinkByProfileIDAndTelegram(
		ctx context.Context,
		profileID string,
	) (*ProfileLinkInfo, error)
	CreateTelegramProfileLink(ctx context.Context, params *CreateProfileLinkParams) error
	SoftDeleteTelegramProfileLink(ctx context.Context, remoteID string) error
	GetMaxProfileLinkOrder(ctx context.Context, profileID string) (int, error)

	// GetProfileSlugByID resolves a profile ID to its slug (for bot status messages).
	GetProfileSlugByID(ctx context.Context, profileID string) (string, error)
}

// Service provides Telegram account linking business logic.
type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator func() string
}

// NewService creates a new Telegram service.
func NewService(
	logger *logfx.Logger,
	repo Repository,
	idGenerator func() string,
) *Service {
	return &Service{
		logger:      logger,
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// GenerateLinkToken creates a crypto-random token for the Telegram deep link flow.
// Returns the raw token string (to be embedded in the deep link).
func (s *Service) GenerateLinkToken(
	ctx context.Context,
	profileID string,
	profileSlug string,
	userID string,
) (string, error) {
	// Check if profile already has a telegram link
	existing, err := s.repo.GetProfileLinkByProfileIDAndTelegram(ctx, profileID)
	if err == nil && existing != nil {
		return "", ErrProfileAlreadyHasTelegram
	}

	// Generate random token
	tokenBytes := make([]byte, TokenLength)

	_, err = rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateToken, err)
	}

	token := hex.EncodeToString(tokenBytes)

	linkToken := &TelegramLinkToken{
		ID:              s.idGenerator(),
		Token:           token,
		ProfileID:       profileID,
		ProfileSlug:     profileSlug,
		CreatedByUserID: userID,
	}

	err = s.repo.CreateLinkToken(ctx, linkToken)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateToken, err)
	}

	s.logger.DebugContext(ctx, "Telegram link token generated",
		slog.String("profile_id", profileID),
		slog.String("profile_slug", profileSlug))

	return token, nil
}

// LinkAccountByToken validates a token and links a Telegram account to a profile.
func (s *Service) LinkAccountByToken( //nolint:funlen
	ctx context.Context,
	token string,
	telegramUserID int64,
	telegramUsername string,
) (*LinkResult, error) {
	// Look up the token
	linkToken, err := s.repo.GetLinkTokenByToken(ctx, token)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	remoteID := strconv.FormatInt(telegramUserID, 10)

	// Check if this Telegram user is already linked to any profile
	existing, err := s.repo.GetProfileLinkByTelegramRemoteID(ctx, remoteID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyLinked
	}

	// Check if the target profile already has a telegram link
	existingForProfile, err := s.repo.GetProfileLinkByProfileIDAndTelegram(ctx, linkToken.ProfileID)
	if err == nil && existingForProfile != nil {
		return nil, ErrProfileAlreadyHasTelegram
	}

	// Get the next order value for this profile's links
	maxOrder, err := s.repo.GetMaxProfileLinkOrder(ctx, linkToken.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToLink, err)
	}

	// Build the Telegram URI
	uri := ""
	if telegramUsername != "" {
		uri = "https://t.me/" + telegramUsername
	}

	// Create the profile link
	err = s.repo.CreateTelegramProfileLink(ctx, &CreateProfileLinkParams{
		ID:               s.idGenerator(),
		ProfileID:        linkToken.ProfileID,
		RemoteID:         remoteID,
		PublicID:         telegramUsername,
		URI:              uri,
		Order:            maxOrder + 1,
		AddedByProfileID: linkToken.ProfileID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToLink, err)
	}

	// Consume the token (mark as used)
	err = s.repo.ConsumeLinkToken(ctx, token)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to consume link token (link was already created)",
			slog.String("token_id", linkToken.ID),
			slog.String("error", err.Error()))
	}

	s.logger.InfoContext(ctx, "Telegram account linked",
		slog.String("profile_id", linkToken.ProfileID),
		slog.String("profile_slug", linkToken.ProfileSlug),
		slog.Int64("telegram_user_id", telegramUserID))

	return &LinkResult{
		ProfileID:        linkToken.ProfileID,
		ProfileSlug:      linkToken.ProfileSlug,
		TelegramUserID:   telegramUserID,
		TelegramUsername: telegramUsername,
	}, nil
}

// GetLinkedProfile returns the profile link for a Telegram user, or nil if not linked.
func (s *Service) GetLinkedProfile(
	ctx context.Context,
	telegramUserID int64,
) (*ProfileLinkInfo, error) {
	remoteID := strconv.FormatInt(telegramUserID, 10)

	info, err := s.repo.GetProfileLinkByTelegramRemoteID(ctx, remoteID)
	if err != nil {
		return nil, ErrNotLinked
	}

	return info, nil
}

// UnlinkAccount removes the Telegram link for a given Telegram user.
func (s *Service) UnlinkAccount(ctx context.Context, telegramUserID int64) error {
	remoteID := strconv.FormatInt(telegramUserID, 10)

	// Verify it exists first
	_, err := s.repo.GetProfileLinkByTelegramRemoteID(ctx, remoteID)
	if err != nil {
		return ErrNotLinked
	}

	err = s.repo.SoftDeleteTelegramProfileLink(ctx, remoteID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnlink, err)
	}

	s.logger.InfoContext(ctx, "Telegram account unlinked",
		slog.Int64("telegram_user_id", telegramUserID))

	return nil
}

// CleanupExpiredTokens removes expired link tokens from the database.
func (s *Service) CleanupExpiredTokens(ctx context.Context) error {
	err := s.repo.CleanupExpiredTokens(ctx)
	if err != nil {
		return fmt.Errorf("cleanup expired tokens: %w", err)
	}

	return nil
}
