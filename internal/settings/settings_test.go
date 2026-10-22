package settings

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultPathUsesXDGConfigHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	path, err := DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath: %v", err)
	}
	want := filepath.Join(dir, "dossier", "config.toml")
	if path != want {
		t.Fatalf("DefaultPath = %q, want %q", path, want)
	}
}

func TestLoadDefaultsToIndexStartView(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, _, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.UI.StartView != "index" {
		t.Fatalf("default ui.start_view = %q, want index", cfg.UI.StartView)
	}
}

func TestLoadAcceptsChangeStartView(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[ui]\nstart_view = \"change\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, _, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.UI.StartView != "change" {
		t.Fatalf("ui.start_view = %q, want change", cfg.UI.StartView)
	}
}

func TestLoadRejectsUppercaseStartView(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[ui]\nstart_view = \"INDEX\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "ui.start_view") {
		t.Fatalf("expected strict start_view error, got %v", err)
	}
}

func TestLoadRejectsInvalidStartView(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[ui]\nstart_view = \"dashboard\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "ui.start_view") || !strings.Contains(err.Error(), "index") || !strings.Contains(err.Error(), "change") {
		t.Fatalf("expected actionable start_view error, got %v", err)
	}
}

func TestLoadMissingDefaultReturnsDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, path, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Theme.Base != "none" {
		t.Fatalf("default theme base = %q, want none", cfg.Theme.Base)
	}
	if len(cfg.Keys.Viewer.Down) == 0 {
		t.Fatal("expected default viewer bindings")
	}
	if path != filepath.Join(dir, "dossier", "config.toml") {
		t.Fatalf("resolved path = %q", path)
	}
}

func TestLoadMissingExplicitPathReturnsError(t *testing.T) {
	_, _, err := Load(filepath.Join(t.TempDir(), "missing.toml"))
	if err == nil || !strings.Contains(err.Error(), "missing.toml") {
		t.Fatalf("expected missing explicit path error, got %v", err)
	}
}

func TestLoadParsesThemeAndSparseKeyOverrides(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := `
[theme]
base = "light"

[theme.ui]
primary_fg = "#112233"

[theme.glamour]
style = "light"

[theme.glamour.elements.h1]
foreground = "#334455"
background = "#ffffff"

[theme.glamour.syntax.keyword]
foreground = "#556677"

[theme.chroma]
style = "github"

[theme.chroma.tokens]
Keyword = "bold #778899"

[keys.viewer]
down = ["n"]
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, _, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Theme.Base != "light" || cfg.Theme.UI.PrimaryFg != "#112233" {
		t.Fatalf("unexpected theme: %+v", cfg.Theme)
	}
	if got := cfg.Theme.Glamour.Elements["h1"].Foreground; got != "#334455" {
		t.Fatalf("h1 foreground = %q", got)
	}
	if got := cfg.Theme.Chroma.Tokens["Keyword"]; got != "bold #778899" {
		t.Fatalf("Keyword token = %q", got)
	}
	if len(cfg.Keys.Viewer.Down) != 1 || cfg.Keys.Viewer.Down[0] != "n" {
		t.Fatalf("viewer.down = %v", cfg.Keys.Viewer.Down)
	}
	if len(cfg.Keys.Viewer.Up) == 0 {
		t.Fatal("expected unspecified bindings to retain defaults")
	}
}

func TestLoadNamedKeyStyleWithSparseOverrides(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := `
[keys]
style = "nvim"

[keys.viewer]
down = ["n"]
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, _, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(cfg.Keys.Viewer.Down, []string{"n"}) {
		t.Fatalf("viewer.down = %v", cfg.Keys.Viewer.Down)
	}
	if !reflect.DeepEqual(cfg.Keys.Viewer.Quit, []string{"Q"}) {
		t.Fatalf("viewer.quit should inherit nvim base: %v", cfg.Keys.Viewer.Quit)
	}
}

func TestLoadWithKeyStyleOverridesConfiguredBase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := "[keys]\nstyle = \"unknown-in-config\"\n[keys.viewer]\ndown = [\"n\"]\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, _, err := LoadWithKeyStyle(path, "NVIM")
	if err != nil {
		t.Fatalf("LoadWithKeyStyle: %v", err)
	}
	if cfg.Keys.Style != "nvim" || !reflect.DeepEqual(cfg.Keys.Viewer.Down, []string{"n"}) {
		t.Fatalf("CLI base did not preserve config overrides: %+v", cfg.Keys)
	}
}

func TestLoadRejectsUnknownKeyStyle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[keys]\nstyle = \"wat\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), `unknown keystyle "wat"`) {
		t.Fatalf("expected unknown keystyle error, got %v", err)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[theme.ui]\nprimry_fg = \"#fff\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "primry_fg") {
		t.Fatalf("expected unknown key error, got %v", err)
	}
}

func TestBuiltinThemesComeFromEmbeddedTOML(t *testing.T) {
	cfg, err := BuiltinTheme("light")
	if err != nil {
		t.Fatalf("BuiltinTheme: %v", err)
	}
	if cfg.UI.PrimaryFg != "0" || cfg.UI.ViewBackground != "#ffffff" {
		t.Fatalf("unexpected light palette: %+v", cfg.UI)
	}
	if cfg.Glamour.Style != "light" || cfg.Chroma.Style != "github" {
		t.Fatalf("unexpected renderer styles: %+v %+v", cfg.Glamour, cfg.Chroma)
	}
}

func TestBuiltinGruvboxLightSoftTheme(t *testing.T) {
	cfg, err := BuiltinTheme("gruvbox-light-soft")
	if err != nil {
		t.Fatalf("BuiltinTheme: %v", err)
	}
	if cfg.UI.PrimaryFg != "#3c3836" || cfg.UI.AccentYellow != "#af3a03" || cfg.UI.ViewBackground != "" {
		t.Fatalf("unexpected Gruvbox palette: %+v", cfg.UI)
	}
	if cfg.Glamour.Style != "light" || cfg.Chroma.Style != "github" {
		t.Fatalf("unexpected renderer styles: %+v %+v", cfg.Glamour, cfg.Chroma)
	}
}

func TestDocumentedExampleConfigsLoad(t *testing.T) {
	for _, name := range []string{"config.example.toml", "config.gruvbox-light-soft.toml"} {
		t.Run(name, func(t *testing.T) {
			cfg, _, err := Load(filepath.Join("..", "..", "examples", name))
			if err != nil {
				t.Fatalf("Load documented example: %v", err)
			}
			if cfg.Theme.Base == "" {
				t.Fatal("example must select a base")
			}
		})
	}
}

func TestDefaultKeysUseNeovimProfile(t *testing.T) {
	keys := DefaultKeys()
	checks := []struct {
		name string
		got  []string
		want []string
	}{
		{name: "viewer quit", got: keys.Viewer.Quit, want: []string{"Q"}},
		{name: "viewer back", got: keys.Viewer.Back, want: []string{"q", "esc"}},
		{name: "viewer info", got: keys.Viewer.Info, want: []string{"?"}},
		{name: "viewer changes", got: append(append([]string{}, keys.Viewer.Previous...), keys.Viewer.Next...), want: []string{"shift+tab", "tab"}},
		{name: "viewer artifacts", got: append(append([]string{}, keys.Viewer.PreviousTab...), keys.Viewer.NextTab...), want: []string{"h", "l"}},
		{name: "viewer horizontal", got: append(append([]string{}, keys.Viewer.Left...), keys.Viewer.Right...), want: []string{"H", "left", "L", "right"}},
		{name: "viewer pages", got: append(append([]string{}, keys.Viewer.PageUp...), keys.Viewer.PageDown...), want: []string{"pgup", "ctrl+u", "pgdown", "ctrl+d"}},
		{name: "index primary", got: keys.Index.Open, want: []string{"enter"}},
		{name: "index toggle", got: keys.Index.Toggle, want: []string{"space"}},
		{name: "index new", got: keys.Index.New, want: []string{"n"}},
		{name: "index lifecycle", got: keys.Index.Lifecycle, want: []string{"a"}},
		{name: "index inspect", got: keys.Index.Inspect, want: []string{"i"}},
		{name: "index edit", got: keys.Index.Edit, want: []string{"e"}},
		{name: "index validate", got: keys.Index.Validate, want: []string{"v"}},
		{name: "index undo", got: keys.Index.Undo, want: []string{"u"}},
		{name: "index back", got: keys.Index.Back, want: []string{"q", "Q", "esc"}},
		{name: "spec quit", got: keys.Spec.Quit, want: []string{"Q"}},
		{name: "spec back", got: keys.Spec.Back, want: []string{"q", "esc"}},
		{name: "config back", got: keys.Config.Back, want: []string{"q", "?", "esc"}},
	}
	for _, check := range checks {
		if !reflect.DeepEqual(check.got, check.want) {
			t.Errorf("%s = %v, want %v", check.name, check.got, check.want)
		}
	}

	root := reflect.ValueOf(keys)
	for i := range root.NumField() {
		context := root.Field(i)
		if context.Kind() != reflect.Struct {
			continue
		}
		contextType := context.Type()
		for j := range context.NumField() {
			for _, key := range context.Field(j).Interface().([]string) {
				if strings.HasPrefix(key, "ctrl+") && key != "ctrl+u" && key != "ctrl+d" {
					t.Errorf("unexpected Ctrl binding %s.%s=%q", root.Type().Field(i).Tag.Get("toml"), contextType.Field(j).Tag.Get("toml"), key)
				}
			}
		}
	}
}

func TestDefaultKeysBindEveryAction(t *testing.T) {
	keys := DefaultKeys()
	root := reflect.ValueOf(keys)
	typeOfRoot := root.Type()
	for i := range root.NumField() {
		context := root.Field(i)
		if context.Kind() != reflect.Struct {
			continue
		}
		contextType := context.Type()
		for j := range context.NumField() {
			bindings := context.Field(j).Interface().([]string)
			if len(bindings) == 0 {
				t.Errorf("%s.%s has no default binding", typeOfRoot.Field(i).Tag.Get("toml"), contextType.Field(j).Tag.Get("toml"))
			}
		}
	}
	if err := ValidateKeys(keys); err != nil {
		t.Fatalf("default keys must be valid: %v", err)
	}
}

func TestValidateKeysRejectsContextCollision(t *testing.T) {
	keys := DefaultKeys()
	keys.Viewer.Down = []string{"n"}
	keys.Viewer.Up = []string{"n"}

	err := ValidateKeys(keys)
	if err == nil || !strings.Contains(err.Error(), "viewer.down") || !strings.Contains(err.Error(), "viewer.up") {
		t.Fatalf("expected contextual collision, got %v", err)
	}
}

func TestValidateKeysAllowsReuseAcrossContexts(t *testing.T) {
	keys := DefaultKeys()
	keys.Viewer.Down = []string{"z"}
	keys.Index.Down = []string{"z"}
	if err := ValidateKeys(keys); err != nil {
		t.Fatalf("expected cross-context reuse to be valid: %v", err)
	}
}

func TestValidateKeysRejectsEmptyAndSequences(t *testing.T) {
	tests := []struct {
		name string
		keys []string
	}{
		{name: "empty list", keys: []string{}},
		{name: "empty key", keys: []string{""}},
		{name: "sequence", keys: []string{"g g"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := DefaultKeys()
			keys.Viewer.Down = tt.keys
			if err := ValidateKeys(keys); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
