package workers

import (
	"time"
)

// YouTubeSyncConfig holds configuration for the YouTube sync workers.
type YouTubeSyncConfig struct {
	FullSyncEnabled         bool          `conf:"full_sync_enabled"         default:"true"`
	FullSyncInterval        time.Duration `conf:"full_sync_interval"        default:"6h"`
	IncrementalSyncEnabled  bool          `conf:"incremental_sync_enabled"  default:"true"`
	IncrementalSyncInterval time.Duration `conf:"incremental_sync_interval" default:"15m"`
	CheckInterval           time.Duration `conf:"check_interval"            default:"1m"`
	BatchSize               int           `conf:"batch_size"                default:"10"`
	StoriesPerLink          int           `conf:"stories_per_link"          default:"50"`
	FullSyncMaxStories      int           `conf:"full_sync_max_stories"     default:"1000"`
	TokenRefreshBuffer      time.Duration `conf:"token_refresh_buffer"      default:"5m"`
}

// SpeakerDeckSyncConfig holds configuration for the SpeakerDeck sync worker.
type SpeakerDeckSyncConfig struct {
	FullSyncEnabled  bool          `conf:"full_sync_enabled"  default:"true"`
	FullSyncInterval time.Duration `conf:"full_sync_interval" default:"6h"`
	CheckInterval    time.Duration `conf:"check_interval"     default:"1m"`
	BatchSize        int           `conf:"batch_size"         default:"10"`
}

// ExternalSiteSyncConfig holds configuration for the external site sync worker.
type ExternalSiteSyncConfig struct {
	FullSyncEnabled  bool          `conf:"full_sync_enabled"  default:"true"`
	FullSyncInterval time.Duration `conf:"full_sync_interval" default:"6h"`
	CheckInterval    time.Duration `conf:"check_interval"     default:"1m"`
	BatchSize        int           `conf:"batch_size"         default:"10"`
}

// DomainSyncConfig holds configuration for the custom domain sync worker.
type DomainSyncConfig struct {
	Enabled      bool          `conf:"enabled"       default:"false"`
	SyncInterval time.Duration `conf:"sync_interval" default:"5m"`
	BaseDomains  string        `conf:"base_domains"`
}

// BulletinConfig holds configuration for the bulletin digest worker.
type BulletinConfig struct {
	Enabled       bool          `conf:"enabled"        default:"false"`
	CheckInterval time.Duration `conf:"check_interval" default:"5m"`
}

// Config holds all worker configurations.
type Config struct {
	YouTubeSync      YouTubeSyncConfig        `conf:"youtube_sync"`
	GitHubSync       GitHubSyncConfig         `conf:"github_sync"`
	SpeakerDeckSync  SpeakerDeckSyncConfig    `conf:"speakerdeck_sync"`
	ExternalSiteSync ExternalSiteSyncConfig   `conf:"external_site_sync"`
	DomainSync       DomainSyncConfig         `conf:"domain_sync"`
	Queue            QueueWorkerConfig        `conf:"queue"`
	TelegramBot      TelegramBotPollingConfig `conf:"telegram_bot"`
	Bulletin         BulletinConfig           `conf:"bulletin"`
}
