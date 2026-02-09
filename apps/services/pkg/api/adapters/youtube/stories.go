package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/linksync"
)

// Sentinel errors for stories operations.
var (
	ErrFailedToFetchVideos  = linksync.ErrFailedToFetchStories
	ErrFailedToRefreshToken = linksync.ErrFailedToRefreshToken
)

// FetchRemoteStories fetches YouTube videos for a channel.
// Implements linksync.RemoteStoryFetcher interface.
func (p *Provider) FetchRemoteStories(
	ctx context.Context,
	accessToken string,
	remoteSourceID string,
	publishedAfter *time.Time,
	maxResults int,
) ([]*linksync.RemoteStoryItem, error) {
	p.logger.DebugContext(ctx, "Fetching YouTube videos",
		slog.String("channel_id", remoteSourceID),
		slog.Int("max_results", maxResults))

	// First, get the uploads playlist ID for the channel
	uploadsPlaylistID, err := p.getUploadsPlaylistID(ctx, accessToken, remoteSourceID)
	if err != nil {
		return nil, err
	}

	// Fetch videos from the uploads playlist
	videos, err := p.fetchPlaylistVideos(
		ctx,
		accessToken,
		uploadsPlaylistID,
		publishedAfter,
		maxResults,
	)
	if err != nil {
		return nil, err
	}

	p.logger.DebugContext(ctx, "Successfully fetched YouTube videos",
		slog.String("channel_id", remoteSourceID),
		slog.Int("video_count", len(videos)))

	return videos, nil
}

// RefreshAccessToken refreshes an expired access token using the refresh token.
// Implements linksync.RemoteStoryFetcher interface.
func (p *Provider) RefreshAccessToken(
	ctx context.Context,
	refreshToken string,
) (*linksync.TokenRefreshResult, error) {
	p.logger.DebugContext(ctx, "Refreshing YouTube access token")

	values := url.Values{
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(values.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to refresh access token",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToRefreshToken, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "Token refresh failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToRefreshToken, resp.StatusCode)
	}

	var tokenResp tokenResponse

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRefreshToken, err)
	}

	if tokenResp.AccessToken == "" {
		return nil, ErrFailedToRefreshToken
	}

	// Calculate token expiry
	var expiresAt *time.Time

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	// Google doesn't return a new refresh token on refresh
	// Only return new refresh token if one was provided
	var newRefreshToken *string
	if tokenResp.RefreshToken != "" {
		newRefreshToken = &tokenResp.RefreshToken
	}

	p.logger.DebugContext(ctx, "Successfully refreshed YouTube access token")

	return &linksync.TokenRefreshResult{
		AccessToken:          tokenResp.AccessToken,
		AccessTokenExpiresAt: expiresAt,
		RefreshToken:         newRefreshToken,
	}, nil
}

// getUploadsPlaylistID fetches the uploads playlist ID for a channel.
func (p *Provider) getUploadsPlaylistID(
	ctx context.Context,
	accessToken string,
	channelID string,
) (string, error) {
	reqURL := "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&id=" + url.QueryEscape(
		channelID,
	)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		p.logger.ErrorContext(ctx, "Failed to get channel details",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))

		return "", fmt.Errorf("%w: status %d", ErrFailedToFetchVideos, resp.StatusCode)
	}

	var channelResp struct {
		Items []struct {
			ContentDetails struct {
				RelatedPlaylists struct {
					Uploads string `json:"uploads"`
				} `json:"relatedPlaylists"`
			} `json:"contentDetails"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &channelResp); err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}

	if len(channelResp.Items) == 0 {
		return "", fmt.Errorf("%w: no channel found", ErrFailedToFetchVideos)
	}

	return channelResp.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

// fetchPlaylistVideos fetches videos from a playlist.
func (p *Provider) fetchPlaylistVideos(
	ctx context.Context,
	accessToken string,
	playlistID string,
	publishedAfter *time.Time,
	maxResults int,
) ([]*linksync.RemoteStoryItem, error) {
	videos := make([]*linksync.RemoteStoryItem, 0, maxResults)
	pageToken := ""

	for len(videos) < maxResults {
		// Build request URL
		reqURL := fmt.Sprintf(
			"https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,contentDetails&playlistId=%s&maxResults=%d",
			url.QueryEscape(playlistID),
			min(50, maxResults-len(videos)), // YouTube API max is 50
		)

		if pageToken != "" {
			reqURL += "&pageToken=" + url.QueryEscape(pageToken)
		}

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusOK {
			p.logger.ErrorContext(ctx, "Failed to fetch playlist items",
				slog.Int("status", resp.StatusCode),
				slog.String("response", string(body)))

			return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchVideos, resp.StatusCode)
		}

		var playlistResp playlistItemsResponse

		if err := json.Unmarshal(body, &playlistResp); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		// Process items
		for i, item := range playlistResp.Items {
			publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

			// Filter by publishedAfter if provided
			if publishedAfter != nil && !publishedAt.After(*publishedAfter) {
				// Items are returned in reverse chronological order
				// So if we hit an item older than publishedAfter, we can stop
				return videos, nil
			}

			// Store the raw playlist item metadata
			var rawItem map[string]any

			err := json.Unmarshal(playlistResp.RawItems[i], &rawItem)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
			}

			video := &linksync.RemoteStoryItem{
				RemoteID:    item.ContentDetails.VideoID,
				PublishedAt: publishedAt,
				Properties: map[string]any{
					"playlistItemMetadata": rawItem,
				},
			}

			videos = append(videos, video)

			if len(videos) >= maxResults {
				break
			}
		}

		// Check for more pages
		if playlistResp.NextPageToken == "" {
			break
		}

		pageToken = playlistResp.NextPageToken
	}

	// Fetch video metadata for the videos we have
	if len(videos) > 0 {
		videoIDs := make([]string, len(videos))
		for i, v := range videos {
			videoIDs[i] = v.RemoteID
		}

		videoMetadata, err := p.fetchVideoMetadata(ctx, accessToken, videoIDs)
		if err != nil {
			// Log but don't fail - video metadata is optional
			p.logger.WarnContext(ctx, "Failed to fetch video metadata",
				slog.String("error", err.Error()))
		} else {
			for _, video := range videos {
				if meta, ok := videoMetadata[video.RemoteID]; ok {
					video.Properties["videoMetadata"] = meta
				}
			}
		}
	}

	return videos, nil
}

// playlistItemsResponse represents the YouTube playlist items API response.
type playlistItemsResponse struct {
	NextPageToken string            `json:"nextPageToken"`
	RawItems      []json.RawMessage `json:"items"`
	Items         []playlistItem
}

// playlistItem holds the minimal parsed fields needed for sync logic.
type playlistItem struct {
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
	} `json:"snippet"`
	ContentDetails struct {
		VideoID string `json:"videoId"`
	} `json:"contentDetails"`
}

// UnmarshalJSON custom unmarshals to preserve raw items while parsing minimal fields.
func (r *playlistItemsResponse) UnmarshalJSON(data []byte) error {
	type Alias struct {
		NextPageToken string            `json:"nextPageToken"`
		RawItems      []json.RawMessage `json:"items"`
	}

	var alias Alias

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	r.NextPageToken = alias.NextPageToken
	r.RawItems = alias.RawItems
	r.Items = make([]playlistItem, len(alias.RawItems))

	for i, raw := range alias.RawItems {
		err := json.Unmarshal(raw, &r.Items[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// fetchVideoMetadata fetches full metadata for multiple videos.
func (p *Provider) fetchVideoMetadata(
	ctx context.Context,
	accessToken string,
	videoIDs []string,
) (map[string]map[string]any, error) {
	result := make(map[string]map[string]any)

	// YouTube API allows up to 50 video IDs per request
	for i := 0; i < len(videoIDs); i += 50 {
		end := i + 50
		if end > len(videoIDs) {
			end = len(videoIDs)
		}

		batch := videoIDs[i:end]
		ids := strings.Join(batch, ",")

		reqURL := "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics,status&id=" + url.QueryEscape(
			ids,
		)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchVideos, resp.StatusCode)
		}

		var videosResp struct {
			Items []json.RawMessage `json:"items"`
		}

		if err := json.Unmarshal(body, &videosResp); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		for _, rawItem := range videosResp.Items {
			var item map[string]any

			err := json.Unmarshal(rawItem, &item)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
			}

			id, _ := item["id"].(string)
			if id != "" {
				result[id] = item
			}
		}
	}

	return result, nil
}
