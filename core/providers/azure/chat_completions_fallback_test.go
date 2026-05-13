package azure

import "testing"

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
