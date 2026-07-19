package ui

import "testing"

func TestBuiltinThemesExist(t *testing.T) {
	names := []string{"dark", "light", "dracula"}
	for _, name := range names {
		theme, ok := LookupTheme(name)
		if !ok {
			t.Errorf("expected theme %q to exist", name)
			continue
		}
		if theme.GlamourStyle == "" {
			t.Errorf("theme %q: GlamourStyle should not be empty", name)
		}
		if theme.ChromaStyle == "" {
			t.Errorf("theme %q: ChromaStyle should not be empty", name)
		}
	}
}

func TestLookupThemeCaseInsensitive(t *testing.T) {
	cases := []string{"DARK", "Dark", "LIGHT", "Light", "DRACULA", "Dracula"}
	for _, name := range cases {
		_, ok := LookupTheme(name)
		if !ok {
			t.Errorf("LookupTheme(%q) should be found (case-insensitive)", name)
		}
	}
}

func TestLookupThemeNotFound(t *testing.T) {
	_, ok := LookupTheme("nonexistent")
	if ok {
		t.Error("LookupTheme for nonexistent theme should return false")
	}
}

func TestDefaultThemeIsDark(t *testing.T) {
	dt := DefaultTheme()
	if dt.Name != "dark" {
		t.Errorf("DefaultTheme().Name = %q, want %q", dt.Name, "dark")
	}
}

func TestDarkThemeConfig(t *testing.T) {
	t.Run("glamour", func(t *testing.T) {
		if DarkTheme.GlamourStyle != "dark" {
			t.Errorf("DarkTheme.GlamourStyle = %q, want %q", DarkTheme.GlamourStyle, "dark")
		}
	})
	t.Run("chroma", func(t *testing.T) {
		if DarkTheme.ChromaStyle != "monokai" {
			t.Errorf("DarkTheme.ChromaStyle = %q, want %q", DarkTheme.ChromaStyle, "monokai")
		}
	})
}

func TestLightThemeConfig(t *testing.T) {
	if LightTheme.GlamourStyle != "light" {
		t.Errorf("LightTheme.GlamourStyle = %q, want %q", LightTheme.GlamourStyle, "light")
	}
	if LightTheme.ChromaStyle != "github" {
		t.Errorf("LightTheme.ChromaStyle = %q, want %q", LightTheme.ChromaStyle, "github")
	}
}

func TestDraculaThemeConfig(t *testing.T) {
	if DraculaTheme.GlamourStyle != "dracula" {
		t.Errorf("DraculaTheme.GlamourStyle = %q, want %q", DraculaTheme.GlamourStyle, "dracula")
	}
	if DraculaTheme.ChromaStyle != "dracula" {
		t.Errorf("DraculaTheme.ChromaStyle = %q, want %q", DraculaTheme.ChromaStyle, "dracula")
	}
}
