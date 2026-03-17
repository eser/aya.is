package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiBaseURL = "https://api.resend.com"

	// httpClientTimeout is the default HTTP client timeout for Resend API requests.
	httpClientTimeout = 30 * time.Second
)

// ErrResendAPI indicates an error response from the Resend API.
var ErrResendAPI = errors.New("resend API error")

// Client is an HTTP client for the Resend email API.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new Resend API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{ //nolint:exhaustruct // std lib type
			Timeout: httpClientTimeout,
		},
	}
}

type sendEmailRequest struct {
	From    string   `json:"from"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	To      []string `json:"to"`
}

type apiError struct {
	Name       string `json:"name"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// SendEmail sends an email via the Resend API.
func (c *Client) SendEmail(
	ctx context.Context,
	from string,
	recipient string,
	subject string,
	htmlBody string,
) error {
	reqBody := sendEmailRequest{
		From:    from,
		To:      []string{recipient},
		Subject: subject,
		HTML:    htmlBody,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiBaseURL+"/emails",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)

		var apiErr apiError

		jsonErr := json.Unmarshal(respBody, &apiErr)
		if jsonErr == nil && apiErr.Message != "" {
			return fmt.Errorf("%w (%d): %s", ErrResendAPI, resp.StatusCode, apiErr.Message)
		}

		return fmt.Errorf("%w (%d): %s", ErrResendAPI, resp.StatusCode, string(respBody))
	}

	return nil
}
