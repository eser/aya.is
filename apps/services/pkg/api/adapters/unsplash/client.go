package unsplash

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

// Sentinel errors.
var (
	ErrFailedToSearchPhotos = errors.New("failed to search Unsplash photos")
	ErrNoAccessKey          = errors.New("unsplash access key is not configured")
)

// HTTPClient interface for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config holds Unsplash API configuration.
type Config struct {
	AccessKey string `conf:"access_key"`
}

// IsConfigured returns true if the Unsplash API is configured.
func (c *Config) IsConfigured() bool {
	return c.AccessKey != ""
}

// PhotoURLs contains various image size URLs from Unsplash.
type PhotoURLs struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

// PhotoUser contains photographer information.
type PhotoUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// Photo represents an Unsplash photo.
type Photo struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	AltDesc     string    `json:"alt_description"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Color       string    `json:"color"`
	URLs        PhotoURLs `json:"urls"`
	User        PhotoUser `json:"user"`
}

// SearchResult represents the search API response.
type SearchResult struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}

// Client provides Unsplash API operations.
type Client struct {
	config     *Config
	logger     *logfx.Logger
	httpClient HTTPClient
}

// NewClient creates a new Unsplash client.
func NewClient(
	config *Config,
	logger *logfx.Logger,
	httpClient HTTPClient,
) *Client {
	return &Client{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

// Config returns the Unsplash configuration.
func (c *Client) Config() *Config {
	return c.config
}

// SearchPhotos searches for photos on Unsplash.
func (c *Client) SearchPhotos(
	ctx context.Context,
	query string,
	page int,
	perPage int,
) (*SearchResult, error) {
	if c.config.AccessKey == "" {
		return nil, ErrNoAccessKey
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("query", query)
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("per_page", strconv.Itoa(perPage))
	queryParams.Set("orientation", "landscape") // Better for cover images

	apiURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "api.unsplash.com",
		Path:     "/search/photos",
		RawQuery: queryParams.Encode(),
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiURL.String(),
		nil,
	)
	req.Header.Set("Authorization", "Client-ID "+c.config.AccessKey)
	req.Header.Set("Accept-Version", "v1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to search Unsplash photos",
			slog.String("error", err.Error()),
			slog.String("query", query))

		return nil, fmt.Errorf("%w: %w", ErrFailedToSearchPhotos, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "Unsplash search failed",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)),
			slog.String("query", query))

		return nil, fmt.Errorf("%w: status %d", ErrFailedToSearchPhotos, resp.StatusCode)
	}

	var result SearchResult

	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.ErrorContext(ctx, "Failed to parse Unsplash response",
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToSearchPhotos, err)
	}

	c.logger.DebugContext(ctx, "Unsplash search successful",
		slog.String("query", query),
		slog.Int("total", result.Total),
		slog.Int("results_count", len(result.Results)))

	return &result, nil
}
