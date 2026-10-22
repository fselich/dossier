package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

func dynamicSpecializedChange(path string) openspec.Change {
	return openspec.Change{
		Name: "custom-work", Path: path, Schema: "custom",
		Artifacts: []openspec.ChangeArtifact{
			{ID: "diagnosis", Outputs: []openspec.ArtifactOutput{{RelativePath: "diagnosis.md", DisplayName: "diagnosis", Content: "# Diagnosis", Present: true}}},
			{ID: "specs", Outputs: []openspec.ArtifactOutput{
				{RelativePath: "specs/auth/spec.md", DisplayName: "auth", Content: "# Auth", Present: true},
				{RelativePath: "specs/cache/spec.md", DisplayName: "cache", Content: "# Cache", Present: true},
			}},
			{ID: "tasks", Outputs: []openspec.ArtifactOutput{{RelativePath: "planning/checklist.md", DisplayName: "checklist", Content: "- [ ] dynamic task\n", Present: true}}},
		},
	}
}

func TestDynamicTasksUseSelectedArtifactOutput(t *testing.T) {
	root := t.TempDir()
	changePath := filepath.Join(root, "openspec", "changes", "custom-work")
	taskPath := filepath.Join(changePath, "planning", "checklist.md")
	if err := os.MkdirAll(filepath.Dir(taskPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(taskPath, []byte("- [ ] dynamic task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	change := dynamicSpecializedChange(changePath)
	m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}, loader: openspec.NewLoader(openspec.OSFS{}), mode: ModeNormal}
	m.selectArtifactOutput("tasks", "planning/checklist.md")
	m.loadTaskItems()
	if len(m.tasks.Items) != 1 || m.tasks.Items[0].Text != "dynamic task" {
		t.Fatalf("dynamic tasks not loaded: %#v", m.tasks.Items)
	}
	_ = m.doToggle()
	got, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "- [x] dynamic task") {
		t.Fatalf("dynamic task output not toggled: %q", got)
	}
}

func TestDynamicSpecSubnavigationUsesArtifactOutputs(t *testing.T) {
	change := dynamicSpecializedChange(t.TempDir())
	m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}, mode: ModeNormal}
	m.selectArtifactOutput("specs", "specs/cache/spec.md")
	if !m.hasSpecSubnav() {
		t.Fatal("expected multi-output specs subnavigation")
	}
	bar := m.renderSpecSubnav()
	if !strings.Contains(bar, "auth") || !strings.Contains(bar, "cache") {
		t.Fatalf("dynamic spec outputs missing: %q", bar)
	}
}

func TestDynamicEditorPathUsesExactSafeOutput(t *testing.T) {
	root := t.TempDir()
	changePath := filepath.Join(root, "openspec", "changes", "custom-work")
	path := filepath.Join(changePath, "diagnosis.md")
	if err := os.MkdirAll(changePath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("diagnosis"), 0o644); err != nil {
		t.Fatal(err)
	}
	change := dynamicSpecializedChange(changePath)
	m := Model{project: &openspec.Project{Changes: []openspec.Change{change}}, mode: ModeNormal}
	m.selectArtifactOutput("diagnosis", "diagnosis.md")
	got, ok := m.selectedArtifactFilePath()
	if !ok || got != path {
		t.Fatalf("selected editor path = %q,%v want %q", got, ok, path)
	}
}

func TestSingleActivePathSchedulesOnlySelectedStatus(t *testing.T) {
	change := dynamicSpecializedChange(t.TempDir())
	client := &trackingStatusClient{}
	m := Model{
		root: change.Path, mode: ModeNormal, singlePath: true,
		project: &openspec.Project{Changes: []openspec.Change{change}}, openSpec: client,
		discoveryFingerprints: discoveryFingerprints([]openspec.Change{change}),
		enrichedFingerprints:  make(map[string]string), pendingEnrichments: make(map[string]string), enrichmentRetryAfter: make(map[string]time.Time),
	}
	if cmd := m.scheduleStatusEnrichment(); cmd == nil {
		t.Fatal("explicit active path did not schedule selected status enrichment")
	}
}

func TestSingleArchivedPathPollsOnlySelectedOutputs(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-custom-work")
	if err := os.MkdirAll(archivePath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archivePath, ".openspec.yaml"), []byte("schema: spike\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archivePath, "research.md"), []byte("# Research"), 0o644); err != nil {
		t.Fatal(err)
	}
	loader := openspec.NewLoader(openspec.OSFS{})
	project, err := loader.LoadFromPath(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	m := NewSinglePath(project, openspec.ProjectConfig{}, archivePath, loader, Theme{}, settings.DefaultKeys(), false)
	if err := os.WriteFile(filepath.Join(archivePath, "decision.md"), []byte("# Decision"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = m.handleTick()
	if _, ok := m.current().ArtifactByID("decision"); !ok {
		t.Fatalf("explicit archived path did not refresh selected outputs: %#v", m.current().Artifacts)
	}
}

func TestSingleArchivedPathStartsReadOnlyDynamicArchiveViewer(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-custom-work")
	change := dynamicSpecializedChange(archivePath)
	project := &openspec.Project{Name: "project", Changes: []openspec.Change{change}}
	m := NewSinglePath(project, openspec.ProjectConfig{}, archivePath, openspec.NewLoader(openspec.OSFS{}), Theme{}, settings.DefaultKeys(), false)
	if m.mode != ModeViewingArchive || m.current() == nil || m.current().Name != "custom-work" {
		t.Fatalf("archived path mode/current = %d/%#v", m.mode, m.current())
	}
	before := m.artifactSelection
	result, cmd := m.updateViewer(tea.KeyPressMsg{Text: "e"})
	updated := result.(Model)
	if cmd != nil || updated.artifactSelection != before {
		t.Fatal("archived explicit path allowed editor mutation")
	}
}
