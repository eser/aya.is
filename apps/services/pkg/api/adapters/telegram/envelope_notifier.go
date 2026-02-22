package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	envelopes "github.com/eser/aya.is/services/pkg/api/business/profile_envelopes"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
)

// EnvelopeNotifier sends Telegram notifications when new envelopes are created.
// It implements the profile_envelopes.EnvelopeNotifier port interface.
type EnvelopeNotifier struct {
	client          *Client
	telegramService *telegrambiz.Service
	logger          *logfx.Logger
}

// NewEnvelopeNotifier creates a new Telegram envelope notifier adapter.
func NewEnvelopeNotifier(
	client *Client,
	telegramService *telegrambiz.Service,
	logger *logfx.Logger,
) *EnvelopeNotifier {
	return &EnvelopeNotifier{
		client:          client,
		telegramService: telegramService,
		logger:          logger,
	}
}

// NotifyNewEnvelope looks up the target profile's linked Telegram account
// and sends a notification message. This is best-effort — errors are logged
// but never propagated.
func (n *EnvelopeNotifier) NotifyNewEnvelope(
	ctx context.Context,
	params *envelopes.EnvelopeNotification,
) {
	// Look up the target profile's telegram link
	link, err := n.telegramService.GetProfileTelegramLink(ctx, params.TargetProfileID)
	if err != nil {
		// No linked Telegram account — nothing to notify, this is normal
		return
	}

	chatID, parseErr := strconv.ParseInt(link.RemoteID, 10, 64)
	if parseErr != nil {
		n.logger.WarnContext(ctx, "Invalid telegram remote_id for envelope notification",
			slog.String("target_profile_id", params.TargetProfileID),
			slog.String("remote_id", link.RemoteID))

		return
	}

	fromLabel := "someone"
	if params.SenderProfileTitle != "" {
		fromLabel = params.SenderProfileTitle
	}

	locale := params.Locale
	if locale == "" {
		locale = "en"
	}

	text := fmt.Sprintf(
		"📬 <b>New message in your Mailbox</b>\n\n"+
			"<b>%s</b>\nFrom: <i>%s</i>\n\n"+
			"<a href=\"https://aya.is/%s/mailbox\">Open Mailbox</a>",
		params.EnvelopeTitle,
		fromLabel,
		locale,
	)

	sendErr := n.client.SendMessage(ctx, chatID, text)
	if sendErr != nil {
		n.logger.WarnContext(ctx, "Failed to send envelope notification via Telegram",
			slog.String("target_profile_id", params.TargetProfileID),
			slog.Int64("chat_id", chatID),
			slog.String("error", sendErr.Error()))
	} else {
		n.logger.InfoContext(ctx, "Envelope notification sent via Telegram",
			slog.String("target_profile_id", params.TargetProfileID),
			slog.Int64("chat_id", chatID))
	}
}
