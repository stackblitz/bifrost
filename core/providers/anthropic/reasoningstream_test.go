package anthropic

import (
	"context"
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
)

func TestToAnthropicResponsesStreamResponse_ReasoningThenTextBlocks(t *testing.T) {
	t.Parallel()

	ctx, cancel := schemas.NewBifrostContextWithCancel(context.Background())
	defer cancel()

	state := schemas.AcquireChatToResponsesStreamState()
	defer schemas.ReleaseChatToResponsesStreamState(state)

	makeChunk := func(role *string, reasoning *string, content *string, finishReason *string) *schemas.BifrostChatResponse {
		return &schemas.BifrostChatResponse{
			ID:    "chatcmpl-reasoning",
			Model: "test-model",
			Choices: []schemas.BifrostResponseChoice{
				{
					FinishReason: finishReason,
					ChatStreamResponseChoice: &schemas.ChatStreamResponseChoice{
						Delta: &schemas.ChatStreamResponseChoiceDelta{
							Role:      role,
							Reasoning: reasoning,
							Content:   content,
						},
					},
				},
			},
		}
	}

	role := string(schemas.ChatMessageRoleAssistant)
	reasoning1 := "The"
	reasoning2 := " answer"
	text1 := "Why"
	text2 := " it works."
	stop := string(schemas.BifrostFinishReasonStop)

	chunks := []*schemas.BifrostChatResponse{
		makeChunk(&role, nil, nil, nil),
		makeChunk(nil, &reasoning1, nil, nil),
		makeChunk(nil, &reasoning2, nil, nil),
		makeChunk(nil, nil, &text1, nil),
		makeChunk(nil, nil, &text2, nil),
		makeChunk(nil, nil, nil, &stop),
	}

	var contentEvents []*AnthropicStreamEvent
	for _, chunk := range chunks {
		for _, bifrostResp := range chunk.ToBifrostResponsesStreamResponse(state) {
			for _, event := range ToAnthropicResponsesStreamResponse(ctx, bifrostResp) {
				switch event.Type {
				case AnthropicStreamEventTypeContentBlockStart,
					AnthropicStreamEventTypeContentBlockDelta,
					AnthropicStreamEventTypeContentBlockStop:
					contentEvents = append(contentEvents, event)
				}
			}
		}
	}

	if len(contentEvents) != 8 {
		t.Fatalf("expected 8 content block events, got %d", len(contentEvents))
	}

	assertIndex := func(event *AnthropicStreamEvent, want int) {
		t.Helper()
		if event.Index == nil || *event.Index != want {
			t.Fatalf("event index = %v, want %d", event.Index, want)
		}
	}

	if contentEvents[0].Type != AnthropicStreamEventTypeContentBlockStart {
		t.Fatalf("event[0] type = %v, want content_block_start", contentEvents[0].Type)
	}
	assertIndex(contentEvents[0], 0)
	if contentEvents[0].ContentBlock == nil || contentEvents[0].ContentBlock.Type != AnthropicContentBlockTypeThinking {
		t.Fatalf("event[0] content block = %#v, want thinking", contentEvents[0].ContentBlock)
	}

	if contentEvents[1].Type != AnthropicStreamEventTypeContentBlockDelta {
		t.Fatalf("event[1] type = %v, want content_block_delta", contentEvents[1].Type)
	}
	assertIndex(contentEvents[1], 0)
	if contentEvents[1].Delta == nil || contentEvents[1].Delta.Type != AnthropicStreamDeltaTypeThinking {
		t.Fatalf("event[1] delta = %#v, want thinking_delta", contentEvents[1].Delta)
	}

	if contentEvents[2].Type != AnthropicStreamEventTypeContentBlockDelta {
		t.Fatalf("event[2] type = %v, want content_block_delta", contentEvents[2].Type)
	}
	assertIndex(contentEvents[2], 0)
	if contentEvents[2].Delta == nil || contentEvents[2].Delta.Type != AnthropicStreamDeltaTypeThinking {
		t.Fatalf("event[2] delta = %#v, want thinking_delta", contentEvents[2].Delta)
	}

	if contentEvents[3].Type != AnthropicStreamEventTypeContentBlockStop {
		t.Fatalf("event[3] type = %v, want content_block_stop", contentEvents[3].Type)
	}
	assertIndex(contentEvents[3], 0)

	if contentEvents[4].Type != AnthropicStreamEventTypeContentBlockStart {
		t.Fatalf("event[4] type = %v, want content_block_start", contentEvents[4].Type)
	}
	assertIndex(contentEvents[4], 1)
	if contentEvents[4].ContentBlock == nil || contentEvents[4].ContentBlock.Type != AnthropicContentBlockTypeText {
		t.Fatalf("event[4] content block = %#v, want text", contentEvents[4].ContentBlock)
	}

	if contentEvents[5].Type != AnthropicStreamEventTypeContentBlockDelta {
		t.Fatalf("event[5] type = %v, want content_block_delta", contentEvents[5].Type)
	}
	assertIndex(contentEvents[5], 1)
	if contentEvents[5].Delta == nil || contentEvents[5].Delta.Type != AnthropicStreamDeltaTypeText {
		t.Fatalf("event[5] delta = %#v, want text_delta", contentEvents[5].Delta)
	}

	if contentEvents[6].Type != AnthropicStreamEventTypeContentBlockDelta {
		t.Fatalf("event[6] type = %v, want content_block_delta", contentEvents[6].Type)
	}
	assertIndex(contentEvents[6], 1)
	if contentEvents[6].Delta == nil || contentEvents[6].Delta.Type != AnthropicStreamDeltaTypeText {
		t.Fatalf("event[6] delta = %#v, want text_delta", contentEvents[6].Delta)
	}

	if contentEvents[7].Type != AnthropicStreamEventTypeContentBlockStop {
		t.Fatalf("event[7] type = %v, want content_block_stop", contentEvents[7].Type)
	}
	assertIndex(contentEvents[7], 1)
}
