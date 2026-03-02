package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
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

	p.logger.WarnContext(ctx, "Successfully fetched YouTube videos",
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

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRefreshToken, err)
	}

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

	unmarshalErr := json.Unmarshal(body, &tokenResp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToRefreshToken, unmarshalErr)
	}

	if tokenResp.AccessToken == "" {
		return nil, ErrFailedToRefreshToken
	}

	result := buildTokenRefreshResult(&tokenResp)

	p.logger.DebugContext(ctx, "Successfully refreshed YouTube access token")

	return result, nil
}

// buildTokenRefreshResult constructs a TokenRefreshResult from a token response.
func buildTokenRefreshResult(tokenResp *tokenResponse) *linksync.TokenRefreshResult {
	var expiresAt *time.Time

	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	var newRefreshToken *string
	if tokenResp.RefreshToken != "" {
		newRefreshToken = &tokenResp.RefreshToken
	}

	return &linksync.TokenRefreshResult{
		AccessToken:          tokenResp.AccessToken,
		AccessTokenExpiresAt: expiresAt,
		RefreshToken:         newRefreshToken,
	}
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}

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

	err = json.Unmarshal(body, &channelResp)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}

	if len(channelResp.Items) == 0 {
		return "", fmt.Errorf("%w: no channel found", ErrFailedToFetchVideos)
	}

	return channelResp.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

// fetchPlaylistVideos fetches videos from a playlist.
func (p *Provider) fetchPlaylistVideos( //nolint:cyclop,funlen // paginated API fetch loop
	ctx context.Context,
	accessToken string,
	playlistID string,
	publishedAfter *time.Time,
	maxResults int,
) ([]*linksync.RemoteStoryItem, error) {
	videos := make([]*linksync.RemoteStoryItem, 0, maxResults)
	pageToken := ""

	const youtubeAPIMaxResults = 50

	for len(videos) < maxResults {
		// Build request URL
		reqURL := fmt.Sprintf(
			"https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,contentDetails,status&playlistId=%s&maxResults=%d",
			url.QueryEscape(playlistID),
			min(youtubeAPIMaxResults, maxResults-len(videos)),
		)

		if pageToken != "" {
			reqURL += "&pageToken=" + url.QueryEscape(pageToken)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		body, readErr := func() ([]byte, error) {
			defer resp.Body.Close() //nolint:errcheck

			return io.ReadAll(resp.Body)
		}()
		if readErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, readErr)
		}

		if resp.StatusCode != http.StatusOK {
			p.logger.ErrorContext(ctx, "Failed to fetch playlist items",
				slog.Int("status", resp.StatusCode),
				slog.String("response", string(body)))

			return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchVideos, resp.StatusCode)
		}

		var playlistResp playlistItemsResponse

		err = json.Unmarshal(body, &playlistResp)
		if err != nil {
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

	// Fetch and attach video metadata for the videos we have
	if len(videos) > 0 {
		p.enrichWithVideoMetadata(ctx, accessToken, videos)
		p.logDiagnosticKeys(ctx, videos)
	}

	return videos, nil
}

// enrichWithVideoMetadata fetches full video metadata and attaches it to each video's properties.
func (p *Provider) enrichWithVideoMetadata(
	ctx context.Context,
	accessToken string,
	videos []*linksync.RemoteStoryItem,
) {
	videoIDs := make([]string, len(videos))
	for index, video := range videos {
		videoIDs[index] = video.RemoteID
	}

	videoMetadata, err := p.fetchVideoMetadata(ctx, accessToken, videoIDs)
	if err != nil {
		p.logger.WarnContext(ctx, "Failed to fetch video metadata",
			slog.String("error", err.Error()))

		return
	}

	for _, video := range videos {
		if meta, ok := videoMetadata[video.RemoteID]; ok {
			video.Properties["videoMetadata"] = meta
		}
	}
}

// logDiagnosticKeys logs property keys of the first video for debugging.
func (p *Provider) logDiagnosticKeys(
	ctx context.Context,
	videos []*linksync.RemoteStoryItem,
) {
	first := videos[0]

	var vmKeys, pmKeys []string

	if videoMeta, ok := first.Properties["videoMetadata"].(map[string]any); ok {
		for key := range videoMeta {
			vmKeys = append(vmKeys, key)
		}
	}

	if playlistMeta, ok := first.Properties["playlistItemMetadata"].(map[string]any); ok {
		for key := range playlistMeta {
			pmKeys = append(pmKeys, key)
		}
	}

	p.logger.WarnContext(ctx, "YouTube API response diagnostic",
		slog.String("remote_id", first.RemoteID),
		slog.Any("videoMetadata_keys", vmKeys),
		slog.Any("playlistItemMetadata_keys", pmKeys))
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
		return fmt.Errorf("failed to unmarshal playlist items: %w", err)
	}

	r.NextPageToken = alias.NextPageToken
	r.RawItems = alias.RawItems
	r.Items = make([]playlistItem, len(alias.RawItems))

	for index, raw := range alias.RawItems {
		itemErr := json.Unmarshal(raw, &r.Items[index])
		if itemErr != nil {
			return fmt.Errorf("failed to unmarshal playlist item: %w", itemErr)
		}
	}

	return nil
}

// youtubeAPIMaxBatchSize is the maximum number of video IDs per YouTube API request.
const youtubeAPIMaxBatchSize = 50

// fetchVideoMetadata fetches full metadata for multiple videos.
func (p *Provider) fetchVideoMetadata(
	ctx context.Context,
	accessToken string,
	videoIDs []string,
) (map[string]map[string]any, error) {
	result := make(map[string]map[string]any)

	// YouTube API allows up to 50 video IDs per request
	for i := 0; i < len(videoIDs); i += youtubeAPIMaxBatchSize {
		end := min(i+youtubeAPIMaxBatchSize, len(videoIDs))

		batch := videoIDs[i:end]
		ids := strings.Join(batch, ",")

		videosURL := "https://www.googleapis.com/youtube/v3/videos" +
			"?part=snippet,contentDetails,statistics,status&id=" +
			url.QueryEscape(ids)

		batchResult, batchErr := p.fetchVideoBatch(ctx, accessToken, videosURL)
		if batchErr != nil {
			return nil, batchErr
		}

		maps.Copy(result, batchResult)
	}

	return result, nil
}

// fetchVideoBatch fetches a single batch of video metadata from the given URL.
func (p *Provider) fetchVideoBatch(
	ctx context.Context,
	accessToken string,
	videosURL string,
) (map[string]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, videosURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}

	body, readErr := func() ([]byte, error) {
		defer resp.Body.Close() //nolint:errcheck

		return io.ReadAll(resp.Body)
	}()
	if readErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, readErr)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFailedToFetchVideos, resp.StatusCode)
	}

	var videosResp struct {
		Items []json.RawMessage `json:"items"`
	}

	err = json.Unmarshal(body, &videosResp)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
	}

	result := make(map[string]map[string]any, len(videosResp.Items))

	for _, rawItem := range videosResp.Items {
		var item map[string]any

		err = json.Unmarshal(rawItem, &item)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToFetchVideos, err)
		}

		id, _ := item["id"].(string)
		if id != "" {
			result[id] = item
		}
	}

	return result, nil
}
