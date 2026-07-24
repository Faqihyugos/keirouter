package connectors

import (
	"fmt"
	"testing"

	"github.com/mydisha/keirouter/backend/internal/core"
	"github.com/stretchr/testify/require"
)

func TestOpenAIResponsesGrokCLIHeaders(t *testing.T) {
	c := NewOpenAIResponses("grok-cli", "https://cli-chat-proxy.grok.com/v1/responses")

	headers := c.headers(core.Credentials{AccessToken: "tok_test"})

	require.Equal(t, "Bearer tok_test", headers["Authorization"])
	require.Equal(t, grokCLIVersion, headers["x-grok-client-version"])
	require.Equal(t, grokCLIClientIdentifier, headers["x-grok-client-identifier"])
	require.Equal(t, grokCLITokenAuth, headers["X-XAI-Token-Auth"])
	require.Equal(t, grokCLIAuthenticateResp, headers["x-authenticateresponse"])
	require.Equal(t, grokCLIClientMode, headers["x-grok-client-mode"])
	require.Equal(t, fmt.Sprintf("grok-shell/%s (linux; x86_64)", grokCLIVersion), headers["User-Agent"])
}

func TestOpenAIResponsesGrokCLIHeadersAllowOverrides(t *testing.T) {
	c := NewOpenAIResponses("grok-cli", "https://cli-chat-proxy.grok.com/v1/responses")

	headers := c.headers(core.Credentials{
		AccessToken: "tok_test",
		Headers: map[string]string{
			"x-grok-client-version": "9.9.9",
			"User-Agent":            "custom-ua",
		},
	})

	require.Equal(t, "Bearer tok_test", headers["Authorization"])
	require.Equal(t, "9.9.9", headers["x-grok-client-version"])
	require.Equal(t, "custom-ua", headers["User-Agent"])
	require.Equal(t, grokCLIClientIdentifier, headers["x-grok-client-identifier"])
}

func TestOpenAIResponsesCodexHeadersUnaffected(t *testing.T) {
	c := NewOpenAIResponses("codex", "https://chatgpt.com/backend-api/codex/responses")

	headers := c.headers(core.Credentials{AccessToken: "tok_test"})

	require.Equal(t, "Bearer tok_test", headers["Authorization"])
	require.Empty(t, headers["x-grok-client-version"])
	require.Empty(t, headers["x-grok-client-identifier"])
	require.Empty(t, headers["X-XAI-Token-Auth"])
	require.Empty(t, headers["x-authenticateresponse"])
	require.Empty(t, headers["x-grok-client-mode"])
	require.Empty(t, headers["User-Agent"])
}
