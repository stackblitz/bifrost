package azure

import (
	"strings"

	"github.com/maximhq/bifrost/core/providers/anthropic"
	"github.com/maximhq/bifrost/core/schemas"
)

// azureChatCompletionsOnlyModelPrefixes lists deployment-name fragments for
// Azure-hosted models that only expose the Chat Completions API. A Responses
// API call against these deployments returns 404/"not supported", so Bifrost
// must internally fall back to Chat Completions. Matching is case-insensitive
// substring against the deployment name.
var azureChatCompletionsOnlyModelPrefixes = []string{
	"kimi",
	"moonshot",
	"minimax",
	"glm",
	"deepseek",
	"qwen",
	"llama",
	"mistral",
	"codestral",
	"cohere",
	"command",
	"jais",
	"phi",
	"ai21",
	"jamba",
	"grok",
}

// isAzureChatCompletionsOnlyModel returns true when the Azure deployment name
// indicates a third-party / open-source model that doesn't expose the OpenAI
// Responses API. The check is intentionally a name-based heuristic — Azure
// deployment names are user-chosen, so callers using non-standard names will
// stay on the default (Responses API) path.
func isAzureChatCompletionsOnlyModel(model string) bool {
	modelLower := strings.ToLower(model)
	for _, p := range azureChatCompletionsOnlyModelPrefixes {
		if strings.Contains(modelLower, p) {
			return true
		}
	}
	return false
}

// getRequestBodyForAnthropicResponses serializes a BifrostResponsesRequest into the Anthropic wire format for Azure.
// It delegates to BuildAnthropicResponsesRequestBody with the Azure provider and the target deployment name.
func getRequestBodyForAnthropicResponses(ctx *schemas.BifrostContext, request *schemas.BifrostResponsesRequest, deployment string, isStreaming bool, shouldSendBackRawRequest bool, shouldSendBackRawResponse bool) ([]byte, *schemas.BifrostError) {
	return anthropic.BuildAnthropicResponsesRequestBody(ctx, request, anthropic.AnthropicRequestBuildConfig{
		Provider:                  schemas.Azure,
		Deployment:                deployment,
		IsStreaming:               isStreaming,
		ValidateTools:             true,
		ShouldSendBackRawRequest:  shouldSendBackRawRequest,
		ShouldSendBackRawResponse: shouldSendBackRawResponse,
	})
}

// getAzureScopes returns the configured scopes or the default scope if none are valid.
// It filters out empty/whitespace-only strings.
func getAzureScopes(configuredScopes []string) []string {
	scopes := []string{DefaultAzureScope}
	if len(configuredScopes) > 0 {
		cleaned := make([]string, 0, len(configuredScopes))
		for _, s := range configuredScopes {
			if strings.TrimSpace(s) != "" {
				cleaned = append(cleaned, strings.TrimSpace(s))
			}
		}
		if len(cleaned) > 0 {
			scopes = cleaned
		}
	}
	return scopes
}
