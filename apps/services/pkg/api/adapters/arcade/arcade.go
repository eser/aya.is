package arcade

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Arcade struct {
	HTTPClient HTTPClient
	Config     Config
}

func New(config Config, httpClient HTTPClient) *Arcade {
	return &Arcade{
		Config:     config,
		HTTPClient: httpClient,
	}
}

func (arcade *Arcade) DoHTTPCall(ctx context.Context, req *http.Request) (_ []byte, err error) {
	res, err := arcade.HTTPClient.Do(req)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	defer func() {
		closeErr := res.Body.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return body, nil
}

func isTweetURL(text string) bool {
	const prefix = "https://twitter.com/"

	const statusPart = "/status/"

	// 1. Check for the main prefix
	if !strings.HasPrefix(text, prefix) {
		return false
	}

	// 2. Find "/status/" after the prefix
	// We search in the substring starting after the prefix.
	statusIndexInSubstring := strings.Index(text[len(prefix):], statusPart)
	if statusIndexInSubstring == -1 {
		// No "/status/" found after "https://twitter.com/"
		return false
	}

	// Calculate the absolute start index of the status ID part within the original string.
	startOfID := len(prefix) + statusIndexInSubstring + len(statusPart)

	// 3. Ensure there's something after "/status/"
	if startOfID >= len(text) {
		// URL ends exactly at "/status/"
		return false
	}

	// 4. Check if all characters after "/status/" are digits
	idPart := text[startOfID:]
	for _, r := range idPart {
		if r < '0' || r > '9' {
			// Found a non-digit character in the supposed ID part
			return false
		}
	}

	// If we passed all checks, it looks like a tweet URL structure.
	return true
}

// endsWithTweetURL checks if the given text ends with a string that is a valid Tweet URL.
func endsWithTweetURL(text string) bool {
	// Split the text into words
	words := strings.Fields(text)
	if len(words) == 0 {
		return false
	}

	// Get the last word
	lastWord := words[len(words)-1]

	// Check if the last word is a Tweet URL
	return isTweetURL(lastWord)
}

func (arcade *Arcade) GetRecentPostsByUsername( //nolint:funlen
	ctx context.Context,
	username string,
	userID string,
) ([]*profiles.ExternalPost, error) {
	url := arcade.Config.URL

	requestData := ExecuteToolRequest{ //nolint:exhaustruct
		Input: ExecuteToolInput{
			Username:   username,
			MaxResults: "100", // TODO(@eser) Hardcoded for now, consider making this configurable
		},
		ToolName: "X.SearchRecentTweetsByUsername@0.1.12", // TODO(@eser) Consider making this configurable
		UserID:   userID,
	}

	payloadBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	payloadReader := bytes.NewReader(payloadBytes)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		payloadReader,
	) // Use payloadReader
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	req.Header.Add("Authorization", "Bearer "+arcade.Config.APIKey)
	req.Header.Add("Content-Type", "application/json")

	result, err := arcade.DoHTTPCall(ctx, req)
	if err != nil {
		return nil, err
	}

	var response ExecuteToolResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err //nolint:wrapcheck
	}

	posts := make([]*profiles.ExternalPost, len(response.Output.Value.Data))
	count := 0

	for _, post := range response.Output.Value.Data {
		if strings.HasPrefix(post.Text, "@") {
			continue
		}

		if strings.HasPrefix(post.Text, "RT @") {
			continue
		}

		if endsWithTweetURL(post.Text) {
			continue
		}

		posts[count] = &profiles.ExternalPost{
			ID:        post.ID,
			Content:   post.Text,
			Permalink: post.TweetURL,
			CreatedAt: nil,
		}
		count++
	}

	return posts[:count], nil
}
