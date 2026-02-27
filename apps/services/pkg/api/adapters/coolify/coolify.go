package coolify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

var (
	ErrBaseDomainMissing = errors.New("base domain missing from update set")
	ErrNotConfigured     = errors.New("coolify adapter not configured")
	ErrAPIRequestFailed  = errors.New("coolify API request failed")
)

// Config holds configuration for the Coolify HTTP API adapter.
type Config struct {
	APIURL           string        `conf:"api_url"             default:"https://cool.acikyazilim.com/api/v1"`
	APIToken         string        `conf:"api_token"`
	WebclientAppUUID string        `conf:"webclient_app_uuid"  default:"xkgccwc8sggsw0s44g04oo0k"`
	RequestTimeout   time.Duration `conf:"request_timeout"     default:"30s"`
	BaseDomains      string        `conf:"base_domains"`
}

// Client implements profiles.WebserverSyncer for Coolify infrastructure.
type Client struct {
	config      *Config
	logger      *logfx.Logger
	httpClient  *http.Client
	baseDomains []string
}

// NewClient creates a new Coolify API client.
func NewClient(config *Config, logger *logfx.Logger) *Client {
	baseDomains := make([]string, 0)

	if config.BaseDomains != "" {
		for _, d := range strings.Split(config.BaseDomains, ",") {
			trimmed := strings.TrimSpace(d)
			if trimmed != "" {
				baseDomains = append(baseDomains, trimmed)
			}
		}
	}

	return &Client{
		config: config,
		logger: logger,
		httpClient: &http.Client{
			Timeout: config.RequestTimeout,
		},
		baseDomains: baseDomains,
	}
}

// BaseDomains returns the configured base domains that must always be present.
func (c *Client) BaseDomains() []string {
	return c.baseDomains
}

// applicationResponse is the partial response from GET /applications/{uuid}.
type applicationResponse struct {
	FQDN string `json:"fqdn"`
}

// GetCurrentDomains retrieves the current domain list from the Coolify API.
func (c *Client) GetCurrentDomains(ctx context.Context) ([]string, error) {
	if c.config.APIToken == "" {
		return nil, ErrNotConfigured
	}

	url := fmt.Sprintf("%s/applications/%s", c.config.APIURL, c.config.WebclientAppUUID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAPIRequestFailed, err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAPIRequestFailed, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("%w: status %d, body: %s", ErrAPIRequestFailed, resp.StatusCode, string(body))
	}

	var appResp applicationResponse

	if err := json.NewDecoder(resp.Body).Decode(&appResp); err != nil {
		return nil, fmt.Errorf("%w: failed to decode response: %w", ErrAPIRequestFailed, err)
	}

	if appResp.FQDN == "" {
		return []string{}, nil
	}

	domains := make([]string, 0)
	for _, d := range strings.Split(appResp.FQDN, ",") {
		trimmed := strings.TrimSpace(d)
		if trimmed != "" {
			domains = append(domains, trimmed)
		}
	}

	return domains, nil
}

// updateRequest is the request body for PATCH /applications/{uuid}.
type updateRequest struct {
	Domains string `json:"domains"`
}

// UpdateDomains updates the domain list on the Coolify webclient app.
// It refuses to execute if any base domain is missing from the provided list.
func (c *Client) UpdateDomains(ctx context.Context, domains []string) error {
	if c.config.APIToken == "" {
		return ErrNotConfigured
	}

	// Safety: verify all base domains are present
	domainSet := make(map[string]bool, len(domains))
	for _, d := range domains {
		domainSet[d] = true
	}

	for _, base := range c.baseDomains {
		if !domainSet[base] {
			return fmt.Errorf("%w: %s", ErrBaseDomainMissing, base)
		}
	}

	domainStr := strings.Join(domains, ",")

	c.logger.WarnContext(ctx, "Updating Coolify domains",
		slog.String("app_uuid", c.config.WebclientAppUUID),
		slog.Int("domain_count", len(domains)),
		slog.String("domains", domainStr))

	body, err := json.Marshal(updateRequest{Domains: domainStr})
	if err != nil {
		return fmt.Errorf("%w: failed to marshal request: %w", ErrAPIRequestFailed, err)
	}

	url := fmt.Sprintf("%s/applications/%s", c.config.APIURL, c.config.WebclientAppUUID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAPIRequestFailed, err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAPIRequestFailed, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("%w: status %d, body: %s", ErrAPIRequestFailed, resp.StatusCode, string(respBody))
	}

	c.logger.WarnContext(ctx, "Successfully updated Coolify domains",
		slog.Int("domain_count", len(domains)))

	return nil
}
