package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
)

var (
	ErrAITranslationNotAvailable = errors.New("AI translation not available")
	ErrAIGenerationFailed        = errors.New("AI generation failed")
	ErrFailedToParseAIResponse   = errors.New("failed to parse AI translation response")
)

// localeToLanguageName maps locale codes to full language names for AI prompts.
// Using full names prevents LLMs from defaulting to English when given short codes.
var localeToLanguageName = map[string]string{ //nolint:gochecknoglobals
	"ar":    "Arabic",
	"de":    "German",
	"en":    "English",
	"es":    "Spanish",
	"fr":    "French",
	"it":    "Italian",
	"ja":    "Japanese",
	"ko":    "Korean",
	"nl":    "Dutch",
	"pt-PT": "Portuguese (Portugal)",
	"ru":    "Russian",
	"tr":    "Turkish",
	"zh-CN": "Chinese (Simplified)",
}

// languageNameForLocale returns the full language name for a locale code.
// Falls back to the locale code itself if not found.
func languageNameForLocale(locale string) string {
	if name, ok := localeToLanguageName[locale]; ok {
		return name
	}

	return locale
}

type translationResult struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

// AIContentTranslator implements stories.ContentTranslator and profiles.ContentTranslator
// using the aifx registry for AI-powered translation.
type AIContentTranslator struct {
	aiModels *aifx.Registry
}

// NewAIContentTranslator creates a new AIContentTranslator.
func NewAIContentTranslator(aiModels *aifx.Registry) *AIContentTranslator {
	return &AIContentTranslator{aiModels: aiModels}
}

// Translate implements the ContentTranslator interface.
func (t *AIContentTranslator) Translate(
	ctx context.Context,
	sourceLocale, targetLocale string,
	title, summary, content string,
) (string, string, string, error) {
	if t.aiModels == nil {
		return "", "", "", ErrAITranslationNotAvailable
	}

	model := t.aiModels.GetDefault()
	if model == nil {
		return "", "", "", ErrAITranslationNotAvailable
	}

	sourceLang := languageNameForLocale(sourceLocale)
	targetLang := languageNameForLocale(targetLocale)

	prompt := fmt.Sprintf(
		`Translate the following content from %s to %s.
The output MUST be written entirely in %s.
Return ONLY a valid JSON object with exactly these three keys: "title", "summary", "content".
Do not include any other text, explanation, or markdown formatting.
Preserve all markdown formatting in the content field.
Do not translate code blocks, URLs, or technical terms.

Input:
Title: %s
Summary: %s
Content:
%s`,
		sourceLang,
		targetLang,
		targetLang,
		title,
		summary,
		content,
	)

	result, err := model.GenerateText(ctx, &aifx.GenerateTextOptions{
		Messages: []aifx.Message{
			aifx.NewTextMessage(aifx.RoleUser, prompt),
		},
		System:    "You are a professional translator. Translate content accurately while preserving formatting.",
		MaxTokens: 8192,
	})
	if err != nil {
		logAIErrorClassification(ctx, "translation", err)

		return "", "", "", fmt.Errorf("%w: %w", ErrAIGenerationFailed, err)
	}

	responseText := extractJSON(result.Text())

	var translated translationResult

	err = json.Unmarshal([]byte(responseText), &translated)
	if err != nil {
		return "", "", "", fmt.Errorf("%w: %w", ErrFailedToParseAIResponse, err)
	}

	return translated.Title, translated.Summary, translated.Content, nil
}

// AIContentGenerator implements profiles.ContentGenerator
// using the aifx registry for AI-powered content generation.
type AIContentGenerator struct {
	aiModels *aifx.Registry
}

// NewAIContentGenerator creates a new AIContentGenerator.
func NewAIContentGenerator(aiModels *aifx.Registry) *AIContentGenerator {
	return &AIContentGenerator{aiModels: aiModels}
}

// GenerateCV implements the ContentGenerator interface.
// It generates a professional CV page from the user's profile data using AI.
func (g *AIContentGenerator) GenerateCV(
	ctx context.Context,
	locale string,
	profileTitle string,
	profileDescription string,
	linkedInURL string,
	links []*profiles.ProfileLinkBrief,
	contributions []*profiles.ProfileMembership,
) (string, string, string, error) {
	if g.aiModels == nil {
		return "", "", "", ErrAITranslationNotAvailable
	}

	model := g.aiModels.GetDefault()
	if model == nil {
		return "", "", "", ErrAITranslationNotAvailable
	}

	targetLang := languageNameForLocale(locale)

	// Build context from links
	var linksContext strings.Builder

	for _, link := range links {
		if link.URI != "" {
			linksContext.WriteString(fmt.Sprintf("- %s: %s", link.Kind, link.URI))

			if link.Title != "" {
				linksContext.WriteString(fmt.Sprintf(" (%s)", link.Title))
			}

			linksContext.WriteString("\n")
		}
	}

	// Build context from contributions (organization memberships)
	var contributionsContext strings.Builder

	for _, membership := range contributions {
		if membership.Profile == nil {
			continue
		}

		role := membership.Kind
		orgName := membership.Profile.Title

		contributionsContext.WriteString(fmt.Sprintf("- %s at %s", role, orgName))

		if membership.StartedAt != nil {
			contributionsContext.WriteString(
				fmt.Sprintf(" (since %s)", membership.StartedAt.Format("January 2006")),
			)
		}

		if membership.FinishedAt != nil {
			contributionsContext.WriteString(
				fmt.Sprintf(" (until %s)", membership.FinishedAt.Format("January 2006")),
			)
		}

		if membership.Profile.Description != "" {
			contributionsContext.WriteString(
				"\n  Organization: " + membership.Profile.Description,
			)
		}

		contributionsContext.WriteString("\n")
	}

	prompt := fmt.Sprintf(
		`Based on the profile information provided below, generate a professional CV page in markdown.
The output MUST be written entirely in %s.

Profile Name: %s
Profile Bio: %s
LinkedIn Profile: %s

Connected Accounts and Links:
%s
Organization Memberships and Contributions:
%s
Instructions:
- Generate a professional CV/resume page in markdown format.
- Structure the CV with ## Experience, ## Education, and ## Certificates sections.
- For experience entries, use ### Role at Company format with date ranges and achievement bullet points.
- Use all available context to generate realistic professional content.
- The user will review and edit this draft afterward, so focus on providing a solid starting structure.
- Do not invent specific details that cannot be inferred from the provided context.
- Keep the tone professional and concise.

Return ONLY a valid JSON object with exactly these three keys: "title", "summary", "content".
Do not include any other text, explanation, or markdown formatting around the JSON.
The "content" field should contain the full markdown CV.`,
		targetLang,
		profileTitle,
		profileDescription,
		linkedInURL,
		linksContext.String(),
		contributionsContext.String(),
	)

	result, err := model.GenerateText(ctx, &aifx.GenerateTextOptions{
		Messages: []aifx.Message{
			aifx.NewTextMessage(aifx.RoleUser, prompt),
		},
		System: "You are a professional CV/resume writer. " +
			"Generate well-structured, professional CV content based on the provided profile information. " +
			"Always respond with valid JSON only.",
		MaxTokens: 8192,
	})
	if err != nil {
		logAIErrorClassification(ctx, "content generation", err)

		return "", "", "", fmt.Errorf("%w: %w", ErrAIGenerationFailed, err)
	}

	responseText := extractJSON(result.Text())

	var generated translationResult

	err = json.Unmarshal([]byte(responseText), &generated)
	if err != nil {
		return "", "", "", fmt.Errorf("%w: %w", ErrFailedToParseAIResponse, err)
	}

	return generated.Title, generated.Summary, generated.Content, nil
}

// logAIErrorClassification emits a structured warning when an AI call fails,
// tagging the log with the classified error category for observability.
func logAIErrorClassification(ctx context.Context, operation string, err error) {
	switch {
	case errors.Is(err, aifx.ErrRateLimited):
		slog.WarnContext(ctx, "AI rate limited", slog.String("operation", operation))
	case errors.Is(err, aifx.ErrAuthFailed):
		slog.ErrorContext(ctx, "AI authentication failed", slog.String("operation", operation))
	case errors.Is(err, aifx.ErrInsufficientCredits):
		slog.ErrorContext(ctx, "AI insufficient credits", slog.String("operation", operation))
	case errors.Is(err, aifx.ErrServiceUnavailable):
		slog.WarnContext(ctx, "AI service unavailable", slog.String("operation", operation))
	}
}

// extractJSON strips markdown code fences from AI responses.
// LLMs commonly wrap JSON in ```json ... ``` despite being told not to.
func extractJSON(text string) string {
	trimmed := strings.TrimSpace(text)

	if strings.HasPrefix(trimmed, "```") {
		// Remove opening fence (```json or ```)
		if idx := strings.Index(trimmed, "\n"); idx != -1 {
			trimmed = trimmed[idx+1:]
		}

		// Remove closing fence
		if idx := strings.LastIndex(trimmed, "```"); idx != -1 {
			trimmed = trimmed[:idx]
		}

		return strings.TrimSpace(trimmed)
	}

	return trimmed
}
