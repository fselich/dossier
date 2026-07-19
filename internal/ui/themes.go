package ui

import (
	"image/color"
	"strings"
)

type Theme struct {
	Name         string
	GlamourStyle string
	ChromaStyle  string
	ViewBg       color.Color
	DiffAddBg    string
	DiffRemoveBg string
}

type ThemeColors struct {
}

var (
	DarkTheme = Theme{
		Name:         "dark",
		GlamourStyle: "dark",
		ChromaStyle:  "monokai",
		ViewBg:       color.RGBA{26, 26, 26, 255},
		DiffAddBg:    "#1a3a1a",
		DiffRemoveBg: "#3a1a1a",
	}

	NoneTheme = Theme{
		Name:         "none",
		GlamourStyle: "dark",
		ChromaStyle:  "monokai",
		ViewBg:       nil,
		DiffAddBg:    "#1a3a1a",
		DiffRemoveBg: "#3a1a1a",
	}

	LightTheme = Theme{
		Name:         "light",
		GlamourStyle: "light",
		ChromaStyle:  "github",
		ViewBg:       color.RGBA{255, 255, 255, 255},
		DiffAddBg:    "#e6ffed",
		DiffRemoveBg: "#ffeef0",
	}

	DraculaTheme = Theme{
		Name:         "dracula",
		GlamourStyle: "dracula",
		ChromaStyle:  "dracula",
		ViewBg:       color.RGBA{40, 42, 54, 255},
		DiffAddBg:    "#1f3425",
		DiffRemoveBg: "#3d1f26",
	}
)

var Themes = map[string]Theme{
	"dark":    DarkTheme,
	"none":    NoneTheme,
	"light":   LightTheme,
	"dracula": DraculaTheme,
}

func DefaultTheme() Theme {
	return DarkTheme
}

func LookupTheme(name string) (Theme, bool) {
	t, ok := Themes[strings.ToLower(name)]
	return t, ok
}
