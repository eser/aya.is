package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
)

// NewEnvelopeAutoRedeemer returns a callback that automatically redeems
// telegram_group invitations when they are accepted — creating a single-use
// invite link and DMing it to the recipient.
func NewEnvelopeAutoRedeemer(
	client *Client,
	telegramService *telegrambiz.Service,
	envelopeService *mailbox.Service,
	logger *logfx.Logger,
) mailbox.OnEnvelopeAcceptedFunc {
	return func(ctx context.Context, envelope *mailbox.Envelope) {
		if envelope.Kind != mailbox.KindInvitation {
			return
		}

		// Parse invitation properties.
		propsJSON, marshalErr := json.Marshal(envelope.Properties)
		if marshalErr != nil {
			logger.WarnContext(ctx, "Failed to marshal envelope properties for auto-redeem",
				slog.String("envelope_id", envelope.ID),
				slog.String("error", marshalErr.Error()))

			return
		}

		var props mailbox.InvitationProperties

		unmarshalErr := json.Unmarshal(propsJSON, &props)
		if unmarshalErr != nil {
			logger.WarnContext(ctx, "Failed to unmarshal invitation properties",
				slog.String("envelope_id", envelope.ID),
				slog.String("error", unmarshalErr.Error()))

			return
		}

		if props.InvitationKind != mailbox.InvitationKindTelegramGroup {
			return
		}

		if props.TelegramChatID == 0 {
			logger.WarnContext(ctx, "Invitation missing telegram_chat_id, skipping auto-redeem",
				slog.String("envelope_id", envelope.ID))

			return
		}

		// Look up the recipient's managed Telegram link.
		link, linkErr := telegramService.GetProfileTelegramLink(ctx, envelope.TargetProfileID)
		if linkErr != nil {
			logger.DebugContext(ctx,
				"No managed Telegram link for invitation target, skipping auto-redeem",
				slog.String("target_profile_id", envelope.TargetProfileID))

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

		// Create a single-use invite link for the group.
		inviteLink, createErr := client.CreateChatInviteLink(
			ctx, props.TelegramChatID, "AYA Invitation", 1,
		)
		if createErr != nil {
			logger.ErrorContext(ctx, "Failed to create chat invite link for auto-redeem",
				slog.String("envelope_id", envelope.ID),
				slog.Int64("chat_id", props.TelegramChatID),
				slog.String("error", createErr.Error()))

			return
		}

		// Transition the envelope to redeemed with the invite link.
		props.InviteLink = &inviteLink.InviteLink

		redeemErr := envelopeService.RedeemEnvelope(ctx, envelope.ID, props)
		if redeemErr != nil {
			logger.ErrorContext(ctx, "Failed to redeem envelope during auto-redeem",
				slog.String("envelope_id", envelope.ID),
				slog.String("error", redeemErr.Error()))
		}

		// DM the invite link to the recipient.
		groupName := props.GroupName
		if groupName == "" && envelope.Message != nil {
			groupName = *envelope.Message
		}

		text := fmt.Sprintf(
			"\U0001F389 <b>Your invitation has been accepted!</b>\n\n"+
				"Here is your single-use invite link for <b>%s</b>:\n\n%s\n\n"+
				"This link can only be used once.",
			groupName,
			inviteLink.InviteLink,
		)

		sendErr := client.SendMessage(ctx, chatID, text)
		if sendErr != nil {
			logger.WarnContext(ctx, "Failed to send auto-redeem invite link DM",
				slog.Int64("chat_id", chatID),
				slog.String("error", sendErr.Error()))
		} else {
			logger.InfoContext(ctx, "Auto-redeemed invitation and sent invite link",
				slog.String("envelope_id", envelope.ID),
				slog.String("target_profile_id", envelope.TargetProfileID),
				slog.Int64("chat_id", chatID))
		}
	}
}
