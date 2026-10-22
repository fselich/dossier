package ui

import (
	"path/filepath"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/aemonge/dossier/internal/settings"
	"github.com/alecthomas/chroma/v2"
)

func TestBuiltinThemesExist(t *testing.T) {
	names := []string{"dark", "light", "dracula", "gruvbox-light-soft"}
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
	cases := []string{"DARK", "Dark", "LIGHT", "Light", "DRACULA", "Dracula", "GRUVBOX-LIGHT-SOFT"}
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

func TestDefaultThemeIsNone(t *testing.T) {
	dt := DefaultTheme()
	if dt.Name != "none" {
		t.Errorf("DefaultTheme().Name = %q, want %q", dt.Name, "none")
	}
	if dt.ViewBg != nil {
		t.Error("DefaultTheme should have nil ViewBg (respect terminal default)")
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

func TestBuildStylesProducesStyles(t *testing.T) {
	styles := BuildStyles(DarkColors)
	if styles.BaseText.GetForeground() == nil {
		t.Error("BaseText should have foreground set")
	}
	if styles.Header.GetForeground() == nil {
		t.Error("Header should have foreground set")
	}
	if styles.Section.GetForeground() == nil {
		t.Error("Section should have foreground set")
	}
	if styles.Help.GetForeground() == nil {
		t.Error("Help should have foreground set")
	}
}

func TestBuildStylesUsesCorrectColors(t *testing.T) {
	styles := BuildStyles(DarkColors)
	if styles.BaseText.GetForeground() != DarkColors.PrimaryFg {
		t.Error("BaseText should use PrimaryFg")
	}
	if styles.Header.GetForeground() != DarkColors.AccentBlue {
		t.Error("Header should use AccentBlue")
	}
	if styles.Section.GetForeground() != DarkColors.AccentYellow {
		t.Error("Section should use AccentYellow")
	}
	if styles.Error.GetForeground() != DarkColors.AccentRed {
		t.Error("Error should use AccentRed")
	}
}

func TestDarkColorsPreserveOriginalValues(t *testing.T) {
	colors := DarkColors
	checkStr := func(name string, got, want any) {
		if got != want {
			t.Errorf("DarkColors.%s: got %v, want %v", name, got, want)
		}
	}
	checkStr("PrimaryFg", colors.PrimaryFg, lipgloss.Color("15"))
	checkStr("MutedFg", colors.MutedFg, lipgloss.Color("8"))
	checkStr("AccentBlue", colors.AccentBlue, lipgloss.Color("12"))
	checkStr("AccentYellow", colors.AccentYellow, lipgloss.Color("11"))
	checkStr("AccentGreen", colors.AccentGreen, lipgloss.Color("2"))
	checkStr("AccentRed", colors.AccentRed, lipgloss.Color("9"))
}

func TestLightColorsAdaptedForWhiteBackground(t *testing.T) {
	colors := LightColors
	checkStr := func(name string, got, want any) {
		if got != want {
			t.Errorf("LightColors.%s: got %v, want %v", name, got, want)
		}
	}
	checkStr("PrimaryFg", colors.PrimaryFg, lipgloss.Color("0"))
	checkStr("AccentYellow", colors.AccentYellow, lipgloss.Color("3"))
	checkStr("AccentRed", colors.AccentRed, lipgloss.Color("1"))
	checkStr("MutedFg", colors.MutedFg, lipgloss.Color("7"))
	checkStr("MidFg", colors.MidFg, lipgloss.Color("8"))
}

func TestNoneThemeInheritsDarkColors(t *testing.T) {
	if NoneTheme.Colors != DarkColors {
		t.Error("NoneTheme should use DarkColors")
	}
}

func TestBuildThemeAppliesUIAndRendererOverrides(t *testing.T) {
	cfg, err := settings.BuiltinTheme("light")
	if err != nil {
		t.Fatal(err)
	}
	cfg.UI.PrimaryFg = "#112233"
	cfg.Glamour.Elements = map[string]settings.ColorOverride{
		"h1": {Foreground: "#334455", Background: "#ffffff"},
	}
	cfg.Glamour.Syntax = map[string]settings.ColorOverride{
		"keyword": {Foreground: "#556677"},
	}
	cfg.Chroma.Tokens = map[string]string{"Keyword": "bold #778899"}

	theme, err := BuildTheme(cfg)
	if err != nil {
		t.Fatalf("BuildTheme: %v", err)
	}
	if theme.Colors.PrimaryFg != lipgloss.Color("#112233") {
		t.Fatalf("PrimaryFg = %v", theme.Colors.PrimaryFg)
	}
	if theme.GlamourStyles.H1.Color == nil || *theme.GlamourStyles.H1.Color != "#334455" {
		t.Fatalf("Glamour H1 color = %v", theme.GlamourStyles.H1.Color)
	}
	if theme.GlamourStyles.CodeBlock.Chroma == nil || theme.GlamourStyles.CodeBlock.Chroma.Keyword.Color == nil || *theme.GlamourStyles.CodeBlock.Chroma.Keyword.Color != "#556677" {
		t.Fatalf("Glamour keyword override not applied")
	}
	entry := theme.Chroma.Get(chroma.Keyword)
	if !entry.Colour.IsSet() || entry.Colour.String() != "#778899" || entry.Bold != chroma.Yes {
		t.Fatalf("Chroma Keyword = %+v", entry)
	}
}

func TestBuildThemeDoesNotMutateBuiltinGlamourStyle(t *testing.T) {
	base, err := settings.BuiltinTheme("light")
	if err != nil {
		t.Fatal(err)
	}
	before, err := BuildTheme(base)
	if err != nil {
		t.Fatal(err)
	}
	if before.GlamourStyles.CodeBlock.Chroma == nil || before.GlamourStyles.CodeBlock.Chroma.Keyword.Color == nil {
		t.Fatal("light style should define a keyword color")
	}
	want := *before.GlamourStyles.CodeBlock.Chroma.Keyword.Color

	custom := base
	custom.Glamour.Syntax = map[string]settings.ColorOverride{"keyword": {Foreground: "#123456"}}
	if _, err := BuildTheme(custom); err != nil {
		t.Fatal(err)
	}

	after, err := BuildTheme(base)
	if err != nil {
		t.Fatal(err)
	}
	if got := *after.GlamourStyles.CodeBlock.Chroma.Keyword.Color; got != want {
		t.Fatalf("built-in Glamour style mutated: got %q, want %q", got, want)
	}
}

func TestBuildThemeFromDocumentedExamples(t *testing.T) {
	tests := []struct {
		name      string
		viewBgNil bool
	}{
		{name: "config.example.toml", viewBgNil: false},
		{name: "config.gruvbox-light-soft.toml", viewBgNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, _, err := settings.Load(filepath.Join("..", "..", "examples", tt.name))
			if err != nil {
				t.Fatal(err)
			}
			theme, err := BuildTheme(cfg.Theme)
			if err != nil {
				t.Fatalf("BuildTheme example: %v", err)
			}
			if (theme.ViewBg == nil) != tt.viewBgNil {
				t.Fatalf("unexpected example background: %v", theme.ViewBg)
			}
		})
	}
}

func TestBuildThemeRejectsUnknownOverrides(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*settings.ThemeConfig)
	}{
		{name: "glamour element", mutate: func(c *settings.ThemeConfig) {
			c.Glamour.Elements = map[string]settings.ColorOverride{"wat": {Foreground: "#fff"}}
		}},
		{name: "glamour syntax", mutate: func(c *settings.ThemeConfig) {
			c.Glamour.Syntax = map[string]settings.ColorOverride{"wat": {Foreground: "#fff"}}
		}},
		{name: "chroma token", mutate: func(c *settings.ThemeConfig) { c.Chroma.Tokens = map[string]string{"Wat": "#fff"} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := settings.BuiltinTheme("light")
			if err != nil {
				t.Fatal(err)
			}
			tt.mutate(&cfg)
			if _, err := BuildTheme(cfg); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestDraculaUsesDarkerYellow(t *testing.T) {
	if DraculaColors.AccentYellow != lipgloss.Color("3") {
		t.Errorf("DraculaColors.AccentYellow = %v, want Color(\"3\")", DraculaColors.AccentYellow)
	}
	if DraculaColors.AccentYellow == DarkColors.AccentYellow {
		t.Error("Dracula AccentYellow should differ from Dark")
	}
}
