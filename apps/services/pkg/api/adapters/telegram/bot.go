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
	case "/join":
		b.handleJoinDirect(ctx, msg)
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
		"/join — List private groups you can join\n" +
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

		if link.LinkVisibility != "" && link.LinkVisibility != "public" {
			label += " (" + link.LinkVisibility + "+)"
		}

		if link.LinkURI != "" {
			builder.WriteString(fmt.Sprintf("  • <a href=\"%s\">%s</a>\n", link.LinkURI, label))
		} else {
			builder.WriteString(fmt.Sprintf("  • %s\n", label))
		}
	}

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "List invite-only groups I can join", CallbackData: "join"},
			},
		},
	}

	_ = b.client.SendMessageWithKeyboard(ctx, msg.Chat.ID, builder.String(), keyboard)
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

func (b *Bot) handleGroupRegister(ctx context.Context, msg *Message) { //nolint:cyclop,funlen
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

	// Verify the bot is an administrator in this group
	botUserID, botErr := b.client.GetBotUserID(ctx)
	if botErr != nil {
		b.logger.WarnContext(ctx, "Failed to get bot user ID",
			slog.String("error", botErr.Error()))

		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Something went wrong. Please try again later.")

		return
	}

	botMember, botMemberErr := b.client.GetChatMember(ctx, msg.Chat.ID, botUserID)
	if botMemberErr != nil || botMember.Status != "administrator" {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"I need to be an <b>administrator</b> in this group to manage invites.\n"+
				"Please promote me to admin, then try /register again.")

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

func (b *Bot) handleCallbackQuery(ctx context.Context, cq *CallbackQuery) { //nolint:varnamelen
	if cq.From == nil || cq.From.IsBot {
		return
	}

	switch {
	case cq.Data == "join":
		b.handleJoin(ctx, cq)
	case strings.HasPrefix(cq.Data, "join:"):
		b.handleJoinGroup(ctx, cq)
	case strings.HasPrefix(cq.Data, "joined:"):
		_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "You're already a member of this group.")

		return
	}

	// Answer callback to remove loading indicator
	_ = b.client.AnswerCallbackQuery(ctx, cq.ID, "")
}

func (b *Bot) handleJoinDirect(ctx context.Context, msg *Message) {
	info, err := b.service.GetLinkedProfile(ctx, msg.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, msg.Chat.ID,
			"Your Telegram account is <b>not linked</b> to any AYA profile.\n\n"+
				"Send /start to get a verification code.")

		return
	}

	b.sendJoinableGroups(ctx, msg.Chat.ID, msg.From.ID, info.ProfileID)
}

func (b *Bot) handleJoin(ctx context.Context, cq *CallbackQuery) { //nolint:varnamelen
	info, err := b.service.GetLinkedProfile(ctx, cq.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, cq.From.ID,
			"Your Telegram account is <b>not linked</b> to any AYA profile.\n\n"+
				"Send /start to get a verification code.")

		return
	}

	chatID := cq.From.ID
	if cq.Message != nil {
		chatID = cq.Message.Chat.ID
	}

	b.sendJoinableGroups(ctx, chatID, cq.From.ID, info.ProfileID)
}

func (b *Bot) sendJoinableGroups( //nolint:cyclop,funlen
	ctx context.Context,
	chatID int64,
	telegramUserID int64,
	profileID string,
) {
	groups, groupsErr := b.service.GetEligibleTelegramGroups(ctx, profileID)
	if groupsErr != nil {
		b.logger.WarnContext(ctx, "Failed to get eligible telegram groups",
			slog.String("profile_id", profileID),
			slog.String("error", groupsErr.Error()))

		_ = b.client.SendMessage(ctx, chatID,
			"Something went wrong. Please try again later.")

		return
	}

	if len(groups) == 0 {
		_ = b.client.SendMessage(ctx, chatID,
			"No Telegram groups are available for your teams.")

		return
	}

	// Build text list with full names and buttons only for joinable groups
	var textBuilder strings.Builder

	textBuilder.WriteString(
		"Pick a button below to get invited to the selected group. Available private groups are:\n\n",
	)

	var buttons [][]InlineKeyboardButton

	for _, group := range groups {
		member, memberErr := b.client.GetChatMember(ctx, group.ChatID, telegramUserID)

		isJoined := memberErr == nil && member != nil &&
			(member.Status == "member" || member.Status == "administrator" || member.Status == "creator")

		label := group.GroupTitle
		if group.ProfileSlug != "" {
			label += " (" + group.ProfileSlug + ")"
		}

		if isJoined {
			textBuilder.WriteString("\u2713 " + label + "\n")
		} else {
			textBuilder.WriteString("\u25CB " + label + "\n")

			buttons = append(buttons, []InlineKeyboardButton{
				{Text: label, CallbackData: "join:" + group.ResourceID},
			})
		}
	}

	textBuilder.WriteString("\n\u2713 = already a member")

	if len(buttons) > 0 {
		keyboard := &InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		}

		_ = b.client.SendMessageWithKeyboard(ctx, chatID, textBuilder.String(), keyboard)
	} else {
		_ = b.client.SendMessage(ctx, chatID, textBuilder.String())
	}
}

func (b *Bot) handleJoinGroup( //nolint:cyclop,funlen
	ctx context.Context,
	cq *CallbackQuery, //nolint:varnamelen
) {
	resourceID := strings.TrimPrefix(cq.Data, "join:")

	info, err := b.service.GetLinkedProfile(ctx, cq.From.ID)
	if err != nil {
		_ = b.client.SendMessage(ctx, cq.From.ID,
			"Your Telegram account is <b>not linked</b> to any AYA profile.\n\n"+
				"Send /start to get a verification code.")

		return
	}

	// Re-fetch eligible groups to validate access
	groups, groupsErr := b.service.GetEligibleTelegramGroups(ctx, info.ProfileID)
	if groupsErr != nil {
		_ = b.client.SendMessage(ctx, cq.From.ID,
			"Something went wrong. Please try again later.")

		return
	}

	// Find the requested group
	var target *telegrambiz.EligibleTelegramGroup

	for i := range groups {
		if groups[i].ResourceID == resourceID {
			target = &groups[i]

			break
		}
	}

	if target == nil {
		_ = b.client.SendMessage(ctx, cq.From.ID,
			"You don't have access to this group, or it no longer exists.")

		return
	}

	// Check if already a member
	member, memberErr := b.client.GetChatMember(ctx, target.ChatID, cq.From.ID)
	if memberErr == nil && member != nil &&
		(member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
		_ = b.client.SendMessage(ctx, cq.From.ID,
			fmt.Sprintf("You're already a member of <b>%s</b>.", target.GroupTitle))

		return
	}

	// Generate single-use invite link
	inviteName := "AYA /join"
	if cq.From.Username != "" {
		inviteName += " for @" + cq.From.Username
	}

	inviteLink, createErr := b.client.CreateChatInviteLink(ctx, target.ChatID, inviteName, 1)
	if createErr != nil {
		b.logger.ErrorContext(ctx, "Failed to create invite link",
			slog.Int64("chat_id", target.ChatID),
			slog.String("resource_id", resourceID),
			slog.String("error", createErr.Error()))

		_ = b.client.SendMessage(ctx, cq.From.ID,
			"Failed to create an invite link. The bot may not be an administrator in this group.")

		return
	}

	// Record audit entry for leak detection
	b.service.RecordInviteLinkGenerated(
		ctx,
		info.ProfileID,
		cq.From.ID,
		target.ChatID,
		target.GroupTitle,
		inviteLink.InviteLink,
	)

	text := fmt.Sprintf(
		"Here is your invite link for <b>%s</b>:\n\n%s\n\n"+
			"This link can only be used once.",
		target.GroupTitle,
		inviteLink.InviteLink,
	)

	_ = b.client.SendMessage(ctx, cq.From.ID, text)
}

func (b *Bot) handleUnknown(ctx context.Context, msg *Message) {
	_ = b.client.SendMessage(ctx, msg.Chat.ID,
		"I don't understand that command. Use /help to see available commands.")
}
