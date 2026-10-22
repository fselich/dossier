package ui

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

func dynamicViewerChange() openspec.Change {
	return openspec.Change{
		Name: "fix-cache",
		Artifacts: []openspec.ChangeArtifact{
			{ID: "proposal", Status: openspec.ArtifactStatusDone, Outputs: []openspec.ArtifactOutput{{RelativePath: "proposal.md", Content: "# Proposal", Present: true}}},
			{ID: "reproduction", Status: openspec.ArtifactStatusDone, Outputs: []openspec.ArtifactOutput{{RelativePath: "reproduction.md", Content: "# Reproduction", Present: true}}},
			{ID: "diagnosis", Status: openspec.ArtifactStatusReady},
			{ID: "specs", Status: openspec.ArtifactStatusDone, Outputs: []openspec.ArtifactOutput{
				{RelativePath: "specs/cache/spec.md", DisplayName: "cache", Content: "# Cache", Present: true},
				{RelativePath: "specs/storage/spec.md", DisplayName: "storage", Content: "# Storage", Present: true},
			}},
			{ID: "tasks", Status: openspec.ArtifactStatusBlocked},
		},
	}
}

func TestDynamicViewerSupportsStandardAndCustomArtifactSets(t *testing.T) {
	tests := []struct {
		name      string
		artifacts []string
	}{
		{name: "feature", artifacts: []string{"proposal", "specs", "design", "tasks"}},
		{name: "bugfix", artifacts: []string{"proposal", "reproduction", "diagnosis", "specs", "tasks"}},
		{name: "spike", artifacts: []string{"proposal", "research", "decision"}},
		{name: "task-only", artifacts: []string{"tasks"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			change := openspec.Change{Name: test.name}
			for _, id := range test.artifacts {
				change.Artifacts = append(change.Artifacts, openspec.ChangeArtifact{ID: id, Status: openspec.ArtifactStatusReady})
			}
			m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}}
			if !m.selectDefaultArtifact() || m.artifactSelection.ArtifactID != test.artifacts[0] {
				t.Fatalf("default selection = %+v, want %q", m.artifactSelection, test.artifacts[0])
			}
		})
	}
}

func TestDynamicViewerNavigatesSchemaArtifactsIncludingUnavailable(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}}
	m.selectDefaultArtifact()
	if got := m.artifactSelection.ArtifactID; got != "proposal" {
		t.Fatalf("default artifact = %q, want proposal", got)
	}
	for _, expected := range []string{"reproduction", "diagnosis", "specs", "tasks", "proposal"} {
		m.moveArtifactSelection(1)
		if got := m.artifactSelection.ArtifactID; got != expected {
			t.Fatalf("next artifact = %q, want %q", got, expected)
		}
	}
}

func TestDynamicViewerKeyNavigationUsesSchemaPositions(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{
		mode: ModeNormal, project: &openspec.Project{Changes: []openspec.Change{change}},
		keyMap: settings.DefaultKeys(), artifactRenderCache: make(map[artifactOutputIdentity]string),
	}
	m.selectDefaultArtifact()
	result, _ := m.updateViewer(tea.KeyPressMsg{Text: "l"})
	m = result.(Model)
	if got := m.artifactSelection.ArtifactID; got != "reproduction" {
		t.Fatalf("next artifact = %q, want reproduction", got)
	}
	result, _ = m.updateViewer(tea.KeyPressMsg{Text: "3"})
	m = result.(Model)
	if got := m.artifactSelection.ArtifactID; got != "diagnosis" {
		t.Fatalf("position 3 artifact = %q, want diagnosis", got)
	}
}

func TestDynamicViewerSelectsExactMultiOutputIdentity(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}}
	if !m.selectArtifactOutput("specs", "specs/storage/spec.md") {
		t.Fatal("expected exact output selection")
	}
	artifact, output, ok := m.selectedArtifactOutput()
	if !ok || artifact.ID != "specs" || output == nil || output.RelativePath != "specs/storage/spec.md" {
		t.Fatalf("unexpected selection: artifact=%+v output=%+v ok=%v", artifact, output, ok)
	}
}

func TestDynamicViewerReconcilesOutputInsertionAndRemovalByIdentity(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}}
	m.selectArtifactOutput("specs", "specs/storage/spec.md")

	m.project.Changes[0].Artifacts[3].Outputs = append([]openspec.ArtifactOutput{
		{RelativePath: "specs/aaa/spec.md", Present: true},
	}, m.project.Changes[0].Artifacts[3].Outputs...)
	m.reconcileArtifactSelection()
	if got := m.artifactSelection.OutputPath; got != "specs/storage/spec.md" {
		t.Fatalf("selection changed after insertion: %q", got)
	}

	m.project.Changes[0].Artifacts[3].Outputs = []openspec.ArtifactOutput{{RelativePath: "specs/cache/spec.md", Present: true}}
	m.reconcileArtifactSelection()
	if got := m.artifactSelection.OutputPath; got != "specs/cache/spec.md" {
		t.Fatalf("selection after removal = %q, want remaining output", got)
	}
}

func TestDynamicViewerHeaderShowsWorkflowSchema(t *testing.T) {
	change := dynamicViewerChange()
	change.Schema = "feature"
	m := Model{mode: ModeNormal, width: 100, project: &openspec.Project{Name: "project", Changes: []openspec.Change{change}}}
	header := m.renderHeader()
	if !strings.Contains(header, "[feature]") {
		t.Fatalf("workflow schema missing from header: %q", header)
	}
}

func TestDynamicViewerTabBarUsesArtifactIDsAndStatuses(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{
		mode: ModeNormal, width: 120,
		project: &openspec.Project{Changes: []openspec.Change{change}},
	}
	m.selectArtifactOutput("diagnosis", "")
	bar := m.renderTabBar()
	for _, expected := range []string{"proposal", "reproduction", "diagnosis", "specs", "tasks"} {
		if !strings.Contains(bar, expected) {
			t.Errorf("expected dynamic tab bar to contain %q: %q", expected, bar)
		}
	}
}

func TestHierarchyInspectOpensExactDynamicOutput(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{
		mode: ModeIndex, project: &openspec.Project{Changes: []openspec.Change{change}},
		loader: openspec.NewLoader(openspec.OSFS{}), root: t.TempDir(),
		index: indexState{
			ExpandedSpecs: make(map[int]bool), ExpandedChanges: map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: map[string]bool{"active:fix-cache:artifact:specs": true},
		},
		artifactRenderCache: make(map[artifactOutputIdentity]string),
	}
	m.buildIndexItems()
	identity := "active:fix-cache:artifact:specs:output:specs/storage/spec.md"
	itemIndex := indexItemByIdentity(m.index.Items, identity)
	if itemIndex < 0 {
		t.Fatalf("missing hierarchy output %q", identity)
	}
	result, _ := m.inspectIndexItem(m.index.Items[itemIndex])
	m = result.(Model)
	if m.mode != ModeNormal || m.artifactSelection.ArtifactID != "specs" || m.artifactSelection.OutputPath != "specs/storage/spec.md" {
		t.Fatalf("unexpected inspected selection: mode=%d selection=%+v", m.mode, m.artifactSelection)
	}
	m.enterIndex()
	if got := m.selectedIndexIdentity(); got != identity {
		t.Fatalf("return selection = %q, want %q", got, identity)
	}
}

func TestDynamicViewerMouseSelectsSchemaArtifact(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{mode: ModeNormal, width: 120, project: &openspec.Project{Changes: []openspec.Change{change}}}
	m.selectDefaultArtifact()
	firstLabelWidth := lipgloss.Width("proposal"+artifactTabStatusSuffix(change.Artifacts[0])) + 2
	clickX := 1 + firstLabelWidth + 1
	result, _ := m.handleMouseClick(tea.MouseClickMsg(tea.Mouse{X: clickX, Y: 2, Button: tea.MouseLeft}))
	m = result.(Model)
	if got := m.artifactSelection.ArtifactID; got != "reproduction" {
		t.Fatalf("mouse selected %q, want reproduction", got)
	}
}

func TestMergeReloadedChangeChoosesNearestArtifactAfterRemoval(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{mode: ModeNormal, project: &openspec.Project{Changes: []openspec.Change{change}}, artifactRenderCache: make(map[artifactOutputIdentity]string)}
	m.selectArtifactOutput("specs", "specs/storage/spec.md")
	fresh := change
	fresh.Artifacts = append([]openspec.ChangeArtifact(nil), change.Artifacts[:3]...)
	fresh.Artifacts = append(fresh.Artifacts, change.Artifacts[4])
	m.mergeReloadedChange(fresh)
	if m.artifactSelection.ArtifactID != "tasks" {
		t.Fatalf("nearest surviving artifact = %+v, want tasks", m.artifactSelection)
	}
}

func TestDynamicViewerTabBarKeepsSelectedArtifactVisibleWhenOverflowing(t *testing.T) {
	change := openspec.Change{Name: "wide"}
	for _, id := range []string{"proposal", "research", "reproduction", "diagnosis", "experiments", "decision", "specs", "design", "tasks"} {
		change.Artifacts = append(change.Artifacts, openspec.ChangeArtifact{
			ID: id, Status: openspec.ArtifactStatusDone,
			Outputs: []openspec.ArtifactOutput{{RelativePath: id + ".md", Present: true}},
		})
	}
	m := Model{mode: ModeNormal, width: 52, project: &openspec.Project{Changes: []openspec.Change{change}}}
	m.selectArtifactOutput("decision", "decision.md")
	bar := m.renderTabBar()
	if !strings.Contains(bar, "decision") {
		t.Fatalf("selected artifact missing from overflowed tabs: %q", bar)
	}
	if got := lipgloss.Width(bar); got > m.width-2 {
		t.Fatalf("tab width = %d, want <= %d: %q", got, m.width-2, bar)
	}
}

func TestDynamicViewerRendersSelectedGenericOutputByIdentity(t *testing.T) {
	change := dynamicViewerChange()
	m := Model{
		mode: ModeNormal, width: 100,
		project: &openspec.Project{Changes: []openspec.Change{change}},
		vp:      viewport.New(viewport.WithWidth(98), viewport.WithHeight(20)), vpReady: true,
		artifactRenderCache: make(map[artifactOutputIdentity]string),
	}
	m.selectArtifactOutput("reproduction", "reproduction.md")
	cmd := m.loadViewportForArtifact()
	if cmd == nil {
		t.Fatal("expected asynchronous artifact render")
	}
	msg, ok := cmd().(artifactRenderedMsg)
	if !ok {
		t.Fatalf("expected artifactRenderedMsg, got %T", cmd())
	}
	if msg.key.ArtifactID != "reproduction" || msg.key.OutputPath != "reproduction.md" {
		t.Fatalf("unexpected render key: %+v", msg.key)
	}
}
