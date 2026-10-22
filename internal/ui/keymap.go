package ui

import (
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/settings"
)

func (m Model) effectiveKeyMap() settings.KeyConfig {
	if m.keyMap.Viewer.Quit == nil {
		return settings.DefaultKeys()
	}
	return m.keyMap
}

func matchesKey(msg tea.KeyPressMsg, bindings []string) bool {
	key := msg.String()
	for _, binding := range bindings {
		if key == binding {
			return true
		}
	}
	return false
}

func primaryKeyLabel(bindings []string) string {
	if len(bindings) == 0 {
		return "?"
	}
	return displayKey(bindings[0])
}

func pairedKeyLabel(first, second []string) string {
	return primaryKeyLabel(first) + "/" + primaryKeyLabel(second)
}

func combinedKeyLabel(groups ...[]string) string {
	seen := map[string]bool{}
	var labels []string
	for _, group := range groups {
		for _, key := range group {
			label := displayKey(key)
			if !seen[label] {
				seen[label] = true
				labels = append(labels, label)
			}
		}
	}
	return strings.Join(labels, "/")
}

func displayKey(key string) string {
	switch key {
	case "space":
		return "Space"
	case "enter":
		return "Enter"
	case "esc":
		return "Esc"
	case "tab":
		return "Tab"
	case "shift+tab":
		return "Shift+Tab"
	case "backspace":
		return "Backspace"
	case "up":
		return "↑"
	case "down":
		return "↓"
	case "left":
		return "←"
	case "right":
		return "→"
	case "pgup":
		return "PgUp"
	case "pgdown":
		return "PgDown"
	}
	if strings.HasPrefix(key, "ctrl+") {
		return "Ctrl+" + strings.ToUpper(strings.TrimPrefix(key, "ctrl+"))
	}
	if utf8.RuneCountInString(key) == 1 {
		return key
	}
	return key
}
