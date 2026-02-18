package speakerdeck

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/siteimporter"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

// presentationIDRegexp extracts the presentation ID from media:content URLs.
// e.g., https://speakerdeck.com/.../presentations/abc123/slide_0.jpg -> abc123.
var presentationIDRegexp = regexp.MustCompile(`/presentations/([^/]+)/`)

// FetchAll fetches all presentations from a SpeakerDeck user's RSS feed.
func (p *Provider) FetchAll(
	ctx context.Context,
	username string,
) ([]*siteimporter.ImportItem, error) {
	var allItems []*siteimporter.ImportItem

	for page := 1; ; page++ {
		rssURL := fmt.Sprintf("%s/%s.rss?page=%d", speakerDeckBaseURL, username, page)

		feed, err := p.parser.ParseURLWithContext(rssURL, ctx)
		if err != nil {
			if page == 1 {
				return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
			}
			// Later pages may 404 when there are no more items
			break
		}

		if len(feed.Items) == 0 {
			break
		}

		for _, item := range feed.Items {
			importItem := p.parseRSSItem(ctx, item)
			if importItem != nil {
				allItems = append(allItems, importItem)
			}
		}

		p.logger.DebugContext(ctx, "Fetched SpeakerDeck RSS page",
			slog.String("username", username),
			slog.Int("page", page),
			slog.Int("items", len(feed.Items)))
	}

	return allItems, nil
}

// parseRSSItem converts a gofeed item to an ImportItem.
func (p *Provider) parseRSSItem(ctx context.Context, item *gofeed.Item) *siteimporter.ImportItem {
	if item == nil {
		return nil
	}

	// Extract remote ID from media:content URL
	remoteID := extractPresentationID(item)
	if remoteID == "" {
		// Fallback: use GUID or link as remote ID
		remoteID = item.GUID
		if remoteID == "" {
			remoteID = item.Link
		}

		if remoteID == "" {
			p.logger.WarnContext(ctx, "Skipping SpeakerDeck item without ID",
				slog.String("title", item.Title))

			return nil
		}
	}

	// Extract published date
	publishedAt := time.Now()
	if item.PublishedParsed != nil {
		publishedAt = *item.PublishedParsed
	}

	// Extract description from content:encoded or description
	description := ""
	if item.Content != "" {
		description = item.Content
	} else if item.Description != "" {
		description = item.Description
	}

	// Strip HTML tags from description
	description = stripHTMLTags(description)

	// Extract thumbnail URL from media:content or enclosures
	thumbnailURL := extractThumbnailURL(item)

	return &siteimporter.ImportItem{
		RemoteID:     remoteID,
		Title:        item.Title,
		Description:  description,
		PublishedAt:  publishedAt,
		Link:         item.Link,
		ThumbnailURL: thumbnailURL,
		StoryKind:    "presentation",
		Properties:   make(map[string]any),
	}
}

// getMediaContentURLs extracts URLs from media:content RSS extensions.
func getMediaContentURLs(item *gofeed.Item) []string {
	if item.Extensions == nil {
		return nil
	}

	media, ok := item.Extensions["media"]
	if !ok {
		return nil
	}

	contents, ok := media["content"]
	if !ok {
		return nil
	}

	return collectExtensionURLs(contents)
}

// collectExtensionURLs extracts non-empty URL attributes from extension elements.
func collectExtensionURLs(elements []ext.Extension) []string {
	var urls []string

	for _, elem := range elements {
		if contentURL, ok := elem.Attrs["url"]; ok && contentURL != "" {
			urls = append(urls, contentURL)
		}
	}

	return urls
}

// extractPresentationID extracts the presentation ID from media extensions.
func extractPresentationID(item *gofeed.Item) string {
	// Try media:content URLs first
	for _, contentURL := range getMediaContentURLs(item) {
		matches := presentationIDRegexp.FindStringSubmatch(contentURL)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// Try enclosures
	for _, enc := range item.Enclosures {
		if enc.URL != "" {
			matches := presentationIDRegexp.FindStringSubmatch(enc.URL)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return ""
}

// extractThumbnailURL extracts the thumbnail URL from media extensions or enclosures.
func extractThumbnailURL(item *gofeed.Item) string {
	// Try media:content URL directly (it's typically a slide image)
	urls := getMediaContentURLs(item)
	if len(urls) > 0 {
		return urls[0]
	}

	// Try enclosures
	for _, enc := range item.Enclosures {
		if enc.URL != "" && strings.Contains(enc.Type, "image") {
			return enc.URL
		}
	}

	// Try image
	if item.Image != nil && item.Image.URL != "" {
		return item.Image.URL
	}

	return ""
}

// stripHTMLTags removes HTML tags from a string.
func stripHTMLTags(input string) string {
	var result strings.Builder

	inTag := false

	for _, char := range input {
		if char == '<' {
			inTag = true

			continue
		}

		if char == '>' {
			inTag = false

			continue
		}

		if !inTag {
			result.WriteRune(char)
		}
	}

	return strings.TrimSpace(result.String())
}
