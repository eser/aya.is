package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/aifx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

const (
	charLimitX        = 280
	charLimitLinkedIn = 3000
	charLimitMastodon = 500
	charLimitBluesky  = 300
)

// Platform character limits for the optimize action.
var platformCharLimits = map[string]int{ //nolint:gochecknoglobals
	"x":        charLimitX,
	"linkedin": charLimitLinkedIn,
	"mastodon": charLimitMastodon,
	"bluesky":  charLimitBluesky,
}

// validShareWizardActions lists the supported AI assist actions.
var validShareWizardActions = map[string]bool{ //nolint:gochecknoglobals
	"summarize":   true,
	"adjust_tone": true,
	"optimize":    true,
	"hashtags":    true,
	"translate":   true,
}

// validShareWizardPlatforms lists the supported social platforms.
var validShareWizardPlatforms = map[string]bool{ //nolint:gochecknoglobals
	"x":        true,
	"linkedin": true,
	"mastodon": true,
	"bluesky":  true,
}

// validShareWizardTones lists the supported tone options.
var validShareWizardTones = map[string]bool{ //nolint:gochecknoglobals
	"professional": true,
	"casual":       true,
	"enthusiastic": true,
	"informative":  true,
}

func RegisterHTTPRoutesForShareWizard( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	aiModels *aifx.Registry,
) {
	// AI assist endpoint for Share Wizard
	routes.Route(
		"POST /{locale}/stories/{slug}/share/ai-assist",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}

			var requestBody struct {
				Action       string `json:"action"`
				Text         string `json:"text"`
				StoryContent string `json:"story_content"`
				Platform     string `json:"platform"`
				Tone         string `json:"tone"`
				TargetLocale string `json:"target_locale"`
			}

			err := ctx.ParseJSONBody(&requestBody)
			if err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Validate action
			if !validShareWizardActions[requestBody.Action] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid action"))
			}

			// Validate platform when relevant
			if requestBody.Platform != "" && !validShareWizardPlatforms[requestBody.Platform] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid platform"))
			}

			// Validate tone when relevant
			if requestBody.Action == "adjust_tone" && !validShareWizardTones[requestBody.Tone] {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid tone"))
			}

			// Validate target_locale for translate action
			if requestBody.Action == "translate" && requestBody.TargetLocale == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("target_locale is required for translate action"),
				)
			}

			// Get AI model
			if aiModels == nil {
				return ctx.Results.Error(
					http.StatusServiceUnavailable,
					httpfx.WithErrorMessage("AI service not available"),
				)
			}

			model := aiModels.GetDefault()
			if model == nil {
				return ctx.Results.Error(
					http.StatusServiceUnavailable,
					httpfx.WithErrorMessage("AI service not available"),
				)
			}

			// Build prompts based on action
			systemPrompt, userPrompt := buildShareWizardPrompts(
				requestBody.Action,
				requestBody.Text,
				requestBody.StoryContent,
				requestBody.Platform,
				requestBody.Tone,
				requestBody.TargetLocale,
			)

			// Extend HTTP write deadline for long-running AI call
			rc := http.NewResponseController(ctx.ResponseWriter)
			_ = rc.SetWriteDeadline(time.Now().Add(aiCallDeadlineSec * time.Second))

			result, err := model.GenerateText(ctx.Request.Context(), &aifx.GenerateTextOptions{
				Messages: []aifx.Message{
					aifx.NewTextMessage(aifx.RoleUser, userPrompt),
				},
				Tools:          nil,
				System:         systemPrompt,
				MaxTokens:      aiMaxTokens,
				Temperature:    nil,
				TopP:           nil,
				StopWords:      nil,
				ToolChoice:     "",
				ResponseFormat: nil,
				ThinkingBudget: nil,
				SafetySettings: nil,
				Extensions:     nil,
			})
			if err != nil {
				logAIErrorClassification(ctx.Request.Context(), "share_wizard", err)

				logger.ErrorContext(ctx.Request.Context(), "Share wizard AI assist failed",
					slog.String("error", err.Error()),
					slog.String("action", requestBody.Action))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("AI generation failed"),
				)
			}

			wrappedResponse := map[string]any{
				"data": map[string]any{
					"result": result.Text(),
				},
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		},
	).
		HasSummary("Share Wizard AI Assist").
		HasDescription("Generate AI-assisted text for sharing stories on social media.").
		HasResponse(http.StatusOK)
}

// buildShareWizardPrompts builds system and user prompts based on the requested action.
func buildShareWizardPrompts( //nolint:funlen
	action, text, storyContent, platform, tone, targetLocale string,
) (string, string) {
	systemPrompt := "You are a social media content assistant. " +
		"Generate concise, engaging text suitable for social media sharing. " +
		"Return only the generated text without any explanation, quotes, or formatting."

	var userPrompt string

	switch action {
	case "summarize":
		userPrompt = fmt.Sprintf(
			"Summarize the following content into a concise social media post suitable for sharing.\n\n"+
				"Story content for context:\n%s\n\n"+
				"Current draft (if any):\n%s",
			storyContent,
			text,
		)

	case "adjust_tone":
		userPrompt = fmt.Sprintf(
			"Rewrite the following text in a %s tone. "+
				"Keep the core message but adjust the style and word choice.\n\n"+
				"Text:\n%s",
			tone,
			text,
		)

	case "optimize":
		charLimit := platformCharLimits[platform]
		if charLimit == 0 {
			charLimit = charLimitX
		}

		userPrompt = fmt.Sprintf(
			"Optimize the following text for %s. "+
				"Keep it within %d characters. "+
				"Make it engaging and platform-appropriate.\n\n"+
				"Text:\n%s",
			platform,
			charLimit,
			text,
		)

	case "hashtags":
		userPrompt = fmt.Sprintf(
			"Generate relevant hashtags for the following content. "+
				"Return only hashtags separated by spaces. "+
				"Generate between 3 and 8 hashtags.\n\n"+
				"Content:\n%s\n\n"+
				"Current text:\n%s",
			storyContent,
			text,
		)

	case "translate":
		targetLang := languageNameForLocale(targetLocale)

		userPrompt = fmt.Sprintf(
			"Translate the following social media post to %s. "+
				"Keep the tone and style appropriate for social media. "+
				"Return only the translated text.\n\n"+
				"Text:\n%s",
			targetLang,
			text,
		)

	default:
		userPrompt = text
	}

	return systemPrompt, userPrompt
}
