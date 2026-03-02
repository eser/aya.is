package bulletin

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// Repository defines the data access port for the bulletin module.
type Repository interface {
	GetActiveSubscriptionsForWindow(ctx context.Context, utcHour int) ([]*Subscription, error)
	GetFollowedProfileStoriesSince(
		ctx context.Context,
		subscriberProfileID string,
		localeCode string,
		since time.Time,
		maxStories int,
	) ([]*DigestStory, error)
	UpdateLastBulletinAt(ctx context.Context, subscriptionID string) error
	CreateBulletinLog(ctx context.Context, log *BulletinLog) error

	GetSubscriptionsByProfileID(ctx context.Context, profileID string) ([]*Subscription, error)
	GetSubscription(ctx context.Context, id string) (*Subscription, error)
	UpsertSubscription(ctx context.Context, sub *Subscription) error
	UpdateSubscriptionPreferences(
		ctx context.Context,
		id string,
		frequency DigestFrequency,
		preferredTime int,
	) error
	DeleteSubscription(ctx context.Context, id string) error
	DeleteSubscriptionsByProfileID(ctx context.Context, profileID string) error
}

// Channel defines the delivery port for sending bulletins.
type Channel interface {
	Kind() ChannelKind
	Send(ctx context.Context, recipientProfileID string, digest *Digest) error
}

// UserEmailResolver resolves a profile ID to a user email address.
type UserEmailResolver interface {
	GetUserEmailByProfileID(ctx context.Context, profileID string) (string, error)
}

// RecordIDGenerator is a function that generates unique record IDs.
type RecordIDGenerator func() string

// Service orchestrates bulletin digest creation and delivery.
type Service struct {
	logger   *logfx.Logger
	config   *Config
	repo     Repository
	channels map[ChannelKind]Channel
	idGen    RecordIDGenerator
}

// NewService creates a new bulletin Service.
func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	channels []Channel,
	idGen RecordIDGenerator,
) *Service {
	channelMap := make(map[ChannelKind]Channel, len(channels))
	for _, ch := range channels {
		channelMap[ch.Kind()] = ch
	}

	return &Service{
		logger:   logger,
		config:   config,
		repo:     repo,
		channels: channelMap,
		idGen:    idGen,
	}
}

// ProcessDigestWindow processes all active subscriptions for the current UTC hour window.
func (s *Service) ProcessDigestWindow(ctx context.Context) error {
	now := time.Now().UTC()
	utcHour := now.Hour()

	s.logger.InfoContext(ctx, "Processing bulletin digest window",
		slog.Int("utc_hour", utcHour))

	subscriptions, err := s.repo.GetActiveSubscriptionsForWindow(ctx, utcHour)
	if err != nil {
		return fmt.Errorf("getting active subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		s.logger.DebugContext(ctx, "No subscriptions for current window")

		return nil
	}

	s.logger.InfoContext(ctx, "Found subscriptions for digest window",
		slog.Int("count", len(subscriptions)),
		slog.Int("utc_hour", utcHour))

	for _, sub := range subscriptions {
		processErr := s.processSubscription(ctx, sub, now)
		if processErr != nil {
			s.logger.WarnContext(ctx, "Failed to process subscription",
				slog.String("subscription_id", sub.ID),
				slog.String("profile_id", sub.ProfileID),
				slog.String("error", processErr.Error()))
		}
	}

	return nil
}

// processSubscription handles a single subscription: fetch stories, summarize, deliver.
func (s *Service) processSubscription(
	ctx context.Context,
	sub *Subscription,
	now time.Time,
) error {
	locale := sub.DefaultLocale
	if locale == "" {
		locale = "en"
	}

	// Determine the "since" cutoff — frequency-aware lookback, or last bulletin time
	lookback := frequencyToLookback(sub.Frequency)

	since := now.Add(-lookback)
	if sub.LastBulletinAt != nil {
		since = *sub.LastBulletinAt
	}

	// Fetch stories from followed profiles
	stories, err := s.repo.GetFollowedProfileStoriesSince(
		ctx, sub.ProfileID, locale, since, s.config.MaxStoriesPerDigest,
	)
	if err != nil {
		return fmt.Errorf("fetching followed stories: %w", err)
	}

	// Check threshold
	if len(stories) < s.config.MinStoryThreshold {
		s.logger.DebugContext(ctx, "Not enough stories for digest",
			slog.String("profile_id", sub.ProfileID),
			slog.Int("story_count", len(stories)),
			slog.Int("threshold", s.config.MinStoryThreshold))

		return nil
	}

	// Group stories by author profile
	digest := s.buildDigest(sub.ProfileID, sub.ProfileSlug, locale, stories, now)

	// Deliver and log result
	return s.deliverAndLog(ctx, sub, digest, len(stories), now)
}

// deliverAndLog sends the digest through the channel adapter and records the result.
func (s *Service) deliverAndLog(
	ctx context.Context,
	sub *Subscription,
	digest *Digest,
	storyCount int,
	now time.Time,
) error {
	channel, ok := s.channels[sub.Channel]
	if !ok {
		s.recordBulletinLog(ctx, sub.ID, storyCount, "failed",
			strPtr(fmt.Sprintf("channel %q not available", sub.Channel)), now)

		return fmt.Errorf("%w: %s", ErrChannelNotAvailable, sub.Channel)
	}

	sendErr := channel.Send(ctx, sub.ProfileID, digest)
	if sendErr != nil {
		s.recordBulletinLog(ctx, sub.ID, storyCount, "failed",
			strPtr(sendErr.Error()), now)

		return fmt.Errorf("sending bulletin via %s: %w", sub.Channel, sendErr)
	}

	updateErr := s.repo.UpdateLastBulletinAt(ctx, sub.ID)
	if updateErr != nil {
		s.logger.WarnContext(ctx, "Failed to update last_bulletin_at",
			slog.String("subscription_id", sub.ID),
			slog.String("error", updateErr.Error()))
	}

	s.recordBulletinLog(ctx, sub.ID, storyCount, "sent", nil, now)

	s.logger.InfoContext(ctx, "Bulletin sent successfully",
		slog.String("profile_id", sub.ProfileID),
		slog.String("channel", string(sub.Channel)),
		slog.Int("story_count", storyCount))

	return nil
}

// recordBulletinLog creates a bulletin log entry (best-effort).
func (s *Service) recordBulletinLog(
	ctx context.Context,
	subscriptionID string,
	storyCount int,
	status string,
	errorMessage *string,
	now time.Time,
) {
	logEntry := &BulletinLog{
		ID:             s.idGen(),
		SubscriptionID: subscriptionID,
		StoryCount:     storyCount,
		Status:         status,
		ErrorMessage:   errorMessage,
		CreatedAt:      now,
	}

	_ = s.repo.CreateBulletinLog(ctx, logEntry)
}

// buildDigest groups stories by author profile into a Digest.
func (s *Service) buildDigest(
	recipientProfileID string,
	recipientSlug string,
	locale string,
	stories []*DigestStory,
	now time.Time,
) *Digest {
	groupMap := make(map[string]*DigestGroup)
	groupOrder := make([]string, 0)

	for _, story := range stories {
		group, exists := groupMap[story.AuthorProfileID]
		if !exists {
			group = &DigestGroup{
				ProfileID:         story.AuthorProfileID,
				ProfilePictureURI: story.AuthorProfilePictureURI,
				Slug:              story.AuthorSlug,
				Title:             story.AuthorTitle,
				Stories:           make([]*DigestStory, 0),
			}

			groupMap[story.AuthorProfileID] = group

			groupOrder = append(groupOrder, story.AuthorProfileID)
		}

		group.Stories = append(group.Stories, story)
	}

	// Preserve insertion order (most recent story's author first)
	groups := make([]*DigestGroup, 0, len(groupOrder))
	for _, profileID := range groupOrder {
		groups = append(groups, groupMap[profileID])
	}

	return &Digest{
		RecipientProfileID: recipientProfileID,
		RecipientSlug:      recipientSlug,
		Locale:             locale,
		Groups:             groups,
		GeneratedAt:        now,
	}
}

// GetSubscriptions returns all active subscriptions for a profile.
func (s *Service) GetSubscriptions(ctx context.Context, profileID string) ([]*Subscription, error) {
	subs, err := s.repo.GetSubscriptionsByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("getting subscriptions: %w", err)
	}

	return subs, nil
}

// Subscribe creates or reactivates a bulletin subscription.
func (s *Service) Subscribe(
	ctx context.Context,
	profileID string,
	channel ChannelKind,
	frequency DigestFrequency,
	preferredTime int,
) (*Subscription, error) {
	sub := &Subscription{
		ID:             s.idGen(),
		ProfileID:      profileID,
		Channel:        channel,
		Frequency:      frequency,
		PreferredTime:  preferredTime,
		CreatedAt:      time.Now(),
		LastBulletinAt: nil,
		UpdatedAt:      nil,
		ProfileSlug:    "",
		DefaultLocale:  "",
	}

	err := s.repo.UpsertSubscription(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("upserting subscription: %w", err)
	}

	return sub, nil
}

// UpdatePreferences updates the frequency and preferred time for a subscription.
func (s *Service) UpdatePreferences(
	ctx context.Context,
	subscriptionID string,
	frequency DigestFrequency,
	preferredTime int,
) error {
	err := s.repo.UpdateSubscriptionPreferences(ctx, subscriptionID, frequency, preferredTime)
	if err != nil {
		return fmt.Errorf("updating preferences: %w", err)
	}

	return nil
}

// UpdateBulletinPreferences atomically updates all bulletin preferences for a profile.
// It upserts subscriptions for selected channels and soft-deletes unselected ones.
func (s *Service) UpdateBulletinPreferences(
	ctx context.Context,
	profileID string,
	frequency DigestFrequency,
	preferredTime int,
	channels []ChannelKind,
) error {
	// If no channels selected, soft-delete all subscriptions ("Don't send")
	if len(channels) == 0 {
		err := s.repo.DeleteSubscriptionsByProfileID(ctx, profileID)
		if err != nil {
			return fmt.Errorf("deleting all subscriptions: %w", err)
		}

		return nil
	}

	// Build a set of requested channels for fast lookup
	requested := make(map[ChannelKind]bool, len(channels))
	for _, channel := range channels {
		requested[channel] = true
	}

	// Get existing subscriptions
	existing, err := s.repo.GetSubscriptionsByProfileID(ctx, profileID)
	if err != nil {
		return fmt.Errorf("getting existing subscriptions: %w", err)
	}

	// Upsert requested channels
	for _, channel := range channels {
		sub := &Subscription{
			ID:             s.idGen(),
			ProfileID:      profileID,
			Channel:        channel,
			Frequency:      frequency,
			PreferredTime:  preferredTime,
			CreatedAt:      time.Now(),
			LastBulletinAt: nil,
			UpdatedAt:      nil,
			ProfileSlug:    "",
			DefaultLocale:  "",
		}

		upsertErr := s.repo.UpsertSubscription(ctx, sub)
		if upsertErr != nil {
			return fmt.Errorf("upserting subscription for %s: %w", channel, upsertErr)
		}
	}

	// Soft-delete channels that are no longer selected
	for _, sub := range existing {
		if !requested[sub.Channel] {
			delErr := s.repo.DeleteSubscription(ctx, sub.ID)
			if delErr != nil {
				return fmt.Errorf("deleting subscription %s: %w", sub.ID, delErr)
			}
		}
	}

	return nil
}

// Unsubscribe soft-deletes a subscription.
func (s *Service) Unsubscribe(ctx context.Context, subscriptionID string) error {
	err := s.repo.DeleteSubscription(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("deleting subscription: %w", err)
	}

	return nil
}

// Lookback duration constants for digest frequencies.
const (
	lookbackBiDaily = 48 * time.Hour
	lookbackDaily   = 24 * time.Hour
	lookbackWeekly  = 7 * lookbackDaily
)

// frequencyToLookback returns the lookback duration for a digest frequency.
func frequencyToLookback(freq DigestFrequency) time.Duration {
	switch freq {
	case FrequencyBiDaily:
		return lookbackBiDaily
	case FrequencyWeekly:
		return lookbackWeekly
	case FrequencyDaily:
		return lookbackDaily
	default:
		return lookbackDaily
	}
}

func strPtr(s string) *string {
	return &s
}
