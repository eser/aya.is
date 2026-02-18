package speakerdeck

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

// jsonLDRegexp extracts JSON-LD script content from HTML.
var jsonLDRegexp = regexp.MustCompile(
	`<script\s+type="application/ld\+json">\s*([\s\S]*?)\s*</script>`,
)

// jsonLDData represents the relevant JSON-LD structure from SpeakerDeck pages.
type jsonLDData struct {
	AssociatedMedia *jsonLDMedia `json:"associatedMedia"`
}

type jsonLDMedia struct {
	ContentURL string `json:"contentUrl"`
}

// FetchPDFURL extracts the PDF download URL from a SpeakerDeck presentation page
// by parsing the JSON-LD structured data embedded in the HTML.
func (p *Provider) FetchPDFURL(ctx context.Context, presentationURL string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, presentationURL, nil)
	if err != nil {
		p.logger.WarnContext(ctx, "Failed to create request for PDF URL",
			slog.String("url", presentationURL),
			slog.Any("error", err))

		return ""
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.WarnContext(ctx, "Failed to fetch presentation page for PDF URL",
			slog.String("url", presentationURL),
			slog.Any("error", err))

		return ""
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		p.logger.WarnContext(ctx, "Unexpected status fetching presentation page",
			slog.String("url", presentationURL),
			slog.Int("status", resp.StatusCode))

		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.WarnContext(ctx, "Failed to read presentation page body",
			slog.String("url", presentationURL),
			slog.Any("error", err))

		return ""
	}

	return extractPDFURLFromHTML(string(body))
}

// extractPDFURLFromHTML parses JSON-LD from HTML to find the PDF content URL.
func extractPDFURLFromHTML(html string) string {
	matches := jsonLDRegexp.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) < 2 { //nolint:mnd
			continue
		}

		var data jsonLDData

		err := json.Unmarshal([]byte(match[1]), &data)
		if err != nil {
			continue
		}

		if data.AssociatedMedia != nil && data.AssociatedMedia.ContentURL != "" {
			return data.AssociatedMedia.ContentURL
		}
	}

	return ""
}

// BuildPDFURL constructs the PDF URL from a presentation ID and slug.
// Format: https://files.speakerdeck.com/presentations/{id}/{slug}.pdf
// This is a fallback when HTML scraping is not possible.
func BuildPDFURL(presentationID string, slug string) string {
	if presentationID == "" || slug == "" {
		return ""
	}

	return fmt.Sprintf(
		"https://files.speakerdeck.com/presentations/%s/%s.pdf",
		presentationID,
		slug,
	)
}
