package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Sentinel errors for live broadcast operations.
var ErrFailedToCheckLiveBroadcasts = errors.New("failed to check live broadcasts")

// LiveBroadcastResult contains information about an active live broadcast.
type LiveBroadcastResult struct {
	StartedAt    *time.Time
	BroadcastID  string
	Title        string
	ThumbnailURL string
	IsLive       bool
}

// CheckLiveBroadcasts checks if the authenticated user has any active live broadcasts.
// Uses the liveBroadcasts.list endpoint with broadcastStatus=active (1 quota unit).
func (p *Provider) CheckLiveBroadcasts( //nolint:funlen
	ctx context.Context,
	accessToken string,
) (*LiveBroadcastResult, error) {
	reqURL := "https://www.googleapis.com/youtube/v3/liveBroadcasts" +
		"?part=snippet,status&broadcastStatus=active&mine=true"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCheckLiveBroadcasts, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to check live broadcasts",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToCheckLiveBroadcasts, err)
	}

	body, readErr := func() ([]byte, error) {
		defer resp.Body.Close() //nolint:errcheck

		return io.ReadAll(resp.Body)
	}()
	if readErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCheckLiveBroadcasts, readErr)
	}

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "Live broadcasts API returned error",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToCheckLiveBroadcasts, resp.StatusCode)
	}

	var broadcastResp struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title           string `json:"title"`
				ActualStartTime string `json:"actualStartTime"`
				LiveChatID      string `json:"liveChatId"`
				Thumbnails      struct {
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			Status struct {
				LifeCycleStatus string `json:"lifeCycleStatus"`
			} `json:"status"`
		} `json:"items"`
	}

	err = json.Unmarshal(body, &broadcastResp)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCheckLiveBroadcasts, err)
	}

	if len(broadcastResp.Items) == 0 {
		return &LiveBroadcastResult{IsLive: false}, nil //nolint:exhaustruct
	}

	// Take the first active broadcast
	item := broadcastResp.Items[0]

	var startedAt *time.Time

	if item.Snippet.ActualStartTime != "" {
		parsed, parseErr := time.Parse(time.RFC3339, item.Snippet.ActualStartTime)
		if parseErr == nil {
			startedAt = &parsed
		}
	}

	return &LiveBroadcastResult{
		IsLive:       true,
		BroadcastID:  item.ID,
		Title:        item.Snippet.Title,
		StartedAt:    startedAt,
		ThumbnailURL: item.Snippet.Thumbnails.Default.URL,
	}, nil
}
