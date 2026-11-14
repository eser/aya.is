package privacy

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/events"
)

const (
	// RedactedText replaces every piece of personal text of a deleted account.
	RedactedText = "**redacted**"

	// TombstoneUserID and TombstoneProfileID identify the shared "deleted user" placeholder that inherits the
	// authorship of content a departing user chose to leave behind. Both are created lazily on the first deletion
	// (never by a migration, because the seed logic treats any existing profile as "already seeded").
	TombstoneUserID      = "0000000000000000DELETEDUSR"
	TombstoneProfileID   = "0000000000000000DELETEDUSR"
	TombstoneProfileSlug = "deleted-user"
	TombstoneTitle       = "Deleted user"

	// RedactedSlugPrefix is prepended to the profile id to build the slug a deleted profile is renamed to,
	// which frees the original username for anyone to register again.
	RedactedSlugPrefix = "deleted-"

	exportContentType = "application/zip"
)

// ExportStatus is the lifecycle state of a personal data export request.
type ExportStatus string

const (
	// ExportStatusNone is reported (never stored) when the user has never requested an export.
	ExportStatusNone       ExportStatus = "none"
	ExportStatusPending    ExportStatus = "pending"
	ExportStatusProcessing ExportStatus = "processing"
	ExportStatusReady      ExportStatus = "ready"
	ExportStatusFailed     ExportStatus = "failed"
	ExportStatusExpired    ExportStatus = "expired"
)

// Config holds the privacy feature configuration.
//
// ExportCooldown is the minimum time between two export requests of the same user. ExportRetention is how long a
// finished archive stays downloadable. DownloadURLTTL is the lifetime of the signed download URL handed to the
// browser. ExportMaxBytes aborts an export whose archive grows beyond this size (it is assembled in memory).
// FrontendURI builds the link in the "your archive is ready" e-mail.
type Config struct {
	ExportCooldown  time.Duration `conf:"export_cooldown"  default:"720h"`
	ExportRetention time.Duration `conf:"export_retention" default:"168h"`
	DownloadURLTTL  time.Duration `conf:"download_url_ttl" default:"10m"`
	ExportMaxBytes  int64         `conf:"export_max_bytes" default:"268435456"`
	FrontendURI     string        `conf:"frontend_uri"     default:""`
}

// DataExport is one request to download the user's personal data.
type DataExport struct {
	RequestedAt  time.Time
	StartedAt    *time.Time
	CompletedAt  *time.Time
	ExpiresAt    *time.Time
	NotifiedAt   *time.Time
	StorageKey   *string
	SizeBytes    *int64
	ErrorMessage *string
	ID           string
	UserID       string
	ProfileID    string
	Status       ExportStatus
}

// Account is the privacy-relevant view of a user and their individual profile.
type Account struct {
	HiddenAt    *time.Time
	UserID      string
	ProfileID   string
	ProfileSlug string
	// Redacted is true once the account has been deleted.
	Redacted bool
}

// Contact is where the "archive is ready" notice is sent.
type Contact struct {
	Email  string
	Locale string
	Slug   string
}

// OwnedProfile is a profile that has no other owner than the departing user.
type OwnedProfile struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Kind  string `json:"kind"`
	Title string `json:"title"`
}

// DeleteOptions are the choices offered when deleting an account.
type DeleteOptions struct {
	// DeleteContent erases the stories, pages, comments and questions the user wrote (their text is replaced and
	// they are removed from the site). When false they stay, attributed to the "deleted user" placeholder.
	DeleteContent bool
}

// ExportState is what the settings page shows about personal data exports.
type ExportState struct {
	RequestedAt   *time.Time   `json:"requested_at"`
	ReadyAt       *time.Time   `json:"ready_at"`
	ExpiresAt     *time.Time   `json:"expires_at"`
	NextRequestAt *time.Time   `json:"next_request_at"`
	ID            *string      `json:"id"`
	Status        ExportStatus `json:"status"`
	SizeBytes     *int64       `json:"size_bytes"`
	// Available is false when the server has no private storage configured for archives.
	Available   bool `json:"available"`
	CanRequest  bool `json:"can_request"`
	CanDownload bool `json:"can_download"`
}

// State is the full privacy state of the signed-in user.
type State struct {
	ProfileSlug string      `json:"profile_slug"`
	Export      ExportState `json:"export"`
	IsHidden    bool        `json:"is_hidden"`
}

// Repository is the storage port of the privacy service.
//
//nolint:interfacebloat // one port: account state, export bookkeeping, dataset streaming and the erasure transaction
type Repository interface {
	GetAccount(ctx context.Context, userID string) (*Account, error)
	GetContact(ctx context.Context, userID string) (*Contact, error)

	// HideProfile hides the profile and ends every active session of the user in one transaction.
	HideProfile(ctx context.Context, userID string, profileID string) (bool, error)
	// UnhideProfile restores a self-hidden profile; it reports whether anything changed.
	UnhideProfile(ctx context.Context, userID string) (bool, error)
	ListSoleOwnedProfiles(ctx context.Context, profileID string) ([]OwnedProfile, error)

	CreateExport(ctx context.Context, export *DataExport) error
	GetLatestExport(ctx context.Context, userID string) (*DataExport, error)
	GetExport(ctx context.Context, id string) (*DataExport, error)
	ClaimExport(ctx context.Context, id string) (bool, error)
	MarkExportReady(
		ctx context.Context,
		id string,
		storageKey string,
		sizeBytes int64,
		expiresAt time.Time,
	) error
	MarkExportFailed(ctx context.Context, id string, message string) error
	MarkExportNotified(ctx context.Context, id string) error
	ListExpiredExports(ctx context.Context, limit int) ([]*DataExport, error)
	MarkExportExpired(ctx context.Context, id string) error
	ListExportsByUser(ctx context.Context, userID string) ([]*DataExport, error)

	// ExportDatasetNames lists, in a stable order, the datasets an archive is made of in this database.
	ExportDatasetNames(ctx context.Context) ([]string, error)
	// StreamDataset calls fn once per JSON document of the dataset.
	StreamDataset(
		ctx context.Context,
		name string,
		userID string,
		profileID string,
		fn func(row json.RawMessage) error,
	) error

	// RedactAccount erases the user's personal data in one transaction. It must be idempotent.
	RedactAccount(ctx context.Context, userID string, profileID string, opts DeleteOptions) error
}

// ObjectStorage is the private object storage that holds finished archives.
type ObjectStorage interface {
	PutObject(ctx context.Context, key string, contentType string, body io.ReadSeeker) error
	PresignedDownloadURL(
		ctx context.Context,
		key string,
		filename string,
		expiresIn time.Duration,
	) (string, error)
	RemoveObject(ctx context.Context, key string) error
}

// Enqueuer schedules background work.
type Enqueuer interface {
	Enqueue(ctx context.Context, params events.QueueEnqueueParams) (string, error)
}

// Mailer notifies a user that their archive is ready.
type Mailer interface {
	SendExportReady(
		ctx context.Context,
		to string,
		locale string,
		settingsURL string,
		retentionDays int,
	) error
}
