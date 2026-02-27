package telegram

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"strconv"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
)

const maxSummaryLen = 120

// BulletinSender sends bulletin digests via Telegram.
type BulletinSender struct {
	client          *Client
	telegramService *telegrambiz.Service
	logger          *logfx.Logger
	siteURI         string
}

// NewBulletinSender creates a new Telegram bulletin channel adapter.
func NewBulletinSender(
	client *Client,
	telegramService *telegrambiz.Service,
	logger *logfx.Logger,
	siteURI string,
) *BulletinSender {
	return &BulletinSender{
		client:          client,
		telegramService: telegramService,
		logger:          logger,
		siteURI:         siteURI,
	}
}

// Kind returns the channel kind identifier.
func (s *BulletinSender) Kind() bulletinbiz.ChannelKind {
	return bulletinbiz.ChannelTelegram
}

// Send delivers a digest to the recipient via Telegram.
func (s *BulletinSender) Send(
	ctx context.Context,
	recipientProfileID string,
	digest *bulletinbiz.Digest,
) error {
	link, err := s.telegramService.GetProfileTelegramLink(ctx, recipientProfileID)
	if err != nil {
		s.logger.DebugContext(ctx, "No managed Telegram link for bulletin recipient",
			slog.String("profile_id", recipientProfileID))

		return fmt.Errorf("getting telegram link: %w", err)
	}

	if link.RemoteID == "" {
		return fmt.Errorf("telegram link has empty remote_id for profile %s", recipientProfileID)
	}

	chatID, parseErr := strconv.ParseInt(link.RemoteID, 10, 64)
	if parseErr != nil {
		return fmt.Errorf("parsing telegram remote_id %q: %w", link.RemoteID, parseErr)
	}

	text, keyboard := s.buildMessage(digest)

	sendErr := s.client.SendMessageWithKeyboard(ctx, chatID, text, keyboard)
	if sendErr != nil {
		return fmt.Errorf("sending telegram message: %w", sendErr)
	}

	s.logger.InfoContext(ctx, "Bulletin sent via Telegram",
		slog.String("profile_id", recipientProfileID),
		slog.Int64("chat_id", chatID))

	return nil
}

// buildMessage constructs the HTML-formatted message and inline keyboard.
func (s *BulletinSender) buildMessage(digest *bulletinbiz.Digest) (string, *InlineKeyboardMarkup) {
	var b strings.Builder

	b.WriteString("\U0001F4EC <b>Your Daily Digest</b>\n")

	buttons := make([][]InlineKeyboardButton, 0)

	for _, group := range digest.Groups {
		b.WriteString(fmt.Sprintf("\n<b>From %s:</b>\n", html.EscapeString(group.Title)))

		for _, story := range group.Stories {
			storyURL := fmt.Sprintf("%s/%s/stories/%s", s.siteURI, digest.Locale, story.Slug)

			b.WriteString(fmt.Sprintf("• <a href=\"%s\">%s</a>\n",
				storyURL,
				html.EscapeString(story.Title)))

			summary := s.pickSummary(story)
			if summary != "" {
				b.WriteString(fmt.Sprintf("  <i>%s</i>\n", html.EscapeString(summary)))
			}

			buttons = append(buttons, []InlineKeyboardButton{
				{Text: "Read: " + truncate(story.Title, 30), URL: storyURL},
			})
		}
	}

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	return b.String(), keyboard
}

// pickSummary returns the AI summary if available, otherwise the human summary, truncated.
func (s *BulletinSender) pickSummary(story *bulletinbiz.DigestStory) string {
	text := story.Summary

	if story.SummaryAI != nil && *story.SummaryAI != "" {
		text = *story.SummaryAI
	}

	return truncate(text, maxSummaryLen)
}

// truncate shortens a string to maxLen characters, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen-1]) + "…"
}
