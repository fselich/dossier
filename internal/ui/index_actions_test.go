package ui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

type fakeUIArchiveClient struct {
	run func()
}

func (c *fakeUIArchiveClient) Archive(_ string, name string) (openspec.ArchiveResult, error) {
	if c.run != nil {
		c.run()
	}
	return openspec.ArchiveResult{ChangeName: name}, nil
}

type fakeActionClient struct {
	createName   string
	createSchema string
	validateName string
	validateKind string
	createErr    error
	validateErr  error
	createRun    func(name, schema string)
}

func (c *fakeActionClient) Status(_ string, name string) (openspec.ChangeStatus, error) {
	return openspec.ChangeStatus{ChangeName: name}, nil
}
func (c *fakeActionClient) Schemas(_ string) ([]openspec.SchemaInfo, error) { return nil, nil }
func (c *fakeActionClient) CreateChange(_ string, name, schema string) (openspec.CreateChangeResult, error) {
	c.createName, c.createSchema = name, schema
	if c.createRun != nil {
		c.createRun(name, schema)
	}
	return openspec.CreateChangeResult{Name: name, Schema: schema}, c.createErr
}
func (c *fakeActionClient) Validate(_ string, name, kind string) (openspec.ValidationResult, error) {
	c.validateName, c.validateKind = name, kind
	return openspec.ValidationResult{Valid: c.validateErr == nil, Summary: "valid"}, c.validateErr
}

func actionTestModel(root string, project *openspec.Project) Model {
	m := Model{
		root: root, loader: openspec.NewLoader(openspec.OSFS{}), project: project,
		mode: ModeIndex, keyMap: settings.DefaultKeys(), openSpec: &fakeActionClient{},
		index: indexState{ExpandedSpecs: make(map[int]bool), ExpandedChanges: make(map[string]bool), ExpandedArtifacts: make(map[string]bool)},
	}
	m.buildIndexItems()
	return m
}

func TestIndexNewChangePromptsAndUsesExplicitSchema(t *testing.T) {
	root := t.TempDir()
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	client := &fakeActionClient{}
	m.openSpec = client
	m.schemaCatalog = []openspec.SchemaInfo{{Name: "bugfix"}}

	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "n"})
	m = result.(Model)
	if m.index.Action.Mode != indexActionName || m.index.Action.Schema != "bugfix" {
		t.Fatalf("unexpected new action state: %+v", m.index.Action)
	}
	for _, r := range "fix-cache" {
		result, _ = m.updateIndex(tea.KeyPressMsg{Text: string(r)})
		m = result.(Model)
	}
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	if cmd == nil || m.index.Action.Mode != indexActionPending {
		t.Fatalf("expected pending create, state=%+v", m.index.Action)
	}
	msg := cmd()
	if client.createName != "fix-cache" || client.createSchema != "bugfix" {
		t.Fatalf("create called with name=%q schema=%q", client.createName, client.createSchema)
	}
	if _, ok := msg.(indexActionResultMsg); !ok {
		t.Fatalf("unexpected result message %T", msg)
	}
}

func TestIndexNewChangeChoosesFilteredSchemaAndSelectsCreatedRow(t *testing.T) {
	root := t.TempDir()
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	client := &fakeActionClient{}
	client.createRun = func(name, schema string) {
		dir := filepath.Join(root, "openspec", "changes", name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("schema: "+schema+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	m.openSpec = client
	m.schemaCatalog = []openspec.SchemaInfo{
		{Name: "feature", Description: "Feature workflow", Source: "project", Artifacts: []string{"proposal", "tasks"}},
		{Name: "bugfix", Description: "Bug workflow", Source: "user", Artifacts: []string{"diagnosis", "tasks"}},
	}
	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "n"})
	m = result.(Model)
	for _, r := range "bug" {
		result, _ = m.updateIndex(tea.KeyPressMsg{Text: string(r)})
		m = result.(Model)
	}
	result, _ = m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	if m.index.Action.Mode != indexActionName || m.index.Action.Schema != "bugfix" {
		t.Fatalf("filtered schema not selected: %+v", m.index.Action)
	}
	for _, r := range "fix-cache" {
		result, _ = m.updateIndex(tea.KeyPressMsg{Text: string(r)})
		m = result.(Model)
	}
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	updatedModel, _ := m.Update(cmd())
	m = updatedModel.(Model)
	if client.createSchema != "bugfix" || m.selectedIndexIdentity() != "active:fix-cache" {
		t.Fatalf("create schema/selection = %q/%q", client.createSchema, m.selectedIndexIdentity())
	}
}

func TestIndexNewChangeReportsCLIRejectionWithoutPartialRow(t *testing.T) {
	root := t.TempDir()
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	client := &fakeActionClient{createErr: errors.New("name must be kebab-case")}
	m.openSpec = client
	m.schemaCatalog = []openspec.SchemaInfo{{Name: "feature"}}
	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "n"})
	m = result.(Model)
	for _, r := range "Bad Name" {
		result, _ = m.updateIndex(tea.KeyPressMsg{Text: string(r)})
		m = result.(Model)
	}
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	updatedModel, _ := m.Update(cmd())
	m = updatedModel.(Model)
	if !m.index.Action.StatusError || !strings.Contains(m.index.Action.Status, "kebab-case") || indexItemByIdentity(m.index.Items, "active:Bad Name") >= 0 {
		t.Fatalf("CLI rejection not preserved without partial row: %+v", m.index.Action)
	}
}

func TestIndexActionPromptCancelIsIsolated(t *testing.T) {
	m := actionTestModel(t.TempDir(), &openspec.Project{Name: "project"})
	m.schemaCatalog = []openspec.SchemaInfo{{Name: "feature"}}
	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "n"})
	m = result.(Model)
	before := m.index.Cursor
	result, _ = m.updateIndex(tea.KeyPressMsg{Text: "j"})
	m = result.(Model)
	if m.index.Cursor != before || m.index.Action.Name != "j" {
		t.Fatalf("ordinary prompt key leaked to navigation: %+v", m.index.Action)
	}
	result, _ = m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = result.(Model)
	if m.index.Action.Mode != indexActionIdle {
		t.Fatalf("escape did not cancel: %+v", m.index.Action)
	}
}

func TestIndexValidateRequirementUsesContainingSpec(t *testing.T) {
	m := actionTestModel(t.TempDir(), &openspec.Project{Name: "project"})
	client := &fakeActionClient{}
	m.openSpec = client
	m.projectSpecs = []openspec.ProjectSpec{{Name: "cache", RequirementNames: []string{"Freshness"}, RequirementCount: 1}}
	m.index.ExpandedSpecs[0] = true
	m.buildIndexItems()
	m.index.Cursor = indexItemByIdentity(m.index.Items, "canonical-spec:cache:requirement:Freshness")
	result, cmd := m.updateIndex(tea.KeyPressMsg{Text: "v"})
	m = result.(Model)
	if cmd == nil || m.index.Action.Mode != indexActionPending {
		t.Fatal("expected pending validation")
	}
	_ = cmd()
	if client.validateName != "cache" || client.validateKind != "spec" {
		t.Fatalf("unexpected validation target %q/%q", client.validateName, client.validateKind)
	}
}

func TestIndexReadOnlyBlocksMutationsButAllowsValidation(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	m := actionTestModel(t.TempDir(), &openspec.Project{Changes: []openspec.Change{change}})
	m.readOnly = true
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")
	for _, key := range []string{"n", "a", "e", "u"} {
		result, cmd := m.updateIndex(tea.KeyPressMsg{Text: key})
		updated := result.(Model)
		if cmd != nil || updated.index.Action.Mode != indexActionIdle {
			t.Fatalf("read-only key %q started mutation", key)
		}
	}
	result, cmd := m.updateIndex(tea.KeyPressMsg{Text: "v"})
	if cmd == nil || result.(Model).index.Action.Mode != indexActionPending {
		t.Fatal("read-only validation should remain available")
	}
}

func TestIndexEditPathCapabilities(t *testing.T) {
	root := t.TempDir()
	changePath := filepath.Join(root, "openspec", "changes", "fix-cache")
	outputPath := filepath.Join(changePath, "diagnosis.md")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outputPath, []byte("diagnosis"), 0o644); err != nil {
		t.Fatal(err)
	}
	change := openspec.Change{Name: "fix-cache", Path: changePath, Artifacts: []openspec.ChangeArtifact{{
		ID: "diagnosis", Outputs: []openspec.ArtifactOutput{{RelativePath: "diagnosis.md", Present: true}},
	}}}
	m := actionTestModel(root, &openspec.Project{Changes: []openspec.Change{change}})
	m.index.ExpandedChanges["active:fix-cache"] = true
	m.buildIndexItems()
	item := m.index.Items[indexItemByIdentity(m.index.Items, "active:fix-cache:artifact:diagnosis")]
	path, ok := m.indexEditPath(item)
	if !ok || path != outputPath {
		t.Fatalf("edit path = %q,%v want %q", path, ok, outputPath)
	}
	archived := item
	archived.archived = true
	if _, ok := m.indexEditPath(archived); ok {
		t.Fatal("archived output must not be editable")
	}
}

func TestIndexEditPathResolvesCanonicalSpecAndRejectsUnsafeOutput(t *testing.T) {
	root := t.TempDir()
	specPath := filepath.Join(root, "openspec", "specs", "cache", "spec.md")
	if err := os.MkdirAll(filepath.Dir(specPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(specPath, []byte("### Requirement: Freshness\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	changePath := filepath.Join(root, "openspec", "changes", "fix-cache")
	if err := os.MkdirAll(changePath, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(root, "outside.md")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(changePath, "linked.md")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	change := openspec.Change{Name: "fix-cache", Path: changePath, Artifacts: []openspec.ChangeArtifact{{
		ID: "linked", Outputs: []openspec.ArtifactOutput{{RelativePath: "linked.md", Present: true}},
	}, {
		ID: "missing", Outputs: []openspec.ArtifactOutput{{RelativePath: "missing.md", Present: false}},
	}}}
	m := actionTestModel(root, &openspec.Project{Changes: []openspec.Change{change}})
	m.projectSpecs = []openspec.ProjectSpec{{Name: "cache", RequirementNames: []string{"Freshness"}, RequirementCount: 1}}
	m.buildIndexItems()
	specItem := m.index.Items[indexItemByIdentity(m.index.Items, "canonical-spec:cache")]
	if path, ok := m.indexEditPath(specItem); !ok || path != specPath {
		t.Fatalf("canonical edit path = %q,%v", path, ok)
	}
	m.index.ExpandedChanges["active:fix-cache"] = true
	m.buildIndexItems()
	for _, identity := range []string{"active:fix-cache:artifact:linked", "active:fix-cache:artifact:missing"} {
		item := m.index.Items[indexItemByIdentity(m.index.Items, identity)]
		if path, ok := m.indexEditPath(item); ok {
			t.Fatalf("unsafe/missing item %s resolved to %q", identity, path)
		}
	}
}

func TestIndexEditorReturnReloadPreservesRequirementIdentity(t *testing.T) {
	root := t.TempDir()
	specPath := filepath.Join(root, "openspec", "specs", "cache", "spec.md")
	if err := os.MkdirAll(filepath.Dir(specPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(specPath, []byte("### Requirement: Freshness\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "openspec", "changes"), 0o755); err != nil {
		t.Fatal(err)
	}
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	m.projectSpecs = []openspec.ProjectSpec{{Name: "cache", RequirementNames: []string{"Freshness"}, RequirementCount: 1}}
	m.index.ExpandedSpecs[0] = true
	m.buildIndexItems()
	identity := "canonical-spec:cache:requirement:Freshness"
	m.index.Cursor = indexItemByIdentity(m.index.Items, identity)
	updatedModel, _ := m.Update(editorReturnMsg{indexIdentity: identity})
	updated := updatedModel.(Model)
	if updated.selectedIndexIdentity() != identity {
		t.Fatalf("editor reload selection = %q, want %q", updated.selectedIndexIdentity(), identity)
	}
}

func TestIndexLifecycleConfirmationAndDuplicateSuppression(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	m := actionTestModel(t.TempDir(), &openspec.Project{Changes: []openspec.Change{change}})
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")
	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "a"})
	m = result.(Model)
	if m.index.Action.Mode != indexActionConfirmArchive {
		t.Fatalf("expected archive confirmation, got %+v", m.index.Action)
	}
	m.lifecycle = openspec.NewLifecycleService(m.root, nil)
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	if cmd == nil || m.index.Action.Mode != indexActionPending {
		t.Fatal("expected one pending lifecycle operation")
	}
	result, duplicate := m.updateIndex(tea.KeyPressMsg{Text: "a"})
	if duplicate != nil || result.(Model).index.Action.Mode != indexActionPending {
		t.Fatal("duplicate lifecycle input was not suppressed")
	}
}

func TestIndexArchiveConfirmationRendersSummaryPendingAndMissingCLI(t *testing.T) {
	root := t.TempDir()
	change := hierarchyChange("fix-cache", "bugfix")
	change.Artifacts[2].Outputs = []openspec.ArtifactOutput{{RelativePath: "tasks.md", Content: "- [x] done\n- [ ] pending\n", Present: true}}
	change.Artifacts = append(change.Artifacts, openspec.ChangeArtifact{ID: "specs-extra", Outputs: []openspec.ArtifactOutput{{RelativePath: "specs/cache/spec.md", Present: true}}})
	m := actionTestModel(root, &openspec.Project{Changes: []openspec.Change{change}})
	m.lifecycle = openspec.NewLifecycleService(root, nil)
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")
	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "a"})
	m = result.(Model)
	help := m.renderHelpBar()
	for _, expected := range []string{"tasks 1/2", "1 incomplete", "cache", "a/Enter: confirm"} {
		if !strings.Contains(help, expected) {
			t.Fatalf("archive summary missing %q: %q", expected, help)
		}
	}
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	if !strings.Contains(m.renderHelpBar(), "archiving change in progress") {
		t.Fatalf("pending progress missing: %q", m.renderHelpBar())
	}
	updatedModel, _ := m.Update(cmd())
	updated := updatedModel.(Model)
	if !updated.index.Action.StatusError || !strings.Contains(updated.index.Action.Status, "CLI is required") {
		t.Fatalf("missing CLI error not persistent: %+v", updated.index.Action)
	}
}

func TestIndexActionRenderingUsesCustomLabelsAndAvailableUndo(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	m := actionTestModel(t.TempDir(), &openspec.Project{Changes: []openspec.Change{change}})
	m.keyMap.Index.New = []string{"c"}
	m.keyMap.Index.Lifecycle = []string{"x"}
	m.keyMap.Index.Validate = []string{"z"}
	m.keyMap.Index.Undo = []string{"r"}
	m.lifecycleUndo = &openspec.LifecycleUndo{Kind: openspec.LifecycleReactivate, ChangeName: "old-change"}
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")
	help := m.renderHelpBar()
	for _, expected := range []string{"c: new", "x: archive", "z: validate", "r: undo"} {
		if !strings.Contains(help, expected) {
			t.Fatalf("custom action help missing %q: %q", expected, help)
		}
	}
	for _, unexpected := range []string{"n: new", "a: archive", "v: validate", "u: undo"} {
		if strings.Contains(help, unexpected) {
			t.Fatalf("default action help remains %q: %q", unexpected, help)
		}
	}
}

func TestIndexReactivationAndUndoTransferSelection(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	if err := os.MkdirAll(archive, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archive, ".openspec.yaml"), []byte("schema: bugfix\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archive, "proposal.md"), []byte("proposal"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	archives, err := m.loader.ListArchiveChangesFrom(root)
	if err != nil {
		t.Fatal(err)
	}
	m.index.ArchiveChanges = archives
	m.lifecycle = openspec.NewLifecycleService(root, nil)
	m.buildIndexItems()
	m.index.Cursor = indexItemByIdentity(m.index.Items, "archive:fix-cache")

	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "a"})
	m = result.(Model)
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	msg := cmd().(indexActionResultMsg)
	updatedModel, _ := m.Update(msg)
	m = updatedModel.(Model)
	if m.selectedIndexIdentity() != "active:fix-cache" || m.lifecycleUndo == nil {
		t.Fatalf("reactivation did not transfer selection or retain undo: identity=%q undo=%v", m.selectedIndexIdentity(), m.lifecycleUndo)
	}

	result, _ = m.updateIndex(tea.KeyPressMsg{Text: "u"})
	m = result.(Model)
	result, cmd = m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	updatedModel, _ = m.Update(cmd())
	m = updatedModel.(Model)
	if m.selectedIndexIdentity() != "archive:fix-cache" || m.lifecycleUndo != nil {
		t.Fatalf("undo did not transfer selection or clear record: identity=%q undo=%v", m.selectedIndexIdentity(), m.lifecycleUndo)
	}
}

func TestIndexArchiveTransfersSelectionAndRetainsUndo(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	if err := os.MkdirAll(active, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(active, ".openspec.yaml"), []byte("schema: bugfix\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(active, "proposal.md"), []byte("proposal"), 0o644); err != nil {
		t.Fatal(err)
	}
	project, err := openspec.NewLoader(openspec.OSFS{}).LoadFrom(root)
	if err != nil {
		t.Fatal(err)
	}
	m := actionTestModel(root, project)
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeUIArchiveClient{}
	client.run = func() {
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	m.lifecycle = openspec.NewLifecycleService(root, client)
	m.index.Cursor = indexItemByIdentity(m.index.Items, "active:fix-cache")
	result, _ := m.updateIndex(tea.KeyPressMsg{Text: "a"})
	m = result.(Model)
	result, cmd := m.updateIndex(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	msg := cmd().(indexActionResultMsg)
	updatedModel, _ := m.Update(msg)
	m = updatedModel.(Model)
	if m.selectedIndexIdentity() != "archive:fix-cache" || m.lifecycleUndo == nil {
		t.Fatalf("archive did not transfer selection or retain undo: identity=%q undo=%v status=%q", m.selectedIndexIdentity(), m.lifecycleUndo, m.index.Action.Status)
	}
}

func makeArchiveUndoRecord(t *testing.T) (string, *openspec.LifecycleService, *openspec.LifecycleUndo) {
	t.Helper()
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	if err := os.MkdirAll(active, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(active, "proposal.md"), []byte("proposal"), 0o644); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeUIArchiveClient{}
	client.run = func() {
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := openspec.NewLifecycleService(root, client)
	record, _, err := service.Archive(openspec.Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	return root, service, record
}

func TestIndexSuccessfulLifecycleReplacesAndDisposesOlderUndo(t *testing.T) {
	root, service, oldRecord := makeArchiveUndoRecord(t)
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	m.lifecycle = service
	m.lifecycleUndo = oldRecord
	newRecord := &openspec.LifecycleUndo{Kind: openspec.LifecycleReactivate, ChangeName: "other-change"}
	m.applyIndexActionResult(indexActionResultMsg{Kind: indexOperationReactivate, Name: "other-change", TargetIdentity: "active:other-change", Undo: newRecord})
	if m.lifecycleUndo != newRecord {
		t.Fatal("new lifecycle record did not replace older record")
	}
	if err := service.Undo(oldRecord); err == nil || !strings.Contains(err.Error(), "backup unavailable") {
		t.Fatalf("older archive backup was not disposed: %v", err)
	}
}

func TestIndexQuitDisposesSessionArchiveBackup(t *testing.T) {
	root, service, record := makeArchiveUndoRecord(t)
	m := actionTestModel(root, &openspec.Project{Name: "project"})
	m.lifecycle = service
	m.lifecycleUndo = record
	_, cmd := m.updateIndex(tea.KeyPressMsg{Text: "q"})
	if cmd == nil {
		t.Fatal("expected quit command")
	}
	if err := service.Undo(record); err == nil || !strings.Contains(err.Error(), "backup unavailable") {
		t.Fatalf("quit did not dispose session backup: %v", err)
	}
}

func TestIndexContextualHelpOmitsInapplicableActions(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	archived := hierarchyChange("old-change", "feature")
	m := actionTestModel(t.TempDir(), &openspec.Project{Changes: []openspec.Change{change}})
	m.index.ArchiveChanges = []openspec.Change{archived}
	m.buildIndexItems()

	m.index.Cursor = indexItemByIdentity(m.index.Items, "section:active-work")
	help := m.renderHelpBar()
	if !strings.Contains(help, "Enter/Space: toggle") || strings.Contains(help, "n: new") || strings.Contains(help, "i: inspect") {
		t.Fatalf("section help has inapplicable actions: %q", help)
	}

	m.index.Cursor = indexItemByIdentity(m.index.Items, "archive:old-change")
	help = m.renderHelpBar()
	if !strings.Contains(help, "a: make active") || strings.Contains(help, "n: new") || strings.Contains(help, "v: validate") || strings.Contains(help, "e: edit") {
		t.Fatalf("archive help has wrong actions: %q", help)
	}
}

func TestIndexSuccessStatusClearsOnlyWhenStillCurrent(t *testing.T) {
	m := actionTestModel(t.TempDir(), &openspec.Project{Name: "project"})
	m.index.Action.Status = "created fix-cache"
	updatedModel, _ := m.Update(indexStatusClearMsg("older status"))
	m = updatedModel.(Model)
	if m.index.Action.Status != "created fix-cache" {
		t.Fatal("stale timer cleared newer status")
	}
	updatedModel, _ = m.Update(indexStatusClearMsg("created fix-cache"))
	m = updatedModel.(Model)
	if m.index.Action.Status != "" {
		t.Fatalf("current success status did not clear: %q", m.index.Action.Status)
	}
}

func TestIndexActionFailureStaysInIndex(t *testing.T) {
	m := actionTestModel(t.TempDir(), &openspec.Project{Name: "project"})
	m.index.Action = indexActionState{Mode: indexActionPending}
	updatedModel, _ := m.Update(indexActionResultMsg{Kind: indexOperationCreate, Err: errors.New("missing CLI")})
	updated := updatedModel.(Model)
	if updated.mode != ModeIndex || updated.index.Action.Mode != indexActionIdle || !strings.Contains(updated.index.Action.Status, "missing CLI") {
		t.Fatalf("failure did not remain actionable in index: %+v", updated.index.Action)
	}
}
