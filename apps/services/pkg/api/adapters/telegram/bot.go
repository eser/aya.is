package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
)

// Bot handles incoming Telegram updates and routes commands.
type Bot struct {
	client  *Client
	service *telegrambiz.Service
	logger  *logfx.Logger
}

// NewBot creates a new bot handler.
func NewBot(client *Client, service *telegrambiz.Service, logger *logfx.Logger) *Bot {
	return &Bot{
		client:  client,
		service: service,
		logger:  logger,
	}
}

// Client returns the underlying Telegram API client.
func (b *Bot) Client() *Client {
	return b.client
}

// HandleUpdate processes a single Telegram update.
func (b *Bot) HandleUpdate(ctx context.Context, update *Update) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	if msg.From == nil || msg.From.IsBot {
		return
	}

	// Only handle private messages
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

func (b *Bot) handleUnknown(ctx context.Context, msg *Message) {
	_ = b.client.SendMessage(ctx, msg.Chat.ID,
		"I don't understand that command. Use /help to see available commands.")
}
