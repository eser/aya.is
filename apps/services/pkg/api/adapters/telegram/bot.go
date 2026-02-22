package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
)

// Bot handles incoming Telegram updates and routes commands.
type Bot struct {
	client          *Client
	service         *telegrambiz.Service
	envelopeService *mailbox.Service
	logger          *logfx.Logger
}

// NewBot creates a new bot handler.
func NewBot(
	client *Client,
	service *telegrambiz.Service,
	envelopeService *mailbox.Service,
	logger *logfx.Logger,
) *Bot {
	return &Bot{
		client:          client,
		service:         service,
		envelopeService: envelopeService,
		logger:          logger,
	}
}

// Client returns the underlying Telegram API client.
func (b *Bot) Client() *Client {
	return b.client
}

// HandleUpdate processes a single Telegram update.
func (b *Bot) HandleUpdate(ctx context.Context, update *Update) { //nolint:cyclop
	// Handle callback queries from inline keyboard buttons
	if update.CallbackQuery != nil {
		b.handleCallbackQuery(ctx, update.CallbackQuery)

		return
	}

	if update.Message == nil {
		return
	}

	msg := update.Message
	if msg.From == nil || msg.From.IsBot {
		return
	}

	// Handle group commands (only /invite is allowed in groups)
	if msg.Chat.Type == "group" || msg.Chat.Type == "supergroup" {
		text := strings.TrimSpace(msg.Text)
		if text == "/invite" || strings.HasPrefix(text, "/invite@") {
			b.handleGroupInvite(ctx, msg)
		}

		return
	}

	// Only handle private messages for all other commands
	if msg.Chat.Type != "private" {
		return
	}

	text := strings.TrimSpace(msg.Text)

	switch text {
	case "/start":
		b.handleStart(ctx, msg)
	case "/help":
		b.handleHelp(ctx, msg)
	case "/status":
		b.handleStatus(ctx, msg)
	case "/groups":
		b.handleGroups(ctx, msg)
	case "/invitations":
		b.handleInvitations(ctx, msg)
	case "/unlink":
		b.handleUnlink(ctx, msg)
	default:
		b.handleUnknown(ctx, msg)
	}
}

func (b *Bot) handleStart(ctx context.Context, msg *Message) {
	// Check if already linked
	existing, existingErr := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if existingErr == nil && existing != nil {
		slug, slugErr := b.service.GetProfileSlugByID(ctx, existing.ProfileID)
		if slugErr != nil || slug == "" {
			slug = existing.ProfileID
		}

		text := fmt.Sprintf(
			"Your Telegram account is already linked to AYA profile <b>@%s</b>.\n\n"+
				"Use /unlink first to disconnect, then try again.",
			slug,
		)

		_ = b.client.SendMessage(ctx, msg.Chat.ID, text)

		return
	}

	// Generate verification code
	code, err := b.service.GenerateVerificationCode(ctx, msg.From.ID, msg.From.Username)
	if err != nil {
		b.logger.WarnContext(ctx, "Failed to generate verification code",
			slog.Int64("telegram_user_id", msg.From.ID),
			slog.String("error", err.Error()))

		if errors.Is(err, telegrambiz.ErrAlreadyLinked) {
			_ = b.client.SendMessage(ctx, msg.Chat.ID,
				"Your Telegram account is already linked to an AYA profile.\n"+
					"Use /unlink first to disconnect, then try again.")
		} else {
			_ = b.client.SendMessage(ctx, msg.Chat.ID,
				"Something went wrong. Please try again later.")
		}

		return
	}

	text := fmt.Sprintf(
		"Your verification code is:\n\n"+
			"<code>%s</code>\n\n"+
			"Go to your profile settings on <b>aya.is</b>, "+
			"click <b>Connect Telegram</b>, and paste this code.\n\n"+
			"This code expires in 10 minutes.",
		code,
	)

	_ = b.client.SendMessage(ctx, msg.Chat.ID, text)
}

func (b *Bot) handleHelp(ctx context.Context, msg *Message) {
	text := "<b>AYA Bot Commands</b>\n\n" +
		"/start — Get a verification code to link your account\n" +
		"/status — Check your linked AYA profile\n" +
		"/groups — List Telegram groups from your profiles\n" +
		"/invitations — View and redeem your accepted invitations\n" +
		"/invite — Generate an invite code (use in a group chat)\n" +
		"/unlink — Disconnect your Telegram account from AYA\n" +
		"/help — Show this help message"

	_ = b.client.SendMessage(ctx, msg.Chat.ID, text)
}

func (b *Bot) handleStatus(ctx context.Context, msg *Message) {
	info, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		text := "Your Telegram account is <b>not linked</b> to any AYA profile.\n\n" +
			"Send /start to get a verification code."
		_ = b.client.SendMessage(ctx, msg.Chat.ID, text)

		return
	}

	slug, slugErr := b.service.GetProfileSlugByID(ctx, info.ProfileID)
	if slugErr != nil || slug == "" {
		slug = info.ProfileID
	}

	text := fmt.Sprintf(
		"Your Telegram account is linked to AYA profile <b>@%s</b>.",
		slug,
	)

	_ = b.client.SendMessage(ctx, msg.Chat.ID, text)
}

func (b *Bot) handleUnlink(ctx context.Context, msg *Message) {
	err := b.service.UnlinkAccount(ctx, msg.From.ID)
	if err != nil {
		if errors.Is(err, telegrambiz.ErrNotLinked) {
			_ = b.client.SendMessage(ctx, msg.Chat.ID,
				"Your Telegram account is not linked to any AYA profile.")
		} else {
			_ = b.client.SendMessage(ctx, msg.Chat.ID,
				"Something went wrong while unlinking your account. Please try again later.")
		}

		return
	}

	_ = b.client.SendMessage(ctx, msg.Chat.ID,
		"Your Telegram account has been disconnected from AYA.")
}

func (b *Bot) handleGroups(ctx context.Context, msg *Message) {
	info, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Your Telegram account is <b>not linked</b> to any AYA profile.\n\n"+
				"Send /start to get a verification code.")

		return
	}

	links, linksErr := b.service.GetGroupTelegramLinks(ctx, info.ProfileID)
	if linksErr != nil {
		b.logger.WarnContext(ctx, "Failed to get group telegram links",
			slog.String("profile_id", info.ProfileID),
			slog.String("error", linksErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong. Please try again later.")

		return
	}

	if len(links) == 0 {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"You don't have access to any Telegram groups yet.")

		return
	}

	var builder strings.Builder

	builder.WriteString("<b>Telegram Groups</b>\n")

	currentSlug := ""

	for _, link := range links {
		if link.ProfileSlug != currentSlug {
			currentSlug = link.ProfileSlug
			builder.WriteString(
				fmt.Sprintf(
					"\n<b>%s</b> (https://aya.is/%s)\n",
					link.ProfileTitle,
					link.ProfileSlug,
				),
			)
		}

		label := link.LinkTitle
		if label == "" && link.LinkPublicID != "" {
			label = "@" + link.LinkPublicID
		} else if label == "" {
			label = "Telegram"
		}

		if link.LinkURI != "" {
			builder.WriteString(fmt.Sprintf("  • <a href=\"%s\">%s</a>\n", link.LinkURI, label))
		} else {
			builder.WriteString(fmt.Sprintf("  • %s\n", label))
		}
	}

	_ = b.client.SendMessageWithOpts(ctx, msg.Chat.ID, builder.String(), SendMessageOpts{
		DisableLinkPreview: true,
	})
}

func (b *Bot) handleInvitations(ctx context.Context, msg *Message) {
	info, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Your Telegram account is <b>not linked</b> to any AYA profile.\n\n"+
				"Send /start to get a verification code.")

		return
	}

	invitations, invErr := b.envelopeService.GetAcceptedInvitations(
		ctx, info.ProfileID, mailbox.InvitationKindTelegramGroup,
	)
	if invErr != nil {
		b.logger.WarnContext(ctx, "Failed to get invitations",
			slog.String("profile_id", info.ProfileID),
			slog.String("error", invErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong. Please try again later.")

		return
	}

	if len(invitations) == 0 {
		noInvMsg := "You have no pending invitations to redeem.\n\n" +
			"Check your Inbox on <b>aya.is</b> for new invitations."

		pendingCount, countErr := b.envelopeService.CountPendingEnvelopes(ctx, info.ProfileID)
		if countErr == nil && pendingCount > 0 {
			noInvMsg = fmt.Sprintf(
				"You have no pending invitations to redeem.\n\n"+
					"Check your Inbox on <b>aya.is</b> for new invitations, "+
					"there are still <b>%d</b> messages not yet responded.",
				pendingCount,
			)
		}

		_ = b.client.SendMessage(ctx, msg.Chat.ID, noInvMsg)

		return
	}

	// Build inline keyboard with one button per invitation
	var (
		rows    [][]InlineKeyboardButton
		details []string
	)

	for i, inv := range invitations {
		// Extract group name from properties for display
		buttonText := inv.Title
		groupName := ""

		if inv.Properties != nil {
			if propsMap, ok := inv.Properties.(map[string]any); ok {
				if gn, gnOk := propsMap["group_name"].(string); gnOk && gn != "" {
					groupName = gn
					buttonText = fmt.Sprintf("%s → %s", inv.Title, gn)
				}
			}
		}

		rows = append(rows, []InlineKeyboardButton{
			{Text: buttonText, CallbackData: "invite:" + inv.ID},
		})

		detail := fmt.Sprintf("%d. <b>%s</b>", i+1, inv.Title)
		if groupName != "" {
			detail += fmt.Sprintf("\n    Group: <i>%s</i>", groupName)
		}

		if inv.SenderProfileTitle != nil && *inv.SenderProfileTitle != "" {
			detail += fmt.Sprintf("\n    From: <i>%s</i>", *inv.SenderProfileTitle)
		} else if inv.SenderProfileID == nil {
			detail += "\n    From: <i>System</i>"
		}

		details = append(details, detail)
	}

	keyboard := &InlineKeyboardMarkup{InlineKeyboard: rows}

	msgText := "<b>Your Invitations</b>\n\n" +
		strings.Join(details, "\n") +
		"\n\nTap an invitation to get your invite link:"

	_ = b.client.SendMessageWithKeyboard(ctx, msg.Chat.ID,
		msgText,
		keyboard,
	)
}

func (b *Bot) handleGroupInvite(ctx context.Context, msg *Message) {
	if msg.From == nil || msg.From.IsBot {
		return
	}

	// Verify the user has a linked AYA profile
	_, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Please link your Telegram account first by sending /start to me in a DM.")

		return
	}

	// Generate a group invite code
	code, codeErr := b.service.GenerateGroupInviteCode(
		ctx, msg.Chat.ID, msg.Chat.Title, msg.From.ID,
	)
	if codeErr != nil {
		b.logger.WarnContext(ctx, "Failed to generate group invite code",
			slog.Int64("chat_id", msg.Chat.ID),
			slog.Int64("telegram_user_id", msg.From.ID),
			slog.String("error", codeErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong. Please try again later.")

		return
	}

	// DM the code to the user
	dmText := fmt.Sprintf(
		"Your group invite code for <b>%s</b> is:\n\n"+
			"<code>%s</code>\n\n"+
			"Paste this code in the invitation form on <b>aya.is</b> within 10 minutes.\n"+
			"Once the invitation is sent, the recipient can accept it at any time.",
		msg.Chat.Title,
		code,
	)

	dmErr := b.client.SendMessage(ctx, msg.From.ID, dmText)
	if dmErr != nil {
		// User may not have started a DM with the bot yet
		_ = b.client.SendMessage(
			ctx,
			msg.Chat.ID,
			"I couldn't send you a DM. Please start a chat with me first by sending /start in a DM.",
		)

		return
	}

	_ = b.client.SendMessage(ctx, msg.Chat.ID, "Invite code sent to your DM.")
}

func (b *Bot) handleCallbackQuery(ctx context.Context, cq *CallbackQuery) { //nolint:funlen,cyclop
	if cq.From == nil || cq.From.IsBot {
		return
	}

	if !strings.HasPrefix(cq.Data, "invite:") {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Unknown action.")

		return
	}

	envelopeID := strings.TrimPrefix(cq.Data, "invite:")

	// Verify the user has a linked profile
	info, err := b.service.GetLinkedProfile(ctx, cq.From.ID)
	if err != nil {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Your account is not linked.")

		return
	}

	// Get the envelope
	envelope, envErr := b.envelopeService.GetEnvelopeByID(ctx, envelopeID)
	if envErr != nil {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Invitation not found.")

		return
	}

	// Verify it belongs to this user's profile and is accepted
	if envelope.TargetProfileID != info.ProfileID {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "This invitation is not for you.")

		return
	}

	if envelope.Status != mailbox.StatusAccepted {
		_ = b.client.AnswerCallbackQuery(
			ctx,
			cq.ID,
			"This invitation has already been used or is no longer valid.",
		)

		return
	}

	// Parse invitation properties to get the Telegram chat ID
	propsJSON, marshalErr := json.Marshal(envelope.Properties)
	if marshalErr != nil {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Invalid invitation data.")

		return
	}

	var props mailbox.InvitationProperties

	unmarshalErr := json.Unmarshal(propsJSON, &props)
	if unmarshalErr != nil {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Invalid invitation data.")

		return
	}

	if props.TelegramChatID == 0 {
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Invalid invitation: missing chat ID.")

		return
	}

	// Create a single-use invite link
	inviteLink, linkErr := b.client.CreateChatInviteLink(
		ctx, props.TelegramChatID, "AYA Invitation", 1,
	)
	if linkErr != nil {
		b.logger.ErrorContext(ctx, "Failed to create chat invite link",
			slog.String("envelope_id", envelopeID),
			slog.Int64("chat_id", props.TelegramChatID),
			slog.String("error", linkErr.Error()))

		_ = b.client.AnswerCallbackQuery(
			ctx,
			cq.ID,
			"Failed to create invite link. Please try again later.",
		)

		return
	}

	// Update properties with the generated invite link and redeem the envelope
	props.InviteLink = &inviteLink.InviteLink

	redeemErr := b.envelopeService.RedeemEnvelope(ctx, envelopeID, props)
	if redeemErr != nil {
		b.logger.ErrorContext(ctx, "Failed to redeem envelope",
			slog.String("envelope_id", envelopeID),
			slog.String("error", redeemErr.Error()))
	}

	// Send the invite link to the user
	chatID := cq.From.ID
	if cq.Message != nil && cq.Message.Chat != nil {
		chatID = cq.Message.Chat.ID
	}

	groupName := props.GroupName
	if groupName == "" {
		groupName = envelope.Title
	}

	text := fmt.Sprintf(
		"Here is your single-use invite link for <b>%s</b>:\n\n%s\n\n"+
			"This link can only be used once.",
		groupName,
		inviteLink.InviteLink,
	)

	_ = b.client.SendMessage(ctx, chatID, text)
	_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "Invite link sent!")
}

func (b *Bot) handleUnknown(ctx context.Context, msg *Message) {
	_ = b.client.SendMessage(ctx, msg.Chat.ID,
		"I don't understand that command. Use /help to see available commands.")
}
