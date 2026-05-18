package ui

import (
	"strings"
	"testing"

	"github.com/fselich/dossier/internal/openspec"
)

func TestConfigToMarkdown(t *testing.T) {
	t.Run("context output", func(t *testing.T) {
		cfg := openspec.ProjectConfig{Context: "Tech stack: Go.\nDomain: TUI."}
		out := configToMarkdown(cfg)
		if !strings.Contains(out, "## Context") {
			t.Error("expected '## Context' heading")
		}
		if !strings.Contains(out, "Tech stack: Go.") {
			t.Error("expected context prose in output")
		}
	})

	t.Run("rules grouping", func(t *testing.T) {
		cfg := openspec.ProjectConfig{
			Context: "Some context.",
			Rules: map[string][]string{
				"proposal": {"Keep it concise", "Include non-goals"},
				"tasks":    {"Small steps"},
			},
		}
		out := configToMarkdown(cfg)
		if !strings.Contains(out, "## Rules") {
			t.Error("expected '## Rules' heading")
		}
		if !strings.Contains(out, "### proposal") {
			t.Error("expected '### proposal' heading")
		}
		if !strings.Contains(out, "### tasks") {
			t.Error("expected '### tasks' heading")
		}
		if !strings.Contains(out, "- Keep it concise") {
			t.Error("expected proposal rule as bullet")
		}
		if !strings.Contains(out, "- Small steps") {
			t.Error("expected task rule as bullet")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		out := configToMarkdown(openspec.ProjectConfig{})
		if out != "" {
			t.Errorf("expected empty string for empty config, got %q", out)
		}
	})

	t.Run("rules sorted alphabetically", func(t *testing.T) {
		cfg := openspec.ProjectConfig{
			Rules: map[string][]string{
				"tasks":    {"step"},
				"proposal": {"concise"},
			},
		}
		out := configToMarkdown(cfg)
		posProposal := strings.Index(out, "### proposal")
		posTasks := strings.Index(out, "### tasks")
		if posProposal > posTasks {
			t.Error("expected 'proposal' to appear before 'tasks' (alphabetical order)")
		}
	})
}
