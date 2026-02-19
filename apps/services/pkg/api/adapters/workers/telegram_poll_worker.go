package workers

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	telegramadapter "github.com/eser/aya.is/services/pkg/api/adapters/telegram"
)

// TelegramPollWorker polls the Telegram Bot API for updates (dev mode alternative to webhooks).
type TelegramPollWorker struct {
	config       *TelegramBotPollingConfig
	logger       *logfx.Logger
	client       *telegramadapter.Client
	bot          *telegramadapter.Bot
	lastUpdateID atomic.Int64
}

// NewTelegramPollWorker creates a new Telegram polling worker.
func NewTelegramPollWorker(
	config *TelegramBotPollingConfig,
	logger *logfx.Logger,
	client *telegramadapter.Client,
	bot *telegramadapter.Bot,
) *TelegramPollWorker {
	return &TelegramPollWorker{
		config: config,
		logger: logger,
		client: client,
		bot:    bot,
	}
}

// Name returns the worker's unique identifier.
func (w *TelegramPollWorker) Name() string {
	return "telegram-bot-poll"
}

// Interval returns how often the worker should execute.
// Returns 0 for continuous execution (long polling blocks internally).
func (w *TelegramPollWorker) Interval() time.Duration {
	return w.config.PollInterval
}

// Execute fetches and processes Telegram updates.
func (w *TelegramPollWorker) Execute(ctx context.Context) error {
	offset := w.lastUpdateID.Load() + 1

	updates, err := w.client.GetUpdates(ctx, offset, w.config.PollTimeout)
	if err != nil {
		w.logger.WarnContext(ctx, "Telegram poll failed",
			slog.String("error", err.Error()))

		return nil // Don't stop the worker on transient errors
	}

	for i := range updates {
		update := &updates[i]
		w.lastUpdateID.Store(update.UpdateID)
		w.bot.HandleUpdate(ctx, update)
	}

	return nil
}
