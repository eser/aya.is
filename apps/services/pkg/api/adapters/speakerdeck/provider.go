package speakerdeck

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/siteimporter"
	"github.com/mmcdole/gofeed"
)

const speakerDeckBaseURL = "https://speakerdeck.com"

// Provider implements siteimporter.SiteProvider for SpeakerDeck.
type Provider struct {
	httpClient *httpclient.Client
	logger     *logfx.Logger
	parser     *gofeed.Parser
}

// NewProvider creates a new SpeakerDeck provider.
func NewProvider(logger *logfx.Logger, httpClient *httpclient.Client) *Provider {
	parser := gofeed.NewParser()
	parser.Client = httpClient.Client

	return &Provider{
		httpClient: httpClient,
		logger:     logger,
		parser:     parser,
	}
}

// Kind returns the provider kind.
func (p *Provider) Kind() string {
	return "speakerdeck"
}

// Check validates a SpeakerDeck URL and returns connection info.
func (p *Provider) Check(ctx context.Context, rawURL string) (*siteimporter.CheckResult, error) {
	username := extractUsername(rawURL)
	if username == "" {
		return nil, fmt.Errorf(
			"%w: could not extract username from URL",
			siteimporter.ErrInvalidURL,
		)
	}

	// Verify by fetching the RSS feed
	rssURL := fmt.Sprintf("%s/%s.rss", speakerDeckBaseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rssURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch feed", siteimporter.ErrSiteNotFound)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: user %q not found", siteimporter.ErrSiteNotFound, username)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"%w: unexpected status %d",
			siteimporter.ErrSiteNotFound,
			resp.StatusCode,
		)
	}

	// Parse the feed to get the title
	feed, err := p.parser.ParseURL(rssURL)
	if err != nil {
		p.logger.WarnContext(ctx, "Failed to parse SpeakerDeck RSS for title",
			slog.String("username", username),
			slog.Any("error", err))
	}

	title := username
	if feed != nil && feed.Title != "" {
		title = feed.Title
	}

	return &siteimporter.CheckResult{
		Valid:    true,
		Username: username,
		URI:      fmt.Sprintf("%s/%s", speakerDeckBaseURL, username),
		Title:    title,
	}, nil
}

// extractUsername parses a SpeakerDeck URL to extract the username.
// Supports: "https://speakerdeck.com/eser", "speakerdeck.com/eser", "eser".
func extractUsername(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	// Remove protocol
	rawURL = strings.TrimPrefix(rawURL, "https://")
	rawURL = strings.TrimPrefix(rawURL, "http://")

	// Remove speakerdeck.com prefix
	rawURL = strings.TrimPrefix(rawURL, "speakerdeck.com/")
	rawURL = strings.TrimPrefix(rawURL, "www.speakerdeck.com/")

	// Remove trailing slash
	rawURL = strings.TrimRight(rawURL, "/")

	// If there's still a slash, take only the first segment
	username, _, _ := strings.Cut(rawURL, "/")

	// Basic validation
	if username == "" || strings.Contains(username, ".") {
		return ""
	}

	return username
}
