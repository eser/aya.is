package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// Sentinel errors.
var (
	ErrBotTokenInvalid    = errors.New("bot token is invalid")
	ErrWebhookSetupFailed = errors.New("failed to set up webhook")
	ErrSendMessageFailed  = errors.New("failed to send message")
	ErrGetUpdatesFailed   = errors.New("failed to get updates")
	ErrAPICallFailed      = errors.New("telegram API call failed")
)

// Config holds Telegram bot configuration.
type Config struct {
	BotToken      string `conf:"bot_token"`
	BotUsername   string `conf:"bot_username"   default:"ayabot"`
	WebhookURL    string `conf:"webhook_url"`
	WebhookSecret string `conf:"webhook_secret"`
	Enabled       bool   `conf:"enabled"        default:"false"`
	UsePolling    bool   `conf:"use_polling"    default:"false"`
}

// IsConfigured returns true if the bot token is present.
func (c *Config) IsConfigured() bool {
	return c.BotToken != ""
}

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// -- Telegram Bot API types --

// Update represents a Telegram update.
type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

// Message represents a Telegram message.
type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
	Date      int64  `json:"date"`
}

// User represents a Telegram user.
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"` // "private", "group", "supergroup", "channel"
}

// BotInfo represents bot identity information.
type BotInfo struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// ChatInviteLink represents a Telegram chat invite link.
type ChatInviteLink struct {
	InviteLink  string `json:"invite_link"`
	Name        string `json:"name"`
	ExpireDate  int64  `json:"expire_date"`
	MemberLimit int    `json:"member_limit"`
}

// apiResponse wraps the standard Telegram Bot API response envelope.
type apiResponse struct {
	OK          bool            `json:"ok"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
}

// Client provides Telegram Bot API operations.
type Client struct {
	config     *Config
	logger     *logfx.Logger
	httpClient HTTPClient
	baseURL    string
}

// NewClient creates a new Telegram Bot API client.
func NewClient(config *Config, logger *logfx.Logger, httpClient HTTPClient) *Client {
	return &Client{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
		baseURL:    "https://api.telegram.org/bot" + config.BotToken,
	}
}

// Config returns the Telegram configuration.
func (c *Client) Config() *Config {
	return c.config
}

// GetMe verifies the bot token and returns bot identity.
func (c *Client) GetMe(ctx context.Context) (*BotInfo, error) {
	result, err := c.callAPI(ctx, "getMe", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBotTokenInvalid, err)
	}

	var info BotInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBotTokenInvalid, err)
	}

	return &info, nil
}

// SetWebhook configures the webhook on Telegram's side.
func (c *Client) SetWebhook(ctx context.Context, url string, secret string) error {
	payload := map[string]any{
		"url":             url,
		"secret_token":    secret,
		"allowed_updates": []string{"message"},
	}

	_, err := c.callAPI(ctx, "setWebhook", payload)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWebhookSetupFailed, err)
	}

	c.logger.InfoContext(ctx, "Telegram webhook configured",
		slog.String("url", url))

	return nil
}

// DeleteWebhook removes the webhook (needed for polling mode).
func (c *Client) DeleteWebhook(ctx context.Context) error {
	_, err := c.callAPI(ctx, "deleteWebhook", nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWebhookSetupFailed, err)
	}

	c.logger.InfoContext(ctx, "Telegram webhook removed")

	return nil
}

// GetUpdates fetches updates using long polling.
func (c *Client) GetUpdates(ctx context.Context, offset int64, timeout int) ([]Update, error) {
	payload := map[string]any{
		"offset":          offset,
		"timeout":         timeout,
		"allowed_updates": []string{"message"},
	}

	result, err := c.callAPI(ctx, "getUpdates", payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetUpdatesFailed, err)
	}

	var updates []Update
	if err := json.Unmarshal(result, &updates); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetUpdatesFailed, err)
	}

	return updates, nil
}

// SendMessage sends a text message to a chat.
func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	payload := map[string]any{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	_, err := c.callAPI(ctx, "sendMessage", payload)
	if err != nil {
		// Log but don't crash on 403 (user blocked bot)
		c.logger.WarnContext(ctx, "Failed to send Telegram message",
			slog.Int64("chat_id", chatID),
			slog.String("error", err.Error()))

		return fmt.Errorf("%w: %w", ErrSendMessageFailed, err)
	}

	return nil
}

// CreateChatInviteLink creates a single-use invite link for a chat.
// This is a stub for future channel invitation functionality.
func (c *Client) CreateChatInviteLink(
	ctx context.Context,
	chatID int64,
	name string,
	memberLimit int,
) (*ChatInviteLink, error) {
	payload := map[string]any{
		"chat_id":      chatID,
		"name":         name,
		"member_limit": memberLimit,
	}

	result, err := c.callAPI(ctx, "createChatInviteLink", payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAPICallFailed, err)
	}

	var link ChatInviteLink
	if err := json.Unmarshal(result, &link); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAPICallFailed, err)
	}

	return &link, nil
}

// DeepLink returns the bot deep link URL with a start parameter.
func (c *Client) DeepLink(token string) string {
	return "https://t.me/" + c.config.BotUsername + "?start=" + token
}

// FormatUserID converts a Telegram user ID to a string for storage.
func FormatUserID(id int64) string {
	return strconv.FormatInt(id, 10)
}

// callAPI makes a POST request to the Telegram Bot API.
func (c *Client) callAPI(ctx context.Context, method string, payload any) (json.RawMessage, error) {
	var body io.Reader

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrAPICallFailed, err)
		}

		body = bytes.NewReader(data)
	}

	reqURL := c.baseURL + "/" + method

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAPICallFailed, err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Telegram API request failed",
			slog.String("method", method),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrAPICallFailed, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, _ := io.ReadAll(resp.Body)

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("%w: failed to parse response: %w", ErrAPICallFailed, err)
	}

	if !apiResp.OK {
		c.logger.ErrorContext(ctx, "Telegram API returned error",
			slog.String("method", method),
			slog.String("description", apiResp.Description),
			slog.Int("status", resp.StatusCode))

		return nil, fmt.Errorf("%w: %s", ErrAPICallFailed, apiResp.Description)
	}

	return apiResp.Result, nil
}
