package workers

import "time"

// GitHubSyncConfig holds configuration for the GitHub sync worker.
type GitHubSyncConfig struct {
	Enabled          bool          `conf:"enabled"            default:"true"`
	CheckInterval    time.Duration `conf:"check_interval"     default:"5m"`
	FullSyncInterval time.Duration `conf:"full_sync_interval" default:"1h"`
	BatchSize        int           `conf:"batch_size"         default:"50"`
}
