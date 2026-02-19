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
func (b *Bot) HandleUpdate(ctx context.Context, update *Update) { //nolint:cyclop
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

	switch {
	case strings.HasPrefix(text, "/start "):
		token := strings.TrimPrefix(text, "/start ")
		b.handleStartWithToken(ctx, msg, strings.TrimSpace(token))
	case text == "/start":
		b.handleStart(ctx, msg)
	case text == "/help":
		b.handleHelp(ctx, msg)
	case text == "/status":
		b.handleStatus(ctx, msg)
	case text == "/unlink":
		b.handleUnlink(ctx, msg)
	default:
		b.handleUnknown(ctx, msg)
	}
}

func (b *Bot) handleStart(ctx context.Context, msg *Message) {
	text := "Welcome to <b>AYA Bot</b>!\n\n" +
		"To connect your Telegram account to your AYA profile, " +
		"go to your profile settings on aya.is and click <b>Connect Telegram</b>.\n\n" +
		"Commands:\n" +
		"/status — Check your linked AYA profile\n" +
		"/unlink — Disconnect your Telegram account\n" +
		"/help — Show this help message"

	_ = b.client.SendMessage(ctx, msg.Chat.ID, text)
}

func (b *Bot) handleStartWithToken(ctx context.Context, msg *Message, token string) {
	if token == "" {
		b.handleStart(ctx, msg)

		return
	}

	result, err := b.service.LinkAccountByToken(
		ctx,
		token,
		msg.From.ID,
		msg.From.Username,
	)
	if err != nil {
		b.logger.WarnContext(ctx, "Telegram link failed",
			slog.Int64("telegram_user_id", msg.From.ID),
			slog.String("error", err.Error()))

		var errMsg string

		switch {
		case errors.Is(err, telegrambiz.ErrTokenNotFound):
			errMsg = "This link has expired or is invalid. Please generate a new one from your AYA profile settings."
		case errors.Is(err, telegrambiz.ErrAlreadyLinked):
			errMsg = "Your Telegram account is already linked to an AYA profile. " +
				"Use /unlink first to disconnect, then try again."
		case errors.Is(err, telegrambiz.ErrProfileAlreadyHasTelegram):
			errMsg = "This AYA profile already has a Telegram account connected. Remove the existing connection first."
		default:
			errMsg = "Something went wrong while linking your account. Please try again later."
		}

		_ = b.client.SendMessage(ctx, msg.Chat.ID, errMsg)

		return
	}

	text := fmt.Sprintf(
		"Successfully linked to AYA profile <b>@%s</b>!\n\n"+
			"Your Telegram account is now connected and managed by AYA.",
		result.ProfileSlug,
	)

	_ = b.client.SendMessage(ctx, msg.Chat.ID, text)
}

func (b *Bot) handleHelp(ctx context.Context, msg *Message) {
	text := "<b>AYA Bot Commands</b>\n\n" +
		"/start — Welcome message and setup instructions\n" +
		"/status — Check your linked AYA profile\n" +
		"/unlink — Disconnect your Telegram account from AYA\n" +
		"/help — Show this help message"

	_ = b.client.SendMessage(ctx, msg.Chat.ID, text)
}

func (b *Bot) handleStatus(ctx context.Context, msg *Message) {
	info, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		text := "Your Telegram account is <b>not linked</b> to any AYA profile.\n\n" +
			"To connect, go to your profile settings on aya.is and click <b>Connect Telegram</b>."
		_ = b.client.SendMessage(ctx, msg.Chat.ID, text)

		return
	}

	text := fmt.Sprintf(
		"Your Telegram account is linked to AYA profile <b>%s</b>.",
		info.ProfileID,
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
