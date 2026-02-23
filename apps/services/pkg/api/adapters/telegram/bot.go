package telegram

import (
	"context"
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
	if update.Message == nil {
		return
	}

	msg := update.Message
	if msg.From == nil || msg.From.IsBot {
		return
	}

	// Handle group commands (/invite and /register are allowed in groups)
	if msg.Chat.Type == "group" || msg.Chat.Type == "supergroup" {
		text := strings.TrimSpace(msg.Text)

		switch {
		case text == "/invite" || strings.HasPrefix(text, "/invite@"):
			b.handleGroupInvite(ctx, msg)
		case text == "/register" || strings.HasPrefix(text, "/register@"):
			b.handleGroupRegister(ctx, msg)
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
		"/invite — Generate an invite code (use in a group chat)\n" +
		"/register — Register a group as a resource (use in a group chat)\n" +
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

func (b *Bot) handleGroupInvite(ctx context.Context, msg *Message) {
	if msg.From == nil || msg.From.IsBot {
		return
	}

	// Verify the user is a group administrator
	member, memberErr := b.client.GetChatMember(ctx, msg.Chat.ID, msg.From.ID)
	if memberErr != nil {
		b.logger.WarnContext(ctx, "Failed to check group membership",
			slog.Int64("chat_id", msg.Chat.ID),
			slog.Int64("user_id", msg.From.ID),
			slog.String("error", memberErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong while checking your permissions. Please try again later.")

		return
	}

	if member.Status != "creator" && member.Status != "administrator" {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Only group administrators can generate invite codes.")

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

func (b *Bot) handleGroupRegister(ctx context.Context, msg *Message) {
	if msg.From == nil || msg.From.IsBot {
		return
	}

	// Verify the user is a group administrator
	member, memberErr := b.client.GetChatMember(ctx, msg.Chat.ID, msg.From.ID)
	if memberErr != nil {
		b.logger.WarnContext(ctx, "Failed to check group membership",
			slog.Int64("chat_id", msg.Chat.ID),
			slog.Int64("user_id", msg.From.ID),
			slog.String("error", memberErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong while checking your permissions. Please try again later.")

		return
	}

	if member.Status != "creator" && member.Status != "administrator" {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Only group administrators can register this group.")

		return
	}

	// Verify the user has a linked AYA profile
	_, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Please link your Telegram account first by sending /start to me in a DM.")

		return
	}

	// Generate a group register code
	code, codeErr := b.service.GenerateGroupRegisterCode(
		ctx, msg.Chat.ID, msg.Chat.Title, msg.From.ID,
	)
	if codeErr != nil {
		b.logger.WarnContext(ctx, "Failed to generate group register code",
			slog.Int64("chat_id", msg.Chat.ID),
			slog.Int64("telegram_user_id", msg.From.ID),
			slog.String("error", codeErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong. Please try again later.")

		return
	}

	// DM the code to the user
	dmText := fmt.Sprintf(
		"Your group registration code for <b>%s</b> is:\n\n"+
			"<code>%s</code>\n\n"+
			"Go to your organization's settings on <b>aya.is</b> → Resources → "+
			"<b>Add Telegram Group</b>, and paste this code.\n\n"+
			"This code expires in 10 minutes.",
		msg.Chat.Title,
		code,
	)

	dmErr := b.client.SendMessage(ctx, msg.From.ID, dmText)
	if dmErr != nil {
		_ = b.client.SendMessage(
			ctx,
			msg.Chat.ID,
			"I couldn't send you a DM. Please start a chat with me first by sending /start in a DM.",
		)

		return
	}

	_ = b.client.SendMessage(ctx, msg.Chat.ID, "Registration code sent to your DM.")
}

func (b *Bot) handleUnknown(ctx context.Context, msg *Message) {
	_ = b.client.SendMessage(ctx, msg.Chat.ID,
		"I don't understand that command. Use /help to see available commands.")
}
