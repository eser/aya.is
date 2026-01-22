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
		for _, item := range playlistResp.Items {
			publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

			// Filter by publishedAfter if provided
			if publishedAfter != nil && !publishedAt.After(*publishedAfter) {
				// Items are returned in reverse chronological order
				// So if we hit an item older than publishedAfter, we can stop
				return videos, nil
			}

			video := &linksync.RemoteStoryItem{
				RemoteID:     item.ContentDetails.VideoID,
				Title:        item.Snippet.Title,
				Description:  item.Snippet.Description,
				PublishedAt:  publishedAt,
				ThumbnailURL: p.getBestThumbnail(item.Snippet.Thumbnails),
				Properties: map[string]any{
					"channel_id":    item.Snippet.ChannelID,
					"channel_title": item.Snippet.ChannelTitle,
					"playlist_id":   playlistID,
					"position":      item.Snippet.Position,
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

	// Fetch video statistics for the videos we have
	if len(videos) > 0 {
		videoIDs := make([]string, len(videos))
		for i, v := range videos {
			videoIDs[i] = v.RemoteID
		}

		stats, err := p.fetchVideoStatistics(ctx, accessToken, videoIDs)
		if err != nil {
			// Log but don't fail - statistics are optional
			p.logger.WarnContext(ctx, "Failed to fetch video statistics",
				slog.String("error", err.Error()))
		} else {
			for _, video := range videos {
				if stat, ok := stats[video.RemoteID]; ok {
					video.ViewCount = stat.ViewCount
					video.LikeCount = stat.LikeCount
					video.Duration = stat.Duration
				}
			}
		}
	}

	return videos, nil
}

// playlistItemsResponse represents the YouTube playlist items API response.
type playlistItemsResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Snippet struct {
			PublishedAt  string `json:"publishedAt"`
			ChannelID    string `json:"channelId"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
			Position     int    `json:"position"`
			Thumbnails   struct {
				Default  *thumbnail `json:"default"`
				Medium   *thumbnail `json:"medium"`
				High     *thumbnail `json:"high"`
				Standard *thumbnail `json:"standard"`
				Maxres   *thumbnail `json:"maxres"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			VideoID string `json:"videoId"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// getBestThumbnail returns the highest quality thumbnail URL available.
func (p *Provider) getBestThumbnail(thumbnails struct {
	Default  *thumbnail `json:"default"`
	Medium   *thumbnail `json:"medium"`
	High     *thumbnail `json:"high"`
	Standard *thumbnail `json:"standard"`
	Maxres   *thumbnail `json:"maxres"`
},
) string {
	if thumbnails.Maxres != nil {
		return thumbnails.Maxres.URL
	}

	if thumbnails.Standard != nil {
		return thumbnails.Standard.URL
	}

	if thumbnails.High != nil {
		return thumbnails.High.URL
	}

	if thumbnails.Medium != nil {
		return thumbnails.Medium.URL
	}

	if thumbnails.Default != nil {
		return thumbnails.Default.URL
	}

	return ""
}

// videoStatistics holds view and like counts for a video.
type videoStatistics struct {
	ViewCount int64
	LikeCount int64
	Duration  string
}

// fetchVideoStatistics fetches statistics for multiple videos.
func (p *Provider) fetchVideoStatistics(
	ctx context.Context,
	accessToken string,
	videoIDs []string,
) (map[string]*videoStatistics, error) {
	result := make(map[string]*videoStatistics)

	// YouTube API allows up to 50 video IDs per request
	for i := 0; i < len(videoIDs); i += 50 {
		end := i + 50
		if end > len(videoIDs) {
			end = len(videoIDs)
		}

		batch := videoIDs[i:end]
		ids := strings.Join(batch, ",")

		reqURL := "https://www.googleapis.com/youtube/v3/videos?part=statistics,contentDetails&id=" + url.QueryEscape(
			ids,
		)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status %d", resp.StatusCode)
		}

		var videosResp struct {
			Items []struct {
				ID         string `json:"id"`
				Statistics struct {
					ViewCount string `json:"viewCount"`
					LikeCount string `json:"likeCount"`
				} `json:"statistics"`
				ContentDetails struct {
					Duration string `json:"duration"`
				} `json:"contentDetails"`
			} `json:"items"`
		}

		if err := json.Unmarshal(body, &videosResp); err != nil {
			return nil, err
		}

		for _, item := range videosResp.Items {
			var viewCount, likeCount int64

			fmt.Sscanf(item.Statistics.ViewCount, "%d", &viewCount)
			fmt.Sscanf(item.Statistics.LikeCount, "%d", &likeCount)

			result[item.ID] = &videoStatistics{
				ViewCount: viewCount,
				LikeCount: likeCount,
				Duration:  item.ContentDetails.Duration,
			}
		}
	}

	return result, nil
}
