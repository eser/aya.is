package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
)

var (
	ErrAITranslationNotAvailable = errors.New("AI translation not available")
	ErrAIGenerationFailed        = errors.New("AI generation failed")
	ErrFailedToParseAIResponse   = errors.New("failed to parse AI translation response")
)

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

	prompt := fmt.Sprintf(
		`Translate the following content from locale "%s" to locale "%s".
Return ONLY a valid JSON object with exactly these three keys: "title", "summary", "content".
Do not include any other text, explanation, or markdown formatting.
Preserve all markdown formatting in the content field.
Do not translate code blocks, URLs, or technical terms.

Input:
Title: %s
Summary: %s
Content:
%s`,
		sourceLocale,
		targetLocale,
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
