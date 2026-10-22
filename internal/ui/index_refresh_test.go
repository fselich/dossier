package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aemonge/dossier/internal/openspec"
)

// ramonaPlanStatus mimics `openspec status` for a ramona-plan change whose
// artifacts generate numbered files (1-idea.md … 4-tasks.md).
func ramonaPlanStatus(root, changeName string) openspec.ChangeStatus {
	dir := filepath.Join(root, "openspec", "changes", changeName)
	status := openspec.ChangeStatus{
		ChangeName:    changeName,
		SchemaName:    "ramona-plan",
		ChangeRoot:    dir,
		ArtifactPaths: map[string]openspec.ArtifactPathStatus{},
	}
	for _, artifact := range []struct{ id, file string }{
		{"idea", "1-idea.md"}, {"plan", "2-plan.md"},
		{"execution", "3-execution.md"}, {"tasks", "4-tasks.md"},
	} {
		absolute := filepath.Join(dir, artifact.file)
		status.ArtifactPaths[artifact.id] = openspec.ArtifactPathStatus{
			OutputPath: artifact.file, ResolvedOutputPath: absolute,
			ExistingOutputPaths: []string{absolute},
		}
		status.Artifacts = append(status.Artifacts,
			openspec.StatusArtifact{ID: artifact.id, OutputPath: artifact.file, Status: "done"})
	}
	return status
}

func writeRamonaPlanFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "openspec", "changes", "demo")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		".openspec.yaml": "schema: ramona-plan\ncreated: 2026-10-20\n",
		"1-idea.md":      "# Idea",
		"2-plan.md":      "# Plan",
		"3-execution.md": "# Execution",
		"4-tasks.md":     "# Tasks\n- [ ] Step 1 — prepare\n- [x] Step 2 — done\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func refreshTestModel(t *testing.T, root string) *Model {
	t.Helper()
	loader := openspec.NewLoader(openspec.OSFS{})
	project, err := loader.LoadFrom(root)
	if err != nil {
		t.Fatal(err)
	}
	m := &Model{
		mode:     ModeIndex,
		root:     root,
		loader:   loader,
		project:  project,
		openSpec: &stubRefreshClient{root: root},
		index: indexState{
			ExpandedSpecs:     map[int]bool{},
			ExpandedChanges:   map[string]bool{},
			ExpandedArtifacts: map[string]bool{},
		},
		vpReady:               true,
		width:                 100,
		height:                40,
		discoveryFingerprints: discoveryFingerprints(project.Changes),
		enrichedFingerprints:  map[string]string{},
		pendingEnrichments:    map[string]string{},
		enrichmentRetryAfter:  map[string]time.Time{},
		lastStatusByName:      map[string]openspec.ChangeStatus{},
		artifactRenderCache:   map[artifactOutputIdentity]string{},
	}
	return m
}

type stubRefreshClient struct {
	root      string
	calls     int
	batchUsed bool
}

func (c *stubRefreshClient) Status(root, changeName string) (openspec.ChangeStatus, error) {
	c.calls++
	return ramonaPlanStatus(root, changeName), nil
}

func (c *stubRefreshClient) StatusAll(root string) ([]openspec.ChangeStatus, error) {
	c.batchUsed = true
	return []openspec.ChangeStatus{ramonaPlanStatus(root, "demo")}, nil
}

func (c *stubRefreshClient) Schemas(string) ([]openspec.SchemaInfo, error) {
	return nil, nil
}

func indexContentHash(t *testing.T, m *Model) string {
	t.Helper()
	content, _ := m.renderIndexContent()
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func cursorOnArtifactRow(t *testing.T, m *Model, artifactID string) {
	t.Helper()
	m.index.ExpandedChanges["active:demo"] = true
	m.buildIndexItems()
	identity := "active:demo:artifact:" + artifactID
	for i, item := range m.index.Items {
		if item.identity == identity {
			m.index.Cursor = i
			return
		}
	}
	t.Fatalf("artifact row %q not found", identity)
}

// enrichInitial runs one status enrichment round synchronously, as the app
// does at startup.
func enrichInitial(t *testing.T, m *Model) {
	t.Helper()
	cmd := m.scheduleStatusEnrichment()
	if cmd == nil {
		t.Fatal("expected enrichment command")
	}
	msg, ok := cmd().(statusEnrichmentMsg)
	if !ok {
		t.Fatalf("expected statusEnrichmentMsg, got %T", cmd())
	}
	m.applyStatusEnrichment(msg)
}

func TestAdoptDiscoveredChangeKeepsSchemaIdentities(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	enrichInitial(t, m)

	if got := m.project.Changes[0].Artifacts[0].ID; got != "idea" {
		t.Fatalf("expected enriched artifact ID idea, got %q", got)
	}

	// Disk changes underneath: adoption must re-apply the cached status so
	// schema identities survive instead of flapping to filename stems.
	if err := os.WriteFile(filepath.Join(root, "openspec", "changes", "demo", "2-plan.md"),
		[]byte("# Plan v2"), 0644); err != nil {
		t.Fatal(err)
	}
	fresh := m.loader.ReloadChange(m.project.Changes[0])
	if !m.adoptDiscoveredChange(0, fresh) {
		t.Fatal("expected adoption after disk change")
	}
	ids := make([]string, 0, 4)
	for _, artifact := range m.project.Changes[0].Artifacts {
		ids = append(ids, artifact.ID)
	}
	want := []string{"idea", "plan", "execution", "tasks"}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("artifact %d: expected schema ID %q, got %q (ids: %v)", i, want[i], ids[i], ids)
		}
	}
	plan, _ := m.project.Changes[0].ArtifactByID("plan")
	if len(plan.Outputs) != 1 || plan.Outputs[0].Content != "# Plan v2" {
		t.Fatalf("adoption lost fresh disk content: %#v", plan.Outputs)
	}
	if _, enriched := m.enrichedFingerprints["demo"]; enriched {
		t.Fatal("adoption must schedule an async authoritative refresh")
	}
}

func TestStatusEnrichmentNoOpDoesNotRefresh(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	enrichInitial(t, m)
	m.artifactRenderCache[artifactOutputIdentity{ChangeName: "demo", ArtifactID: "idea", OutputPath: "1-idea.md"}] = "rendered"

	fingerprint := m.discoveryFingerprints["demo"]
	msg := statusEnrichmentMsg{Results: []statusEnrichmentResult{{
		ChangeName: "demo", Fingerprint: fingerprint,
		Status: ramonaPlanStatus(root, "demo"),
	}}}

	before := indexContentHash(t, m)
	indexRefreshed, currentChanged := m.applyStatusEnrichment(msg)
	if indexRefreshed || currentChanged {
		t.Fatalf("identical status must not refresh: index=%t current=%t", indexRefreshed, currentChanged)
	}
	if after := indexContentHash(t, m); after != before {
		t.Fatal("index content changed for an identical status result")
	}
	if len(m.artifactRenderCache) != 1 {
		t.Fatalf("render cache must survive a no-op enrichment, got %d entries", len(m.artifactRenderCache))
	}

	// The viewer must not reload markdown either.
	m.mode = ModeNormal
	m.changeIdx = 0
	updated, cmd := m.Update(msg)
	if cmd != nil {
		t.Fatal("no-op enrichment must not trigger a viewport reload")
	}
	if updated.(Model).mode != ModeNormal {
		t.Fatal("unexpected mode change")
	}
}

func TestEditCycleSingleRefreshStableIdentity(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	enrichInitial(t, m)
	m.refreshIndexViewport()
	cursorOnArtifactRow(t, m, "plan")
	selected := m.selectedIndexIdentity()
	if selected != "active:demo:artifact:plan" {
		t.Fatalf("expected cursor on plan row, got %q", selected)
	}
	baseline := indexContentHash(t, m)

	if err := os.WriteFile(filepath.Join(root, "openspec", "changes", "demo", "4-tasks.md"),
		[]byte("# Tasks\n- [ ] Step 1 — prepare\n- [x] Step 2 — done\n- [ ] Step 3 — extra\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Tick 1: discovery adoption rebuilds once with schema identities intact.
	m.pollIndexMode()
	if got := m.selectedIndexIdentity(); got != selected {
		t.Fatalf("adoption flapped cursor identity: %q -> %q", selected, got)
	}
	afterAdopt := indexContentHash(t, m)
	if afterAdopt == baseline {
		t.Fatal("expected the edit to change index content")
	}

	// Tick 2: authoritative status arrives and must be a no-op (same state).
	cmd := m.scheduleStatusEnrichment()
	if cmd == nil {
		t.Fatal("expected async authoritative refresh to be scheduled")
	}
	msg := cmd().(statusEnrichmentMsg)
	indexRefreshed, _ := m.applyStatusEnrichment(msg)
	if indexRefreshed {
		t.Fatal("identical authoritative status rebuilt the index (double refresh)")
	}
	if got := indexContentHash(t, m); got != afterAdopt {
		t.Fatal("identical authoritative status changed visible content")
	}
	if got := m.selectedIndexIdentity(); got != selected {
		t.Fatalf("cursor identity did not survive the cycle: %q", got)
	}
}

func TestRebuildIndexCursorSurvivesIdentitySwap(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	m.refreshIndexViewport()
	cursorOnArtifactRow(t, m, "1-idea") // discovered filename-stem ID
	if got := m.selectedIndexIdentity(); got != "active:demo:artifact:1-idea" {
		t.Fatalf("expected discovered row, got %q", got)
	}

	// Startup enrichment swaps discovery IDs for schema IDs.
	enriched := m.loader.EnrichChange(m.project.Changes[0], ramonaPlanStatus(root, "demo"))
	m.project.Changes[0] = enriched
	m.rebuildIndexPreservingCursor()

	if got := m.selectedIndexIdentity(); got != "active:demo:artifact:idea" {
		t.Fatalf("cursor must follow the row across the ID swap, got %q", got)
	}
}

func TestViewerModeTickDoesNotStompEnrichment(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	enrichInitial(t, m)

	// Enter the change viewer with a warm cache and rendered content.
	m.mode = ModeNormal
	m.changeIdx = 0
	m.selectDefaultArtifact()
	m.artifactRenderCache[artifactOutputIdentity{ChangeName: "demo", ArtifactID: "idea", OutputPath: "1-idea.md"}] = "rendered"

	// Steady state: a tick must not touch the enriched change at all —
	// comparing an enriched change against a discovered reload used to
	// mismatch every 500ms and re-render the viewer.
	cmd := m.pollNormalModeContent()
	if cmd != nil {
		t.Fatal("steady-state tick reloaded the viewport")
	}
	if got := m.project.Changes[0].Artifacts[0].ID; got != "idea" {
		t.Fatalf("steady-state tick stomped enrichment: %q", got)
	}
	if len(m.artifactRenderCache) != 1 {
		t.Fatalf("steady-state tick invalidated the render cache, got %d entries", len(m.artifactRenderCache))
	}

	// Real disk change: exactly one adoption, still schema-identified, with
	// fresh content — never passing through the discovered-ID state.
	if err := os.WriteFile(filepath.Join(root, "openspec", "changes", "demo", "2-plan.md"),
		[]byte("# Plan v2"), 0644); err != nil {
		t.Fatal(err)
	}
	_ = m.pollNormalModeContent()
	if got := m.project.Changes[0].Artifacts[1].ID; got != "plan" {
		t.Fatalf("disk-change reload lost schema identity: %q", got)
	}
	plan, _ := m.project.Changes[0].ArtifactByID("plan")
	if len(plan.Outputs) != 1 || plan.Outputs[0].Content != "# Plan v2" {
		t.Fatalf("disk-change reload lost fresh content: %#v", plan.Outputs)
	}
}

func TestBootstrapEnrichmentPrePaintsSchemaIdentities(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)

	m.bootstrapEnrichment()
	if got := m.project.Changes[0].Artifacts[0].ID; got != "idea" {
		t.Fatalf("expected schema identities before first paint, got %q", got)
	}
	if _, ok := m.lastStatusByName["demo"]; !ok {
		t.Fatal("bootstrap must cache the status result")
	}
	if m.enrichedFingerprints["demo"] != m.discoveryFingerprints["demo"] {
		t.Fatal("bootstrap must mark the change enriched at its discovery fingerprint")
	}
}

type failingStatusClient struct{}

func (failingStatusClient) Status(string, string) (openspec.ChangeStatus, error) {
	return openspec.ChangeStatus{}, errors.New("openspec unavailable")
}

func (failingStatusClient) Schemas(string) ([]openspec.SchemaInfo, error) {
	return nil, nil
}

func TestBootstrapEnrichmentDegradesWhenCLIFails(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	m.openSpec = failingStatusClient{}

	m.bootstrapEnrichment()
	if got := m.project.Changes[0].Artifacts[0].ID; got != "1-idea" {
		t.Fatalf("expected discovered degradation, got %q", got)
	}
	if len(m.lastStatusByName) != 0 {
		t.Fatal("failed status lookups must not be cached")
	}
}

func TestEnterIndexRestoresCursorAcrossIdentitySwap(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	enrichInitial(t, m)
	m.refreshIndexViewport()
	cursorOnArtifactRow(t, m, "idea")
	item := m.index.Items[m.rawCursorIndex()]

	// Leaving the index captures identity plus rebuild-stable anchor.
	leaving := *m
	leaving.returnIndexIdentity = item.identity
	leaving.returnIndexAnchor = m.cursorAnchor(m.rawCursorIndex())

	// While away, enrichment is dropped (e.g. degraded CLI): rows swap back
	// to discovered IDs, so the exact identity no longer exists.
	fresh := m.loader.ReloadChange(m.project.Changes[0])
	leaving.project.Changes[0] = fresh
	leaving.mode = ModeNormal
	leaving.enterIndex()

	if got := leaving.selectedIndexIdentity(); got != "active:demo:artifact:1-idea" {
		t.Fatalf("cursor must re-anchor across the ID swap, got %q", got)
	}
}

func stripAnsiText(s string) string {
	out := make([]rune, 0, len(s))
	inEscape := false
	for _, r := range s {
		switch {
		case r == '\x1b':
			inEscape = true
		case inEscape:
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
		default:
			out = append(out, r)
		}
	}
	return string(out)
}

func TestEnrichedArtifactsShowDiskFileLabels(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	enrichInitial(t, m)
	m.refreshIndexViewport()
	cursorOnArtifactRow(t, m, "idea")

	// Rows keep the schema-declared artifact ID as the primary label, with
	// the numbered disk file stem as dimmed secondary text.
	content, _ := m.renderIndexContent()
	plain := stripAnsiText(content)
	for _, want := range []string{"1-idea", "2-plan", "3-execution", "4-tasks"} {
		if !strings.Contains(plain, want) {
			t.Fatalf("index must show the disk file stem as secondary label, missing %q", want)
		}
	}
	if got := m.selectedIndexIdentity(); got != "active:demo:artifact:idea" {
		t.Fatalf("identity must stay schema-based, got %q", got)
	}

	// Viewer tab bar shows the stable schema IDs, not the disk stems.
	m.mode = ModeNormal
	m.changeIdx = 0
	m.selectDefaultArtifact()
	bar := stripAnsiText(m.renderTabBar())
	for _, want := range []string{"idea", "plan", "execution", "tasks"} {
		if !strings.Contains(bar, want) {
			t.Fatalf("tab bar missing %q: %q", want, bar)
		}
	}

	// Filtering by the disk name matches the enriched row.
	m.mode = ModeIndex
	m.index.FilterText = "1-idea"
	m.applyFilter()
	if len(m.index.FilterIndices) == 0 {
		t.Fatal("filter by disk file name must match the enriched artifact row")
	}
}

func TestMultiOutputArtifactKeepsSchemaIDLabel(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	specsDir := filepath.Join(root, "openspec", "changes", "demo", "specs", "auth")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specsDir, "spec.md"), []byte("# Delta"), 0644); err != nil {
		t.Fatal(err)
	}
	m := refreshTestModel(t, root)
	enrichInitial(t, m)

	// "tasks" resolves to one file (4-tasks.md): ID stays primary, stem secondary.
	tasks, ok := m.project.Changes[0].ArtifactByID("tasks")
	if !ok {
		t.Fatal("expected tasks artifact")
	}
	if primary, secondary := artifactLabels(tasks); primary != "tasks" || secondary != "4-tasks" {
		t.Fatalf("mono-output labels = (%q, %q), want (tasks, 4-tasks)", primary, secondary)
	}
}

func TestBootstrapUsesBatchStatusClient(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	client := m.openSpec.(*stubRefreshClient)

	m.bootstrapEnrichment()
	if !client.batchUsed {
		t.Fatal("bootstrap must prefer the single-invocation batch status pass")
	}
	if client.calls != 0 {
		t.Fatalf("batch pass must not spawn per-change lookups, got %d", client.calls)
	}
	if got := m.project.Changes[0].Artifacts[0].ID; got != "idea" {
		t.Fatalf("expected schema identities, got %q", got)
	}
}

type perChangeOnlyClient struct {
	calls int
}

func (c *perChangeOnlyClient) Status(root, changeName string) (openspec.ChangeStatus, error) {
	c.calls++
	return ramonaPlanStatus(root, changeName), nil
}

func (perChangeOnlyClient) Schemas(string) ([]openspec.SchemaInfo, error) {
	return nil, nil
}

func TestBootstrapFallsBackWithoutBatchClient(t *testing.T) {
	root := writeRamonaPlanFixture(t)
	m := refreshTestModel(t, root)
	client := &perChangeOnlyClient{}
	m.openSpec = client

	m.bootstrapEnrichment()
	if client.calls != 1 {
		t.Fatalf("expected one per-change fallback lookup, got %d", client.calls)
	}
	if got := m.project.Changes[0].Artifacts[0].ID; got != "idea" {
		t.Fatalf("fallback pass must still enrich, got %q", got)
	}
}
