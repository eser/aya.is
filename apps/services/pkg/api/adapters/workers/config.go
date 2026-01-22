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
	BatchSize               int           `conf:"batch_size"                default:"10"`
	StoriesPerLink          int           `conf:"stories_per_link"          default:"50"`
	TokenRefreshBuffer      time.Duration `conf:"token_refresh_buffer"      default:"5m"`
}

// Config holds all worker configurations.
type Config struct {
	YouTubeSync YouTubeSyncConfig `conf:"youtube_sync"`
}
