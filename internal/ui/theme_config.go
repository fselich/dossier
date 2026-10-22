package ui

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	glamourstyles "charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
	"github.com/aemonge/dossier/internal/settings"
	"github.com/alecthomas/chroma/v2"
	chromastyles "github.com/alecthomas/chroma/v2/styles"
)

var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func BuildTheme(cfg settings.ThemeConfig) (Theme, error) {
	baseName := cfg.Base
	if baseName == "" {
		baseName = "none"
	}
	base, err := settings.BuiltinTheme(baseName)
	if err != nil {
		return Theme{}, err
	}
	cfg = settings.MergeTheme(base, cfg)
	if err := validateUIColors(cfg.UI); err != nil {
		return Theme{}, err
	}

	glamourStyle, ok := glamourstyles.DefaultStyles[strings.ToLower(cfg.Glamour.Style)]
	if !ok {
		return Theme{}, fmt.Errorf("unknown Glamour style %q", cfg.Glamour.Style)
	}
	resolvedGlamour := *glamourStyle
	if resolvedGlamour.CodeBlock.Chroma != nil {
		chromaCopy := *resolvedGlamour.CodeBlock.Chroma
		resolvedGlamour.CodeBlock.Chroma = &chromaCopy
	}
	if err := applyGlamourOverrides(&resolvedGlamour, cfg.Glamour); err != nil {
		return Theme{}, err
	}

	resolvedChroma, err := buildChromaStyle(cfg.Chroma)
	if err != nil {
		return Theme{}, err
	}

	colors := ThemeColors{
		PrimaryFg:     lipgloss.Color(cfg.UI.PrimaryFg),
		MutedFg:       lipgloss.Color(cfg.UI.MutedFg),
		MidFg:         lipgloss.Color(cfg.UI.MidFg),
		AccentBlue:    lipgloss.Color(cfg.UI.AccentBlue),
		AccentYellow:  lipgloss.Color(cfg.UI.AccentYellow),
		AccentCyan:    lipgloss.Color(cfg.UI.AccentCyan),
		AccentGreen:   lipgloss.Color(cfg.UI.AccentGreen),
		AccentRed:     lipgloss.Color(cfg.UI.AccentRed),
		AccentMagenta: lipgloss.Color(cfg.UI.AccentMagenta),
		ActiveBg:      lipgloss.Color(cfg.UI.ActiveBg),
		ActiveFg:      lipgloss.Color(cfg.UI.ActiveFg),
	}

	theme := Theme{
		Name:          cfg.Base,
		GlamourStyle:  cfg.Glamour.Style,
		GlamourStyles: resolvedGlamour,
		ChromaStyle:   cfg.Chroma.Style,
		Chroma:        resolvedChroma,
		DiffAddBg:     cfg.UI.DiffAddBackground,
		DiffRemoveBg:  cfg.UI.DiffRemoveBackground,
		Colors:        colors,
		Styles:        BuildStyles(colors),
	}
	if cfg.UI.ViewBackground != "" {
		theme.ViewBg = lipgloss.Color(cfg.UI.ViewBackground)
	}
	return theme, nil
}

func mustBuiltinTheme(name string) Theme {
	cfg, err := settings.BuiltinTheme(name)
	if err != nil {
		panic(err)
	}
	theme, err := BuildTheme(cfg)
	if err != nil {
		panic(err)
	}
	return theme
}

func validateUIColors(colors settings.UIColors) error {
	value := reflect.ValueOf(colors)
	typeOfValue := value.Type()
	for i := range value.NumField() {
		name := typeOfValue.Field(i).Tag.Get("toml")
		color := value.Field(i).String()
		if name == "view_background" && color == "" {
			continue
		}
		if color == "" {
			return fmt.Errorf("theme.ui.%s must be set", name)
		}
		if err := validateColor(color); err != nil {
			return fmt.Errorf("theme.ui.%s: %w", name, err)
		}
	}
	return nil
}

func validateColor(value string) error {
	if hexColorPattern.MatchString(value) {
		return nil
	}
	n, err := strconv.Atoi(value)
	if err == nil && n >= 0 && n <= 255 {
		return nil
	}
	return fmt.Errorf("invalid color %q; use #RRGGBB or an ANSI color number 0-255", value)
}

func applyGlamourOverrides(style any, cfg settings.GlamourConfig) error {
	root := reflect.ValueOf(style)
	if root.Kind() != reflect.Pointer || root.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("invalid Glamour style target")
	}
	root = root.Elem()
	for name, override := range cfg.Elements {
		field, ok := findNormalizedField(root, name)
		if !ok {
			return fmt.Errorf("unknown Glamour element %q", name)
		}
		if err := applyColorOverride(field, override); err != nil {
			return fmt.Errorf("glamour element %s: %w", name, err)
		}
	}

	codeBlock := root.FieldByName("CodeBlock")
	chromaField := codeBlock.FieldByName("Chroma")
	if chromaField.IsNil() {
		chromaField.Set(reflect.New(chromaField.Type().Elem()))
	}
	chromaValue := chromaField.Elem()
	for name, override := range cfg.Syntax {
		field, ok := findNormalizedField(chromaValue, name)
		if !ok {
			return fmt.Errorf("unknown Glamour syntax token %q", name)
		}
		if err := applyColorOverride(field, override); err != nil {
			return fmt.Errorf("glamour syntax token %s: %w", name, err)
		}
	}
	return nil
}

func findNormalizedField(value reflect.Value, name string) (reflect.Value, bool) {
	want := normalizeName(name)
	typeOfValue := value.Type()
	for i := range value.NumField() {
		if normalizeName(typeOfValue.Field(i).Name) == want {
			return value.Field(i), true
		}
	}
	return reflect.Value{}, false
}

func normalizeName(value string) string {
	return strings.Map(func(r rune) rune {
		if r == '_' || r == '-' || unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, value)
}

func applyColorOverride(target reflect.Value, override settings.ColorOverride) error {
	if override.Foreground != "" {
		if err := validateColor(override.Foreground); err != nil {
			return fmt.Errorf("foreground: %w", err)
		}
		field := target.FieldByName("Color")
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("foreground is not configurable")
		}
		value := override.Foreground
		field.Set(reflect.ValueOf(&value))
	}
	if override.Background != "" {
		if err := validateColor(override.Background); err != nil {
			return fmt.Errorf("background: %w", err)
		}
		field := target.FieldByName("BackgroundColor")
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("background is not configurable")
		}
		value := override.Background
		field.Set(reflect.ValueOf(&value))
	}
	return nil
}

func buildChromaStyle(cfg settings.ChromaConfig) (*chroma.Style, error) {
	base := chromastyles.Get(cfg.Style)
	if base == nil {
		return nil, fmt.Errorf("unknown Chroma style %q", cfg.Style)
	}
	builder := base.Builder()
	for name, descriptor := range cfg.Tokens {
		tokenType, err := chroma.TokenTypeString(name)
		if err != nil {
			return nil, fmt.Errorf("unknown Chroma token %q", name)
		}
		builder.Add(tokenType, descriptor)
	}
	style, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("chroma style: %w", err)
	}
	return style, nil
}
