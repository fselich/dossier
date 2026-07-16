package ui

import (
	"testing"

	"charm.land/glamour/v2"
)

func TestGlamourStyleName(t *testing.T) {
	t.Run("defaults to dark when empty", func(t *testing.T) {
		t.Setenv(glamourStyleEnv, "")

		got := glamourStyleName()

		if got != "dark" {
			t.Errorf("expected dark, got %q", got)
		}
	})

	t.Run("uses environment override", func(t *testing.T) {
		t.Setenv(glamourStyleEnv, "light")

		got := glamourStyleName()

		if got != "light" {
			t.Errorf("expected light, got %q", got)
		}
	})
}

func TestEnsureRendererUsesConfiguredGlamourStyle(t *testing.T) {
	t.Setenv(glamourStyleEnv, "light")
	raw := "# Title\n\nBody"

	m := &Model{}
	m.ensureRenderer(80)
	if m.glamourRenderer == nil {
		t.Fatal("expected renderer to be initialized")
	}
	got, err := m.glamourRenderer.Render(raw)
	if err != nil {
		t.Fatalf("render with configured style: %v", err)
	}

	wantRenderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("light"),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("create expected light renderer: %v", err)
	}
	want, err := wantRenderer.Render(raw)
	if err != nil {
		t.Fatalf("render with expected light style: %v", err)
	}

	if got != want {
		t.Error("expected renderer output to match light style")
	}
}
