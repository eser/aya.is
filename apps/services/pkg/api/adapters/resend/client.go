package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiBaseURL = "https://api.resend.com"

// Client is an HTTP client for the Resend email API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Resend API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type sendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type sendEmailResponse struct {
	ID string `json:"id"`
}

type apiError struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

// SendEmail sends an email via the Resend API.
func (c *Client) SendEmail(
	ctx context.Context,
	from string,
	to string,
	subject string,
	htmlBody string,
) error {
	reqBody := sendEmailRequest{
		From:    from,
		To:      []string{to},
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
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)

		var apiErr apiError

		jsonErr := json.Unmarshal(respBody, &apiErr)
		if jsonErr == nil && apiErr.Message != "" {
			return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, apiErr.Message)
		}

		return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
