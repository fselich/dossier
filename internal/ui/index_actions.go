package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/openspec"
)

func indexOperationLabel(operation indexOperation) string {
	switch operation {
	case indexOperationCreate:
		return "creating change"
	case indexOperationValidate:
		return "validating"
	case indexOperationArchive:
		return "archiving change"
	case indexOperationReactivate:
		return "making change active"
	case indexOperationUndo:
		return "undoing lifecycle action"
	default:
		return "working"
	}
}

type openSpecActionClient interface {
	CreateChange(root, name, schema string) (openspec.CreateChangeResult, error)
	Validate(root, name, kind string) (openspec.ValidationResult, error)
}

func (m Model) handleIndexActionInput(msg tea.KeyPressMsg) (Model, tea.Cmd, bool) {
	keys := m.effectiveKeyMap()
	action := m.index.Action
	switch action.Mode {
	case indexActionIdle:
		return m, nil, false
	case indexActionPending:
		if matchesKey(msg, keys.Index.New) || matchesKey(msg, keys.Index.Lifecycle) || matchesKey(msg, keys.Index.Open) ||
			matchesKey(msg, keys.Index.Edit) || matchesKey(msg, keys.Index.Validate) || matchesKey(msg, keys.Index.Undo) {
			return m, nil, true
		}
		return m, nil, false
	case indexActionSchema:
		if matchesKey(msg, keys.Filter.Cancel) {
			m.index.Action = indexActionState{}
			return m, nil, true
		}
		visible := m.filteredSchemas()
		switch {
		case matchesKey(msg, keys.Index.Down):
			if action.SchemaCursor < len(visible)-1 {
				action.SchemaCursor++
			}
		case matchesKey(msg, keys.Index.Up):
			if action.SchemaCursor > 0 {
				action.SchemaCursor--
			}
		case matchesKey(msg, keys.Filter.Backspace):
			if action.SchemaFilter != "" {
				_, size := utf8.DecodeLastRuneInString(action.SchemaFilter)
				action.SchemaFilter = action.SchemaFilter[:len(action.SchemaFilter)-size]
				action.SchemaCursor = 0
			}
		case matchesKey(msg, keys.Index.Open), matchesKey(msg, keys.Filter.Accept):
			if len(visible) > 0 {
				action.Mode = indexActionName
				action.Schema = visible[action.SchemaCursor].Name
				action.Name = ""
			}
		default:
			if msg.Text != "" && utf8.RuneCountInString(msg.Text) == 1 {
				action.SchemaFilter += msg.Text
				action.SchemaCursor = 0
			}
		}
		m.index.Action = action
		return m, nil, true
	case indexActionName:
		switch {
		case matchesKey(msg, keys.Filter.Cancel):
			m.index.Action = indexActionState{}
		case matchesKey(msg, keys.Filter.Backspace):
			if action.Name != "" {
				_, size := utf8.DecodeLastRuneInString(action.Name)
				action.Name = action.Name[:len(action.Name)-size]
				m.index.Action = action
			}
		case matchesKey(msg, keys.Index.Open), matchesKey(msg, keys.Filter.Accept):
			name := strings.TrimSpace(action.Name)
			if name == "" {
				action.Status = "change name cannot be empty"
				action.StatusError = true
				m.index.Action = action
				return m, nil, true
			}
			client, ok := m.openSpec.(openSpecActionClient)
			if !ok {
				action.Status = "OpenSpec CLI is required to create a change"
				action.StatusError = true
				m.index.Action = action
				return m, nil, true
			}
			schema := action.Schema
			root := m.root
			m.index.Action = indexActionState{Mode: indexActionPending, Operation: indexOperationCreate, Name: name, Schema: schema, TargetIdentity: "active:" + name}
			return m, func() tea.Msg {
				result, err := client.CreateChange(root, name, schema)
				if result.Name != "" {
					name = result.Name
				}
				return indexActionResultMsg{Kind: indexOperationCreate, Name: name, TargetIdentity: "active:" + name, Err: err}
			}, true
		default:
			if msg.Text != "" && utf8.RuneCountInString(msg.Text) == 1 {
				action.Name += msg.Text
				m.index.Action = action
			}
		}
		return m, nil, true
	case indexActionConfirmArchive, indexActionConfirmReactivate:
		if matchesKey(msg, keys.Filter.Cancel) {
			m.index.Action = indexActionState{}
			return m, nil, true
		}
		if !matchesKey(msg, keys.Index.Lifecycle) && !matchesKey(msg, keys.Index.Open) {
			return m, nil, true
		}
		if m.readOnly {
			m.index.Action = indexActionState{}
			return m, nil, true
		}
		item, ok := m.indexItemByStableIdentity(action.TargetIdentity)
		if !ok || m.lifecycle == nil {
			m.index.Action = indexActionState{Status: "selected lifecycle item is no longer available", StatusError: true}
			return m, nil, true
		}
		service := m.lifecycle
		if action.Mode == indexActionConfirmArchive {
			change, ok := m.indexItemChange(item)
			if !ok {
				m.index.Action = indexActionState{Status: "active change is no longer available", StatusError: true}
				return m, nil, true
			}
			m.index.Action = indexActionState{Mode: indexActionPending, Operation: indexOperationArchive, Name: change.Name, TargetIdentity: action.TargetIdentity}
			return m, func() tea.Msg {
				record, result, err := service.Archive(change)
				return indexActionResultMsg{Kind: indexOperationArchive, Name: change.Name, TargetIdentity: "archive:" + change.Name, Archive: result, Undo: record, Err: err}
			}, true
		}
		change, ok := m.indexItemChange(item)
		if !ok {
			m.index.Action = indexActionState{Status: "archived change is no longer available", StatusError: true}
			return m, nil, true
		}
		m.index.Action = indexActionState{Mode: indexActionPending, Operation: indexOperationReactivate, Name: change.Name, TargetIdentity: action.TargetIdentity}
		return m, func() tea.Msg {
			record, err := service.Reactivate(change.Path)
			return indexActionResultMsg{Kind: indexOperationReactivate, Name: change.Name, TargetIdentity: "active:" + change.Name, Undo: record, Err: err}
		}, true
	case indexActionConfirmUndo:
		if matchesKey(msg, keys.Filter.Cancel) {
			m.index.Action = indexActionState{}
			return m, nil, true
		}
		if !matchesKey(msg, keys.Index.Undo) && !matchesKey(msg, keys.Index.Open) {
			return m, nil, true
		}
		if m.readOnly || m.lifecycle == nil || m.lifecycleUndo == nil {
			m.index.Action = indexActionState{}
			return m, nil, true
		}
		record := m.lifecycleUndo
		service := m.lifecycle
		target := "active:" + record.ChangeName
		if record.Kind == openspec.LifecycleReactivate {
			target = "archive:" + record.ChangeName
		}
		m.index.Action = indexActionState{Mode: indexActionPending, Operation: indexOperationUndo, Name: record.ChangeName, TargetIdentity: target}
		return m, func() tea.Msg {
			err := service.Undo(record)
			return indexActionResultMsg{Kind: indexOperationUndo, Name: record.ChangeName, TargetIdentity: target, Err: err}
		}, true
	}
	return m, nil, false
}

func (m Model) filteredSchemas() []openspec.SchemaInfo {
	filter := strings.ToLower(m.index.Action.SchemaFilter)
	if filter == "" {
		return append([]openspec.SchemaInfo(nil), m.schemaCatalog...)
	}
	var result []openspec.SchemaInfo
	for _, schema := range m.schemaCatalog {
		haystack := strings.ToLower(schema.Name + " " + schema.Description + " " + schema.Source + " " + strings.Join(schema.Artifacts, " "))
		if strings.Contains(haystack, filter) {
			result = append(result, schema)
		}
	}
	return result
}

func (m Model) indexItemByStableIdentity(identity string) (indexItem, bool) {
	index := indexItemByIdentity(m.index.Items, identity)
	if index < 0 {
		return indexItem{}, false
	}
	return m.index.Items[index], true
}

func (m Model) startNewChange() Model {
	if m.readOnly {
		return m
	}
	if len(m.schemaCatalog) == 0 {
		message := "no OpenSpec workflow schemas available"
		if m.schemaCatalogErr != "" {
			message = m.schemaCatalogErr
		}
		m.index.Action = indexActionState{Status: message, StatusError: true}
		return m
	}
	if len(m.schemaCatalog) == 1 {
		m.index.Action = indexActionState{Mode: indexActionName, Schema: m.schemaCatalog[0].Name}
		return m
	}
	m.index.Action = indexActionState{Mode: indexActionSchema}
	return m
}

func (m Model) archiveConfirmationSummary(identity string) string {
	item, ok := m.indexItemByStableIdentity(identity)
	if !ok {
		return strings.TrimPrefix(identity, "active:")
	}
	change, ok := m.indexItemChange(item)
	if !ok {
		return strings.TrimPrefix(identity, "active:")
	}
	done, total := taskCounts(change)
	text := change.Name
	if total > 0 {
		text += fmt.Sprintf(" — tasks %d/%d", done, total)
		if done < total {
			text += fmt.Sprintf(" (%d incomplete)", total-done)
		}
	}
	var specs []string
	for _, artifact := range change.Artifacts {
		for _, output := range artifact.Outputs {
			parts := strings.Split(filepath.ToSlash(output.RelativePath), "/")
			if len(parts) >= 3 && parts[0] == "specs" && parts[1] != "" {
				specs = append(specs, parts[1])
			}
		}
	}
	if len(specs) > 0 {
		text += " — specs: " + strings.Join(specs, ", ")
	}
	return text
}

func (m Model) startLifecycle(item indexItem) Model {
	if m.readOnly {
		return m
	}
	switch item.kind {
	case indexKindActive:
		m.index.Action = indexActionState{Mode: indexActionConfirmArchive, Operation: indexOperationArchive, TargetIdentity: item.identity}
	case indexKindArchived:
		m.index.Action = indexActionState{Mode: indexActionConfirmReactivate, Operation: indexOperationReactivate, TargetIdentity: item.identity}
	}
	return m
}

func (m Model) startValidation(item indexItem) (Model, tea.Cmd) {
	client, ok := m.openSpec.(openSpecActionClient)
	if !ok {
		m.index.Action = indexActionState{Status: "OpenSpec CLI is required to validate", StatusError: true}
		return m, nil
	}
	name, kind, ok := m.indexValidationTarget(item)
	if !ok {
		return m, nil
	}
	root := m.root
	identity := item.identity
	m.index.Action = indexActionState{Mode: indexActionPending, Operation: indexOperationValidate, Name: name, TargetIdentity: identity}
	return m, func() tea.Msg {
		result, err := client.Validate(root, name, kind)
		if err == nil && !result.Valid {
			detail := result.Summary
			if len(result.Errors) > 0 {
				detail = strings.Join(result.Errors, "; ")
			}
			err = fmt.Errorf("validation failed: %s", detail)
		}
		return indexActionResultMsg{Kind: indexOperationValidate, Name: name, TargetIdentity: identity, Validation: result, Err: err}
	}
}

func (m Model) indexValidationTarget(item indexItem) (string, string, bool) {
	switch item.kind {
	case indexKindActive:
		if m.project != nil && item.idx >= 0 && item.idx < len(m.project.Changes) {
			return m.project.Changes[item.idx].Name, "change", true
		}
	case indexKindSpec, indexKindRequirement:
		if item.idx >= 0 && item.idx < len(m.projectSpecs) {
			return m.projectSpecs[item.idx].Name, "spec", true
		}
	}
	return "", "", false
}

func (m Model) indexEditPath(item indexItem) (string, bool) {
	if m.readOnly || item.archived {
		return "", false
	}
	if item.kind == indexKindSpec || item.kind == indexKindRequirement {
		if item.idx < 0 || item.idx >= len(m.projectSpecs) {
			return "", false
		}
		path := filepath.Join(m.root, "openspec", "specs", m.projectSpecs[item.idx].Name, "spec.md")
		return safeIndexFile(filepath.Join(m.root, "openspec", "specs"), path)
	}
	change, ok := m.indexItemChange(item)
	if !ok {
		return "", false
	}
	var output *openspec.ArtifactOutput
	switch item.kind {
	case indexKindActive:
		for artifactIndex := range change.Artifacts {
			for outputIndex := range change.Artifacts[artifactIndex].Outputs {
				candidate := &change.Artifacts[artifactIndex].Outputs[outputIndex]
				if candidate.Present {
					output = candidate
					break
				}
			}
			if output != nil {
				break
			}
		}
	case indexKindArtifact:
		if item.artifactIdx >= 0 && item.artifactIdx < len(change.Artifacts) {
			for outputIndex := range change.Artifacts[item.artifactIdx].Outputs {
				candidate := &change.Artifacts[item.artifactIdx].Outputs[outputIndex]
				if candidate.Present {
					output = candidate
					break
				}
			}
		}
	case indexKindArtifactOutput:
		if item.artifactIdx >= 0 && item.artifactIdx < len(change.Artifacts) && item.outputIdx >= 0 && item.outputIdx < len(change.Artifacts[item.artifactIdx].Outputs) {
			candidate := &change.Artifacts[item.artifactIdx].Outputs[item.outputIdx]
			if candidate.Present {
				output = candidate
			}
		}
	}
	if output == nil {
		return "", false
	}
	return safeIndexFile(change.Path, filepath.Join(change.Path, filepath.FromSlash(output.RelativePath)))
}

func safeIndexFile(root, path string) (string, bool) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", false
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	relative, err := filepath.Rel(absoluteRoot, absolutePath)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return "", false
	}
	rootResolved, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return "", false
	}
	pathResolved, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", false
	}
	resolvedRelative, err := filepath.Rel(rootResolved, pathResolved)
	if err != nil || resolvedRelative == ".." || strings.HasPrefix(resolvedRelative, ".."+string(filepath.Separator)) || filepath.IsAbs(resolvedRelative) {
		return "", false
	}
	info, err := os.Stat(pathResolved)
	if err != nil || !info.Mode().IsRegular() {
		return "", false
	}
	return absolutePath, true
}

func (m Model) editIndexItem(item indexItem) (Model, tea.Cmd) {
	if m.readOnly {
		return m, nil
	}
	path, ok := m.indexEditPath(item)
	if !ok {
		return m, nil
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	command := exec.Command(editor, path)
	identity := item.identity
	return m, tea.ExecProcess(command, func(err error) tea.Msg {
		if err != nil {
			return indexActionResultMsg{Kind: indexOperationNone, TargetIdentity: identity, Err: fmt.Errorf("editor: %w", err)}
		}
		return editorReturnMsg{indexIdentity: identity}
	})
}

func (m *Model) reloadIndex(targetIdentity string) {
	if project, err := m.loader.LoadFrom(m.root); err == nil {
		// Fingerprints must capture the discovered state (pre-enrichment),
		// otherwise every later adoption comparison mismatches.
		discovered := discoveryFingerprints(project.Changes)
		// Editor return must not swap schema-declared artifact rows back to
		// filename-stem discovery IDs: re-apply cached status synchronously.
		for i := range project.Changes {
			project.Changes[i] = m.applyCachedStatus(project.Changes[i])
		}
		m.project = project
		m.discoveryFingerprints = discovered
		m.enrichedFingerprints = make(map[string]string)
		m.pendingEnrichments = make(map[string]string)
	}
	if archives, err := m.loader.ListArchiveChangesFrom(m.root); err == nil {
		m.index.ArchiveChanges = archives
	}
	if specs, err := m.loader.LoadProjectSpecsFrom(m.root); err == nil {
		m.projectSpecs = specs
	}
	m.expandIndexIdentityAncestors(targetIdentity)
	m.buildIndexItems()
	m.applyFilter()
	m.selectVisibleIndexIdentity(targetIdentity)
	if m.vpReady {
		m.refreshIndexViewport()
	}
}

func (m *Model) expandIndexIdentityAncestors(identity string) {
	if requirementAt := strings.Index(identity, ":requirement:"); requirementAt >= 0 && strings.HasPrefix(identity, "canonical-spec:") {
		specName := strings.TrimPrefix(identity[:requirementAt], "canonical-spec:")
		for index := range m.projectSpecs {
			if m.projectSpecs[index].Name == specName {
				if m.index.ExpandedSpecs == nil {
					m.index.ExpandedSpecs = make(map[int]bool)
				}
				m.index.ExpandedSpecs[index] = true
				break
			}
		}
	}
	if artifactAt := strings.Index(identity, ":artifact:"); artifactAt >= 0 {
		if m.index.ExpandedChanges == nil {
			m.index.ExpandedChanges = make(map[string]bool)
		}
		m.index.ExpandedChanges[identity[:artifactAt]] = true
	}
	if outputAt := strings.Index(identity, ":output:"); outputAt >= 0 {
		if m.index.ExpandedArtifacts == nil {
			m.index.ExpandedArtifacts = make(map[string]bool)
		}
		m.index.ExpandedArtifacts[identity[:outputAt]] = true
	}
}

func (m *Model) selectVisibleIndexIdentity(identity string) {
	rawIndex := indexItemByIdentity(m.index.Items, identity)
	if rawIndex < 0 {
		m.index.Cursor = min(m.index.Cursor, max(0, m.visibleItemCount()-1))
		return
	}
	if m.index.FilterIndices == nil {
		m.index.Cursor = rawIndex
		return
	}
	for visibleIndex, filteredRawIndex := range m.index.FilterIndices {
		if filteredRawIndex == rawIndex {
			m.index.Cursor = visibleIndex
			return
		}
	}
}

func (m *Model) applyIndexActionResult(msg indexActionResultMsg) {
	status := ""
	if msg.Err != nil {
		status = msg.Err.Error()
		m.index.Action = indexActionState{Status: status, StatusError: true}
		target := msg.TargetIdentity
		switch msg.Kind {
		case indexOperationArchive:
			target = "active:" + msg.Name
		case indexOperationReactivate:
			target = "archive:" + msg.Name
		}
		m.reloadIndex(target)
		return
	}
	switch msg.Kind {
	case indexOperationCreate:
		status = "created " + msg.Name
	case indexOperationValidate:
		status = msg.Validation.Summary
		if status == "" {
			status = msg.Name + " is valid"
		}
	case indexOperationArchive:
		status = "archived " + msg.Name
		if m.lifecycleUndo != nil && m.lifecycle != nil {
			m.lifecycle.Cleanup(m.lifecycleUndo)
		}
		m.lifecycleUndo = msg.Undo
	case indexOperationReactivate:
		status = "made " + msg.Name + " active"
		if m.lifecycleUndo != nil && m.lifecycle != nil {
			m.lifecycle.Cleanup(m.lifecycleUndo)
		}
		m.lifecycleUndo = msg.Undo
	case indexOperationUndo:
		status = "undid lifecycle action for " + msg.Name
		m.lifecycleUndo = nil
	default:
		status = "index refreshed"
	}
	m.index.Action = indexActionState{Status: status}
	m.reloadIndex(msg.TargetIdentity)
}
