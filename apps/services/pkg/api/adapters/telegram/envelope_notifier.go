package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
)

// NewEnvelopeNotifier returns a callback that sends a Telegram DM to the
// envelope recipient when they have a linked Telegram account.
func NewEnvelopeNotifier(
	client *Client,
	telegramService *telegrambiz.Service,
	logger *logfx.Logger,
) mailbox.OnEnvelopeCreatedFunc {
	return func(ctx context.Context, envelope *mailbox.Envelope, params *mailbox.SendMessageParams) {
		link, err := telegramService.GetProfileTelegramLink(ctx, envelope.TargetProfileID)
		if err != nil {
			logger.DebugContext(
				ctx,
				"No managed Telegram link for envelope target, skipping notification",
				slog.String("target_profile_id", envelope.TargetProfileID),
			)

			return
		}

		if link.RemoteID == "" {
			logger.WarnContext(ctx, "Telegram link has empty remote_id",
				slog.String("target_profile_id", envelope.TargetProfileID))

			return
		}

		chatID, parseErr := strconv.ParseInt(link.RemoteID, 10, 64)
		if parseErr != nil {
			logger.WarnContext(ctx, "Invalid telegram remote_id",
				slog.String("remote_id", link.RemoteID))

			return
		}

		fromLabel := params.SenderProfileTitle
		if fromLabel == "" {
			fromLabel = "someone"
		}

		locale := params.Locale
		if locale == "" {
			locale = "en"
		}

		text := fmt.Sprintf(
			"\U0001F4EC <b>New message in your Mailbox</b>\n\n"+
				"<b>%s</b>\nFrom: <i>%s</i>\n\n"+
				"<a href=\"https://aya.is/%s/mailbox\">Open Mailbox</a>",
			envelope.Title,
			fromLabel,
			locale,
		)

		sendErr := client.SendMessage(ctx, chatID, text)
		if sendErr != nil {
			logger.WarnContext(ctx, "Failed to send envelope notification",
				slog.Int64("chat_id", chatID),
				slog.String("error", sendErr.Error()))
		} else {
			logger.InfoContext(ctx, "Envelope notification sent via Telegram",
				slog.String("target_profile_id", envelope.TargetProfileID),
				slog.Int64("chat_id", chatID))
		}
	}
}
