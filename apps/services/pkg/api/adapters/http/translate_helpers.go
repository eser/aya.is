package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
)

var ErrAITranslationNotAvailable = errors.New("AI translation not available")

type translationResult struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

func translateContent(
	ctx context.Context,
	aiModels *aifx.Registry,
	sourceLocale, targetLocale, title, summary, content string,
) (string, string, string, error) {
	if aiModels == nil {
		return "", "", "", ErrAITranslationNotAvailable
	}

	model := aiModels.GetDefault()
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
		return "", "", "", fmt.Errorf("AI generation failed: %w", err)
	}

	responseText := extractJSON(result.Text())

	var translated translationResult

	err = json.Unmarshal([]byte(responseText), &translated)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse AI translation response: %w", err)
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
