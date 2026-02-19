package telegram

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

const (
	// CodeLength is the number of characters in a verification code.
	CodeLength = 6
	// CodeExpiryMinutes is the TTL for a verification code.
	CodeExpiryMinutes = 10
)

// codeChars are the characters used for verification codes.
// Removed confusing characters: 0/O, 1/I/L to improve readability.
var codeChars = []byte("ABCDEFGHJKMNPQRSTUVWXYZ23456789") //nolint:gochecknoglobals

// Sentinel errors.
var (
	ErrCodeNotFound              = errors.New("verification code not found or expired")
	ErrCodeConsumed              = errors.New("verification code already consumed")
	ErrAlreadyLinked             = errors.New("telegram account already linked to a profile")
	ErrProfileAlreadyHasTelegram = errors.New("profile already has a telegram link")
	ErrFailedToCreateCode        = errors.New("failed to create verification code")
	ErrFailedToLink              = errors.New("failed to link telegram account")
	ErrFailedToUnlink            = errors.New("failed to unlink telegram account")
	ErrNotLinked                 = errors.New("telegram account is not linked")
	ErrGroupInviteCodeNotFound   = errors.New("group invite code not found or expired")
	ErrGroupInviteCodeConsumed   = errors.New("group invite code already consumed")
	ErrFailedToCreateInviteCode  = errors.New("failed to create group invite code")
)

// Repository defines storage operations for the Telegram service.
type Repository interface { //nolint:interfacebloat
	CreateVerificationCode(ctx context.Context, code *TelegramVerificationCode) error
	GetVerificationCodeByCode(ctx context.Context, code string) (*TelegramVerificationCode, error)
	ConsumeVerificationCode(ctx context.Context, code string) error
	CleanupExpiredCodes(ctx context.Context) error

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

	// GetMemberProfileTelegramLinks returns all telegram links from non-individual profiles
	// that the given member profile belongs to (visibility filtering happens in the service).
	GetMemberProfileTelegramLinks(
		ctx context.Context,
		memberProfileID string,
	) ([]RawGroupTelegramLink, error)

	// Group invite code operations.
	CreateGroupInviteCode(ctx context.Context, code *TelegramGroupInviteCode) error
	GetGroupInviteCodeByCode(ctx context.Context, code string) (*TelegramGroupInviteCode, error)
	ConsumeGroupInviteCode(ctx context.Context, code string) error
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

// GenerateVerificationCode creates a short code for the Telegram account linking flow.
// Called by the bot when a user sends /start. Returns the code string to send back to the user.
func (s *Service) GenerateVerificationCode(
	ctx context.Context,
	telegramUserID int64,
	telegramUsername string,
) (string, error) {
	remoteID := strconv.FormatInt(telegramUserID, 10)

	// Check if this Telegram user is already linked to a profile
	existing, err := s.repo.GetProfileLinkByTelegramRemoteID(ctx, remoteID)
	if err == nil && existing != nil {
		return "", ErrAlreadyLinked
	}

	// Generate random code
	code, err := generateCode()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateCode, err)
	}

	now := time.Now()

	verificationCode := &TelegramVerificationCode{
		ID:               s.idGenerator(),
		Code:             code,
		TelegramUserID:   telegramUserID,
		TelegramUsername: telegramUsername,
		CreatedAt:        now,
		ExpiresAt:        now.Add(CodeExpiryMinutes * time.Minute),
		ConsumedAt:       nil,
	}

	err = s.repo.CreateVerificationCode(ctx, verificationCode)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateCode, err)
	}

	s.logger.DebugContext(ctx, "Telegram verification code generated",
		slog.Int64("telegram_user_id", telegramUserID))

	return code, nil
}

// VerifyCodeAndLink validates a verification code and links the Telegram account to a profile.
// Called by the web UI when the user pastes the code.
func (s *Service) VerifyCodeAndLink( //nolint:funlen
	ctx context.Context,
	code string,
	profileID string,
	profileSlug string,
	userID string,
) (*LinkResult, error) {
	// Look up the code
	verificationCode, err := s.repo.GetVerificationCodeByCode(ctx, code)
	if err != nil {
		return nil, ErrCodeNotFound
	}

	remoteID := strconv.FormatInt(verificationCode.TelegramUserID, 10)

	// Check if this Telegram user is already linked to any profile
	existing, err := s.repo.GetProfileLinkByTelegramRemoteID(ctx, remoteID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyLinked
	}

	// Check if the target profile already has a managed telegram link
	existingForProfile, err := s.repo.GetProfileLinkByProfileIDAndTelegram(ctx, profileID)
	if err == nil && existingForProfile != nil {
		return nil, ErrProfileAlreadyHasTelegram
	}

	// Get the next order value for this profile's links
	maxOrder, err := s.repo.GetMaxProfileLinkOrder(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToLink, err)
	}

	// Build the Telegram URI
	uri := ""
	if verificationCode.TelegramUsername != "" {
		uri = "https://t.me/" + verificationCode.TelegramUsername
	}

	// Create the profile link
	err = s.repo.CreateTelegramProfileLink(ctx, &CreateProfileLinkParams{
		ID:               s.idGenerator(),
		ProfileID:        profileID,
		RemoteID:         remoteID,
		PublicID:         verificationCode.TelegramUsername,
		URI:              uri,
		Order:            maxOrder + 1,
		AddedByProfileID: profileID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToLink, err)
	}

	// Consume the code (mark as used)
	err = s.repo.ConsumeVerificationCode(ctx, code)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to consume verification code (link was already created)",
			slog.String("code_id", verificationCode.ID),
			slog.String("error", err.Error()))
	}

	s.logger.InfoContext(ctx, "Telegram account linked via verification code",
		slog.String("profile_id", profileID),
		slog.String("profile_slug", profileSlug),
		slog.Int64("telegram_user_id", verificationCode.TelegramUserID))

	return &LinkResult{
		ProfileID:        profileID,
		ProfileSlug:      profileSlug,
		TelegramUserID:   verificationCode.TelegramUserID,
		TelegramUsername: verificationCode.TelegramUsername,
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

// GetProfileSlugByID resolves a profile ID to its slug.
func (s *Service) GetProfileSlugByID(ctx context.Context, profileID string) (string, error) {
	slug, err := s.repo.GetProfileSlugByID(ctx, profileID)
	if err != nil {
		return "", fmt.Errorf("get profile slug: %w", err)
	}

	return slug, nil
}

// CleanupExpiredCodes removes expired verification codes from the database.
func (s *Service) CleanupExpiredCodes(ctx context.Context) error {
	err := s.repo.CleanupExpiredCodes(ctx)
	if err != nil {
		return fmt.Errorf("cleanup expired codes: %w", err)
	}

	return nil
}

// GetGroupTelegramLinks returns Telegram links from non-individual profiles
// that the user is a member of, filtered by the user's membership-based visibility access.
func (s *Service) GetGroupTelegramLinks(
	ctx context.Context,
	memberProfileID string,
) ([]GroupTelegramLink, error) {
	rawLinks, err := s.repo.GetMemberProfileTelegramLinks(ctx, memberProfileID)
	if err != nil {
		return nil, fmt.Errorf("get member profile telegram links: %w", err)
	}

	result := make([]GroupTelegramLink, 0, len(rawLinks))

	for _, raw := range rawLinks {
		memberLevel := profiles.MembershipKindLevel[profiles.MembershipKind(raw.MembershipKind)]
		requiredMembership := profiles.MinMembershipForVisibility[profiles.LinkVisibility(raw.LinkVisibility)]
		requiredLevel := profiles.MembershipKindLevel[requiredMembership]

		// Public links (requiredMembership == "") are always visible
		if requiredMembership != "" && memberLevel < requiredLevel {
			continue
		}

		result = append(result, GroupTelegramLink{
			ProfileSlug:  raw.ProfileSlug,
			ProfileTitle: raw.ProfileTitle,
			LinkTitle:    raw.LinkTitle,
			LinkURI:      raw.LinkURI,
			LinkPublicID: raw.LinkPublicID,
		})
	}

	return result, nil
}

// GenerateGroupInviteCode creates a short code for the Telegram group invite flow.
// Called by the bot when a lead types /invite in a group chat.
func (s *Service) GenerateGroupInviteCode(
	ctx context.Context,
	chatID int64,
	chatTitle string,
	telegramUserID int64,
) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateInviteCode, err)
	}

	now := time.Now()

	inviteCode := &TelegramGroupInviteCode{
		ID:                      s.idGenerator(),
		Code:                    code,
		TelegramChatID:          chatID,
		TelegramChatTitle:       chatTitle,
		CreatedByTelegramUserID: telegramUserID,
		CreatedAt:               now,
		ExpiresAt:               now.Add(CodeExpiryMinutes * time.Minute),
		ConsumedAt:              nil,
	}

	err = s.repo.CreateGroupInviteCode(ctx, inviteCode)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToCreateInviteCode, err)
	}

	s.logger.DebugContext(ctx, "Telegram group invite code generated",
		slog.Int64("chat_id", chatID),
		slog.String("chat_title", chatTitle),
		slog.Int64("telegram_user_id", telegramUserID))

	return code, nil
}

// ResolveGroupInviteCode looks up a group invite code and returns its data if valid.
// Does NOT consume the code — consumption happens after envelope creation.
func (s *Service) ResolveGroupInviteCode(
	ctx context.Context,
	code string,
) (*TelegramGroupInviteCode, error) {
	inviteCode, err := s.repo.GetGroupInviteCodeByCode(ctx, code)
	if err != nil {
		return nil, ErrGroupInviteCodeNotFound
	}

	return inviteCode, nil
}

// ConsumeGroupInviteCode marks an invite code as consumed after envelope creation.
func (s *Service) ConsumeGroupInviteCode(ctx context.Context, code string) error {
	err := s.repo.ConsumeGroupInviteCode(ctx, code)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrGroupInviteCodeConsumed, err)
	}

	return nil
}

// generateCode creates a random verification code using crypto/rand.
func generateCode() (string, error) {
	randomBytes := make([]byte, CodeLength)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	code := make([]byte, CodeLength)
	charsLen := len(codeChars)

	for i, b := range randomBytes {
		code[i] = codeChars[int(b)%charsLen]
	}

	return string(code), nil
}
