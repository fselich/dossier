package ui

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

func TestCustomViewerNavigationReplacesDefaultBinding(t *testing.T) {
	keys := settings.DefaultKeys()
	keys.Viewer.Down = []string{"n"}
	m := Model{
		mode: ModeNormal, keyMap: keys, artifactSelection: artifactSelection{ArtifactID: "tasks", OutputPath: "tasks.md"},
		project: &openspec.Project{Changes: []openspec.Change{{Artifacts: []openspec.ChangeArtifact{{
			ID: "tasks", Outputs: []openspec.ArtifactOutput{{RelativePath: "tasks.md", Present: true}},
		}}}}},
		tasks: taskState{
			Items: []openspec.TaskItem{
				{Kind: openspec.KindTask, Text: "one"},
				{Kind: openspec.KindTask, Text: "two"},
			},
		},
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "n"})
	m = result.(Model)
	if m.tasks.Cursor != 1 {
		t.Fatalf("custom down key did not move cursor: %d", m.tasks.Cursor)
	}

	m.tasks.Cursor = 0
	result, _ = m.dispatchKey(tea.KeyPressMsg{Text: "j"})
	m = result.(Model)
	if m.tasks.Cursor != 0 {
		t.Fatalf("replaced default down key should be inactive: %d", m.tasks.Cursor)
	}
}

func TestCustomKeymapUpdatesHelp(t *testing.T) {
	keys := settings.DefaultKeys()
	keys.Viewer.ToggleTask = []string{"x"}
	keys.Viewer.Open = []string{"o"}
	m := Model{mode: ModeNormal, artifactSelection: artifactSelection{ArtifactID: "tasks"}, keyMap: keys}

	help := m.renderHelpBar()
	if !strings.Contains(help, "x: toggle") || !strings.Contains(help, "o: edit") {
		t.Fatalf("custom help labels missing: %q", help)
	}
	if strings.Contains(help, "Space: toggle") || strings.Contains(help, "e: edit") {
		t.Fatalf("default help labels should be replaced: %q", help)
	}
}

func TestDefaultHelpShowsVisibleArtifactsAndHiddenChanges(t *testing.T) {
	m := Model{mode: ModeNormal, artifactSelection: artifactSelection{ArtifactID: "proposal"}}
	help := m.renderHelpBar()
	if !strings.Contains(help, "Shift+Tab/Tab: change") {
		t.Fatalf("change navigation help missing: %q", help)
	}
	if !strings.Contains(help, "h/l: artifact") {
		t.Fatalf("artifact navigation help missing: %q", help)
	}

	m.viewingCode = true
	m.gitState.ShowingDiff = true
	help = m.renderHelpBar()
	if !strings.Contains(help, "H/L: horizontal") {
		t.Fatalf("diff horizontal help missing: %q", help)
	}
	if strings.Contains(help, "Shift+Tab/Tab: horizontal") {
		t.Fatalf("change keys must not be advertised as horizontal: %q", help)
	}
}

func TestReadOnlyCustomHelpOmitsMutationBindings(t *testing.T) {
	keys := settings.DefaultKeys()
	keys.Viewer.ToggleTask = []string{"x"}
	keys.Viewer.Open = []string{"o"}
	m := Model{mode: ModeNormal, artifactSelection: artifactSelection{ArtifactID: "tasks"}, keyMap: keys, readOnly: true}

	help := m.renderHelpBar()
	if strings.Contains(help, "x: toggle") || strings.Contains(help, "o: edit") {
		t.Fatalf("read-only help should omit mutation actions: %q", help)
	}
}

func TestFilterUsesCustomEditingBindings(t *testing.T) {
	keys := settings.DefaultKeys()
	keys.Filter.Accept = []string{"ctrl+j"}
	m := Model{
		mode:   ModeIndex,
		keyMap: keys,
		index: indexState{
			FilterActive:  true,
			FilterText:    "abc",
			ExpandedSpecs: make(map[int]bool),
		},
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	if !m.index.FilterActive {
		t.Fatal("replaced default Enter binding should not accept filter")
	}

	result, _ = m.dispatchKey(tea.KeyPressMsg{Text: "ctrl+j"})
	m = result.(Model)
	if m.index.FilterActive {
		t.Fatal("custom ctrl+j binding should accept filter")
	}
}
