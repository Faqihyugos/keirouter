package connectors

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mydisha/keirouter/backend/internal/core"
)

func TestResolveCodexModel(t *testing.T) {
	cases := []struct {
		in     string
		model  string
		effort string
	}{
		// The reported failure: dash typo of a dotted version, with no exact
		// upstream id — alias lands on the 5.6 flagship.
		{"gpt-5-6", "gpt-5.6-sol", ""},
		{"gpt-5.6", "gpt-5.6-sol", ""},
		{"gpt-5-6-high", "gpt-5.6-sol", "high"},
		// Known catalog ids pass through untouched.
		{"gpt-5.6-terra", "gpt-5.6-terra", ""},
		{"gpt-5-codex", "gpt-5-codex", ""},
		{"gpt-5.3-codex-spark", "gpt-5.3-codex-spark", ""},
		// Effort suffixes are stripped and returned separately.
		{"gpt-5.3-codex-xhigh", "gpt-5.3-codex", "xhigh"},
		{"gpt-5.3-codex-high", "gpt-5.3-codex", "high"},
		{"gpt-5.1-codex-mini-high", "gpt-5.1-codex-mini", "high"},
		{"gpt-5.5-none", "gpt-5.5", "none"},
		// gpt-5.6 family -max/-ultra alias suffixes (dash typo tolerated); a
		// bare "-max" elsewhere stays part of the model id (gpt-5.1-codex-max).
		{"gpt-5.6-sol-ultra", "gpt-5.6-sol", "ultra"},
		{"gpt-5-6-terra-max", "gpt-5.6-terra", "max"},
		{"gpt-5.1-codex-max", "gpt-5.1-codex-max", ""},
		// Unknown ids pass through (normalized) for upstream to decide.
		{"gpt-5.9-nova", "gpt-5.9-nova", ""},
	}
	for _, tc := range cases {
		model, effort := resolveCodexModel(tc.in)
		require.Equal(t, tc.model, model, "model for %q", tc.in)
		require.Equal(t, tc.effort, effort, "effort for %q", tc.in)
	}
}

func TestPrepareCodexRequest(t *testing.T) {
	// Effort suffix overrides a client-injected reasoning default.
	req := &core.ChatRequest{Model: "gpt-5.3-codex-xhigh", Reasoning: &core.ReasoningConfig{Effort: "medium"}}
	out := prepareCodexRequest(req)
	require.Equal(t, "gpt-5.3-codex", out.Model)
	require.Equal(t, "xhigh", out.Reasoning.Effort)
	// The caller's request is untouched (fallback to another provider must
	// dispatch with the original model id).
	require.Equal(t, "gpt-5.3-codex-xhigh", req.Model)
	require.Equal(t, "medium", req.Reasoning.Effort)

	// No suffix + no client reasoning -> default medium effort.
	out = prepareCodexRequest(&core.ChatRequest{Model: "gpt-5.5"})
	require.Equal(t, "gpt-5.5", out.Model)
	require.Equal(t, "medium", out.Reasoning.Effort)

	// Client-provided effort is preserved when the model has no suffix.
	out = prepareCodexRequest(&core.ChatRequest{Model: "gpt-5.5", Reasoning: &core.ReasoningConfig{Effort: "low"}})
	require.Equal(t, "low", out.Reasoning.Effort)

	// Claude Code's adaptive thinking (the reported 400) is normalized to the
	// CLI default instead of leaking an invalid wire value upstream.
	out = prepareCodexRequest(&core.ChatRequest{Model: "gpt-5-6", Reasoning: &core.ReasoningConfig{Effort: "adaptive"}})
	require.Equal(t, "gpt-5.6-sol", out.Model)
	require.Equal(t, "medium", out.Reasoning.Effort)

	// Ultra alias suffix resolves to the gpt-5.6 wire top value "max".
	out = prepareCodexRequest(&core.ChatRequest{Model: "gpt-5.6-sol-ultra"})
	require.Equal(t, "gpt-5.6-sol", out.Model)
	require.Equal(t, "max", out.Reasoning.Effort)
}

func TestNormalizeCodexEffort(t *testing.T) {
	cases := []struct {
		model, effort, want string
	}{
		// gpt-5.6 family: enum is none/minimal/low/medium/high/max (no xhigh).
		{"gpt-5.6-sol", "adaptive", "medium"},
		{"gpt-5.6-sol", "xhigh", "max"},
		{"gpt-5.6-terra", "ultra", "max"},
		{"gpt-5.6-luna", "max", "max"},
		{"gpt-5.6-sol", "minimal", "minimal"},
		{"gpt-5.6-sol", "off", "none"},
		{"gpt-5.6-sol", "bogus", "medium"},
		// Legacy codex models: enum tops out at xhigh (no max).
		{"gpt-5.3-codex", "max", "xhigh"},
		{"gpt-5.3-codex", "xhigh", "xhigh"},
		{"gpt-5.5", "adaptive", "medium"},
		{"gpt-5.5", "auto", "medium"},
		{"gpt-5.5", "", "medium"},
		{"gpt-5.5", "HIGH", "high"},
		{"gpt-5-codex", "bogus", "medium"},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, normalizeCodexEffort(tc.model, tc.effort), "model=%s effort=%s", tc.model, tc.effort)
	}
}

func TestApplyCodexHeaders(t *testing.T) {
	h := map[string]string{}
	applyCodexHeaders(h, core.Credentials{
		AccountID: "acc-1",
		Extra:     map[string]string{"chatgpt_account_id": "wk-123"},
	})
	require.Equal(t, "codex_cli_rs", h["originator"])
	require.Equal(t, "responses=experimental", h["OpenAI-Beta"])
	require.Equal(t, codexClientVersion, h["version"])
	require.Equal(t, "wk-123", h["chatgpt-account-id"])
	require.Equal(t, "acc-1", h["session_id"])

	// Imported accounts may carry the CLI-style key instead.
	h = map[string]string{}
	applyCodexHeaders(h, core.Credentials{Extra: map[string]string{"chatgptAccountId": "wk-456"}})
	require.Equal(t, "wk-456", h["chatgpt-account-id"])
}
