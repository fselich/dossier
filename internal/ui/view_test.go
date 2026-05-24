package ui

import (
	"strings"
	"testing"

	"github.com/fselich/dossier/internal/openspec"
)

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
		result := extractRequirement(raw, "Login")
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
		result := extractRequirement(raw, "Nonexistent")
		if result != "" {
			t.Errorf("expected empty result, got %q", result)
		}
	})

	t.Run("last requirement in document", func(t *testing.T) {
		result := extractRequirement(raw, "Logout")
		if !strings.Contains(result, "### Requirement: Logout") {
			t.Error("expected requirement header")
		}
		if !strings.Contains(result, "session ends") {
			t.Error("expected full block for last requirement")
		}
	})

	t.Run("requirement with no following header", func(t *testing.T) {
		single := "### Requirement: Only\nJust one requirement."
		result := extractRequirement(single, "Only")
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
		expandedSpecs: make(map[int]bool),
		archiveChanges: []openspec.Change{
			{Name: "old-feat", DisplayDate: "01/05/2026"},
		},
		projectSpecs: []openspec.ProjectSpec{
			{Name: "auth", RequirementCount: 1, RequirementNames: []string{"Login"}},
		},
	}

	t.Run("with active changes specs and archived", func(t *testing.T) {
		m.buildIndexItems()
		if len(m.indexItems) != 4 {
			t.Fatalf("expected 4 index items (2 active + 1 spec + 1 archive), got %d", len(m.indexItems))
		}
		if m.indexItems[0].kind != indexKindActive {
			t.Error("expected first item to be active change")
		}
		if m.indexItems[2].kind != indexKindSpec {
			t.Error("expected third item to be spec")
		}
		if m.indexItems[3].kind != indexKindArchived {
			t.Error("expected fourth item to be archived")
		}
	})

	t.Run("empty index", func(t *testing.T) {
		empty := &Model{
			project:       &openspec.Project{},
			expandedSpecs: make(map[int]bool),
		}
		empty.buildIndexItems()
		if len(empty.indexItems) != 0 {
			t.Errorf("expected 0 items, got %d", len(empty.indexItems))
		}
	})
}

func TestRenderTasksContent(t *testing.T) {
	t.Run("with task cursor", func(t *testing.T) {
		m := &Model{
			width: 80,
			taskItems: []openspec.TaskItem{
				{Kind: openspec.KindSection, Text: "Section 1", LineNum: 0},
				{Kind: openspec.KindTask, Text: "do thing", Done: false, LineNum: 1},
				{Kind: openspec.KindTask, Text: "another thing", Done: false, LineNum: 2},
			},
			taskCursor: 1,
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
			width:      80,
			taskItems:  nil,
			taskCursor: 0,
		}
		content, _ := m.renderTasksContent()
		if content != "" {
			t.Errorf("expected empty content, got %q", content)
		}
	})
}
