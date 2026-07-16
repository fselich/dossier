package ui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/fselich/dossier/internal/openspec"
)

var sgrRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func forceColor(t *testing.T) {
	t.Helper()
	t.Setenv("NO_COLOR", "0")
	t.Setenv("CLICOLOR_FORCE", "1")
	t.Setenv("TERM", "xterm-256color")
}

func activeSGRBefore(t *testing.T, content, needle string) string {
	t.Helper()
	idx := strings.Index(content, needle)
	if idx < 0 {
		t.Fatalf("expected %q in rendered content, got %q", needle, content)
	}

	active := ""
	for _, loc := range sgrRe.FindAllStringIndex(content[:idx], -1) {
		seq := content[loc[0]:loc[1]]
		switch seq {
		case "\x1b[m", "\x1b[0m":
			active = ""
		default:
			active = seq
		}
	}
	return active
}

func TestRenderTasksContentKeepsFocusedRowsVisible(t *testing.T) {
	forceColor(t)

	items := []openspec.TaskItem{
		{Kind: openspec.KindSection, Text: "1. Done section"},
		{Kind: openspec.KindTask, Text: "completed item", Done: true},
		{Kind: openspec.KindSection, Text: "2. Active section"},
		{Kind: openspec.KindTask, Text: "pending item"},
	}

	t.Run("focused section keeps section text", func(t *testing.T) {
		m := &Model{width: 100, tasks: taskState{Items: items, Cursor: 2}}

		content, _ := m.renderTasksContent()

		if got, want := activeSGRBefore(t, content, "2. Active section"), extractOpeningEscape(sectionStyle); got != want {
			t.Fatalf("expected focused section to keep section style %q, got %q in %q", want, got, content)
		}
	})

	t.Run("focused pending task keeps checkbox and text", func(t *testing.T) {
		m := &Model{width: 100, tasks: taskState{Items: items, Cursor: 3}}

		content, _ := m.renderTasksContent()

		if got, want := activeSGRBefore(t, content, "[ ]"), extractOpeningEscape(taskPendingStyle); got != want {
			t.Fatalf("expected focused pending checkbox to keep pending style %q, got %q in %q", want, got, content)
		}
	})

	t.Run("focused done task keeps checkbox and text", func(t *testing.T) {
		m := &Model{width: 100, tasks: taskState{Items: items, Cursor: 1}}

		content, _ := m.renderTasksContent()

		if got, want := activeSGRBefore(t, content, "[x]"), extractOpeningEscape(taskDoneStyle); got != want {
			t.Fatalf("expected focused done checkbox to keep done style %q, got %q in %q", want, got, content)
		}
	})
}
