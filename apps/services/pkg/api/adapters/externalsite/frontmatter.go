package externalsite

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// ParsedFrontmatter holds metadata extracted from a markdown file's frontmatter.
type ParsedFrontmatter struct {
	Title       string
	Date        *time.Time
	Slug        string
	Tags        []string
	Language    string // Detected locale code (e.g., "en", "tr")
	Description string
	Extra       map[string]any // All other fields
}

// ParseMarkdownFile splits frontmatter from body and parses metadata.
// Supports TOML (+++) and YAML (---) delimiters.
// filePath is used for slug/language inference when frontmatter lacks them.
func ParseMarkdownFile(content string, filePath string) (*ParsedFrontmatter, string, error) {
	fm, body, err := splitFrontmatter(content)
	if err != nil {
		return nil, "", err
	}

	if fm == nil {
		// No frontmatter — return empty metadata with full content as body
		slug := slugFromPath(filePath)

		return &ParsedFrontmatter{Slug: slug}, strings.TrimSpace(content), nil
	}

	result := mapToFrontmatter(fm)

	// Infer slug from file path if not set in frontmatter
	if result.Slug == "" {
		result.Slug = slugFromPath(filePath)
	}

	// Infer language from file path if not set in frontmatter
	if result.Language == "" {
		result.Language = languageFromPath(filePath)
	}

	return result, strings.TrimSpace(body), nil
}

// splitFrontmatter detects delimiter type, parses the frontmatter block,
// and returns the raw map + remaining body.
func splitFrontmatter(content string) (map[string]any, string, error) {
	trimmed := strings.TrimSpace(content)

	switch {
	case strings.HasPrefix(trimmed, "+++"):
		return parseTOMLFrontmatter(trimmed)
	case strings.HasPrefix(trimmed, "---"):
		return parseYAMLFrontmatter(trimmed)
	default:
		return nil, content, nil
	}
}

func parseTOMLFrontmatter(content string) (map[string]any, string, error) {
	// Remove opening +++
	rest := content[3:]

	idx := strings.Index(rest, "\n+++")
	if idx < 0 {
		return nil, content, nil
	}

	fmBlock := rest[:idx]
	body := rest[idx+4:] // skip \n+++

	var raw map[string]any
	if _, err := toml.Decode(fmBlock, &raw); err != nil {
		return nil, "", err
	}

	return raw, body, nil
}

func parseYAMLFrontmatter(content string) (map[string]any, string, error) {
	// Remove opening ---
	rest := content[3:]

	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return nil, content, nil
	}

	fmBlock := rest[:idx]
	body := rest[idx+4:] // skip \n---

	var raw map[string]any

	err := yaml.Unmarshal([]byte(fmBlock), &raw)
	if err != nil {
		return nil, "", err
	}

	return raw, body, nil
}

// mapToFrontmatter extracts well-known fields from a raw frontmatter map.
func mapToFrontmatter(raw map[string]any) *ParsedFrontmatter {
	fm := &ParsedFrontmatter{
		Extra: make(map[string]any),
	}

	for key, val := range raw {
		switch strings.ToLower(key) {
		case "title":
			fm.Title = toString(val)
		case "date":
			fm.Date = toTime(val)
		case "slug":
			fm.Slug = toString(val)
		case "description", "summary":
			fm.Description = toString(val)
		case "tags", "categories":
			fm.Tags = toStringSlice(val)
		case "language", "lang":
			fm.Language = toString(val)
		default:
			fm.Extra[key] = val
		}
	}

	// Check Zola's [extra] section for language if not found at top level
	if fm.Language == "" {
		if extra, ok := fm.Extra["extra"].(map[string]any); ok {
			if lang := toString(extra["language"]); lang != "" {
				fm.Language = lang
			} else if lang := toString(extra["lang"]); lang != "" {
				fm.Language = lang
			}
		}
	}

	// Infer language from tags when no explicit language field is set.
	// Common in Hugo blogs that don't use multilingual mode.
	if fm.Language == "" {
		fm.Language = languageFromTags(fm.Tags)
	}

	return fm
}

// toString converts an interface value to string.
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}

// toTime parses a time value from various formats.
func toTime(v any) *time.Time {
	switch val := v.(type) {
	case time.Time:
		return &val
	case string:
		for _, layout := range []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02",
		} {
			if t, err := time.Parse(layout, val); err == nil {
				return &t
			}
		}
	}

	return nil
}

// toStringSlice converts an interface value to []string.
func toStringSlice(v any) []string {
	switch val := v.(type) {
	case []any:
		result := make([]string, 0, len(val))

		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}

		return result
	case []string:
		return val
	}

	return nil
}

// slugFromPath extracts a slug from a file path.
// For Hugo page bundles (index.md inside a directory), uses the parent directory name.
// e.g., "content/posts/my-cool-post.md" → "my-cool-post"
// e.g., "content/posts/2024-01-02-hello-hugo/index.md" → "2024-01-02-hello-hugo".
func slugFromPath(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// If this is a page bundle index file, use parent directory name as slug
	if strings.ToLower(name) == "index" {
		dir := filepath.Dir(filePath)
		parent := filepath.Base(dir)

		if parent != "." && parent != "/" {
			return parent
		}
	}

	return name
}

// languageFromPath detects locale from Zola-style path conventions.
// e.g., "content/posts/my-post.tr.md" → "tr"
// e.g., "content/posts/my-post.md" → "" (no locale detected).
func languageFromPath(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	// Check for locale suffix like ".tr", ".en", ".de"
	parts := strings.Split(nameWithoutExt, ".")
	if len(parts) >= 2 {
		candidate := parts[len(parts)-1]
		if len(candidate) == 2 || len(candidate) == 5 { // "tr" or "pt-BR"
			return candidate
		}
	}

	return ""
}

// tagLanguageMap maps common tag values to ISO 639-1 locale codes.
// Keys must be lowercase.
var tagLanguageMap = map[string]string{
	"turkce":   "tr",
	"türkçe":   "tr",
	"turkish":  "tr",
	"english":  "en",
	"deutsch":  "de",
	"german":   "de",
	"français": "fr",
	"french":   "fr",
	"español":  "es",
	"spanish":  "es",
	"japanese": "ja",
	"korean":   "ko",
	"arabic":   "ar",
	"russian":  "ru",
	"chinese":  "zh",
}

// languageFromTags infers locale from well-known language tag values.
// Returns "" if no language tag is found.
func languageFromTags(tags []string) string {
	for _, tag := range tags {
		if locale, ok := tagLanguageMap[strings.ToLower(tag)]; ok {
			return locale
		}
	}

	return ""
}
