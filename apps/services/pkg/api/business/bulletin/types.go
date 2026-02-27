package bulletin

import "time"

// ChannelKind identifies the delivery channel for a bulletin.
type ChannelKind string

const (
	ChannelTelegram ChannelKind = "telegram"
	ChannelEmail    ChannelKind = "email"
)

// DigestFrequency identifies how often bulletins are sent.
type DigestFrequency string

const (
	FrequencyDaily   DigestFrequency = "daily"
	FrequencyBiDaily DigestFrequency = "bidaily"
	FrequencyWeekly  DigestFrequency = "weekly"
)

// Subscription represents a user's bulletin subscription for a specific channel.
type Subscription struct {
	CreatedAt      time.Time
	LastBulletinAt *time.Time
	UpdatedAt      *time.Time
	ID             string
	ProfileID      string
	ProfileSlug    string
	Channel        ChannelKind
	Frequency      DigestFrequency
	DefaultLocale  string
	PreferredTime  int // UTC hour (0-23)
}

// DigestStory holds the data needed to render a single story in the digest.
type DigestStory struct {
	PublishedAt             *time.Time
	StoryPictureURI         *string
	SummaryAI               *string
	StoryID                 string
	Slug                    string
	Kind                    string
	LocaleCode              string
	Title                   string
	Summary                 string
	AuthorProfileID         string
	AuthorSlug              string
	AuthorTitle             string
	AuthorProfilePictureURI *string
}

// DigestGroup groups stories by the followed profile (author).
type DigestGroup struct {
	ProfilePictureURI *string
	ProfileID         string
	Slug              string
	Title             string
	Stories           []*DigestStory
}

// Digest is the fully composed bulletin ready for delivery.
type Digest struct {
	GeneratedAt        time.Time
	RecipientProfileID string
	RecipientSlug      string
	Locale             string
	Groups             []*DigestGroup
}

// BulletinLog records a bulletin delivery attempt.
type BulletinLog struct {
	CreatedAt      time.Time
	ErrorMessage   *string
	ID             string
	SubscriptionID string
	Status         string // "sent" or "failed"
	StoryCount     int
}
