package ui

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aemonge/dossier/internal/openspec"
)

type trackingStatusClient struct {
	active    int32
	maxActive int32
	calls     int32
	delay     time.Duration
	schemas   []openspec.SchemaInfo
}

func (c *trackingStatusClient) Status(_ string, changeName string) (openspec.ChangeStatus, error) {
	active := atomic.AddInt32(&c.active, 1)
	defer atomic.AddInt32(&c.active, -1)
	atomic.AddInt32(&c.calls, 1)
	for {
		maxActive := atomic.LoadInt32(&c.maxActive)
		if active <= maxActive || atomic.CompareAndSwapInt32(&c.maxActive, maxActive, active) {
			break
		}
	}
	time.Sleep(c.delay)
	return openspec.ChangeStatus{
		ChangeName: changeName,
		SchemaName: "feature",
		Artifacts:  []openspec.StatusArtifact{{ID: "proposal", Status: "done"}},
	}, nil
}

func (c *trackingStatusClient) Schemas(_ string) ([]openspec.SchemaInfo, error) {
	return append([]openspec.SchemaInfo(nil), c.schemas...), nil
}

func discoveredChange(name, content string) openspec.Change {
	return openspec.Change{
		Name:   name,
		Path:   "/project/openspec/changes/" + name,
		Schema: "feature",
		Artifacts: []openspec.ChangeArtifact{{
			ID:     "proposal",
			Status: openspec.ArtifactStatusDone,
			Source: openspec.ArtifactSourceDiscovered,
			Outputs: []openspec.ArtifactOutput{{
				RelativePath: "proposal.md",
				Content:      content,
				Present:      true,
			}},
		}},
	}
}

func TestSchemaCatalogPreservesOpenSpecOrder(t *testing.T) {
	m := Model{}
	updatedModel, _ := m.Update(schemaCatalogMsg{Schemas: []openspec.SchemaInfo{
		{Name: "spike"},
		{Name: "feature"},
	}})
	updated := updatedModel.(Model)
	if len(updated.schemaCatalog) != 2 || updated.schemaCatalog[0].Name != "spike" {
		t.Fatalf("expected OpenSpec order preserved, got %#v", updated.schemaCatalog)
	}
}

func TestStatusEnrichmentIsBounded(t *testing.T) {
	changes := make([]openspec.Change, 10)
	for i := range changes {
		changes[i] = discoveredChange(fmt.Sprintf("change-%02d", i), fmt.Sprintf("content-%d", i))
	}
	client := &trackingStatusClient{delay: 5 * time.Millisecond}
	m := Model{
		root:                  "/project",
		project:               &openspec.Project{Changes: changes},
		openSpec:              client,
		discoveryFingerprints: discoveryFingerprints(changes),
		enrichedFingerprints:  make(map[string]string),
		pendingEnrichments:    make(map[string]string),
	}

	cmd := m.scheduleStatusEnrichment()
	if cmd == nil {
		t.Fatal("expected enrichment command")
	}
	msg, ok := cmd().(statusEnrichmentMsg)
	if !ok {
		t.Fatalf("expected statusEnrichmentMsg, got %T", cmd())
	}
	if len(msg.Results) != len(changes) {
		t.Fatalf("expected %d results, got %d", len(changes), len(msg.Results))
	}
	if maxActive := atomic.LoadInt32(&client.maxActive); maxActive > maxConcurrentStatus {
		t.Fatalf("expected at most %d concurrent calls, got %d", maxConcurrentStatus, maxActive)
	}
	if maxActive := atomic.LoadInt32(&client.maxActive); maxActive < 2 {
		t.Fatalf("expected concurrent enrichment, got max %d", maxActive)
	}
}

func TestStatusEnrichmentUsesFingerprintCache(t *testing.T) {
	change := discoveredChange("cached", "content")
	fingerprint := changeDiscoveryFingerprint(change)
	client := &trackingStatusClient{}
	m := Model{
		root:                  "/project",
		project:               &openspec.Project{Changes: []openspec.Change{change}},
		openSpec:              client,
		discoveryFingerprints: map[string]string{"cached": fingerprint},
		enrichedFingerprints:  map[string]string{"cached": fingerprint},
		pendingEnrichments:    make(map[string]string),
	}

	if cmd := m.scheduleStatusEnrichment(); cmd != nil {
		t.Fatal("expected cache hit to avoid command")
	}
	if calls := atomic.LoadInt32(&client.calls); calls != 0 {
		t.Fatalf("expected no calls, got %d", calls)
	}
}

func TestStatusEnrichmentIgnoresStaleResult(t *testing.T) {
	change := discoveredChange("changing", "new content")
	m := Model{
		project:               &openspec.Project{Changes: []openspec.Change{change}},
		loader:                openspec.NewLoader(openspec.OSFS{}),
		discoveryFingerprints: map[string]string{"changing": "new-fingerprint"},
		enrichedFingerprints:  make(map[string]string),
		pendingEnrichments:    map[string]string{"changing": "old-fingerprint"},
	}
	msg := statusEnrichmentMsg{Results: []statusEnrichmentResult{{
		ChangeName:  "changing",
		Fingerprint: "old-fingerprint",
		Status: openspec.ChangeStatus{
			ChangeName: "changing",
			SchemaName: "bugfix",
			Artifacts:  []openspec.StatusArtifact{{ID: "diagnosis", Status: "done"}},
		},
	}}}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)
	if updated.project.Changes[0].Schema != "feature" {
		t.Fatalf("stale result changed schema to %q", updated.project.Changes[0].Schema)
	}
	if _, ok := updated.project.Changes[0].ArtifactByID("diagnosis"); ok {
		t.Fatal("stale result added diagnosis")
	}
}

func TestStatusEnrichmentAppliesCurrentResult(t *testing.T) {
	change := discoveredChange("changing", "content")
	fingerprint := changeDiscoveryFingerprint(change)
	m := Model{
		project:               &openspec.Project{Changes: []openspec.Change{change}},
		loader:                openspec.NewLoader(openspec.OSFS{}),
		discoveryFingerprints: map[string]string{"changing": fingerprint},
		enrichedFingerprints:  make(map[string]string),
		pendingEnrichments:    map[string]string{"changing": fingerprint},
	}
	msg := statusEnrichmentMsg{Results: []statusEnrichmentResult{{
		ChangeName:  "changing",
		Fingerprint: fingerprint,
		Status: openspec.ChangeStatus{
			ChangeName: "changing",
			SchemaName: "bugfix",
			Artifacts:  []openspec.StatusArtifact{{ID: "diagnosis", Status: "ready"}},
		},
	}}}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)
	if updated.project.Changes[0].Schema != "bugfix" {
		t.Fatalf("expected enriched schema, got %q", updated.project.Changes[0].Schema)
	}
	if artifact, ok := updated.project.Changes[0].ArtifactByID("diagnosis"); !ok || artifact.Source != openspec.ArtifactSourceStatus {
		t.Fatalf("expected enriched diagnosis, got %#v", artifact)
	}
	if updated.enrichedFingerprints["changing"] != fingerprint {
		t.Fatal("expected enrichment fingerprint cache")
	}
	if _, pending := updated.pendingEnrichments["changing"]; pending {
		t.Fatal("expected pending entry cleared")
	}
}

func TestAdoptDiscoveredChangeInvalidatesOnlyChangedFingerprint(t *testing.T) {
	original := discoveredChange("change", "before")
	originalFingerprint := changeDiscoveryFingerprint(original)
	enriched := original
	enriched.Artifacts[0].Source = openspec.ArtifactSourceStatus
	m := Model{
		project:               &openspec.Project{Changes: []openspec.Change{enriched}},
		discoveryFingerprints: map[string]string{"change": originalFingerprint},
		enrichedFingerprints:  map[string]string{"change": originalFingerprint},
		pendingEnrichments:    map[string]string{"change": originalFingerprint},
		enrichmentRetryAfter:  map[string]time.Time{"change": time.Now().Add(time.Minute)},
	}

	fresh := discoveredChange("change", "after")
	if !m.adoptDiscoveredChange(0, fresh) {
		t.Fatal("expected changed discovery to be adopted")
	}
	if m.project.Changes[0].Artifacts[0].Outputs[0].Content != "after" {
		t.Fatal("expected fresh discovered content")
	}
	if _, ok := m.enrichedFingerprints["change"]; ok {
		t.Fatal("expected status cache invalidated")
	}
	if _, ok := m.pendingEnrichments["change"]; ok {
		t.Fatal("expected stale pending request invalidated")
	}
	if _, ok := m.enrichmentRetryAfter["change"]; ok {
		t.Fatal("expected retry delay invalidated")
	}
}

func TestAdoptDiscoveredChangeKeepsEnrichmentWhenFingerprintUnchanged(t *testing.T) {
	discovered := discoveredChange("change", "same")
	fingerprint := changeDiscoveryFingerprint(discovered)
	enriched := discoveredChange("change", "same")
	enriched.Artifacts[0].Source = openspec.ArtifactSourceStatus
	m := Model{
		project:               &openspec.Project{Changes: []openspec.Change{enriched}},
		discoveryFingerprints: map[string]string{"change": fingerprint},
		enrichedFingerprints:  map[string]string{"change": fingerprint},
		pendingEnrichments:    make(map[string]string),
		enrichmentRetryAfter:  make(map[string]time.Time),
	}

	if m.adoptDiscoveredChange(0, discovered) {
		t.Fatal("expected unchanged discovery to be ignored")
	}
	if m.project.Changes[0].Artifacts[0].Source != openspec.ArtifactSourceStatus {
		t.Fatal("expected authoritative enrichment retained")
	}
	if m.enrichedFingerprints["change"] != fingerprint {
		t.Fatal("expected status cache retained")
	}
}

func TestAdoptDiscoveredProjectPreservesUnchangedEnrichment(t *testing.T) {
	discovered := discoveredChange("existing", "same")
	fingerprint := changeDiscoveryFingerprint(discovered)
	enriched := discoveredChange("existing", "same")
	enriched.Artifacts[0].Source = openspec.ArtifactSourceStatus
	m := Model{
		project:               &openspec.Project{Name: "project", Changes: []openspec.Change{enriched}},
		discoveryFingerprints: map[string]string{"existing": fingerprint},
		enrichedFingerprints:  map[string]string{"existing": fingerprint},
		pendingEnrichments:    make(map[string]string),
		enrichmentRetryAfter:  make(map[string]time.Time),
	}
	freshProject := &openspec.Project{Name: "project", Changes: []openspec.Change{
		discovered,
		discoveredChange("new", "new content"),
	}}

	m.adoptDiscoveredProject(freshProject)
	if m.project.Changes[0].Artifacts[0].Source != openspec.ArtifactSourceStatus {
		t.Fatal("expected unchanged authoritative enrichment preserved")
	}
	if m.project.Changes[1].Name != "new" {
		t.Fatalf("expected new change, got %#v", m.project.Changes)
	}
	if m.discoveryFingerprints["new"] == "" {
		t.Fatal("expected new change fingerprint")
	}
}

func TestAdoptDiscoveredProjectRemovesDeletedChangeState(t *testing.T) {
	change := discoveredChange("deleted", "content")
	fingerprint := changeDiscoveryFingerprint(change)
	m := Model{
		project:               &openspec.Project{Changes: []openspec.Change{change}},
		discoveryFingerprints: map[string]string{"deleted": fingerprint},
		enrichedFingerprints:  map[string]string{"deleted": fingerprint},
		pendingEnrichments:    map[string]string{"deleted": fingerprint},
		enrichmentRetryAfter:  map[string]time.Time{"deleted": time.Now().Add(time.Minute)},
	}

	m.adoptDiscoveredProject(&openspec.Project{Name: "project"})
	if len(m.project.Changes) != 0 {
		t.Fatal("expected deleted change removed")
	}
	if len(m.discoveryFingerprints) != 0 || len(m.enrichedFingerprints) != 0 || len(m.pendingEnrichments) != 0 || len(m.enrichmentRetryAfter) != 0 {
		t.Fatalf("expected deleted change caches removed: discovery=%v enriched=%v pending=%v retry=%v",
			m.discoveryFingerprints, m.enrichedFingerprints, m.pendingEnrichments, m.enrichmentRetryAfter)
	}
}

func TestStatusEnrichmentRetriesAfterDelay(t *testing.T) {
	change := discoveredChange("retry", "content")
	fingerprint := changeDiscoveryFingerprint(change)
	client := &trackingStatusClient{}
	m := Model{
		root:                  "/project",
		project:               &openspec.Project{Changes: []openspec.Change{change}},
		openSpec:              client,
		discoveryFingerprints: map[string]string{"retry": fingerprint},
		enrichedFingerprints:  make(map[string]string),
		pendingEnrichments:    map[string]string{"retry": fingerprint},
		enrichmentRetryAfter:  make(map[string]time.Time),
	}
	m.applyStatusEnrichment(statusEnrichmentMsg{Results: []statusEnrichmentResult{{
		ChangeName: "retry", Fingerprint: fingerprint, Err: fmt.Errorf("CLI unavailable"),
	}}})
	if cmd := m.scheduleStatusEnrichment(); cmd != nil {
		t.Fatal("expected retry delay to suppress immediate command")
	}
	m.enrichmentRetryAfter["retry"] = time.Now().Add(-time.Second)
	if cmd := m.scheduleStatusEnrichment(); cmd == nil {
		t.Fatal("expected command after retry delay")
	}
}

func TestDiscoveryFingerprintChangesWithArtifactContent(t *testing.T) {
	before := discoveredChange("same", "before")
	after := discoveredChange("same", "after")
	if changeDiscoveryFingerprint(before) == changeDiscoveryFingerprint(after) {
		t.Fatal("expected content change to invalidate fingerprint")
	}
	beforeCopy := before
	if changeDiscoveryFingerprint(before) != changeDiscoveryFingerprint(beforeCopy) {
		t.Fatal("expected deterministic fingerprint")
	}
}
