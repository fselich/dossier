package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func testLoader() *openspec.Loader {
	return openspec.NewLoader(openspec.OSFS{})
}

func TestExtractRequirement(t *testing.T) {
	raw := `# Spec

### Requirement: Login
The system SHALL authenticate users.

#### Scenario: Success
- **WHEN** valid credentials
- **THEN** user is logged in

### Requirement: Logout
The system SHALL log out users.

#### Scenario: Click logout
- **WHEN** user clicks logout
- **THEN** session ends
`

	t.Run("name found returns block", func(t *testing.T) {
		result := openspec.ExtractRequirement(raw, "Login")
		if !strings.Contains(result, "### Requirement: Login") {
			t.Error("expected requirement header in result")
		}
		if !strings.Contains(result, "The system SHALL authenticate users") {
			t.Error("expected requirement body in result")
		}
		if strings.Contains(result, "### Requirement: Logout") {
			t.Error("expected block to stop at next requirement")
		}
	})

	t.Run("name not found returns empty", func(t *testing.T) {
		result := openspec.ExtractRequirement(raw, "Nonexistent")
		if result != "" {
			t.Errorf("expected empty result, got %q", result)
		}
	})

	t.Run("last requirement in document", func(t *testing.T) {
		result := openspec.ExtractRequirement(raw, "Logout")
		if !strings.Contains(result, "### Requirement: Logout") {
			t.Error("expected requirement header")
		}
		if !strings.Contains(result, "session ends") {
			t.Error("expected full block for last requirement")
		}
	})

	t.Run("requirement with no following header", func(t *testing.T) {
		single := "### Requirement: Only\nJust one requirement."
		result := openspec.ExtractRequirement(single, "Only")
		if result != single {
			t.Errorf("expected full content, got %q", result)
		}
	})
}

func TestFirstAvailableTab(t *testing.T) {
	t.Run("change with all tabs", func(t *testing.T) {
		ch := openspec.Change{
			Proposal: openspec.Artifact{Present: true},
			Design:   openspec.Artifact{Present: true},
			Specs:    openspec.Artifact{Present: true},
			Tasks:    openspec.Artifact{Present: true},
		}
		if got := firstAvailableTab(ch); got != TabProposal {
			t.Errorf("expected TabProposal, got %d", got)
		}
	})

	t.Run("change with only proposal and tasks", func(t *testing.T) {
		ch := openspec.Change{
			Proposal: openspec.Artifact{Present: true},
			Tasks:    openspec.Artifact{Present: true},
		}
		if got := firstAvailableTab(ch); got != TabProposal {
			t.Errorf("expected TabProposal, got %d", got)
		}
	})

	t.Run("change with no artifacts", func(t *testing.T) {
		ch := openspec.Change{}
		if got := firstAvailableTab(ch); got != TabProposal {
			t.Errorf("expected TabProposal as default, got %d", got)
		}
	})
}

func TestBuildIndexItems(t *testing.T) {
	m := &Model{
		project: &openspec.Project{
			Changes: []openspec.Change{
				{Name: "feat-a"},
				{Name: "feat-b"},
			},
		},
		projectSpecs: []openspec.ProjectSpec{
			{Name: "auth", RequirementCount: 1, RequirementNames: []string{"Login"}},
		},
	}

	t.Run("with active changes specs and archived", func(t *testing.T) {
		m.index.ExpandedSpecs = make(map[int]bool)
		m.index.ArchiveChanges = []openspec.Change{
			{Name: "old-feat", DisplayDate: "01/05/2026"},
		}
		m.buildIndexItems()
		if len(m.index.Items) != 4 {
			t.Fatalf("expected 4 index items (2 active + 1 spec + 1 archive), got %d", len(m.index.Items))
		}
		if m.index.Items[0].kind != indexKindActive {
			t.Error("expected first item to be active change")
		}
		if m.index.Items[2].kind != indexKindSpec {
			t.Error("expected third item to be spec")
		}
		if m.index.Items[3].kind != indexKindArchived {
			t.Error("expected fourth item to be archived")
		}
	})

	t.Run("empty index", func(t *testing.T) {
		empty := &Model{
			project: &openspec.Project{},
		}
		empty.index.ExpandedSpecs = make(map[int]bool)
		empty.buildIndexItems()
		if len(empty.index.Items) != 0 {
			t.Errorf("expected 0 items, got %d", len(empty.index.Items))
		}
	})
}

func TestRenderTasksContent(t *testing.T) {
	t.Run("with task cursor", func(t *testing.T) {
		m := &Model{
			width: 80,
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindSection, Text: "Section 1", LineNum: 0},
					{Kind: openspec.KindTask, Text: "do thing", Done: false, LineNum: 1},
					{Kind: openspec.KindTask, Text: "another thing", Done: false, LineNum: 2},
				},
				Cursor: 1,
			},
		}
		content, cursorLine := m.renderTasksContent()
		if cursorLine == 0 {
			t.Error("expected non-zero cursor line")
		}
		if !strings.Contains(content, "▶") {
			t.Error("expected cursor indicator (▶) in content")
		}
	})

	t.Run("empty task list", func(t *testing.T) {
		m := &Model{
			width: 80,
			tasks: taskState{
				Items:  nil,
				Cursor: 0,
			},
		}
		content, _ := m.renderTasksContent()
		if content != "" {
			t.Errorf("expected empty content, got %q", content)
		}
	})
}

func TestUpdateKeyPresses(t *testing.T) {
	t.Run("q quits normal mode", func(t *testing.T) {
		m := Model{mode: ModeNormal}
		msg := tea.KeyPressMsg{Text: "q"}
		result, cmd := m.dispatchKey(msg)
		if _, ok := result.(Model); !ok {
			t.Error("expected Model result")
		}
		if cmd == nil {
			t.Error("expected quit command")
		}
	})

	t.Run("i enters config mode", func(t *testing.T) {
		m := Model{mode: ModeNormal, width: 80, height: 24}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		msg := tea.KeyPressMsg{Text: "i"}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.mode != ModeViewingConfig {
			t.Errorf("expected ModeViewingConfig, got %d", updated.mode)
		}
	})

	t.Run("a enters index mode", func(t *testing.T) {
		m := Model{mode: ModeNormal, width: 80, height: 24, project: &openspec.Project{}, loader: testLoader()}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		msg := tea.KeyPressMsg{Text: "a"}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.mode != ModeIndex {
			t.Errorf("expected ModeIndex, got %d", updated.mode)
		}
	})

	t.Run("esc in index mode quits", func(t *testing.T) {
		m := Model{mode: ModeIndex}
		msg := tea.KeyPressMsg{Code: tea.KeyEsc}
		_, cmd := m.dispatchKey(msg)
		if cmd == nil {
			t.Error("expected quit command")
		}
	})

	t.Run("j moves cursor down in index", func(t *testing.T) {
		m := Model{
			mode:  ModeIndex,
			width: 80,
			index: indexState{
				Items:  []indexItem{{kind: indexKindActive, idx: 0}, {kind: indexKindActive, idx: 1}},
				Cursor: 0,
			},
			project: &openspec.Project{},
		}
		msg := tea.KeyPressMsg{Text: "j"}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.index.Cursor != 1 {
			t.Errorf("expected cursor 1, got %d", updated.index.Cursor)
		}
	})

	t.Run("k moves cursor up in index", func(t *testing.T) {
		m := Model{
			mode:  ModeIndex,
			width: 80,
			index: indexState{
				Items:  []indexItem{{kind: indexKindActive, idx: 0}, {kind: indexKindActive, idx: 1}},
				Cursor: 1,
			},
			project: &openspec.Project{},
		}
		msg := tea.KeyPressMsg{Text: "k"}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.index.Cursor != 0 {
			t.Errorf("expected cursor 0, got %d", updated.index.Cursor)
		}
	})
}

func TestToggleTask(t *testing.T) {
	t.Run("toggle pending to done writes to disk", func(t *testing.T) {
		dir := t.TempDir()
		content := "- [ ] do thing"
		if err := os.WriteFile(filepath.Join(dir, "tasks.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		ch := openspec.Change{Name: "test", Path: dir, Tasks: openspec.Artifact{Present: true, Content: content}}
		m := &Model{
			loader:  testLoader(),
			project: &openspec.Project{Changes: []openspec.Change{ch}},
			tasks: taskState{
				Items:  []openspec.TaskItem{{Kind: openspec.KindTask, Text: "do thing", Done: false, LineNum: 0}},
				Cursor: 0,
			},
			width: 80,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		cmd := m.doToggle()
		data, _ := os.ReadFile(filepath.Join(dir, "tasks.md"))
		if !strings.Contains(string(data), "[x]") {
			t.Errorf("expected [x] in tasks.md, got: %s", string(data))
		}
		if cmd != nil {
			t.Error("expected nil command for successful toggle")
		}
	})

	t.Run("toggle on empty items returns nil", func(t *testing.T) {
		m := &Model{tasks: taskState{Items: nil, Cursor: 0}}
		cmd := m.doToggle()
		if cmd != nil {
			t.Error("expected nil command for empty items")
		}
	})
}

func TestLoadViewportDispatch(t *testing.T) {
	t.Run("not ready returns nil", func(t *testing.T) {
		m := &Model{vpReady: false}
		cmd := m.loadViewport()
		if cmd != nil {
			t.Error("expected nil cmd when vp not ready")
		}
	})

	t.Run("index mode returns nil", func(t *testing.T) {
		m := &Model{
			vpReady: true,
			mode:    ModeIndex,
			width:   80,
			project: &openspec.Project{},
			index:   indexState{ExpandedSpecs: make(map[int]bool)},
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		cmd := m.loadViewport()
		if cmd != nil {
			t.Error("expected nil cmd for index mode")
		}
	})

	t.Run("tasks tab returns nil", func(t *testing.T) {
		m := &Model{
			vpReady: true,
			mode:    ModeNormal,
			tab:     TabTasks,
			width:   80,
			project: &openspec.Project{Changes: []openspec.Change{{Name: "test", Tasks: openspec.Artifact{Present: true}}}},
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		cmd := m.loadViewport()
		if cmd != nil {
			t.Error("expected nil cmd for tasks tab")
		}
	})

	t.Run("cache hit returns nil", func(t *testing.T) {
		m := &Model{
			vpReady:     true,
			mode:        ModeNormal,
			tab:         TabProposal,
			width:       80,
			renderCache: map[Tab]string{TabProposal: "cached content"},
			project:     &openspec.Project{Changes: []openspec.Change{{Name: "test", Proposal: openspec.Artifact{Present: true, Content: "content"}}}},
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		cmd := m.loadViewport()
		if cmd != nil {
			t.Error("expected nil cmd for cache hit")
		}
	})

	t.Run("config mode returns glamour cmd", func(t *testing.T) {
		m := &Model{
			vpReady:       true,
			mode:          ModeViewingConfig,
			width:         80,
			projectConfig: openspec.ProjectConfig{Context: "test context"},
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		cmd := m.loadViewport()
		if cmd == nil {
			t.Error("expected non-nil cmd for config mode")
		}
	})
}

func TestHandleTick(t *testing.T) {
	t.Run("viewing archive returns nil", func(t *testing.T) {
		m := &Model{mode: ModeViewingArchive}
		cmd := m.handleTick()
		if cmd != nil {
			t.Error("expected nil cmd for viewing archive")
		}
	})

	t.Run("viewing spec returns nil", func(t *testing.T) {
		m := &Model{mode: ModeViewingSpec}
		cmd := m.handleTick()
		if cmd != nil {
			t.Error("expected nil cmd for viewing spec")
		}
	})

	t.Run("normal mode with no changes returns nil", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "openspec", "changes"), 0755); err != nil {
			t.Fatal(err)
		}
		m := &Model{
			mode:       ModeNormal,
			root:       dir,
			singlePath: true,
			loader:     testLoader(),
			project:    &openspec.Project{Changes: []openspec.Change{{Name: "test"}}},
		}
		cmd := m.handleTick()
		if cmd != nil {
			t.Error("expected nil cmd for normal mode with singlePath")
		}
	})
}

func TestRenderTabBar(t *testing.T) {
	t.Run("active tab highlighted", func(t *testing.T) {
		m := &Model{
			tab:   TabProposal,
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{
					Name:     "test",
					Proposal: openspec.Artifact{Present: true},
				}},
			},
			mode: ModeNormal,
		}
		m.renderCache = make(map[Tab]string)
		result := m.renderTabBar()
		if !strings.Contains(result, "proposal") {
			t.Error("expected 'proposal' in tab bar")
		}
	})

	t.Run("disabled tab shows low style", func(t *testing.T) {
		m := &Model{
			tab:   TabProposal,
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{
					Name:     "test",
					Proposal: openspec.Artifact{Present: true},
				}},
			},
			mode: ModeNormal,
		}
		m.renderCache = make(map[Tab]string)
		result := m.renderTabBar()
		if !strings.Contains(result, "proposal") {
			t.Error("expected proposal in tab bar")
		}
	})
}
