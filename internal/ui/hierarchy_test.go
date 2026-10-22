package ui

import (
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

func hierarchyChange(name, schema string) openspec.Change {
	return openspec.Change{
		Name:   name,
		Schema: schema,
		Artifacts: []openspec.ChangeArtifact{
			{
				ID:     "proposal",
				Status: openspec.ArtifactStatusDone,
				Source: openspec.ArtifactSourceStatus,
				Outputs: []openspec.ArtifactOutput{{
					RelativePath: "proposal.md", DisplayName: "proposal", Present: true,
				}},
			},
			{
				ID:       "specs",
				Status:   openspec.ArtifactStatusDone,
				Requires: []string{"proposal"},
				Source:   openspec.ArtifactSourceStatus,
				Outputs: []openspec.ArtifactOutput{
					{RelativePath: "specs/auth/spec.md", DisplayName: "auth", Present: true},
					{RelativePath: "specs/payments/spec.md", DisplayName: "payments", Present: true},
				},
			},
			{
				ID:       "tasks",
				Status:   openspec.ArtifactStatusReady,
				Requires: []string{"specs", "design"},
				Source:   openspec.ArtifactSourceStatus,
			},
		},
	}
}

func TestBuildSchemaAwareHierarchy(t *testing.T) {
	active := hierarchyChange("fix-cache", "bugfix")
	archived := openspec.Change{
		Name: "old-spike", Schema: "spike", DisplayDate: "01/05/2026",
		Artifacts: []openspec.ChangeArtifact{{
			ID: "research", Status: openspec.ArtifactStatusDone, Source: openspec.ArtifactSourceDiscovered,
			Outputs: []openspec.ArtifactOutput{{RelativePath: "research.md", DisplayName: "research", Present: true}},
		}},
	}
	m := &Model{
		project:      &openspec.Project{Changes: []openspec.Change{active}},
		projectSpecs: []openspec.ProjectSpec{{Name: "canonical-auth", RequirementNames: []string{"Login"}}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true, "archive:old-spike": true},
			ExpandedArtifacts: map[string]bool{"active:fix-cache:artifact:specs": true},
			ArchiveChanges:    []openspec.Change{archived},
		},
		width: 100,
	}

	m.buildIndexItems()
	want := []struct {
		kind     indexItemKind
		identity string
	}{
		{indexKindSection, "section:active-work"},
		{indexKindActive, "active:fix-cache"},
		{indexKindArtifact, "active:fix-cache:artifact:proposal"},
		{indexKindArtifact, "active:fix-cache:artifact:specs"},
		{indexKindArtifactOutput, "active:fix-cache:artifact:specs:output:specs/auth/spec.md"},
		{indexKindArtifactOutput, "active:fix-cache:artifact:specs:output:specs/payments/spec.md"},
		{indexKindArtifact, "active:fix-cache:artifact:tasks"},
		{indexKindSection, "section:canonical-specs"},
		{indexKindSpec, "canonical-spec:canonical-auth"},
		{indexKindSection, "section:history"},
		{indexKindArchived, "archive:old-spike"},
		{indexKindArtifact, "archive:old-spike:artifact:research"},
	}
	if len(m.index.Items) != len(want) {
		t.Fatalf("expected %d hierarchy rows, got %d: %#v", len(want), len(m.index.Items), m.index.Items)
	}
	for index, expected := range want {
		item := m.index.Items[index]
		if item.kind != expected.kind || item.identity != expected.identity {
			t.Errorf("row %d: expected kind=%d identity=%q, got kind=%d identity=%q", index, expected.kind, expected.identity, item.kind, item.identity)
		}
	}
}

func TestHierarchyIdentitySurvivesSiblingInsertion(t *testing.T) {
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: make(map[string]bool),
		},
	}
	m.buildIndexItems()
	identity := "active:fix-cache:artifact:specs"
	before := indexItemByIdentity(m.index.Items, identity)
	if before < 0 {
		t.Fatal("expected specs artifact before insertion")
	}

	m.project.Changes = append([]openspec.Change{hierarchyChange("add-auth", "feature")}, m.project.Changes...)
	m.buildIndexItems()
	after := indexItemByIdentity(m.index.Items, identity)
	if after < 0 {
		t.Fatal("expected specs artifact after insertion")
	}
	if m.index.Items[after].identity != m.index.Items[before+1].identity {
		// Numeric position changes; semantic identity must not.
		if m.index.Items[after].identity != identity {
			t.Fatalf("expected stable identity %q, got %q", identity, m.index.Items[after].identity)
		}
	}
}

func TestRebuildHierarchyPreservesCursorIdentity(t *testing.T) {
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: make(map[string]bool),
		},
	}
	m.buildIndexItems()
	identity := "active:fix-cache:artifact:specs"
	m.index.Cursor = indexItemByIdentity(m.index.Items, identity)
	if m.index.Cursor < 0 {
		t.Fatal("expected initial cursor identity")
	}

	m.project.Changes = append([]openspec.Change{hierarchyChange("add-auth", "feature")}, m.project.Changes...)
	m.rebuildIndexPreservingCursor()
	if m.index.Cursor < 0 || m.index.Cursor >= len(m.index.Items) {
		t.Fatalf("invalid restored cursor %d", m.index.Cursor)
	}
	if got := m.index.Items[m.index.Cursor].identity; got != identity {
		t.Fatalf("expected cursor identity %q, got %q", identity, got)
	}
}

func TestRenderHierarchyDisambiguatesAndAlignsPlanningStatus(t *testing.T) {
	change := hierarchyChange("add-index-workflow-actions", "spec-driven")
	change.Artifacts[2].Outputs = []openspec.ArtifactOutput{{RelativePath: "tasks.md", Content: "- [ ] implement behavior\n", Present: true}}
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{change}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:add-index-workflow-actions": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		theme: Theme{Styles: BuildStyles(DarkColors)},
		width: 82,
	}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	lines := strings.Split(content, "\n")
	var changeLine, artifactLine string
	for _, line := range lines {
		switch {
		case strings.Contains(line, "add-index-workflow-actions"):
			changeLine = line
		case strings.Contains(line, "proposal"):
			artifactLine = line
		}
	}
	if !strings.Contains(changeLine, "planning 2/3") || !strings.Contains(changeLine, "tasks 0/1") {
		t.Fatalf("expected separate planning and task progress, got %q", changeLine)
	}
	if !strings.Contains(changeLine, m.theme.Styles.ProgressDone.Render("[spec-driven]")) {
		t.Fatalf("expected colored schema badge, got %q", changeLine)
	}
	authoredBadge := m.theme.Styles.ProgressComplete.Bold(true).Reverse(true).Render("[authored]")
	if !strings.Contains(artifactLine, authoredBadge) {
		t.Fatalf("expected high-contrast authored planning-document badge, got %q", artifactLine)
	}
	for name, line := range map[string]string{"change": changeLine, "artifact": artifactLine} {
		if got := lipgloss.Width(line); got != m.width-2 {
			t.Errorf("%s row width = %d, want %d: %q", name, got, m.width-2, line)
		}
	}
}

func TestRenderArtifactRowsAlignDependencyColumn(t *testing.T) {
	change := hierarchyChange("add-custom-themes", "spec-driven")
	change.Artifacts = append(change.Artifacts, openspec.ChangeArtifact{
		ID: "design", Status: openspec.ArtifactStatusDone, Requires: []string{"proposal"}, Source: openspec.ArtifactSourceStatus,
	})
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{change}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:add-custom-themes": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		theme: Theme{Styles: BuildStyles(DarkColors)},
		width: 108,
	}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	var dependencyColumns []int
	stripANSI := regexp.MustCompile(`\x1b\[[0-9;:]*m`)
	for _, line := range strings.Split(content, "\n") {
		plain := stripANSI.ReplaceAllString(line, "")
		if strings.Contains(plain, "requires ") {
			dependencyColumns = append(dependencyColumns, strings.Index(plain, "requires "))
		}
	}
	if len(dependencyColumns) < 3 {
		t.Fatalf("expected at least three dependency rows, got %v:\n%s", dependencyColumns, content)
	}
	for _, column := range dependencyColumns[1:] {
		if column != dependencyColumns[0] {
			t.Fatalf("expected dependencies in one stable column, got %v:\n%s", dependencyColumns, content)
		}
	}
}

func TestRenderActiveRowsUseStableTableColumns(t *testing.T) {
	first := hierarchyChange("add-index-workflow-actions", "spec-driven")
	first.Artifacts[2].Outputs = []openspec.ArtifactOutput{{RelativePath: "tasks.md", Content: "- [ ] one\n", Present: true}}
	second := hierarchyChange("add-custom-themes", "spec-driven")
	second.Artifacts[2].Outputs = []openspec.ArtifactOutput{{RelativePath: "tasks.md", Content: strings.Repeat("- [x] done\n", 9) + "- [ ] ten\n", Present: true}}
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{first, second}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   make(map[string]bool),
			ExpandedArtifacts: make(map[string]bool),
		},
		theme: Theme{Styles: BuildStyles(DarkColors)},
		width: 108,
	}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	var schemaColumns []int
	for _, line := range strings.Split(content, "\n") {
		plain := regexp.MustCompile(`\x1b\[[0-9;:]*m`).ReplaceAllString(line, "")
		if strings.Contains(plain, "add-index-workflow-actions") || strings.Contains(plain, "add-custom-themes") {
			schemaColumns = append(schemaColumns, strings.Index(plain, "[spec-driven]"))
		}
	}
	if len(schemaColumns) != 2 || schemaColumns[0] != schemaColumns[1] {
		t.Fatalf("expected schema badges in one stable column, got %v:\n%s", schemaColumns, content)
	}
}

func TestRenderSchemaAwareHierarchy(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	change.StatusEnriched = true
	change.Diagnostic = "degraded status"
	m := &Model{
		project:      &openspec.Project{Changes: []openspec.Change{change}},
		projectSpecs: []openspec.ProjectSpec{{Name: "auth", RequirementCount: 2}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		width: 100,
	}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	for _, expected := range []string{"Active Work", "Canonical Specs", "History", "fix-cache", "bugfix", "planning 2/3", "proposal", "tasks", "ready to author", "requires specs, design", "degraded status"} {
		if !strings.Contains(content, expected) {
			t.Errorf("expected rendered hierarchy to contain %q:\n%s", expected, content)
		}
	}
}

func TestHierarchyPrimaryToggleAndInspect(t *testing.T) {
	m := Model{
		mode:    ModeIndex,
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		keyMap:  settings.DefaultKeys(),
	}
	m.buildIndexItems()
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")

	result, _ := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	if !m.index.ExpandedChanges["active:fix-cache"] || m.mode != ModeIndex {
		t.Fatal("expected Enter to expand change without leaving index")
	}

	result, _ = m.updateIndex(tea.KeyPressMsg{Text: "i"})
	m = result.(Model)
	if m.mode != ModeNormal {
		t.Fatalf("expected i to inspect active change, got mode %d", m.mode)
	}
}

func TestHierarchyClickUsesPrimaryAction(t *testing.T) {
	m := Model{
		mode:    ModeIndex,
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		keyMap:  settings.DefaultKeys(),
	}
	m.buildIndexItems()
	index := indexItemByIdentity(m.index.Items, "active:fix-cache")
	result, _ := m.clickIndexItem(index)
	m = result.(Model)
	if !m.index.ExpandedChanges["active:fix-cache"] || m.mode != ModeIndex {
		t.Fatal("expected selected change click to expand through primary action")
	}
}

func TestHierarchyContextualHelp(t *testing.T) {
	m := Model{
		mode:    ModeIndex,
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		keyMap:  settings.DefaultKeys(),
	}
	m.buildIndexItems()
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")
	help := m.renderHelpBar()
	if !strings.Contains(help, "Enter/Space: toggle") || !strings.Contains(help, "i: inspect") {
		t.Fatalf("expected expandable contextual help, got %q", help)
	}
}

func TestHierarchyFilterRevealsMatchingDescendantAndAncestors(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	change.Artifacts = append(change.Artifacts, openspec.ChangeArtifact{
		ID: "diagnosis", Status: openspec.ArtifactStatusDone, Source: openspec.ArtifactSourceStatus,
		Outputs: []openspec.ArtifactOutput{{RelativePath: "diagnosis.md", DisplayName: "diagnosis", Present: true}},
	})
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{change}},
		index: indexState{
			ExpandedChanges:   make(map[string]bool),
			ExpandedArtifacts: make(map[string]bool),
			ExpandedSpecs:     make(map[int]bool),
			FilterText:        "diagnosis",
		},
	}
	m.buildIndexItems()
	m.applyFilter()
	visible := make([]string, 0, len(m.index.FilterIndices))
	for _, rawIndex := range m.index.FilterIndices {
		visible = append(visible, m.index.Items[rawIndex].identity)
	}
	for _, expected := range []string{"section:active-work", "active:fix-cache", "active:fix-cache:artifact:diagnosis"} {
		found := false
		for _, identity := range visible {
			if identity == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected filtered hierarchy ancestor %q, got %v", expected, visible)
		}
	}
	if m.index.ExpandedChanges["active:fix-cache"] {
		t.Fatal("filter must not overwrite stored expansion state")
	}
}

func TestRenderSchemaAwareHierarchyEmptyStates(t *testing.T) {
	m := &Model{project: &openspec.Project{}, width: 80}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	for _, expected := range []string{"No active work", "No canonical specifications available", "No history"} {
		if !strings.Contains(content, expected) {
			t.Errorf("expected empty hierarchy to contain %q:\n%s", expected, content)
		}
	}
}
