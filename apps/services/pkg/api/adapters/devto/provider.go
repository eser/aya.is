package devto

import (
  "context"
  "fmt"
  "log/slog"
  "strings"

  "github.com/eser/aya.is/services/pkg/ajan/httpclient"
  "github.com/eser/aya.is/services/pkg/ajan/logfx"
  "github.com/eser/aya.is/services/pkg/api/business/siteimporter"
)

const devtoBaseURL = "https://dev.to"

// Provider implements siteimporter.SiteProvider for Dev.to.
type Provider struct {
  logger     *logfx.Logger
  httpClient *httpclient.Client
}

// NewProvider creates a new Dev.to provider.
func NewProvider(logger *logfx.Logger, httpClient *httpclient.Client) *Provider {
  return &Provider{
    logger:     logger,
    httpClient: httpClient,
  }
}

// Kind returns the provider kind.
func (p *Provider) Kind() string {
  return "devto"
}

// Check validates a Dev.to URL and returns connection info.
func (p *Provider) Check(ctx context.Context, rawURL string) (*siteimporter.CheckResult, error) {
  username := extractUsername(rawURL)
  if username == "" {
    return nil, fmt.Errorf(
      "%w: could not extract username from URL",
      siteimporter.ErrInvalidURL,
    )
  }

  articles, err := fetchArticlesPage(ctx, p.httpClient, username, 1, 1)
  if err != nil {
    return nil, fmt.Errorf("%w: %w", siteimporter.ErrSiteNotFound, err)
  }

  title := username
  if len(articles) > 0 && articles[0].User.Name != "" {
    title = articles[0].User.Name
  }

  p.logger.DebugContext(ctx, "Validated Dev.to profile",
    slog.String("username", username))

  return &siteimporter.CheckResult{
    Valid:    true,
    Username: username,
    URI:      fmt.Sprintf("%s/%s", devtoBaseURL, username),
    Title:    title,
  }, nil
}

// extractUsername parses a Dev.to URL to extract the username.
// Supports: "https://dev.to/eser", "dev.to/eser", "@eser".
func extractUsername(rawURL string) string {
  rawURL = strings.TrimSpace(rawURL)
  if rawURL == "" {
    return ""
  }

  rawURL = strings.TrimPrefix(rawURL, "https://")
  rawURL = strings.TrimPrefix(rawURL, "http://")

  rawURL = strings.TrimPrefix(rawURL, "dev.to/")
  rawURL = strings.TrimPrefix(rawURL, "www.dev.to/")

  if idx := strings.IndexAny(rawURL, "?#"); idx != -1 {
    rawURL = rawURL[:idx]
  }

  rawURL = strings.TrimRight(rawURL, "/")

  username, _, _ := strings.Cut(rawURL, "/")
  username = strings.TrimPrefix(username, "@")
  if username == "" {
    return ""
  }

  for _, char := range username {
    if !isValidUsernameChar(char) {
      return ""
    }
  }

  return username
}

func isValidUsernameChar(char rune) bool {
  if char >= 'a' && char <= 'z' {
    return true
  }

  if char >= 'A' && char <= 'Z' {
    return true
  }

  if char >= '0' && char <= '9' {
    return true
  }

  switch char {
  case '-', '_', '.':
    return true
  default:
    return false
  }
}
