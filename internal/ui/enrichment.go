package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/openspec"
)

const (
	maxConcurrentStatus = 4
	statusRetryDelay    = 5 * time.Second
)

type openSpecClient interface {
	Status(root, changeName string) (openspec.ChangeStatus, error)
	Schemas(root string) ([]openspec.SchemaInfo, error)
}

type startEnrichmentMsg struct{}

type schemaCatalogMsg struct {
	Schemas []openspec.SchemaInfo
	Err     error
}

type statusEnrichmentResult struct {
	ChangeName  string
	Fingerprint string
	Status      openspec.ChangeStatus
	Err         error
}

type statusEnrichmentMsg struct {
	Results []statusEnrichmentResult
}

type statusEnrichmentRequest struct {
	ChangeName  string
	Fingerprint string
}

func discoveryFingerprints(changes []openspec.Change) map[string]string {
	fingerprints := make(map[string]string, len(changes))
	for _, change := range changes {
		fingerprints[change.Name] = changeDiscoveryFingerprint(change)
	}
	return fingerprints
}

func changeDiscoveryFingerprint(change openspec.Change) string {
	hash := sha256.New()
	_, _ = fmt.Fprintf(hash, "name=%s\npath=%s\nschema=%s\ncreated=%s\n", change.Name, change.Path, change.Schema, change.Created)
	for _, artifact := range change.Artifacts {
		_, _ = fmt.Fprintf(hash, "artifact=%s\npattern=%s\nsource=%s\n", artifact.ID, artifact.OutputPattern, artifact.Source)
		for _, output := range artifact.Outputs {
			_, _ = fmt.Fprintf(hash, "output=%s\npresent=%t\ncontent=%s\n", output.RelativePath, output.Present, output.Content)
		}
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (m *Model) scheduleStatusEnrichment() tea.Cmd {
	if m.openSpec == nil || m.project == nil || m.singlePath && m.mode != ModeNormal {
		return nil
	}
	if m.discoveryFingerprints == nil {
		m.discoveryFingerprints = discoveryFingerprints(m.project.Changes)
	}
	if m.enrichedFingerprints == nil {
		m.enrichedFingerprints = make(map[string]string)
	}
	if m.pendingEnrichments == nil {
		m.pendingEnrichments = make(map[string]string)
	}
	if m.enrichmentRetryAfter == nil {
		m.enrichmentRetryAfter = make(map[string]time.Time)
	}
	if m.lastStatusByName == nil {
		m.lastStatusByName = make(map[string]openspec.ChangeStatus)
	}

	now := time.Now()
	requests := make([]statusEnrichmentRequest, 0, len(m.project.Changes))
	for _, change := range m.project.Changes {
		fingerprint := m.discoveryFingerprints[change.Name]
		if fingerprint == "" {
			fingerprint = changeDiscoveryFingerprint(change)
			m.discoveryFingerprints[change.Name] = fingerprint
		}
		if m.enrichedFingerprints[change.Name] == fingerprint || m.pendingEnrichments[change.Name] == fingerprint {
			continue
		}
		if retryAt := m.enrichmentRetryAfter[change.Name]; retryAt.After(now) {
			continue
		}
		requests = append(requests, statusEnrichmentRequest{ChangeName: change.Name, Fingerprint: fingerprint})
		m.pendingEnrichments[change.Name] = fingerprint
	}
	if len(requests) == 0 {
		return nil
	}

	root := m.root
	client := m.openSpec
	return func() tea.Msg {
		results := make([]statusEnrichmentResult, len(requests))
		jobs := make(chan int)
		workers := min(maxConcurrentStatus, len(requests))
		var wg sync.WaitGroup
		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				for index := range jobs {
					request := requests[index]
					status, err := client.Status(root, request.ChangeName)
					results[index] = statusEnrichmentResult{
						ChangeName:  request.ChangeName,
						Fingerprint: request.Fingerprint,
						Status:      status,
						Err:         err,
					}
				}
			}()
		}
		for index := range requests {
			jobs <- index
		}
		close(jobs)
		wg.Wait()
		return statusEnrichmentMsg{Results: results}
	}
}

func (m Model) loadSchemaCatalog() tea.Cmd {
	if m.openSpec == nil || m.singlePath {
		return nil
	}
	root := m.root
	client := m.openSpec
	return func() tea.Msg {
		schemas, err := client.Schemas(root)
		return schemaCatalogMsg{Schemas: schemas, Err: err}
	}
}

// applyStatusEnrichment merges authoritative status results into the model.
// It reports whether the index was rebuilt and whether the currently viewed
// artifact changed, so callers avoid redundant viewport reloads: a status
// result that reproduces the exact enriched state refreshes nothing.
func (m *Model) applyStatusEnrichment(msg statusEnrichmentMsg) (indexRefreshed, currentArtifactChanged bool) {
	if m.project == nil {
		return false, false
	}
	changed := false
	for _, result := range msg.Results {
		if m.pendingEnrichments[result.ChangeName] == result.Fingerprint {
			delete(m.pendingEnrichments, result.ChangeName)
		}
		if m.discoveryFingerprints[result.ChangeName] != result.Fingerprint {
			continue
		}
		index := changeIndexByName(m.project.Changes, result.ChangeName)
		if index < 0 {
			continue
		}
		if result.Err != nil {
			appendModelChangeDiagnostic(&m.project.Changes[index], result.Err.Error())
			m.enrichmentRetryAfter[result.ChangeName] = time.Now().Add(statusRetryDelay)
			continue
		}
		if m.lastStatusByName == nil {
			m.lastStatusByName = make(map[string]openspec.ChangeStatus)
		}
		m.lastStatusByName[result.ChangeName] = result.Status
		before := changeDiscoveryFingerprint(m.project.Changes[index])
		enriched := m.loader.EnrichChange(m.project.Changes[index], result.Status)
		m.enrichedFingerprints[result.ChangeName] = result.Fingerprint
		delete(m.enrichmentRetryAfter, result.ChangeName)
		if changeDiscoveryFingerprint(enriched) == before {
			continue
		}
		m.project.Changes[index] = enriched
		m.invalidateArtifactRenderCache(result.ChangeName)
		if m.mode == ModeNormal && index == m.changeIdx {
			m.reconcileArtifactSelection()
			currentArtifactChanged = true
		}
		changed = true
	}
	if changed && m.mode == ModeIndex && m.vpReady {
		m.rebuildIndexPreservingCursor()
		m.refreshIndexViewport()
		indexRefreshed = true
	}
	return indexRefreshed, currentArtifactChanged
}

func changeIndexByName(changes []openspec.Change, name string) int {
	for index := range changes {
		if changes[index].Name == name {
			return index
		}
	}
	return -1
}

func appendModelChangeDiagnostic(change *openspec.Change, diagnostic string) {
	if change.Diagnostic == "" {
		change.Diagnostic = diagnostic
		return
	}
	if change.Diagnostic == diagnostic {
		return
	}
	change.Diagnostic += "; " + diagnostic
}

func (m *Model) adoptDiscoveredChange(index int, fresh openspec.Change) bool {
	if m.project == nil || index < 0 || index >= len(m.project.Changes) {
		return false
	}
	if m.discoveryFingerprints == nil {
		m.discoveryFingerprints = make(map[string]string)
	}
	if m.enrichedFingerprints == nil {
		m.enrichedFingerprints = make(map[string]string)
	}
	if m.pendingEnrichments == nil {
		m.pendingEnrichments = make(map[string]string)
	}
	if m.enrichmentRetryAfter == nil {
		m.enrichmentRetryAfter = make(map[string]time.Time)
	}
	if m.lastStatusByName == nil {
		m.lastStatusByName = make(map[string]openspec.ChangeStatus)
	}
	fingerprint := changeDiscoveryFingerprint(fresh)
	if m.discoveryFingerprints[fresh.Name] == fingerprint {
		return false
	}
	oldName := m.project.Changes[index].Name
	// Capture the selection position before replacing so ModeNormal keeps
	// the nearest surviving artifact instead of resetting to the default.
	oldChange := m.project.Changes[index]
	oldArtifactIndex, oldOutputIndex := artifactSelectionPosition(&oldChange, m.artifactSelection)
	// Re-apply the last known status synchronously so schema-declared
	// artifact identities survive disk reloads; the async refresh still
	// runs (enrichedFingerprints is cleared below) to correct stale status.
	fresh = m.applyCachedStatus(fresh)
	m.project.Changes[index] = fresh
	m.invalidateArtifactRenderCache(oldName)
	if m.mode == ModeNormal && index == m.changeIdx {
		m.reconcileArtifactSelectionNear(oldArtifactIndex, oldOutputIndex)
	}
	if oldName != fresh.Name {
		delete(m.discoveryFingerprints, oldName)
		delete(m.enrichedFingerprints, oldName)
		delete(m.pendingEnrichments, oldName)
		delete(m.enrichmentRetryAfter, oldName)
		delete(m.lastStatusByName, oldName)
	}
	m.discoveryFingerprints[fresh.Name] = fingerprint
	delete(m.enrichedFingerprints, fresh.Name)
	delete(m.pendingEnrichments, fresh.Name)
	delete(m.enrichmentRetryAfter, fresh.Name)
	return true
}

func (m *Model) applyCachedStatus(change openspec.Change) openspec.Change {
	if m.loader == nil {
		return change
	}
	status, ok := m.lastStatusByName[change.Name]
	if !ok {
		return change
	}
	return m.loader.EnrichChange(change, status)
}

func (m *Model) adoptDiscoveredProject(fresh *openspec.Project) {
	if fresh == nil {
		return
	}
	if m.discoveryFingerprints == nil {
		m.discoveryFingerprints = make(map[string]string)
	}
	if m.enrichedFingerprints == nil {
		m.enrichedFingerprints = make(map[string]string)
	}
	if m.pendingEnrichments == nil {
		m.pendingEnrichments = make(map[string]string)
	}
	if m.enrichmentRetryAfter == nil {
		m.enrichmentRetryAfter = make(map[string]time.Time)
	}
	if m.lastStatusByName == nil {
		m.lastStatusByName = make(map[string]openspec.ChangeStatus)
	}
	oldByName := make(map[string]openspec.Change)
	if m.project != nil {
		for _, change := range m.project.Changes {
			oldByName[change.Name] = change
		}
	}

	changes := make([]openspec.Change, len(fresh.Changes))
	newFingerprints := make(map[string]string, len(fresh.Changes))
	for index, change := range fresh.Changes {
		fingerprint := changeDiscoveryFingerprint(change)
		newFingerprints[change.Name] = fingerprint
		if old, ok := oldByName[change.Name]; ok && m.discoveryFingerprints[change.Name] == fingerprint {
			changes[index] = old
			continue
		}
		// Disk state changed (or the change is new): keep schema identities
		// stable by re-applying the cached status synchronously.
		changes[index] = m.applyCachedStatus(change)
		delete(m.enrichedFingerprints, change.Name)
		delete(m.pendingEnrichments, change.Name)
		delete(m.enrichmentRetryAfter, change.Name)
	}
	for name := range m.discoveryFingerprints {
		if _, ok := newFingerprints[name]; ok {
			continue
		}
		delete(m.enrichedFingerprints, name)
		delete(m.pendingEnrichments, name)
		delete(m.enrichmentRetryAfter, name)
		delete(m.lastStatusByName, name)
	}
	m.project = &openspec.Project{Name: fresh.Name, Changes: changes}
	m.discoveryFingerprints = newFingerprints
}

func (m *Model) reconcileDiscoveryFingerprints() {
	if m.project == nil {
		return
	}
	currentNames := make(map[string]bool, len(m.project.Changes))
	for _, change := range m.project.Changes {
		currentNames[change.Name] = true
		if _, ok := m.discoveryFingerprints[change.Name]; !ok {
			m.discoveryFingerprints[change.Name] = changeDiscoveryFingerprint(change)
		}
	}
	for name := range m.discoveryFingerprints {
		if currentNames[name] {
			continue
		}
		delete(m.discoveryFingerprints, name)
		delete(m.enrichedFingerprints, name)
		delete(m.pendingEnrichments, name)
		delete(m.enrichmentRetryAfter, name)
		delete(m.lastStatusByName, name)
	}
}
