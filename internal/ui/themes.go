package ui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

type ThemeColors struct {
	PrimaryFg     color.Color
	MutedFg       color.Color
	MidFg         color.Color
	AccentBlue    color.Color
	AccentYellow  color.Color
	AccentCyan    color.Color
	AccentGreen   color.Color
	AccentRed     color.Color
	AccentMagenta color.Color
	ActiveBg      color.Color
	ActiveFg      color.Color
}

type ThemeStyles struct {
	BaseText         lipgloss.Style
	Header           lipgloss.Style
	TabActive        lipgloss.Style
	TabInactive      lipgloss.Style
	TabDisabled      lipgloss.Style
	IndexActive      lipgloss.Style
	Section          lipgloss.Style
	TaskCursorMark   lipgloss.Style
	TaskDone         lipgloss.Style
	TaskPending      lipgloss.Style
	Help             lipgloss.Style
	Error            lipgloss.Style
	ProgressDone     lipgloss.Style
	ProgressComplete lipgloss.Style
	ProgressEmpty    lipgloss.Style
	Separator        lipgloss.Style
	GitModified      lipgloss.Style
	GitAdded         lipgloss.Style
	GitDeleted       lipgloss.Style
	GitRenamed       lipgloss.Style
	GitUntracked     lipgloss.Style
	GitDot           lipgloss.Style
	DiffRemoved      lipgloss.Style
	TaskCodeCyan     lipgloss.Style
	TaskCodeDone     lipgloss.Style
}

type Theme struct {
	Name         string
	GlamourStyle string
	ChromaStyle  string
	ViewBg       color.Color
	DiffAddBg    string
	DiffRemoveBg string
	Colors       ThemeColors
	Styles       ThemeStyles
}

func BuildStyles(c ThemeColors) ThemeStyles {
	return ThemeStyles{
		BaseText: lipgloss.NewStyle().
			Foreground(c.PrimaryFg),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(c.AccentBlue),

		TabActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(c.ActiveFg).
			Background(c.ActiveBg).
			Padding(0, 1),

		TabInactive: lipgloss.NewStyle().
			Foreground(c.PrimaryFg).
			Padding(0, 1),

		TabDisabled: lipgloss.NewStyle().
			Foreground(c.MutedFg).
			Padding(0, 1),

		IndexActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(c.ActiveFg).
			Background(c.ActiveBg),

		Section: lipgloss.NewStyle().
			Bold(true).
			Foreground(c.AccentYellow),

		TaskCursorMark: lipgloss.NewStyle().
			Bold(true).
			Foreground(c.PrimaryFg),

		TaskDone: lipgloss.NewStyle().
			Foreground(c.MutedFg),

		TaskPending: lipgloss.NewStyle().
			Foreground(c.MidFg),

		Help: lipgloss.NewStyle().
			Foreground(c.MutedFg),

		Error: lipgloss.NewStyle().
			Foreground(c.AccentRed).
			Bold(true),

		ProgressDone: lipgloss.NewStyle().
			Foreground(c.AccentCyan),

		ProgressComplete: lipgloss.NewStyle().
			Foreground(c.AccentGreen),

		ProgressEmpty: lipgloss.NewStyle().
			Foreground(c.MutedFg),

		Separator: lipgloss.NewStyle().
			Foreground(c.MutedFg),

		GitModified: lipgloss.NewStyle().
			Foreground(c.AccentYellow),

		GitAdded: lipgloss.NewStyle().
			Foreground(c.AccentGreen),

		GitDeleted: lipgloss.NewStyle().
			Foreground(c.MutedFg),

		GitRenamed: lipgloss.NewStyle().
			Foreground(c.AccentCyan),

		GitUntracked: lipgloss.NewStyle().
			Foreground(c.AccentMagenta),

		GitDot: lipgloss.NewStyle().
			Foreground(c.MutedFg),

		DiffRemoved: lipgloss.NewStyle().
			Foreground(c.AccentRed),

		TaskCodeCyan: lipgloss.NewStyle().
			Foreground(c.AccentCyan),

		TaskCodeDone: lipgloss.NewStyle().
			Underline(true).
			Foreground(c.MutedFg),
	}
}

var (
	DarkColors = ThemeColors{
		PrimaryFg:     lipgloss.Color("15"),
		MutedFg:       lipgloss.Color("8"),
		MidFg:         lipgloss.Color("7"),
		AccentBlue:    lipgloss.Color("12"),
		AccentYellow:  lipgloss.Color("11"),
		AccentCyan:    lipgloss.Color("6"),
		AccentGreen:   lipgloss.Color("2"),
		AccentRed:     lipgloss.Color("9"),
		AccentMagenta: lipgloss.Color("5"),
		ActiveBg:      lipgloss.Color("4"),
		ActiveFg:      lipgloss.Color("15"),
	}

	LightColors = ThemeColors{
		PrimaryFg:     lipgloss.Color("0"),
		MutedFg:       lipgloss.Color("7"),
		MidFg:         lipgloss.Color("8"),
		AccentBlue:    lipgloss.Color("4"),
		AccentYellow:  lipgloss.Color("3"),
		AccentCyan:    lipgloss.Color("6"),
		AccentGreen:   lipgloss.Color("2"),
		AccentRed:     lipgloss.Color("1"),
		AccentMagenta: lipgloss.Color("5"),
		ActiveBg:      lipgloss.Color("4"),
		ActiveFg:      lipgloss.Color("15"),
	}

	DraculaColors = ThemeColors{
		PrimaryFg:     lipgloss.Color("15"),
		MutedFg:       lipgloss.Color("8"),
		MidFg:         lipgloss.Color("7"),
		AccentBlue:    lipgloss.Color("12"),
		AccentYellow:  lipgloss.Color("3"),
		AccentCyan:    lipgloss.Color("6"),
		AccentGreen:   lipgloss.Color("2"),
		AccentRed:     lipgloss.Color("9"),
		AccentMagenta: lipgloss.Color("5"),
		ActiveBg:      lipgloss.Color("4"),
		ActiveFg:      lipgloss.Color("15"),
	}

	DarkTheme = Theme{
		Name:         "dark",
		GlamourStyle: "dark",
		ChromaStyle:  "monokai",
		ViewBg:       color.RGBA{26, 26, 26, 255},
		DiffAddBg:    "#1a3a1a",
		DiffRemoveBg: "#3a1a1a",
		Colors:       DarkColors,
		Styles:       BuildStyles(DarkColors),
	}

	NoneTheme = Theme{
		Name:         "none",
		GlamourStyle: "dark",
		ChromaStyle:  "monokai",
		ViewBg:       nil,
		DiffAddBg:    "#1a3a1a",
		DiffRemoveBg: "#3a1a1a",
		Colors:       DarkColors,
		Styles:       BuildStyles(DarkColors),
	}

	LightTheme = Theme{
		Name:         "light",
		GlamourStyle: "light",
		ChromaStyle:  "github",
		ViewBg:       color.RGBA{255, 255, 255, 255},
		DiffAddBg:    "#e6ffed",
		DiffRemoveBg: "#ffeef0",
		Colors:       LightColors,
		Styles:       BuildStyles(LightColors),
	}

	DraculaTheme = Theme{
		Name:         "dracula",
		GlamourStyle: "dracula",
		ChromaStyle:  "dracula",
		ViewBg:       color.RGBA{40, 42, 54, 255},
		DiffAddBg:    "#1f3425",
		DiffRemoveBg: "#3d1f26",
		Colors:       DraculaColors,
		Styles:       BuildStyles(DraculaColors),
	}
)

var Themes = map[string]Theme{
	"dark":    DarkTheme,
	"none":    NoneTheme,
	"light":   LightTheme,
	"dracula": DraculaTheme,
}

func DefaultTheme() Theme {
	return NoneTheme
}

func LookupTheme(name string) (Theme, bool) {
	t, ok := Themes[strings.ToLower(name)]
	return t, ok
}
