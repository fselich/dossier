package openspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		yaml := `schema: spec-driven
context: |
  Tech stack: Go.
  Domain: TUI tool.
rules:
  proposal:
    - Keep it concise
  tasks:
    - Small steps
`
		if err := os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte(yaml), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := LoadConfig()
		if cfg.Context == "" {
			t.Error("expected non-empty Context")
		}
		if len(cfg.Rules) != 2 {
			t.Errorf("expected 2 rule groups, got %d", len(cfg.Rules))
		}
		if len(cfg.Rules["proposal"]) != 1 {
			t.Errorf("expected 1 proposal rule, got %d", len(cfg.Rules["proposal"]))
		}
	})

	t.Run("missing file", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		cfg := LoadConfig()
		if cfg.Context != "" {
			t.Error("expected empty Context for missing file")
		}
		if len(cfg.Rules) != 0 {
			t.Error("expected empty Rules for missing file")
		}
	})

	t.Run("empty context", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		yaml := `schema: spec-driven
rules:
  proposal:
    - Keep it concise
`
		if err := os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte(yaml), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := LoadConfig()
		if cfg.Context != "" {
			t.Errorf("expected empty Context, got %q", cfg.Context)
		}
		if len(cfg.Rules["proposal"]) != 1 {
			t.Errorf("expected 1 proposal rule, got %d", len(cfg.Rules["proposal"]))
		}
	})

	t.Run("missing rules", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		yaml := `schema: spec-driven
context: Just context.
`
		if err := os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"), []byte(yaml), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := LoadConfig()
		if cfg.Context != "Just context." {
			t.Errorf("expected context, got %q", cfg.Context)
		}
		if len(cfg.Rules) != 0 {
			t.Errorf("expected empty Rules, got %v", cfg.Rules)
		}
	})
}
