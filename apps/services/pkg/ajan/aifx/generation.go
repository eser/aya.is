package aifx

import (
	"context"
	"strings"
	"sync"
)

// ToolChoice controls how the model selects tools.
type ToolChoice string

const (
	ToolChoiceAuto     ToolChoice = "auto"
	ToolChoiceNone     ToolChoice = "none"
	ToolChoiceRequired ToolChoice = "required"
)

// StopReason indicates why generation stopped.
type StopReason string

const (
	StopReasonEndTurn   StopReason = "end_turn"
	StopReasonMaxTokens StopReason = "max_tokens"
	StopReasonToolUse   StopReason = "tool_use"
	StopReasonStop      StopReason = "stop"
)

// ProviderCapability represents a capability supported by a provider.
type ProviderCapability string

const (
	CapabilityTextGeneration  ProviderCapability = "text_generation"
	CapabilityStreaming       ProviderCapability = "streaming"
	CapabilityToolCalling     ProviderCapability = "tool_calling"
	CapabilityVision          ProviderCapability = "vision"
	CapabilityAudio           ProviderCapability = "audio"
	CapabilityBatchProcessing ProviderCapability = "batch_processing"
	CapabilityStructuredOut   ProviderCapability = "structured_output"
	CapabilityReasoning       ProviderCapability = "reasoning"
)

// ResponseFormat enables structured outputs with JSON schemas.
type ResponseFormat struct {
	Type       string // "json_schema", "json_object", "text"
	Name       string // schema name (for OpenAI structured output)
	JSONSchema []byte // JSON Schema definition
}

// SafetySetting configures content safety thresholds (Google providers).
type SafetySetting struct {
	Category  string // e.g. "HARM_CATEGORY_HARASSMENT"
	Threshold string // e.g. "BLOCK_NONE", "BLOCK_LOW_AND_ABOVE"
}

// Usage contains token usage statistics for a generation.
type Usage struct {
	InputTokens    int
	OutputTokens   int
	TotalTokens    int
	ThinkingTokens int // reasoning model tokens (separate from output)
}

// GenerateTextOptions configures a text generation request.
type GenerateTextOptions struct {
	Messages       []Message
	Tools          []ToolDefinition
	System         string
	MaxTokens      int
	Temperature    *float64
	TopP           *float64
	StopWords      []string
	ToolChoice     ToolChoice      // "auto", "none", "required"
	ResponseFormat *ResponseFormat // structured output (JSON schema)
	ThinkingBudget *int            // reasoning/thinking token budget
	SafetySettings []SafetySetting // Google providers: harm category thresholds
	Extensions     map[string]any  // provider-specific extra options
}

// StreamTextOptions configures a streaming text generation request.
// Same as GenerateTextOptions but used for StreamText to clarify intent.
type StreamTextOptions = GenerateTextOptions

// GenerateTextResult contains the result of a text generation request.
type GenerateTextResult struct {
	Content     []ContentBlock
	StopReason  StopReason
	Usage       Usage
	ModelID     string
	RawRequest  any // provider-specific raw request (for debugging)
	RawResponse any // provider-specific raw response (for debugging)
}

// Text returns concatenated text from all text content blocks.
func (r *GenerateTextResult) Text() string {
	var builder strings.Builder

	for _, block := range r.Content {
		if block.Type == ContentBlockText {
			builder.WriteString(block.Text)
		}
	}

	return builder.String()
}

// ToolCalls extracts all tool calls from the result.
func (r *GenerateTextResult) ToolCalls() []ToolCall {
	var calls []ToolCall

	for _, block := range r.Content {
		if block.Type == ContentBlockToolCall && block.ToolCall != nil {
			calls = append(calls, *block.ToolCall)
		}
	}

	return calls
}

// StreamEventType identifies the type of stream event.
type StreamEventType string

const (
	StreamEventContentDelta  StreamEventType = "content_delta"
	StreamEventToolCallDelta StreamEventType = "tool_call_delta"
	StreamEventMessageDone   StreamEventType = "message_done"
	StreamEventError         StreamEventType = "error"
)

// StreamEvent represents a single event in a streaming response.
type StreamEvent struct {
	Type       StreamEventType
	TextDelta  string
	ToolCall   *ToolCall
	StopReason StopReason
	Usage      *Usage
	Error      error
}

// StreamIterator provides a pull-based API for consuming stream events.
// Follows the Go iterator pattern (similar to sql.Rows).
type StreamIterator struct {
	eventCh <-chan StreamEvent
	cancel  context.CancelFunc
	current StreamEvent
	err     error
	done    bool
	mu      sync.Mutex
}

// NewStreamIterator creates a new stream iterator from an event channel.
func NewStreamIterator(eventCh <-chan StreamEvent, cancel context.CancelFunc) *StreamIterator {
	return &StreamIterator{
		eventCh: eventCh,
		cancel:  cancel,
	}
}

// Next advances the iterator to the next event.
// Returns true if there is a new event, false when the stream is done or errored.
func (iter *StreamIterator) Next() bool {
	iter.mu.Lock()
	defer iter.mu.Unlock()

	if iter.done {
		return false
	}

	event, ok := <-iter.eventCh
	if !ok {
		iter.done = true

		return false
	}

	iter.current = event

	if event.Type == StreamEventError {
		iter.err = event.Error
		iter.done = true

		return false
	}

	if event.Type == StreamEventMessageDone {
		iter.done = true
	}

	return true
}

// Current returns the most recently read stream event.
func (iter *StreamIterator) Current() StreamEvent {
	iter.mu.Lock()
	defer iter.mu.Unlock()

	return iter.current
}

// Err returns the first error encountered during streaming.
func (iter *StreamIterator) Err() error {
	iter.mu.Lock()
	defer iter.mu.Unlock()

	return iter.err
}

// Close cancels the stream and releases resources.
func (iter *StreamIterator) Close() error {
	iter.mu.Lock()
	defer iter.mu.Unlock()

	if iter.cancel != nil {
		iter.cancel()
	}

	iter.done = true

	return nil
}

// Collect drains the entire stream and returns a combined result.
func (iter *StreamIterator) Collect() (*GenerateTextResult, error) {
	var textBuilder strings.Builder

	var (
		toolCalls  []ContentBlock
		stopReason StopReason
		usage      Usage
	)

	for iter.Next() {
		event := iter.Current()

		switch event.Type {
		case StreamEventContentDelta:
			textBuilder.WriteString(event.TextDelta)
		case StreamEventToolCallDelta:
			if event.ToolCall != nil {
				toolCalls = append(toolCalls, ContentBlock{
					Type:     ContentBlockToolCall,
					ToolCall: event.ToolCall,
				})
			}
		case StreamEventMessageDone:
			stopReason = event.StopReason

			if event.Usage != nil {
				usage = *event.Usage
			}
		case StreamEventError:
			// handled by iter.Err()
		}
	}

	err := iter.Err()
	if err != nil {
		return nil, err
	}

	content := make([]ContentBlock, 0, 1+len(toolCalls))

	if textBuilder.Len() > 0 {
		content = append(content, ContentBlock{
			Type: ContentBlockText,
			Text: textBuilder.String(),
		})
	}

	content = append(content, toolCalls...)

	return &GenerateTextResult{
		Content:    content,
		StopReason: stopReason,
		Usage:      usage,
	}, nil
}
