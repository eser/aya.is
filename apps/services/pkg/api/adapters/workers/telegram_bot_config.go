package workers

import "time"

// TelegramBotPollingConfig holds configuration for the Telegram bot polling worker.
type TelegramBotPollingConfig struct {
	Enabled      bool          `conf:"enabled"       default:"false"`
	PollInterval time.Duration `conf:"poll_interval" default:"1s"`
	PollTimeout  int           `conf:"poll_timeout"  default:"30"` // Telegram long poll timeout in seconds
}
