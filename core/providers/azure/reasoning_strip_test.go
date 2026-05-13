package azure

import (
	"testing"

	"github.com/maximhq/bifrost/core/providers/openai"
	"github.com/maximhq/bifrost/core/schemas"
)

// TestStripUnsupportedReasoningChat_NonReasoningModel verifies that Azure clears
// reasoning.effort for non-reasoning deployments (e.g. kimi-k2.6) that route through
// the OpenAI Chat Completions API.
func TestStripUnsupportedReasoningChat_NonReasoningModel(t *testing.T) {
	req := &openai.OpenAIChatRequest{
		Model: "kimi-k2.6",
		ChatParameters: schemas.ChatParameters{
			Reasoning: &schemas.ChatReasoning{Effort: schemas.Ptr("high")},
		},
	}
	stripUnsupportedReasoningChat(req, req.Model)
	if req.ChatParameters.Reasoning != nil {
		t.Fatalf("expected Reasoning to be cleared for non-reasoning Azure model, got %+v", req.ChatParameters.Reasoning)
	}
}

// TestStripUnsupportedReasoningChat_ReasoningModel verifies that reasoning is preserved
// for Azure deployments backed by an OpenAI reasoning model (e.g. gpt-5 or o3).
func TestStripUnsupportedReasoningChat_ReasoningModel(t *testing.T) {
	for _, model := range []string{"o3-mini", "gpt-5.1", "gpt-oss-7b"} {
		req := &openai.OpenAIChatRequest{
			Model: model,
			ChatParameters: schemas.ChatParameters{
				Reasoning: &schemas.ChatReasoning{Effort: schemas.Ptr("high")},
			},
		}
		stripUnsupportedReasoningChat(req, req.Model)
		if req.ChatParameters.Reasoning == nil {
			t.Fatalf("expected Reasoning to be preserved for reasoning model %q", model)
		}
	}
}

// TestStripUnsupportedReasoningResponses_NonReasoningModel mirrors the chat test for
// the Responses API path.
func TestStripUnsupportedReasoningResponses_NonReasoningModel(t *testing.T) {
	req := &openai.OpenAIResponsesRequest{
		Model: "kimi-k2.6",
		ResponsesParameters: schemas.ResponsesParameters{
			Reasoning: &schemas.ResponsesParametersReasoning{Effort: schemas.Ptr("high")},
		},
	}
	stripUnsupportedReasoningResponses(req, req.Model)
	if req.ResponsesParameters.Reasoning != nil {
		t.Fatalf("expected Reasoning to be cleared for non-reasoning Azure model, got %+v", req.ResponsesParameters.Reasoning)
	}
}

func TestStripUnsupportedReasoningResponses_ReasoningModel(t *testing.T) {
	req := &openai.OpenAIResponsesRequest{
		Model: "o3-mini",
		ResponsesParameters: schemas.ResponsesParameters{
			Reasoning: &schemas.ResponsesParametersReasoning{Effort: schemas.Ptr("high")},
		},
	}
	stripUnsupportedReasoningResponses(req, req.Model)
	if req.ResponsesParameters.Reasoning == nil {
		t.Fatalf("expected Reasoning to be preserved for o3-mini")
	}
}

// TestStripUnsupportedReasoning_NilSafe ensures the helpers tolerate a nil request,
// which is what ToOpenAI*Request returns for an invalid input.
func TestStripUnsupportedReasoning_NilSafe(t *testing.T) {
	stripUnsupportedReasoningChat(nil, "kimi-k2.6")
	stripUnsupportedReasoningResponses(nil, "kimi-k2.6")
}

// TestIsAzureChatCompletionsOnlyModel asserts which Azure deployment names are
// classified as "Responses API not supported" so the provider falls back to
// Chat Completions.
func TestIsAzureChatCompletionsOnlyModel(t *testing.T) {
	cases := []struct {
		model string
		want  bool
	}{
		// Third-party deployments that only expose Chat Completions on Azure.
		{"kimi-k2.6", true},
		{"Kimi-K2.6", true},
		{"moonshot-v1", true},
		{"minimax-m1", true},
		{"glm-4.6", true},
		{"deepseek-r1", true},
		{"qwen2.5-coder", true},
		{"meta-llama-3.3", true},
		{"mistral-large", true},
		{"codestral-25.08", true},
		{"phi-4", true},
		{"cohere-command-r", true},
		{"jais-30b", true},
		{"ai21-jamba-1.5", true},
		{"grok-4", true},

		// Native OpenAI deployments — keep the Responses API path.
		{"gpt-4o", false},
		{"gpt-4.1-mini", false},
		{"gpt-5", false},
		{"o3-mini", false},
		{"o4-mini", false},
		{"chatgpt-4o-latest", false},
		{"computer-use-preview", false},
	}
	for _, tc := range cases {
		t.Run(tc.model, func(t *testing.T) {
			got := isAzureChatCompletionsOnlyModel(tc.model)
			if got != tc.want {
				t.Fatalf("isAzureChatCompletionsOnlyModel(%q) = %v, want %v", tc.model, got, tc.want)
			}
		})
	}
}
