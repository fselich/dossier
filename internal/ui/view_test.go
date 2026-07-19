package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/git"
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
		if len(m.index.Items) != 7 {
			t.Fatalf("expected 7 items (3 sections + 2 active + 1 spec + 1 archive), got %d", len(m.index.Items))
		}
		if m.index.Items[0].kind != indexKindSection {
			t.Error("expected first item to be section")
		}
		if m.index.Items[1].kind != indexKindActive {
			t.Error("expected second item to be active change")
		}
		if m.index.Items[4].kind != indexKindSpec {
			t.Error("expected fifth item to be spec")
		}
		if m.index.Items[6].kind != indexKindArchived {
			t.Error("expected seventh item to be archived")
		}
	})

	t.Run("empty index", func(t *testing.T) {
		empty := &Model{
			project: &openspec.Project{},
		}
		empty.index.ExpandedSpecs = make(map[int]bool)
		empty.buildIndexItems()
		if len(empty.index.Items) != 3 {
			t.Errorf("expected 3 section items, got %d", len(empty.index.Items))
		}
		for _, it := range empty.index.Items {
			if it.kind != indexKindSection {
				t.Error("expected all items to be section kinds")
			}
		}
	})
}

func TestCollapsedSections(t *testing.T) {
	t.Run("collapsing section hides its children from item list", func(t *testing.T) {
		m := &Model{
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "feature-a"}, {Name: "feature-b"}},
			},
			projectSpecs: []openspec.ProjectSpec{
				{Name: "auth", RequirementCount: 1, RequirementNames: []string{"Login"}},
			},
			index: indexState{
				ExpandedSpecs: make(map[int]bool),
				ArchiveChanges: []openspec.Change{
					{Name: "old-feat", DisplayDate: "01/05/2026"},
				},
			},
		}
		m.buildIndexItems()
		if len(m.index.Items) != 7 {
			t.Fatalf("expected 7 items (3 sections + 4 children), got %d", len(m.index.Items))
		}

		m.index.CollapsedSections[sectionActive] = true
		m.buildIndexItems()
		if len(m.index.Items) != 5 {
			t.Fatalf("expected 5 items after collapsing active (3 sections + 1 spec + 1 archive), got %d", len(m.index.Items))
		}
		for _, it := range m.index.Items {
			if it.kind == indexKindActive {
				t.Error("expected no active items when section is collapsed")
			}
		}
	})

	t.Run("Space on section toggles collapse state", func(t *testing.T) {
		m := Model{
			mode:  ModeIndex,
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "feature-a"}},
			},
			index: indexState{
				ExpandedSpecs: make(map[int]bool),
			},
		}
		m.buildIndexItems()
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true
		if m.index.CollapsedSections[sectionActive] {
			t.Fatal("expected section to start expanded")
		}

		m.index.Cursor = 0
		result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeySpace})
		updated := result.(Model)
		if !updated.index.CollapsedSections[sectionActive] {
			t.Error("expected section to be collapsed after Space")
		}
		if len(updated.index.Items) != 3 {
			t.Errorf("expected 3 items (3 sections only), got %d", len(updated.index.Items))
		}
	})

	t.Run("Space on spec still expands requirements", func(t *testing.T) {
		m := Model{
			mode:  ModeIndex,
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "feature-a"}},
			},
			projectSpecs: []openspec.ProjectSpec{
				{Name: "auth", RequirementCount: 2, RequirementNames: []string{"Login", "Logout"}},
			},
			index: indexState{
				ExpandedSpecs: make(map[int]bool),
			},
		}
		m.buildIndexItems()
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true

		specIdx := -1
		for i, it := range m.index.Items {
			if it.kind == indexKindSpec {
				specIdx = i
				break
			}
		}
		if specIdx < 0 {
			t.Fatal("expected a spec item in the list")
		}
		m.index.Cursor = specIdx

		result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeySpace})
		updated := result.(Model)
		if !updated.index.ExpandedSpecs[0] {
			t.Error("expected spec to be expanded after Space on spec item")
		}
		reqCount := 0
		for _, it := range updated.index.Items {
			if it.kind == indexKindRequirement {
				reqCount++
			}
		}
		if reqCount != 2 {
			t.Errorf("expected 2 requirement items, got %d", reqCount)
		}
	})

	t.Run("filter respects collapsed section", func(t *testing.T) {
		m := &Model{
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "data-export"}, {Name: "user-auth"}},
			},
			index: indexState{
				ExpandedSpecs:     make(map[int]bool),
				CollapsedSections: [3]bool{true, false, false},
			},
		}
		m.buildIndexItems()
		if len(m.index.Items) != 3 {
			t.Fatalf("expected 3 items (3 sections only), got %d", len(m.index.Items))
		}

		m.index.FilterText = "data"
		m.applyFilter()
		inactiveFound := false
		for _, fi := range m.index.FilterIndices {
			if m.index.Items[fi].kind == indexKindActive {
				inactiveFound = true
				break
			}
		}
		if inactiveFound {
			t.Error("expected no active items in filtered results when section is collapsed")
		}
	})

	t.Run("cursor navigates through sections and items", func(t *testing.T) {
		m := Model{
			mode:  ModeIndex,
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "feat-a"}, {Name: "feat-b"}},
			},
			index: indexState{
				ExpandedSpecs: make(map[int]bool),
			},
		}
		m.buildIndexItems()
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true

		expectedKinds := []indexItemKind{
			indexKindSection, indexKindActive, indexKindActive,
			indexKindSection, indexKindSection,
		}
		if len(m.index.Items) < len(expectedKinds) {
			t.Fatalf("expected at least %d items, got %d", len(expectedKinds), len(m.index.Items))
		}
		for i, ek := range expectedKinds {
			if m.index.Items[i].kind != ek {
				t.Errorf("item[%d]: expected kind %d, got %d", i, ek, m.index.Items[i].kind)
			}
		}

		positions := []int{0, 1, 2, 3, 4, 3, 2, 1, 0}
		expectedKindsAtPos := []indexItemKind{
			indexKindSection, // 0: Active section
			indexKindActive,  // 1: feat-a
			indexKindActive,  // 2: feat-b
			indexKindSection, // 3: Specs section
			indexKindSection, // 4: Archived section
			indexKindSection, // 3: Specs section (back)
			indexKindActive,  // 2: feat-b (back)
			indexKindActive,  // 1: feat-a (back)
			indexKindSection, // 0: Active section (back)
		}

		for step, expectedPos := range positions {
			if m.index.Cursor != expectedPos {
				t.Fatalf("step %d: expected cursor at %d, got %d", step, expectedPos, m.index.Cursor)
			}
			item := m.index.Items[m.visibleItemIdx(m.index.Cursor)]
			if item.kind != expectedKindsAtPos[step] {
				t.Errorf("step %d: expected kind %d at cursor, got %d", step, expectedKindsAtPos[step], item.kind)
			}
			if step < len(positions)-1 {
				nextPos := positions[step+1]
				key := "j"
				if nextPos < expectedPos {
					key = "k"
				}
				result, _ := m.dispatchKey(tea.KeyPressMsg{Text: key})
				m = result.(Model)
			}
		}
	})

	t.Run("Enter on section header does nothing", func(t *testing.T) {
		m := Model{
			mode:  ModeIndex,
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "feat-a"}},
			},
			index: indexState{
				ExpandedSpecs: make(map[int]bool),
			},
		}
		m.buildIndexItems()
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true

		m.index.Cursor = 0
		result, cmd := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEnter})
		updated := result.(Model)
		if cmd != nil {
			t.Error("expected nil cmd (no navigation)")
		}
		if updated.mode != ModeIndex {
			t.Errorf("expected mode to stay ModeIndex, got %d", updated.mode)
		}
		if updated.index.Cursor != 0 {
			t.Errorf("expected cursor to stay at 0, got %d", updated.index.Cursor)
		}
	})

	t.Run("collapsed section shows ellipsis in rendered content", func(t *testing.T) {
		m := &Model{
			width: 80,
			project: &openspec.Project{
				Changes: []openspec.Change{{Name: "feat-a"}},
			},
			index: indexState{
				ExpandedSpecs:     make(map[int]bool),
				CollapsedSections: [3]bool{true, false, false},
			},
		}
		m.buildIndexItems()
		content, _ := m.renderIndexContent()
		if !strings.Contains(content, "…") {
			t.Errorf("expected ellipsis in collapsed section header, got:\n%s", content)
		}

		m.index.CollapsedSections[sectionActive] = false
		m.buildIndexItems()
		content, _ = m.renderIndexContent()
		if strings.Contains(content, "▼") {
			t.Errorf("expected no ▼ indicator in expanded section header, got:\n%s", content)
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

	t.Run("pgdown scrolls a full page in proposal tab", func(t *testing.T) {
		m := Model{
			mode: ModeNormal,
			tab:  TabProposal,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(10))
		m.vp.SetContent(strings.Repeat("line\n", 100))
		msg := tea.KeyPressMsg{Code: tea.KeyPgDown}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.vp.YOffset() <= 1 {
			t.Errorf("expected pgdown to scroll a full page, got offset %d", updated.vp.YOffset())
		}
	})

	t.Run("pgup scrolls a full page back up", func(t *testing.T) {
		m := Model{
			mode: ModeNormal,
			tab:  TabProposal,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(10))
		m.vp.SetContent(strings.Repeat("line\n", 100))
		m.vp.PageDown()
		m.vp.PageDown()
		offsetBefore := m.vp.YOffset()
		msg := tea.KeyPressMsg{Code: tea.KeyPgUp}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.vp.YOffset() >= offsetBefore {
			t.Errorf("expected pgup to scroll back, offset before %d, after %d", offsetBefore, updated.vp.YOffset())
		}
	})

	t.Run("pgdown is a no-op on task list", func(t *testing.T) {
		m := Model{
			mode: ModeNormal,
			tab:  TabTasks,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(10))
		m.vp.SetContent(strings.Repeat("line\n", 100))
		msg := tea.KeyPressMsg{Code: tea.KeyPgDown}
		result, _ := m.dispatchKey(msg)
		updated := result.(Model)
		if updated.vp.YOffset() != 0 {
			t.Errorf("expected task tab pgdown to be a no-op, got offset %d", updated.vp.YOffset())
		}
	})
}

func TestMoveCursorOnSections(t *testing.T) {
	t.Run("moveCursorUp goes to section header", func(t *testing.T) {
		m := &Model{
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindSection, Text: "Section 1"},
					{Kind: openspec.KindTask, Text: "do thing"},
				},
				Cursor: 1,
			},
		}
		m.moveCursorUp()
		if m.tasks.Cursor != 0 {
			t.Errorf("expected cursor at section header (0), got %d", m.tasks.Cursor)
		}
	})

	t.Run("moveCursorDown goes to section header", func(t *testing.T) {
		m := &Model{
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindTask, Text: "do thing"},
					{Kind: openspec.KindSection, Text: "Section 2"},
				},
				Cursor: 0,
			},
		}
		m.moveCursorDown()
		if m.tasks.Cursor != 1 {
			t.Errorf("expected cursor at section header (1), got %d", m.tasks.Cursor)
		}
	})

	t.Run("moveCursorUp stops at first item", func(t *testing.T) {
		m := &Model{
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindSection, Text: "Section 1"},
					{Kind: openspec.KindTask, Text: "do thing"},
				},
				Cursor: 0,
			},
		}
		m.moveCursorUp()
		if m.tasks.Cursor != 0 {
			t.Errorf("expected cursor to stay at 0, got %d", m.tasks.Cursor)
		}
	})

	t.Run("moveCursorDown stops at last item", func(t *testing.T) {
		m := &Model{
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindTask, Text: "do thing"},
					{Kind: openspec.KindSection, Text: "Section 2"},
				},
				Cursor: 1,
			},
		}
		m.moveCursorDown()
		if m.tasks.Cursor != 1 {
			t.Errorf("expected cursor to stay at 1, got %d", m.tasks.Cursor)
		}
	})

	t.Run("moveCursorUp then down navigates through sections and tasks", func(t *testing.T) {
		m := &Model{
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindSection, Text: "S1"},
					{Kind: openspec.KindTask, Text: "T1"},
					{Kind: openspec.KindTask, Text: "T2"},
					{Kind: openspec.KindSection, Text: "S2"},
					{Kind: openspec.KindTask, Text: "T3"},
				},
				Cursor: 2,
			},
		}
		m.moveCursorDown() // T2 -> S2
		if m.tasks.Cursor != 3 {
			t.Errorf("expected cursor at section S2 (3), got %d", m.tasks.Cursor)
		}
		m.moveCursorDown() // S2 -> T3
		if m.tasks.Cursor != 4 {
			t.Errorf("expected cursor at T3 (4), got %d", m.tasks.Cursor)
		}
		m.moveCursorUp() // T3 -> S2
		if m.tasks.Cursor != 3 {
			t.Errorf("expected cursor back at section S2 (3), got %d", m.tasks.Cursor)
		}
		m.moveCursorUp() // S2 -> T2
		if m.tasks.Cursor != 2 {
			t.Errorf("expected cursor back at T2 (2), got %d", m.tasks.Cursor)
		}
	})
}

func TestRenderCursorOnSectionHeader(t *testing.T) {
	t.Run("cursor on section shows ▶ mark", func(t *testing.T) {
		m := &Model{
			width: 80,
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindSection, Text: "Section 1"},
					{Kind: openspec.KindTask, Text: "do thing", Done: false},
				},
				Cursor: 0,
			},
		}
		content, cursorLine := m.renderTasksContent()
		if !strings.Contains(content, "▶") {
			t.Error("expected cursor indicator (▶) in content for section header")
		}
		if cursorLine < 0 {
			t.Errorf("expected non-negative cursor line, got %d", cursorLine)
		}
	})

	t.Run("cursor on task still works", func(t *testing.T) {
		m := &Model{
			width: 80,
			tasks: taskState{
				Items: []openspec.TaskItem{
					{Kind: openspec.KindSection, Text: "Section 1"},
					{Kind: openspec.KindTask, Text: "do thing", Done: false},
				},
				Cursor: 1,
			},
		}
		content, cursorLine := m.renderTasksContent()
		if cursorLine == 0 {
			t.Error("expected non-zero cursor line for task")
		}
		if !strings.Contains(content, "▶") {
			t.Error("expected cursor indicator (▶) in content for task")
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

	t.Run("toggle on section header returns nil", func(t *testing.T) {
		m := &Model{
			tasks: taskState{
				Items:  []openspec.TaskItem{{Kind: openspec.KindSection, Text: "Section 1"}},
				Cursor: 0,
			},
		}
		cmd := m.doToggle()
		if cmd != nil {
			t.Error("expected nil command for section header")
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

func TestMatchesFilter(t *testing.T) {
	m := &Model{
		project: &openspec.Project{
			Changes: []openspec.Change{{Name: "data-export"}, {Name: "auth-module"}},
		},
		projectSpecs: []openspec.ProjectSpec{
			{Name: "mouse-navigation", RequirementNames: []string{"Wheel events", "Click select"}},
		},
		index: indexState{
			ExpandedSpecs: make(map[int]bool),
			ArchiveChanges: []openspec.Change{
				{Name: "refactor-tick", DisplayDate: "30/05/2026"},
			},
		},
	}
	m.buildIndexItems()

	t.Run("active change matches name", func(t *testing.T) {
		item := m.index.Items[1]
		if !m.matchesFilter(item, "data") {
			t.Error("expected 'data-export' to match 'data'")
		}
	})

	t.Run("active change case insensitive", func(t *testing.T) {
		item := m.index.Items[2]
		if !m.matchesFilter(item, "auth") {
			t.Error("expected 'auth-module' to match 'auth'")
		}
	})

	t.Run("active change no match", func(t *testing.T) {
		item := m.index.Items[1]
		if m.matchesFilter(item, "xyz") {
			t.Error("expected 'data-export' not to match 'xyz'")
		}
	})

	t.Run("spec matches name", func(t *testing.T) {
		item := m.index.Items[4]
		if !m.matchesFilter(item, "mouse") {
			t.Error("expected spec name to match 'mouse'")
		}
	})

	t.Run("requirement matches name", func(t *testing.T) {
		item := indexItem{kind: indexKindRequirement, idx: 0, reqIdx: 0}
		if !m.matchesFilter(item, "wheel") {
			t.Error("expected requirement to match 'wheel'")
		}
	})

	t.Run("archived change matches name", func(t *testing.T) {
		item := m.index.Items[6]
		if !m.matchesFilter(item, "tick") {
			t.Error("expected archived name to match 'tick'")
		}
	})

	t.Run("substring partial match", func(t *testing.T) {
		item := m.index.Items[1]
		if !m.matchesFilter(item, "port") {
			t.Error("expected 'data-export' to match substring 'port'")
		}
	})
}

func TestApplyFilter(t *testing.T) {
	items := []indexItem{
		{kind: indexKindActive, idx: 0},
		{kind: indexKindActive, idx: 1},
		{kind: indexKindSpec, idx: 0},
	}
	m := &Model{
		project: &openspec.Project{
			Changes: []openspec.Change{{Name: "data-export"}, {Name: "auth-module"}},
		},
		projectSpecs: []openspec.ProjectSpec{
			{Name: "data-pipeline"},
		},
		index: indexState{
			Items:         items,
			Cursor:        2,
			ExpandedSpecs: make(map[int]bool),
		},
	}

	t.Run("filter set builds FilterIndices", func(t *testing.T) {
		m.index.FilterText = "data"
		m.applyFilter()
		if m.index.FilterIndices == nil {
			t.Fatal("expected non-nil FilterIndices")
		}
		if len(m.index.FilterIndices) != 2 {
			t.Errorf("expected 2 matching items, got %d (%v)", len(m.index.FilterIndices), m.index.FilterIndices)
		}
	})

	t.Run("cursor clamped to visible count", func(t *testing.T) {
		m.index.Cursor = 5
		m.applyFilter()
		if m.index.Cursor != 0 {
			t.Errorf("expected cursor 0 after clamping, got %d", m.index.Cursor)
		}
	})

	t.Run("empty filter clears FilterIndices", func(t *testing.T) {
		m.index.FilterText = ""
		m.applyFilter()
		if m.index.FilterIndices != nil {
			t.Error("expected nil FilterIndices when FilterText is empty")
		}
	})
}

func TestIndexFilterKeypresses(t *testing.T) {
	m := Model{
		mode:  ModeIndex,
		width: 80,
		project: &openspec.Project{
			Changes: []openspec.Change{{Name: "data-export"}, {Name: "auth-module"}},
		},
		index: indexState{
			ExpandedSpecs: make(map[int]bool),
			Items: []indexItem{
				{kind: indexKindActive, idx: 0},
				{kind: indexKindActive, idx: 1},
			},
			Cursor: 0,
		},
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	t.Run("/ enters filter mode", func(t *testing.T) {
		result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "/"})
		updated := result.(Model)
		if !updated.index.FilterActive {
			t.Error("expected FilterActive after pressing /")
		}
	})

	t.Run("typing during filter mode updates FilterText", func(t *testing.T) {
		m.index.FilterActive = true
		m.index.FilterText = ""
		result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "d"})
		updated := result.(Model)
		if updated.index.FilterText != "d" {
			t.Errorf("expected FilterText 'd', got %q", updated.index.FilterText)
		}
	})

	t.Run("backspace during filter removes char", func(t *testing.T) {
		m.index.FilterActive = true
		m.index.FilterText = "da"
		result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyBackspace})
		updated := result.(Model)
		if updated.index.FilterText != "d" {
			t.Errorf("expected FilterText 'd', got %q", updated.index.FilterText)
		}
	})

	t.Run("enter in filter mode confirms", func(t *testing.T) {
		m.index.FilterActive = true
		m.index.FilterText = "data"
		result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEnter})
		updated := result.(Model)
		if updated.index.FilterActive {
			t.Error("expected FilterActive false after Enter")
		}
		if updated.index.FilterText != "data" {
			t.Errorf("expected FilterText 'data' to persist, got %q", updated.index.FilterText)
		}
	})

	t.Run("esc in filter mode cancels and reverts", func(t *testing.T) {
		m.index.FilterActive = true
		m.index.PrevFilterText = ""
		m.index.FilterText = "foo"
		result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEsc})
		updated := result.(Model)
		if updated.index.FilterActive {
			t.Error("expected FilterActive false after Esc in filter mode")
		}
		if updated.index.FilterText != "" {
			t.Errorf("expected FilterText reverted to '', got %q", updated.index.FilterText)
		}
	})

	t.Run("esc with filter clears it", func(t *testing.T) {
		m.index.FilterActive = false
		m.index.FilterText = "data"
		result, cmd := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEsc})
		updated := result.(Model)
		if cmd != nil {
			t.Error("expected nil cmd (filter cleared, not quit)")
		}
		if updated.index.FilterText != "" {
			t.Errorf("expected empty FilterText after Esc, got %q", updated.index.FilterText)
		}
	})

	t.Run("esc without filter quits", func(t *testing.T) {
		m.index.FilterActive = false
		m.index.FilterText = ""
		_, cmd := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEsc})
		if cmd == nil {
			t.Error("expected quit cmd when no filter active")
		}
	})
}

func TestIndexFilterNoMatchMessage(t *testing.T) {
	m := &Model{
		width: 80,
		project: &openspec.Project{
			Changes: []openspec.Change{{Name: "data-export"}},
		},
		index: indexState{
			ExpandedSpecs: make(map[int]bool),
			FilterText:    "nonexistent",
		},
	}
	m.buildIndexItems()
	m.applyFilter()

	content, _ := m.renderIndexContent()
	if !strings.Contains(content, "No items match 'nonexistent'") {
		t.Errorf("expected no-match message in filtered content, got:\n%s", content)
	}
}

func TestIndexFilteredNavigation(t *testing.T) {
	m := Model{
		mode:  ModeIndex,
		width: 80,
		project: &openspec.Project{
			Changes: []openspec.Change{{Name: "data-export"}, {Name: "user-auth"}},
		},
		index: indexState{
			ExpandedSpecs: make(map[int]bool),
			Items: []indexItem{
				{kind: indexKindActive, idx: 0},
				{kind: indexKindActive, idx: 1},
			},
			Cursor:       0,
			FilterText:   "data",
			FilterActive: false,
		},
	}
	m.applyFilter()

	t.Run("visibleItemCount returns filtered count", func(t *testing.T) {
		if n := m.visibleItemCount(); n != 1 {
			t.Errorf("expected 1 visible item, got %d", n)
		}
	})

	t.Run("j/k navigate filtered list", func(t *testing.T) {
		// Try to move past the only visible item
		result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "j"})
		updated := result.(Model)
		if updated.index.Cursor != 0 {
			t.Errorf("expected cursor 0 (only 1 visible), got %d", updated.index.Cursor)
		}
	})

	t.Run("visibleItemIdx maps through filter", func(t *testing.T) {
		idx := m.visibleItemIdx(0)
		if idx != 0 {
			t.Errorf("expected raw index 0 (data-export), got %d", idx)
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

func TestGitCursorLandsOnDeletedFiles(t *testing.T) {
	t.Run("moveGitCursorDown lands on deleted file", func(t *testing.T) {
		m := &Model{
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "a.go", X: ' ', Y: 'M'},
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
					{Path: "b.go", X: ' ', Y: 'M'},
				},
				Cursor: 0,
			},
		}
		m.moveGitCursorDown()
		if m.gitState.Cursor != 1 {
			t.Errorf("expected cursor on deleted file (1), got %d", m.gitState.Cursor)
		}
	})

	t.Run("moveGitCursorUp lands on deleted file", func(t *testing.T) {
		m := &Model{
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "a.go", X: ' ', Y: 'M'},
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
					{Path: "b.go", X: ' ', Y: 'M'},
				},
				Cursor: 2,
			},
		}
		m.moveGitCursorUp()
		if m.gitState.Cursor != 1 {
			t.Errorf("expected cursor on deleted file (1), got %d", m.gitState.Cursor)
		}
	})

	t.Run("moveGitDiffCursorDown skips deleted files", func(t *testing.T) {
		m := &Model{
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "a.go", X: ' ', Y: 'M'},
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
					{Path: "b.go", X: ' ', Y: 'M'},
				},
				Cursor: 0,
			},
		}
		m.moveGitDiffCursorDown()
		if m.gitState.Cursor != 2 {
			t.Errorf("expected cursor to skip to b.go (2), got %d", m.gitState.Cursor)
		}
	})

	t.Run("moveGitDiffCursorUp skips deleted files", func(t *testing.T) {
		m := &Model{
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "a.go", X: ' ', Y: 'M'},
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
					{Path: "b.go", X: ' ', Y: 'M'},
				},
				Cursor: 2,
			},
		}
		m.moveGitDiffCursorUp()
		if m.gitState.Cursor != 0 {
			t.Errorf("expected cursor to skip to a.go (0), got %d", m.gitState.Cursor)
		}
	})
}

func TestGitDeletedEnterEDNoop(t *testing.T) {
	t.Run("d on deleted file does nothing", func(t *testing.T) {
		m := Model{
			mode: ModeNormal,
			tab:  TabGit,
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
				},
				Cursor: 0,
			},
			width: 80, height: 24,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true
		result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "d"})
		updated := result.(Model)
		if updated.gitState.ShowingDiff {
			t.Error("expected diff not to show for deleted file")
		}
	})

	t.Run("Enter on deleted file does nothing", func(t *testing.T) {
		m := Model{
			mode: ModeNormal,
			tab:  TabGit,
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
				},
				Cursor: 0,
			},
			width: 80, height: 24,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true
		result, _ := m.dispatchKey(tea.KeyPressMsg{Code: tea.KeyEnter})
		updated := result.(Model)
		if updated.gitState.ShowingDiff {
			t.Error("expected diff not to show for deleted file")
		}
	})

	t.Run("e on deleted file does nothing", func(t *testing.T) {
		m := Model{
			mode: ModeNormal,
			tab:  TabGit,
			gitState: gitState{
				Files: []git.FileStatus{
					{Path: "del.go", X: ' ', Y: 'D', IsDeleted: true},
				},
				Cursor: 0,
			},
			width: 80, height: 24,
		}
		m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
		m.vpReady = true
		result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "e"})
		updated := result.(Model)
		if updated.gitState.ShowingDiff {
			t.Error("expected diff not to show for deleted file")
		}
	})
}

// ── git tab s toggle tests ─────────────────────────────────────────────────────

func skipIfNoGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found on PATH")
	}
}

func gitInit(t *testing.T, dir string) {
	t.Helper()
	mustGit(t, dir, "init")
	mustGit(t, dir, "config", "user.email", "test@test")
	mustGit(t, dir, "config", "user.name", "Test")
}

func mustGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestGitSStageModified(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	gitInit(t, dir)
	mustGit(t, dir, "commit", "--allow-empty", "-m", "init")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v1"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "a.txt")
	mustGit(t, dir, "commit", "-m", "add a")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v2"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := git.Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	m := Model{
		mode:      ModeNormal,
		tab:       TabGit,
		root:      dir,
		gitRoot:   dir,
		isGitRepo: true,
		gitState: gitState{
			Files:  files,
			Cursor: 0,
		},
		width: 80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if updated.gitState.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", updated.gitState.ErrMsg)
	}
	if len(updated.gitState.Files) == 0 {
		t.Fatal("expected files after stage")
	}
	f := updated.gitState.Files[0]
	if f.X != 'M' || f.Y != ' ' {
		t.Errorf("expected staged (M ), got %c%c", f.X, f.Y)
	}
}

func TestGitSUnstageStaged(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	gitInit(t, dir)
	mustGit(t, dir, "commit", "--allow-empty", "-m", "init")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v1"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "a.txt")
	mustGit(t, dir, "commit", "-m", "add a")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v2"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "a.txt")

	files, err := git.Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	m := Model{
		mode:      ModeNormal,
		tab:       TabGit,
		root:      dir,
		gitRoot:   dir,
		isGitRepo: true,
		gitState: gitState{
			Files:  files,
			Cursor: 0,
		},
		width: 80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if updated.gitState.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", updated.gitState.ErrMsg)
	}
	if len(updated.gitState.Files) == 0 {
		t.Fatal("expected files after unstage")
	}
	f := updated.gitState.Files[0]
	if f.X != ' ' || f.Y != 'M' {
		t.Errorf("expected unstaged ( M), got %c%c", f.X, f.Y)
	}
}

func TestGitSMixedMM(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	gitInit(t, dir)
	mustGit(t, dir, "commit", "--allow-empty", "-m", "init")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v1"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "a.txt")
	mustGit(t, dir, "commit", "-m", "add a")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v2"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "a.txt")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v3"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := git.Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) == 0 || files[0].X != 'M' || files[0].Y != 'M' {
		t.Fatalf("expected MM, got %+v", files)
	}

	m := Model{
		mode:      ModeNormal,
		tab:       TabGit,
		root:      dir,
		gitRoot:   dir,
		isGitRepo: true,
		gitState: gitState{
			Files:  files,
			Cursor: 0,
		},
		width: 80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if updated.gitState.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", updated.gitState.ErrMsg)
	}
	if len(updated.gitState.Files) == 0 {
		t.Fatal("expected files after stage")
	}
	f := updated.gitState.Files[0]
	if f.X != 'M' || f.Y != ' ' {
		t.Errorf("expected `M ` after staging MM (worktree side), got %c%c", f.X, f.Y)
	}
}

func TestGitSStageDeleted(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	gitInit(t, dir)
	mustGit(t, dir, "commit", "--allow-empty", "-m", "init")
	if err := os.WriteFile(filepath.Join(dir, "del.txt"), []byte("bye"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "del.txt")
	mustGit(t, dir, "commit", "-m", "add del")
	if err := os.Remove(filepath.Join(dir, "del.txt")); err != nil {
		t.Fatal(err)
	}

	files, err := git.Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	m := Model{
		mode:      ModeNormal,
		tab:       TabGit,
		root:      dir,
		gitRoot:   dir,
		isGitRepo: true,
		gitState: gitState{
			Files:  files,
			Cursor: 0,
		},
		width: 80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if updated.gitState.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", updated.gitState.ErrMsg)
	}
	if len(updated.gitState.Files) == 0 {
		t.Fatal("expected files after stage")
	}
	f := updated.gitState.Files[0]
	if f.X != 'D' || f.Y != ' ' {
		t.Errorf("expected staged delete (D ), got %c%c", f.X, f.Y)
	}
}

func TestGitSErrorPath(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	gitInit(t, dir)
	mustGit(t, dir, "commit", "--allow-empty", "-m", "init")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v1"), 0644); err != nil {
		t.Fatal(err)
	}
	mustGit(t, dir, "add", "a.txt")
	mustGit(t, dir, "commit", "-m", "add a")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v2"), 0644); err != nil {
		t.Fatal(err)
	}

	gitDir := filepath.Join(dir, ".git")
	lockFile := filepath.Join(gitDir, "index.lock")
	if err := os.WriteFile(lockFile, []byte("lock"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(lockFile) })

	files, err := git.Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	m := Model{
		mode:      ModeNormal,
		tab:       TabGit,
		root:      dir,
		gitRoot:   dir,
		isGitRepo: true,
		gitState: gitState{
			Files:  files,
			Cursor: 0,
		},
		width: 80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if updated.gitState.ErrMsg == "" {
		t.Fatal("expected error message when index.lock is held")
	}
}

func TestGitSInactiveInDiffView(t *testing.T) {
	m := Model{
		mode: ModeNormal,
		tab:  TabGit,
		gitState: gitState{
			Files: []git.FileStatus{
				{Path: "a.txt", X: ' ', Y: 'M'},
			},
			ShowingDiff: true,
			Cursor:      0,
		},
		isGitRepo: true,
		width:     80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	beforeFiles := len(m.gitState.Files)
	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if len(updated.gitState.Files) != beforeFiles {
		t.Error("expected s to be a no-op in diff view")
	}
}

func TestGitSInactiveCleanTree(t *testing.T) {
	m := Model{
		mode: ModeNormal,
		tab:  TabGit,
		gitState: gitState{
			Files:  nil,
			Cursor: 0,
		},
		isGitRepo: true,
		width:     80, height: 24,
	}
	m.vp = viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	m.vpReady = true

	result, _ := m.dispatchKey(tea.KeyPressMsg{Text: "s"})
	updated := result.(Model)
	if updated.gitState.Files != nil {
		t.Error("expected s to be a no-op on clean tree")
	}
}

func TestGetChromaStyleValid(t *testing.T) {
	cs := getChromaStyle("monokai")
	if cs == nil {
		t.Fatal("getChromaStyle(monokai) should not return nil")
	}
}

func TestGetChromaStyleInvalidFallback(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("getChromaStyle should not panic on invalid name: %v", r)
		}
	}()
	cs := getChromaStyle("nonexistent-style")
	if cs == nil {
		t.Fatal("getChromaStyle should return Fallback, not nil, for invalid name")
	}
}

func TestEnsureRendererUsesThemeStyle(t *testing.T) {
	m := &Model{
		theme: DarkTheme,
	}
	m.ensureRenderer(80)
	if m.glamourRenderer == nil {
		t.Fatal("ensureRenderer should create a renderer for dark theme")
	}
}

func TestEnsureRendererFallsBackToDarkOnEmptyStyle(t *testing.T) {
	m := &Model{
		theme: Theme{},
	}
	m.ensureRenderer(80)
	if m.glamourRenderer == nil {
		t.Fatal("ensureRenderer should create a renderer with fallback to dark style when theme has empty GlamourStyle")
	}
}
