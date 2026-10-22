package settings

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed themes/*.toml
var builtinThemes embed.FS

type Config struct {
	UI    UIConfig    `toml:"ui"`
	Theme ThemeConfig `toml:"theme"`
	Keys  KeyConfig   `toml:"keys"`
}

type UIConfig struct {
	StartView string `toml:"start_view"`
}

type ThemeConfig struct {
	Base    string        `toml:"base"`
	UI      UIColors      `toml:"ui"`
	Glamour GlamourConfig `toml:"glamour"`
	Chroma  ChromaConfig  `toml:"chroma"`
}

type UIColors struct {
	PrimaryFg            string `toml:"primary_fg"`
	MutedFg              string `toml:"muted_fg"`
	MidFg                string `toml:"mid_fg"`
	AccentBlue           string `toml:"accent_blue"`
	AccentYellow         string `toml:"accent_yellow"`
	AccentCyan           string `toml:"accent_cyan"`
	AccentGreen          string `toml:"accent_green"`
	AccentRed            string `toml:"accent_red"`
	AccentMagenta        string `toml:"accent_magenta"`
	ActiveBg             string `toml:"active_bg"`
	ActiveFg             string `toml:"active_fg"`
	ViewBackground       string `toml:"view_background"`
	DiffAddBackground    string `toml:"diff_add_background"`
	DiffRemoveBackground string `toml:"diff_remove_background"`
}

type ColorOverride struct {
	Foreground string `toml:"foreground"`
	Background string `toml:"background"`
}

type GlamourConfig struct {
	Style    string                   `toml:"style"`
	Elements map[string]ColorOverride `toml:"elements"`
	Syntax   map[string]ColorOverride `toml:"syntax"`
}

type ChromaConfig struct {
	Style  string            `toml:"style"`
	Tokens map[string]string `toml:"tokens"`
}

type KeyConfig struct {
	Style  string        `toml:"style"`
	Viewer KeyViewer     `toml:"viewer"`
	Index  KeyIndex      `toml:"index"`
	Filter KeyFilter     `toml:"filter"`
	Spec   KeySpec       `toml:"spec"`
	Config KeyConfigView `toml:"config"`
}

type KeyViewer struct {
	Quit         []string `toml:"quit"`
	Stage        []string `toml:"stage"`
	Info         []string `toml:"info"`
	Back         []string `toml:"back"`
	Previous     []string `toml:"previous"`
	Next         []string `toml:"next"`
	Right        []string `toml:"right"`
	Left         []string `toml:"left"`
	NextDiff     []string `toml:"next_diff"`
	PreviousDiff []string `toml:"previous_diff"`
	// Legacy action names are retained for configuration compatibility; the
	// viewer dispatches them positionally to schema artifacts 1 through 4.
	ProposalTab []string `toml:"proposal_tab"`
	DesignTab   []string `toml:"design_tab"`
	SpecsTab    []string `toml:"specs_tab"`
	TasksTab    []string `toml:"tasks_tab"`
	GitTab      []string `toml:"git_tab"`
	NextTab     []string `toml:"next_tab"`
	PreviousTab []string `toml:"previous_tab"`
	ViewDiff    []string `toml:"view_diff"`
	Down        []string `toml:"down"`
	PageDown    []string `toml:"page_down"`
	Up          []string `toml:"up"`
	PageUp      []string `toml:"page_up"`
	ToggleTask  []string `toml:"toggle_task"`
	Open        []string `toml:"open"`
}

type KeyIndex struct {
	Filter    []string `toml:"filter"`
	Info      []string `toml:"info"`
	New       []string `toml:"new"`
	Lifecycle []string `toml:"lifecycle"`
	Inspect   []string `toml:"inspect"`
	Edit      []string `toml:"edit"`
	Validate  []string `toml:"validate"`
	Undo      []string `toml:"undo"`
	Back      []string `toml:"back"`
	Down      []string `toml:"down"`
	Up        []string `toml:"up"`
	Open      []string `toml:"open"`
	Toggle    []string `toml:"toggle"`
	Sort      []string `toml:"sort"`
}

type KeyFilter struct {
	Cancel    []string `toml:"cancel"`
	Accept    []string `toml:"accept"`
	Backspace []string `toml:"backspace"`
}

type KeySpec struct {
	Quit                []string `toml:"quit"`
	Back                []string `toml:"back"`
	Down                []string `toml:"down"`
	Up                  []string `toml:"up"`
	PageDown            []string `toml:"page_down"`
	PageUp              []string `toml:"page_up"`
	PreviousRequirement []string `toml:"previous_requirement"`
	NextRequirement     []string `toml:"next_requirement"`
}

type KeyConfigView struct {
	Back     []string `toml:"back"`
	Down     []string `toml:"down"`
	Up       []string `toml:"up"`
	PageDown []string `toml:"page_down"`
	PageUp   []string `toml:"page_up"`
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config directory: %w", err)
	}
	return filepath.Join(dir, "dossier", "config.toml"), nil
}

func Load(path string) (Config, string, error) {
	return load(path, "")
}

func LoadWithKeyStyle(path, keyStyle string) (Config, string, error) {
	return load(path, keyStyle)
}

func load(path, selectedKeyStyle string) (Config, string, error) {
	explicit := path != ""
	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return Config{}, "", err
		}
	}

	keyStyle := selectedKeyStyle
	if keyStyle == "" {
		keyStyle = "nvim"
	}
	baseKeys, err := BuiltinKeyStyle(keyStyle)
	if err != nil {
		return Config{}, path, fmt.Errorf("config %s: %w", path, err)
	}
	cfg := Config{UI: UIConfig{StartView: "index"}, Theme: ThemeConfig{Base: "none"}, Keys: baseKeys}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && !explicit {
			return cfg, path, nil
		}
		return Config{}, path, fmt.Errorf("config %s: %w", path, err)
	}

	var fileCfg Config
	meta, err := toml.Decode(string(data), &fileCfg)
	if err != nil {
		return Config{}, path, fmt.Errorf("config %s: %w", path, err)
	}
	if unknown := meta.Undecoded(); len(unknown) > 0 {
		names := make([]string, len(unknown))
		for i, key := range unknown {
			names[i] = key.String()
		}
		return Config{}, path, fmt.Errorf("config %s: unknown keys: %s", path, strings.Join(names, ", "))
	}
	if fileCfg.UI.StartView != "" {
		cfg.UI.StartView = fileCfg.UI.StartView
	}
	if cfg.UI.StartView != "index" && cfg.UI.StartView != "change" {
		return Config{}, path, fmt.Errorf("config %s: ui.start_view must be \"index\" or \"change\", got %q", path, fileCfg.UI.StartView)
	}
	if fileCfg.Theme.Base == "" {
		fileCfg.Theme.Base = "none"
	}
	cfg.Theme = fileCfg.Theme
	if selectedKeyStyle == "" && fileCfg.Keys.Style != "" {
		keyStyle = fileCfg.Keys.Style
		baseKeys, err = BuiltinKeyStyle(keyStyle)
		if err != nil {
			return Config{}, path, fmt.Errorf("config %s: %w", path, err)
		}
	}
	cfg.Keys = MergeKeys(baseKeys, fileCfg.Keys)
	cfg.Keys.Style = strings.ToLower(keyStyle)
	if err := ValidateKeys(cfg.Keys); err != nil {
		return Config{}, path, fmt.Errorf("config %s: %w", path, err)
	}
	return cfg, path, nil
}

func BuiltinTheme(name string) (ThemeConfig, error) {
	name = strings.ToLower(name)
	data, err := builtinThemes.ReadFile("themes/" + name + ".toml")
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("unknown theme %q; available: dark, none, light, dracula, gruvbox-light-soft", name)
	}
	var cfg Config
	meta, err := toml.Decode(string(data), &cfg)
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("embedded theme %s: %w", name, err)
	}
	if unknown := meta.Undecoded(); len(unknown) > 0 {
		return ThemeConfig{}, fmt.Errorf("embedded theme %s has unknown keys: %v", name, unknown)
	}
	return cfg.Theme, nil
}

func MergeTheme(base, override ThemeConfig) ThemeConfig {
	merged := base
	mergeStrings(reflect.ValueOf(&merged).Elem(), reflect.ValueOf(override))
	merged.Glamour.Elements = mergeColorMap(base.Glamour.Elements, override.Glamour.Elements)
	merged.Glamour.Syntax = mergeColorMap(base.Glamour.Syntax, override.Glamour.Syntax)
	merged.Chroma.Tokens = mergeStringMap(base.Chroma.Tokens, override.Chroma.Tokens)
	return merged
}

func mergeStrings(dst, src reflect.Value) {
	for i := range dst.NumField() {
		df, sf := dst.Field(i), src.Field(i)
		switch df.Kind() {
		case reflect.String:
			if sf.String() != "" {
				df.SetString(sf.String())
			}
		case reflect.Struct:
			mergeStrings(df, sf)
		}
	}
}

func mergeColorMap(base, override map[string]ColorOverride) map[string]ColorOverride {
	merged := make(map[string]ColorOverride, len(base)+len(override))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range override {
		current := merged[key]
		if value.Foreground != "" {
			current.Foreground = value.Foreground
		}
		if value.Background != "" {
			current.Background = value.Background
		}
		merged[key] = current
	}
	return merged
}

func mergeStringMap(base, override map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(override))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range override {
		merged[key] = value
	}
	return merged
}

func BuiltinKeyStyle(name string) (KeyConfig, error) {
	name = strings.ToLower(name)
	if name != "nvim" {
		return KeyConfig{}, fmt.Errorf("unknown keystyle %q; available: nvim", name)
	}
	return nvimKeys(), nil
}

func DefaultKeys() KeyConfig {
	return nvimKeys()
}

func nvimKeys() KeyConfig {
	return KeyConfig{
		Style: "nvim",
		Viewer: KeyViewer{
			Quit: []string{"Q"}, Stage: []string{"s"}, Info: []string{"?"}, Back: []string{"q", "esc"},
			Previous: []string{"shift+tab"}, Next: []string{"tab"}, Right: []string{"L", "right"}, Left: []string{"H", "left"},
			NextDiff: []string{"]"}, PreviousDiff: []string{"["}, ProposalTab: []string{"1"}, DesignTab: []string{"2"},
			SpecsTab: []string{"3"}, TasksTab: []string{"4"}, GitTab: []string{"5"}, NextTab: []string{"l"},
			PreviousTab: []string{"h"}, ViewDiff: []string{"d", "enter"}, Down: []string{"j", "down"},
			PageDown: []string{"pgdown", "ctrl+d"}, Up: []string{"k", "up"}, PageUp: []string{"pgup", "ctrl+u"},
			ToggleTask: []string{"space"}, Open: []string{"e"},
		},
		Index: KeyIndex{
			Filter: []string{"/"}, Info: []string{"?"}, New: []string{"n"}, Lifecycle: []string{"a"},
			Inspect: []string{"i"}, Edit: []string{"e"}, Validate: []string{"v"}, Undo: []string{"u"},
			Back: []string{"q", "Q", "esc"}, Down: []string{"j", "down"}, Up: []string{"k", "up"},
			Open: []string{"enter"}, Toggle: []string{"space"}, Sort: []string{"s"},
		},
		Filter: KeyFilter{Cancel: []string{"esc"}, Accept: []string{"enter"}, Backspace: []string{"backspace"}},
		Spec: KeySpec{
			Quit: []string{"Q"}, Back: []string{"q", "esc"}, Down: []string{"j", "down"}, Up: []string{"k", "up"},
			PageDown: []string{"pgdown", "ctrl+d"}, PageUp: []string{"pgup", "ctrl+u"}, PreviousRequirement: []string{"h"}, NextRequirement: []string{"l"},
		},
		Config: KeyConfigView{
			Back: []string{"q", "?", "esc"}, Down: []string{"j", "down"}, Up: []string{"k", "up"},
			PageDown: []string{"pgdown", "ctrl+d"}, PageUp: []string{"pgup", "ctrl+u"},
		},
	}
}

func MergeKeys(base, override KeyConfig) KeyConfig {
	merged := base
	mergeKeyValues(reflect.ValueOf(&merged).Elem(), reflect.ValueOf(override))
	return merged
}

func mergeKeyValues(dst, src reflect.Value) {
	for i := range dst.NumField() {
		df, sf := dst.Field(i), src.Field(i)
		switch df.Kind() {
		case reflect.String:
			if sf.String() != "" {
				df.SetString(sf.String())
			}
		case reflect.Struct:
			mergeKeyValues(df, sf)
		case reflect.Slice:
			if !sf.IsNil() {
				df.Set(sf)
			}
		}
	}
}

func ValidateKeys(keys KeyConfig) error {
	root := reflect.ValueOf(keys)
	typeOfRoot := root.Type()
	for i := range root.NumField() {
		context := root.Field(i)
		if context.Kind() != reflect.Struct {
			continue
		}
		contextName := typeOfRoot.Field(i).Tag.Get("toml")
		contextType := context.Type()
		seen := map[string]string{}
		for j := range context.NumField() {
			action := contextType.Field(j).Tag.Get("toml")
			qualified := contextName + "." + action
			bindings := context.Field(j).Interface().([]string)
			if bindings == nil {
				continue
			}
			if len(bindings) == 0 {
				return fmt.Errorf("keys.%s must contain at least one key", qualified)
			}
			for _, key := range bindings {
				if key == "" {
					return fmt.Errorf("keys.%s contains an empty key", qualified)
				}
				if strings.TrimSpace(key) != key || strings.ContainsAny(key, " \t\r\n") {
					return fmt.Errorf("keys.%s contains unsupported key sequence %q", qualified, key)
				}
				if previous, ok := seen[key]; ok {
					return fmt.Errorf("key %q conflicts between keys.%s and keys.%s", key, previous, qualified)
				}
				seen[key] = qualified
			}
		}
	}
	return nil
}
