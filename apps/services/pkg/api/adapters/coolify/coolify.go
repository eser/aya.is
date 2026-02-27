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
	ErrNotConfigured    = errors.New("coolify adapter not configured")
	ErrAPIRequestFailed = errors.New("coolify API request failed")
)

// Config holds configuration for the Coolify HTTP API adapter.
type Config struct {
	APIURL           string        `conf:"api_url"            default:"https://cool.acikyazilim.com/api/v1"`
	APIToken         string        `conf:"api_token"`
	WebclientAppUUID string        `conf:"webclient_app_uuid" default:"xkgccwc8sggsw0s44g04oo0k"`
	RequestTimeout   time.Duration `conf:"request_timeout"    default:"30s"`
}

// Client implements profiles.WebserverSyncer for Coolify infrastructure.
type Client struct {
	config     *Config
	logger     *logfx.Logger
	httpClient *http.Client
}

// NewClient creates a new Coolify API client.
func NewClient(config *Config, logger *logfx.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
		httpClient: &http.Client{
			Timeout: config.RequestTimeout,
		},
	}
}

// applicationResponse is the partial response from GET /applications/{uuid}.
type applicationResponse struct {
	FQDN string `json:"fqdn"`
}

// GetCurrentDomains retrieves the current domain list from the Coolify API.
// Returns bare domains (without https:// prefix).
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

		return nil, fmt.Errorf(
			"%w: status %d, body: %s",
			ErrAPIRequestFailed,
			resp.StatusCode,
			string(body),
		)
	}

	var appResp applicationResponse

	if err := json.NewDecoder(resp.Body).Decode(&appResp); err != nil {
		return nil, fmt.Errorf("%w: failed to decode response: %w", ErrAPIRequestFailed, err)
	}

	if appResp.FQDN == "" {
		return []string{}, nil
	}

	// Coolify stores domains as "https://domain1,https://domain2"
	// Strip the protocol prefix to return bare domains
	domains := make([]string, 0)

	for _, d := range strings.Split(appResp.FQDN, ",") {
		trimmed := strings.TrimSpace(d)
		if trimmed != "" {
			trimmed = strings.TrimPrefix(trimmed, "https://")
			trimmed = strings.TrimPrefix(trimmed, "http://")
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
// Accepts bare domains (e.g. "aya.is") and adds https:// prefix for the Coolify API.
func (c *Client) UpdateDomains(ctx context.Context, domains []string) error {
	if c.config.APIToken == "" {
		return ErrNotConfigured
	}

	// Add https:// prefix for Coolify API
	fqdns := make([]string, 0, len(domains))

	for _, d := range domains {
		fqdns = append(fqdns, "https://"+d)
	}

	domainStr := strings.Join(fqdns, ",")

	c.logger.WarnContext(ctx, "Updating Coolify domains",
		slog.String("app_uuid", c.config.WebclientAppUUID),
		slog.Int("domain_count", len(fqdns)),
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

		return fmt.Errorf(
			"%w: status %d, body: %s",
			ErrAPIRequestFailed,
			resp.StatusCode,
			string(respBody),
		)
	}

	c.logger.WarnContext(ctx, "Successfully updated Coolify domains",
		slog.Int("domain_count", len(fqdns)))

	return nil
}

// RestartApplication triggers a container restart on Coolify to apply domain changes.
// This regenerates Traefik labels without a full rebuild/redeploy.
func (c *Client) RestartApplication(ctx context.Context) error {
	if c.config.APIToken == "" {
		return ErrNotConfigured
	}

	url := fmt.Sprintf("%s/applications/%s/restart", c.config.APIURL, c.config.WebclientAppUUID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAPIRequestFailed, err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAPIRequestFailed, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"%w: restart status %d, body: %s",
			ErrAPIRequestFailed,
			resp.StatusCode,
			string(respBody),
		)
	}

	c.logger.WarnContext(ctx, "Triggered Coolify application restart to apply domain changes",
		slog.String("app_uuid", c.config.WebclientAppUUID))

	return nil
}
