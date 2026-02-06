package aifx

import (
	"context"
	"encoding/json"
	"errors"

	"google.golang.org/genai"
)

// genaiRoleUser is the role name for user messages in the genai SDK.
const genaiRoleUser = "user"

// genaiRoleModel is the role name for assistant/model messages in the genai SDK.
const genaiRoleModel = "model"

// Sentinel errors shared across Google adapters.
var ErrGenAINilResponse = errors.New("genai returned nil response")

// mapMessagesToGenAI converts unified messages to genai Content slices.
// Gemini uses "user" and "model" as role names; "assistant" is mapped to "model".
// System messages are excluded here (they go into SystemInstruction).
// Tool result messages are mapped as "user" role with FunctionResponse parts.
func mapMessagesToGenAI(messages []Message) []*genai.Content {
	contents := make([]*genai.Content, 0, len(messages))

	for _, msg := range messages {
		// Skip system messages — they are handled separately via SystemInstruction.
		if msg.Role == RoleSystem {
			continue
		}

		role := mapRoleToGenAI(msg.Role)
		parts := mapContentBlocksToGenAIParts(msg.Content)

		if len(parts) > 0 {
			contents = append(contents, &genai.Content{
				Role:  role,
				Parts: parts,
			})
		}
	}

	return contents
}

// mapRoleToGenAI converts a unified role to the genai role string.
func mapRoleToGenAI(role Role) string {
	switch role {
	case RoleUser, RoleTool:
		return genaiRoleUser
	case RoleAssistant:
		return genaiRoleModel
	default:
		return genaiRoleUser
	}
}

// mapContentBlocksToGenAIParts converts content blocks to genai Part slices.
func mapContentBlocksToGenAIParts(blocks []ContentBlock) []*genai.Part {
	parts := make([]*genai.Part, 0, len(blocks))

	for idx := range blocks {
		block := &blocks[idx]

		switch block.Type {
		case ContentBlockText:
			parts = append(parts, &genai.Part{
				Text: block.Text,
			})

		case ContentBlockImage:
			if block.Image != nil {
				parts = append(parts, mapImagePartToGenAI(block.Image))
			}

		case ContentBlockAudio:
			if block.Audio != nil {
				parts = append(parts, mapAudioPartToGenAI(block.Audio))
			}

		case ContentBlockFile:
			if block.File != nil {
				parts = append(parts, mapFilePartToGenAI(block.File))
			}

		case ContentBlockToolCall:
			if block.ToolCall != nil {
				parts = append(parts, mapToolCallToGenAIPart(block.ToolCall))
			}

		case ContentBlockToolResult:
			if block.ToolResult != nil {
				parts = append(parts, mapToolResultToGenAIPart(block.ToolResult))
			}
		}
	}

	return parts
}

// mapImagePartToGenAI converts an ImagePart to a genai Part.
// Data URLs and raw bytes produce InlineData; HTTP URLs produce FileData.
func mapImagePartToGenAI(img *ImagePart) *genai.Part {
	// If raw data is already available, use InlineData.
	if len(img.Data) > 0 {
		mimeType := img.MIMEType
		if mimeType == "" {
			mimeType = "image/png"
		}

		return &genai.Part{
			InlineData: &genai.Blob{
				MIMEType: mimeType,
				Data:     img.Data,
			},
		}
	}

	// If the URL is a data: URI, decode it.
	if IsDataURL(img.URL) {
		mimeType, data, err := DecodeDataURL(img.URL)
		if err == nil {
			return &genai.Part{
				InlineData: &genai.Blob{
					MIMEType: mimeType,
					Data:     data,
				},
			}
		}
	}

	// Fall back to FileData for HTTP/GCS URIs.
	mimeType := img.MIMEType
	if mimeType == "" {
		mimeType = DetectMIMEFromURL(img.URL)
	}

	return &genai.Part{
		FileData: &genai.FileData{
			MIMEType: mimeType,
			FileURI:  img.URL,
		},
	}
}

// mapAudioPartToGenAI converts an AudioPart to a genai Part.
func mapAudioPartToGenAI(audio *AudioPart) *genai.Part {
	if len(audio.Data) > 0 {
		mimeType := audio.MIMEType
		if mimeType == "" {
			mimeType = "audio/mpeg"
		}

		return &genai.Part{
			InlineData: &genai.Blob{
				MIMEType: mimeType,
				Data:     audio.Data,
			},
		}
	}

	if IsDataURL(audio.URL) {
		mimeType, data, err := DecodeDataURL(audio.URL)
		if err == nil {
			return &genai.Part{
				InlineData: &genai.Blob{
					MIMEType: mimeType,
					Data:     data,
				},
			}
		}
	}

	mimeType := audio.MIMEType
	if mimeType == "" {
		mimeType = DetectMIMEFromURL(audio.URL)
	}

	return &genai.Part{
		FileData: &genai.FileData{
			MIMEType: mimeType,
			FileURI:  audio.URL,
		},
	}
}

// mapFilePartToGenAI converts a FilePart to a genai Part.
func mapFilePartToGenAI(file *FilePart) *genai.Part {
	mimeType := file.MIMEType
	if mimeType == "" {
		mimeType = DetectMIMEFromURL(file.URI)
	}

	return &genai.Part{
		FileData: &genai.FileData{
			MIMEType: mimeType,
			FileURI:  file.URI,
		},
	}
}

// mapToolCallToGenAIPart converts a ToolCall to a genai FunctionCall Part.
func mapToolCallToGenAIPart(tc *ToolCall) *genai.Part {
	var args map[string]any
	if len(tc.Arguments) > 0 {
		_ = json.Unmarshal(tc.Arguments, &args)
	}

	return &genai.Part{
		FunctionCall: &genai.FunctionCall{
			Name: tc.Name,
			Args: args,
		},
	}
}

// mapToolResultToGenAIPart converts a ToolResult to a genai FunctionResponse Part.
func mapToolResultToGenAIPart(tr *ToolResult) *genai.Part {
	response := map[string]any{
		"result": tr.Content,
	}

	if tr.IsError {
		response["error"] = tr.Content
	}

	return &genai.Part{
		FunctionResponse: &genai.FunctionResponse{
			Name:     tr.ToolCallID,
			Response: response,
		},
	}
}

// mapToolsToGenAI converts unified tool definitions to genai Tool format.
// Uses ParametersJsonSchema for JSON Schema parameters.
func mapToolsToGenAI(tools []ToolDefinition) []*genai.Tool {
	if len(tools) == 0 {
		return nil
	}

	declarations := make([]*genai.FunctionDeclaration, 0, len(tools))

	for idx := range tools {
		tool := &tools[idx]
		decl := &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
		}

		// Convert JSON Schema parameters to the map[string]any format
		// expected by ParametersJsonSchema.
		if len(tool.Parameters) > 0 {
			var schemaMap map[string]any

			err := json.Unmarshal(tool.Parameters, &schemaMap)
			if err == nil {
				decl.ParametersJsonSchema = schemaMap
			}
		}

		declarations = append(declarations, decl)
	}

	return []*genai.Tool{
		{
			FunctionDeclarations: declarations,
		},
	}
}

// mapSafetySettings converts unified safety settings to genai SafetySetting format.
func mapSafetySettings(settings []SafetySetting) []*genai.SafetySetting {
	if len(settings) == 0 {
		return nil
	}

	genaiSettings := make([]*genai.SafetySetting, 0, len(settings))

	for idx := range settings {
		setting := &settings[idx]

		genaiSettings = append(genaiSettings, &genai.SafetySetting{
			Category:  genai.HarmCategory(setting.Category),
			Threshold: genai.HarmBlockThreshold(setting.Threshold),
		})
	}

	return genaiSettings
}

// buildSystemInstruction creates a genai.Content for the system instruction.
func buildSystemInstruction(system string) *genai.Content {
	if system == "" {
		return nil
	}

	return &genai.Content{
		Parts: []*genai.Part{
			{Text: system},
		},
	}
}

// buildGenerateContentConfig constructs the genai.GenerateContentConfig from unified options.
func buildGenerateContentConfig(opts *GenerateTextOptions) *genai.GenerateContentConfig {
	config := &genai.GenerateContentConfig{}

	// System instruction.
	config.SystemInstruction = buildSystemInstruction(opts.System)

	// Generation parameters.
	if opts.MaxTokens > 0 {
		config.MaxOutputTokens = int32(opts.MaxTokens) //nolint:gosec
	}

	if opts.Temperature != nil {
		config.Temperature = genai.Ptr(float32(*opts.Temperature))
	}

	if opts.TopP != nil {
		config.TopP = genai.Ptr(float32(*opts.TopP))
	}

	if len(opts.StopWords) > 0 {
		config.StopSequences = opts.StopWords
	}

	// Tools.
	if len(opts.Tools) > 0 {
		config.Tools = mapToolsToGenAI(opts.Tools)
	}

	// Safety settings.
	if len(opts.SafetySettings) > 0 {
		config.SafetySettings = mapSafetySettings(opts.SafetySettings)
	}

	// Structured output (JSON schema).
	if opts.ResponseFormat != nil {
		switch opts.ResponseFormat.Type {
		case "json_schema", "json_object":
			config.ResponseMIMEType = "application/json"

			if len(opts.ResponseFormat.JSONSchema) > 0 {
				var schemaMap map[string]any

				err := json.Unmarshal(opts.ResponseFormat.JSONSchema, &schemaMap)
				if err == nil {
					config.ResponseSchema = &genai.Schema{}
					config.ResponseJsonSchema = schemaMap
				}
			}
		case "text":
			config.ResponseMIMEType = "text/plain"
		}
	}

	// Thinking budget (reasoning).
	if opts.ThinkingBudget != nil {
		config.ThinkingConfig = &genai.ThinkingConfig{
			ThinkingBudget: genai.Ptr(int32(*opts.ThinkingBudget)), //nolint:gosec
		}
	}

	return config
}

// mapGenAIResponse converts a genai GenerateContentResponse to the unified GenerateTextResult.
func mapGenAIResponse(resp *genai.GenerateContentResponse) (*GenerateTextResult, error) {
	if resp == nil {
		return nil, ErrGenAINilResponse
	}

	result := &GenerateTextResult{
		RawResponse: resp,
	}

	// Map usage metadata.
	if resp.UsageMetadata != nil {
		result.Usage = Usage{
			InputTokens:  int(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:  int(resp.UsageMetadata.TotalTokenCount),
		}

		if resp.UsageMetadata.ThoughtsTokenCount > 0 {
			result.Usage.ThinkingTokens = int(resp.UsageMetadata.ThoughtsTokenCount)
		}
	}

	// Extract content from first candidate.
	if len(resp.Candidates) == 0 {
		return result, nil
	}

	candidate := resp.Candidates[0]

	// Map finish reason.
	result.StopReason = mapGenAIFinishReason(candidate.FinishReason)

	// Map content parts.
	if candidate.Content != nil {
		result.Content = mapGenAIPartsToContentBlocks(candidate.Content.Parts)
	}

	return result, nil
}

// mapGenAIFinishReason converts a genai FinishReason to unified StopReason.
func mapGenAIFinishReason(reason genai.FinishReason) StopReason {
	switch reason {
	case genai.FinishReasonStop:
		return StopReasonEndTurn
	case genai.FinishReasonMaxTokens:
		return StopReasonMaxTokens
	default:
		return StopReasonStop
	}
}

// mapGenAIPartsToContentBlocks converts genai Parts to unified ContentBlock slices.
func mapGenAIPartsToContentBlocks(parts []*genai.Part) []ContentBlock {
	blocks := make([]ContentBlock, 0, len(parts))

	for _, part := range parts {
		if part == nil {
			continue
		}

		if part.Text != "" {
			blocks = append(blocks, ContentBlock{
				Type: ContentBlockText,
				Text: part.Text,
			})
		}

		if part.FunctionCall != nil {
			argsJSON, err := json.Marshal(part.FunctionCall.Args)
			if err != nil {
				argsJSON = []byte("{}")
			}

			blocks = append(blocks, ContentBlock{
				Type: ContentBlockToolCall,
				ToolCall: &ToolCall{
					ID:        "call_" + part.FunctionCall.Name,
					Name:      part.FunctionCall.Name,
					Arguments: argsJSON,
				},
			})
		}

		if part.InlineData != nil {
			blocks = append(blocks, ContentBlock{
				Type: ContentBlockImage,
				Image: &ImagePart{
					MIMEType: part.InlineData.MIMEType,
					Data:     part.InlineData.Data,
				},
			})
		}
	}

	return blocks
}

// emitStreamEventsFromResponse extracts stream events from a genai response chunk
// and sends them to the event channel. Shared by Gemini and Vertex AI streaming.
func emitStreamEventsFromResponse(
	ctx context.Context,
	eventCh chan<- StreamEvent,
	resp *genai.GenerateContentResponse,
) {
	if resp == nil {
		return
	}

	if len(resp.Candidates) == 0 {
		return
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil {
		return
	}

	for _, part := range candidate.Content.Parts {
		if part == nil {
			continue
		}

		if part.Text != "" {
			sendStreamEvent(ctx, eventCh, StreamEvent{
				Type:      StreamEventContentDelta,
				TextDelta: part.Text,
			})
		}

		if part.FunctionCall != nil {
			argsJSON, err := json.Marshal(part.FunctionCall.Args)
			if err != nil {
				argsJSON = []byte("{}")
			}

			sendStreamEvent(ctx, eventCh, StreamEvent{
				Type: StreamEventToolCallDelta,
				ToolCall: &ToolCall{
					ID:        "call_" + part.FunctionCall.Name,
					Name:      part.FunctionCall.Name,
					Arguments: argsJSON,
				},
			})
		}
	}

	// If usage metadata is present in this chunk, attach it.
	if resp.UsageMetadata != nil && candidate.FinishReason != "" {
		usage := &Usage{
			InputTokens:  int(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:  int(resp.UsageMetadata.TotalTokenCount),
		}

		if resp.UsageMetadata.ThoughtsTokenCount > 0 {
			usage.ThinkingTokens = int(resp.UsageMetadata.ThoughtsTokenCount)
		}

		sendStreamEvent(ctx, eventCh, StreamEvent{
			Type:       StreamEventMessageDone,
			StopReason: mapGenAIFinishReason(candidate.FinishReason),
			Usage:      usage,
		})
	}
}

// sendStreamEvent sends an event to the channel, respecting context cancellation.
func sendStreamEvent(ctx context.Context, eventCh chan<- StreamEvent, event StreamEvent) {
	select {
	case eventCh <- event:
	case <-ctx.Done():
	}
}
