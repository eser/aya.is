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
	UpsertStorySummaryAI(
		ctx context.Context,
		storyID string,
		localeCode string,
		summaryAI string,
	) error

	GetSubscriptionsByProfileID(ctx context.Context, profileID string) ([]*Subscription, error)
	GetSubscription(ctx context.Context, id string) (*Subscription, error)
	UpsertSubscription(ctx context.Context, sub *Subscription) error
	UpdateSubscriptionPreferences(ctx context.Context, id string, preferredTime int) error
	DeleteSubscription(ctx context.Context, id string) error
}

// Channel defines the delivery port for sending bulletins.
type Channel interface {
	Kind() ChannelKind
	Send(ctx context.Context, recipientProfileID string, digest *Digest) error
}

// StorySummarizer defines the AI summarization port.
type StorySummarizer interface {
	// SummarizeBatch generates AI summaries for the given stories in the target locale.
	// Returns a map of storyID → AI summary text.
	SummarizeBatch(
		ctx context.Context,
		stories []*DigestStory,
		localeCode string,
	) (map[string]string, error)
}

// UserEmailResolver resolves a profile ID to a user email address.
type UserEmailResolver interface {
	GetUserEmailByProfileID(ctx context.Context, profileID string) (string, error)
}

// RecordIDGenerator is a function that generates unique record IDs.
type RecordIDGenerator func() string

// Service orchestrates bulletin digest creation and delivery.
type Service struct {
	logger     *logfx.Logger
	config     *Config
	repo       Repository
	channels   map[ChannelKind]Channel
	summarizer StorySummarizer // optional — nil means skip AI summaries
	idGen      RecordIDGenerator
}

// NewService creates a new bulletin Service.
func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	channels []Channel,
	summarizer StorySummarizer,
	idGen RecordIDGenerator,
) *Service {
	channelMap := make(map[ChannelKind]Channel, len(channels))
	for _, ch := range channels {
		channelMap[ch.Kind()] = ch
	}

	return &Service{
		logger:     logger,
		config:     config,
		repo:       repo,
		channels:   channelMap,
		summarizer: summarizer,
		idGen:      idGen,
	}
}

// ProcessDigestWindow processes all active subscriptions for the current UTC hour window.
func (s *Service) ProcessDigestWindow(ctx context.Context) error { //nolint:cyclop,funlen
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
) error { //nolint:cyclop,funlen
	locale := sub.DefaultLocale
	if locale == "" {
		locale = "en"
	}

	// Determine the "since" cutoff — last bulletin time or 24h ago for first-time
	since := now.Add(-24 * time.Hour)
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

	// AI summarization for stories missing summary_ai
	if s.summarizer != nil {
		storiesNeedingSummary := make([]*DigestStory, 0)

		for _, story := range stories {
			if story.SummaryAI == nil || *story.SummaryAI == "" {
				storiesNeedingSummary = append(storiesNeedingSummary, story)
			}
		}

		if len(storiesNeedingSummary) > 0 {
			summaries, sumErr := s.summarizer.SummarizeBatch(ctx, storiesNeedingSummary, locale)
			if sumErr != nil {
				s.logger.WarnContext(
					ctx,
					"AI summarization failed, proceeding without AI summaries",
					slog.String("error", sumErr.Error()),
				)
			} else {
				// Persist summaries and update the in-memory story objects
				for _, story := range storiesNeedingSummary {
					if summary, ok := summaries[story.StoryID]; ok {
						persistErr := s.repo.UpsertStorySummaryAI(ctx, story.StoryID, story.LocaleCode, summary)
						if persistErr != nil {
							s.logger.WarnContext(ctx, "Failed to persist AI summary",
								slog.String("story_id", story.StoryID),
								slog.String("error", persistErr.Error()))
						}

						story.SummaryAI = &summary
					}
				}
			}
		}
	}

	// Group stories by author profile
	digest := s.buildDigest(sub.ProfileID, locale, stories, now)

	// Look up the channel adapter
	channel, ok := s.channels[sub.Channel]
	if !ok {
		logEntry := &BulletinLog{
			ID:             s.idGen(),
			SubscriptionID: sub.ID,
			StoryCount:     len(stories),
			Status:         "failed",
			ErrorMessage:   strPtr(fmt.Sprintf("channel %q not available", sub.Channel)),
			CreatedAt:      now,
		}

		_ = s.repo.CreateBulletinLog(ctx, logEntry)

		return fmt.Errorf("%w: %s", ErrChannelNotAvailable, sub.Channel)
	}

	// Deliver
	sendErr := channel.Send(ctx, sub.ProfileID, digest)
	if sendErr != nil {
		logEntry := &BulletinLog{
			ID:             s.idGen(),
			SubscriptionID: sub.ID,
			StoryCount:     len(stories),
			Status:         "failed",
			ErrorMessage:   strPtr(sendErr.Error()),
			CreatedAt:      now,
		}

		_ = s.repo.CreateBulletinLog(ctx, logEntry)

		return fmt.Errorf("sending bulletin via %s: %w", sub.Channel, sendErr)
	}

	// Success — update last_bulletin_at and log
	updateErr := s.repo.UpdateLastBulletinAt(ctx, sub.ID)
	if updateErr != nil {
		s.logger.WarnContext(ctx, "Failed to update last_bulletin_at",
			slog.String("subscription_id", sub.ID),
			slog.String("error", updateErr.Error()))
	}

	logEntry := &BulletinLog{
		ID:             s.idGen(),
		SubscriptionID: sub.ID,
		StoryCount:     len(stories),
		Status:         "sent",
		CreatedAt:      now,
	}

	_ = s.repo.CreateBulletinLog(ctx, logEntry)

	s.logger.InfoContext(ctx, "Bulletin sent successfully",
		slog.String("profile_id", sub.ProfileID),
		slog.String("channel", string(sub.Channel)),
		slog.Int("story_count", len(stories)))

	return nil
}

// buildDigest groups stories by author profile into a Digest.
func (s *Service) buildDigest(
	recipientProfileID string,
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
				ProfileID: story.AuthorProfileID,
				Slug:      story.AuthorSlug,
				Title:     story.AuthorTitle,
				Stories:   make([]*DigestStory, 0),
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
	preferredTime int,
) (*Subscription, error) {
	sub := &Subscription{
		ID:            s.idGen(),
		ProfileID:     profileID,
		Channel:       channel,
		PreferredTime: preferredTime,
	}

	err := s.repo.UpsertSubscription(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("upserting subscription: %w", err)
	}

	return sub, nil
}

// UpdatePreferences updates the preferred time for a subscription.
func (s *Service) UpdatePreferences(
	ctx context.Context,
	subscriptionID string,
	preferredTime int,
) error {
	err := s.repo.UpdateSubscriptionPreferences(ctx, subscriptionID, preferredTime)
	if err != nil {
		return fmt.Errorf("updating preferences: %w", err)
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

func strPtr(s string) *string {
	return &s
}
