package externalsite

import (
	neturl "net/url"
	"path"
	"regexp"
	"strings"
)

// Hugo shortcode patterns:
//   {{< shortcode args >}}       — raw content shortcode
//   {{% shortcode args %}}       — markdown-processed shortcode
//   {{< shortcode args />}}      — self-closing variant
//   {{< shortcode args >}} body {{< /shortcode >}}  — paired shortcode

var (
	// Matches {{< tweet user="X" id="Y" >}} or {{< twitter user="X" id="Y" >}}.
	reTweet = regexp.MustCompile(
		`\{\{[<%]\s*(?:tweet|twitter)\s+` +
			`(?:user="([^"]+)"\s+)?` +
			`(?:id="?(\d+)"?)` +
			`\s*[/%]?[>%]\}\}`,
	)

	// Matches {{< youtube ID >}} or {{< youtube id="ID" >}}.
	reYouTube = regexp.MustCompile(
		`\{\{[<%]\s*youtube\s+` +
			`(?:id="?)?([a-zA-Z0-9_-]+)"?` +
			`\s*[/%]?[>%]\}\}`,
	)

	// Matches {{< vimeo ID >}} or {{< vimeo id="ID" >}}.
	reVimeo = regexp.MustCompile(
		`\{\{[<%]\s*vimeo\s+` +
			`(?:id="?)?(\d+)"?` +
			`\s*[/%]?[>%]\}\}`,
	)

	// Matches {{< gist user ID >}} or {{< gist user ID file >}}.
	reGist = regexp.MustCompile(
		`\{\{[<%]\s*gist\s+"?([a-zA-Z0-9_-]+)"?\s+"?([a-fA-F0-9]+)"?` +
			`(?:\s+"?([^"}\s]+)"?)?` +
			`\s*[/%]?[>%]\}\}`,
	)

	// Matches {{< figure src="..." alt="..." caption="..." >}}.
	reFigureSrc = regexp.MustCompile(
		`\{\{[<%]\s*figure\s+[^}>]*?src="([^"]+)"[^}>]*?[/%]?[>%]\}\}`,
	)
	reFigureAlt = regexp.MustCompile(`alt="([^"]*)"`)
	reFigureCap = regexp.MustCompile(`caption="([^"]*)"`)

	// Catch-all: any remaining Hugo shortcodes (paired or self-closing).
	reShortcodePaired = regexp.MustCompile(
		`\{\{[<%]\s*/?\s*\w+[^}>]*[>%]\}\}`,
	)

	// Zola shortcodes: {{ shortcode(args) }}.
	reZolaShortcode = regexp.MustCompile(
		`\{\{\s*\w+\([^)]*\)\s*\}\}`,
	)
)

// SanitizeContent transforms Hugo/Zola shortcodes in markdown content
// into MDX-compatible equivalents (links, embeds, images, or removal).
func SanitizeContent(content string) string {
	result := content

	// Convert tweet shortcodes → Twitter embed link
	result = reTweet.ReplaceAllStringFunc(result, func(match string) string {
		parts := reTweet.FindStringSubmatch(match)
		if len(parts) >= 3 && parts[2] != "" {
			user := parts[1]
			if user == "" {
				user = "x"
			}

			return "https://twitter.com/" + user + "/status/" + parts[2]
		}

		return ""
	})

	// Convert youtube shortcodes → YouTube link
	result = reYouTube.ReplaceAllString(result,
		"https://www.youtube.com/watch?v=$1")

	// Convert vimeo shortcodes → Vimeo link
	result = reVimeo.ReplaceAllString(result,
		"https://vimeo.com/$1")

	// Convert gist shortcodes → Gist link
	result = reGist.ReplaceAllStringFunc(result, func(match string) string {
		parts := reGist.FindStringSubmatch(match)
		if len(parts) >= 3 {
			url := "https://gist.github.com/" + parts[1] + "/" + parts[2]

			return url
		}

		return ""
	})

	// Convert figure shortcodes → markdown image
	result = reFigureSrc.ReplaceAllStringFunc(result, func(match string) string {
		srcParts := reFigureSrc.FindStringSubmatch(match)
		if len(srcParts) < 2 {
			return ""
		}

		src := srcParts[1]

		alt := ""
		if altMatch := reFigureAlt.FindStringSubmatch(match); len(altMatch) >= 2 {
			alt = altMatch[1]
		}

		if caption := reFigureCap.FindStringSubmatch(match); alt == "" && len(caption) >= 2 {
			alt = caption[1]
		}

		return "![" + alt + "](" + src + ")"
	})

	// Remove any remaining Hugo shortcodes (paired closing tags, unknown shortcodes)
	result = reShortcodePaired.ReplaceAllString(result, "")

	// Remove Zola shortcodes
	result = reZolaShortcode.ReplaceAllString(result, "")

	// Clean up: collapse 3+ consecutive blank lines into 2
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")

	return strings.TrimSpace(result)
}

// SanitizeDescription strips HTML tags and markdown syntax from a description string,
// producing clean plain text suitable for display in story cards.
func SanitizeDescription(s string) string {
	// Strip HTML tags
	inTag := false

	var result strings.Builder

	for _, ch := range s {
		if ch == '<' {
			inTag = true

			continue
		}

		if ch == '>' {
			inTag = false

			result.WriteRune(' ')

			continue
		}

		if !inTag {
			result.WriteRune(ch)
		}
	}

	text := result.String()

	// Strip markdown heading markers (# ## ### etc. and === --- underlines)
	text = reMarkdownHeading.ReplaceAllString(text, "$1")
	text = reMarkdownUnderlineHeading.ReplaceAllString(text, "")

	// Strip markdown image/link syntax but keep alt text
	text = reMarkdownImageLink.ReplaceAllString(text, "$1")

	// Collapse whitespace
	text = reMultipleSpaces.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

var (
	reMarkdownHeading          = regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`)
	reMarkdownUnderlineHeading = regexp.MustCompile(`(?m)^[=\-]{3,}\s*$`)
	reMarkdownImageLink        = regexp.MustCompile(`!?\[([^\]]*)\]\([^)]+\)`)
	reMultipleSpaces           = regexp.MustCompile(`\s{2,}`)
)

// Patterns for image references in markdown.
var (
	// Matches markdown images: ![alt](url).
	reMarkdownImage = regexp.MustCompile(`(!\[[^\]]*\]\()([^)]+)(\))`)

	// Matches HTML img tags: <img ... src="url" ...>.
	reHTMLImgSrc = regexp.MustCompile(`(<img\s[^>]*?src=")([^"]+)("[^>]*>)`)
)

// ResolveRelativeImages rewrites relative image paths in markdown content
// to absolute raw.githubusercontent.com URLs.
//
// Parameters:
//   - content: the markdown body
//   - ownerRepo: e.g. "Ardakilic/arda.pw"
//   - branch: e.g. "master"
//   - filePath: the source file path in the repo, e.g. "content/posts/my-post/index.md"
func ResolveRelativeImages(content, ownerRepo, branch, filePath string) string {
	fileDir := path.Dir(filePath) // e.g. "content/posts/my-post"
	baseURL := "https://raw.githubusercontent.com/" + escapeOwnerRepo(
		ownerRepo,
	) + "/" + neturl.PathEscape(
		branch,
	) + "/"

	resolver := func(imageURL string) string {
		// Skip absolute URLs and protocol-relative URLs
		if strings.HasPrefix(imageURL, "http://") ||
			strings.HasPrefix(imageURL, "https://") ||
			strings.HasPrefix(imageURL, "//") {
			return imageURL
		}

		// Skip data URIs
		if strings.HasPrefix(imageURL, "data:") {
			return imageURL
		}

		// Resolve relative path against the file's directory
		resolved := path.Join(fileDir, imageURL) // path.Join cleans ./  and ../

		return baseURL + resolved
	}

	result := content

	// Rewrite markdown images: ![alt](relative/path)
	result = reMarkdownImage.ReplaceAllStringFunc(result, func(match string) string {
		parts := reMarkdownImage.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		return parts[1] + resolver(parts[2]) + parts[3]
	})

	// Rewrite HTML img src: <img src="relative/path">
	result = reHTMLImgSrc.ReplaceAllStringFunc(result, func(match string) string {
		parts := reHTMLImgSrc.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		return parts[1] + resolver(parts[2]) + parts[3]
	})

	return result
}
