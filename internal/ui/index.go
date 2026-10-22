package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/aemonge/dossier/internal/openspec"
)

func (m *Model) handleTick() tea.Cmd {
	m.pollGitStatus()
	if m.mode == ModeViewingSpec {
		return nil
	}
	if m.mode == ModeViewingArchive {
		if m.singlePath {
			return m.pollSingleArchivedContent()
		}
		return nil
	}

	var pollCmd tea.Cmd
	if m.mode == ModeIndex {
		pollCmd = m.pollIndexMode()
	} else {
		if !m.singlePath {
			pollCmd = m.pollNormalModeChanges()
		}
		if pollCmd == nil {
			pollCmd = m.pollNormalModeContent()
		}
	}
	m.reconcileDiscoveryFingerprints()
	return tea.Batch(pollCmd, m.scheduleStatusEnrichment())
}

func (m *Model) pollIndexMode() tea.Cmd {
	diskChanges, err := m.loader.ListChangeNamesFrom(m.root)
	if err != nil {
		return nil
	}
	diskArchives, err := m.loader.ListArchiveNamesFrom(m.root)
	if err != nil {
		return nil
	}
	diskSpecs, err := m.loader.ListSpecNamesFrom(m.root)
	if err != nil {
		return nil
	}

	archiveNames := make([]string, len(m.index.ArchiveChanges))
	for i, ch := range m.index.ArchiveChanges {
		archiveNames[i] = filepath.Base(ch.Path)
	}
	specNames := make([]string, len(m.projectSpecs))
	for i, ps := range m.projectSpecs {
		specNames[i] = ps.Name
	}

	if sameNames(m.project.Changes, diskChanges) &&
		sameStrings(archiveNames, diskArchives) &&
		sameStrings(specNames, diskSpecs) {
		needsRefresh := false
		for i := range m.project.Changes {
			fresh := m.loader.ReloadChange(m.project.Changes[i])
			if m.adoptDiscoveredChange(i, fresh) {
				needsRefresh = true
			}
		}
		if needsRefresh {
			m.rebuildIndexPreservingCursor()
			m.refreshIndexViewport()
		}
		return nil
	}

	if p, err := m.loader.LoadFrom(m.root); err == nil {
		m.adoptDiscoveredProject(p)
	}
	var archiveErr error
	m.index.ArchiveChanges, archiveErr = m.loader.ListArchiveChangesFrom(m.root)
	if archiveErr != nil {
		m.errMsg = "error loading archive changes: " + archiveErr.Error()
	}
	var specErr error
	m.projectSpecs, specErr = m.loader.LoadProjectSpecsFrom(m.root)
	if specErr != nil {
		m.errMsg = "error loading project specs: " + specErr.Error()
	}
	m.rebuildIndexPreservingCursor()
	m.refreshIndexViewport()
	return nil
}

func (m *Model) pollSingleArchivedContent() tea.Cmd {
	ch := m.currentArchive()
	if ch == nil {
		return nil
	}
	oldArtifactIndex, oldOutputIndex := artifactSelectionPosition(ch, m.artifactSelection)
	fresh := m.loader.ReloadChange(*ch)
	if changeDiscoveryFingerprint(*ch) == changeDiscoveryFingerprint(fresh) {
		return nil
	}
	m.index.ArchiveChanges[m.index.ArchiveCursor] = fresh
	m.invalidateArtifactRenderCache(fresh.Name)
	m.reconcileArtifactSelectionNear(oldArtifactIndex, oldOutputIndex)
	return m.loadViewport()
}

func (m *Model) pollNormalModeChanges() tea.Cmd {
	diskNames, err := m.loader.ListChangeNamesFrom(m.root)
	if err != nil {
		return nil
	}
	if !sameNames(m.project.Changes, diskNames) {
		currentName := ""
		if ch := m.current(); ch != nil {
			currentName = ch.Name
		}
		if p, err := m.loader.LoadFrom(m.root); err == nil {
			m.adoptDiscoveredProject(p)
			m.changeIdx = 0
			for i, ch := range p.Changes {
				if ch.Name == currentName {
					m.changeIdx = i
					break
				}
			}
			if len(p.Changes) == 0 {
				return nil
			}
			m.artifactRenderCache = make(map[artifactOutputIdentity]string)
			m.selectDefaultArtifact()
			m.loadTaskItems()
			return m.loadViewport()
		}
	}
	return nil
}

func (m *Model) pollNormalModeContent() tea.Cmd {
	ch := m.current()
	if ch == nil {
		return nil
	}
	var cursorText string
	if m.tasks.Cursor < len(m.tasks.Items) && m.tasks.Items[m.tasks.Cursor].Kind == openspec.KindTask {
		cursorText = m.tasks.Items[m.tasks.Cursor].Text
	}
	fresh := m.loader.ReloadChange(*ch)
	tasksChanged, viewportDirty := m.mergeReloadedChange(fresh)
	m.adoptDiscoveredChange(m.changeIdx, fresh)

	if tasksChanged {
		if cursorText != "" {
			m.tasks.Cursor = openspec.FindCursorByText(m.tasks.Items, cursorText)
		}
		if m.isTasksView() {
			m.refreshTasksViewport()
		}
	}
	if viewportDirty {
		return m.loadViewport()
	}
	return nil
}

func (m *Model) enterIndex() {
	if len(m.index.ArchiveChanges) == 0 {
		var archiveErr error
		m.index.ArchiveChanges, archiveErr = m.loader.ListArchiveChangesFrom(m.root)
		if archiveErr != nil {
			m.errMsg = "error loading archive changes: " + archiveErr.Error()
		}
	}
	var specErr error
	m.projectSpecs, specErr = m.loader.LoadProjectSpecsFrom(m.root)
	if specErr != nil {
		m.errMsg = "error loading project specs: " + specErr.Error()
	}
	m.index.ExpandedSpecs = make(map[int]bool)
	m.buildIndexItems()
	m.index.Cursor = 0
	restored := -1
	if m.returnIndexIdentity != "" {
		restored = indexItemByIdentity(m.index.Items, m.returnIndexIdentity)
		if restored < 0 {
			// Artifact IDs may have swapped (discovery ↔ schema) while
			// away: re-anchor by (change, output path) before dropping to 0.
			restored = m.indexItemByAnchor(m.returnIndexAnchor)
		}
	}
	if restored >= 0 {
		m.index.Cursor = restored
	}
	m.mode = ModeIndex
	m.vp.SetHeight(m.contentHeight())
	m.refreshIndexViewport()
}

func (m *Model) visibleItemIdx(rawIdx int) int {
	if m.index.FilterIndices != nil {
		return m.index.FilterIndices[rawIdx]
	}
	return rawIdx
}

func (m *Model) visibleItemCount() int {
	if m.index.FilterIndices != nil {
		return len(m.index.FilterIndices)
	}
	return len(m.index.Items)
}

func (m *Model) matchesFilter(item indexItem, lowerQuery string) bool {
	contains := func(values ...string) bool {
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), lowerQuery) {
				return true
			}
		}
		return false
	}
	switch item.kind {
	case indexKindSection:
		return true
	case indexKindActive:
		if m.project != nil && item.idx < len(m.project.Changes) {
			change := m.project.Changes[item.idx]
			return contains(change.Name, change.Schema)
		}
	case indexKindArchived:
		if item.idx < len(m.index.ArchiveChanges) {
			change := m.index.ArchiveChanges[item.idx]
			return contains(change.Name, change.Schema)
		}
	case indexKindArtifact:
		if change, ok := m.indexItemChange(item); ok && item.artifactIdx < len(change.Artifacts) {
			artifact := change.Artifacts[item.artifactIdx]
			outputNames := make([]string, 0, len(artifact.Outputs))
			for _, output := range artifact.Outputs {
				outputNames = append(outputNames, output.RelativePath)
			}
			_, secondary := artifactLabels(artifact)
			return contains(artifact.ID, secondary, artifact.OutputPattern, strings.Join(artifact.Requires, " "), strings.Join(outputNames, " "))
		}
	case indexKindArtifactOutput:
		if change, ok := m.indexItemChange(item); ok && item.artifactIdx < len(change.Artifacts) {
			artifact := change.Artifacts[item.artifactIdx]
			if item.outputIdx < len(artifact.Outputs) {
				output := artifact.Outputs[item.outputIdx]
				return contains(output.DisplayName, output.RelativePath)
			}
		}
	case indexKindSpec:
		if item.idx < len(m.projectSpecs) {
			return contains(m.projectSpecs[item.idx].Name)
		}
	case indexKindRequirement:
		if item.idx < len(m.projectSpecs) && item.reqIdx < len(m.projectSpecs[item.idx].RequirementNames) {
			return contains(m.projectSpecs[item.idx].RequirementNames[item.reqIdx])
		}
	}
	return false
}

func (m *Model) isItemVisible(idx int) bool {
	if m.index.FilterText == "" || m.index.FilterIndices == nil {
		return true
	}
	for _, visibleIndex := range m.index.FilterIndices {
		if visibleIndex == idx {
			return true
		}
	}
	return false
}

func (m *Model) applyFilter() {
	if m.index.FilterText == "" {
		m.index.FilterIndices = nil
		return
	}
	lower := strings.ToLower(m.index.FilterText)
	included := make(map[int]bool)
	for index, item := range m.index.Items {
		if !m.matchesFilter(item, lower) {
			continue
		}
		included[index] = true
		for _, ancestor := range indexIdentityAncestors(item.identity) {
			if ancestorIndex := indexItemByIdentity(m.index.Items, ancestor); ancestorIndex >= 0 {
				included[ancestorIndex] = true
			}
		}
	}
	m.index.FilterIndices = nil
	for index := range m.index.Items {
		if included[index] {
			m.index.FilterIndices = append(m.index.FilterIndices, index)
		}
	}
	if m.index.Cursor >= len(m.index.FilterIndices) {
		m.index.Cursor = 0
	}
}

func indexIdentityAncestors(identity string) []string {
	var ancestors []string
	switch {
	case strings.HasPrefix(identity, "active:"):
		ancestors = append(ancestors, "section:active-work")
	case strings.HasPrefix(identity, "archive:"):
		ancestors = append(ancestors, "section:history")
	case strings.HasPrefix(identity, "canonical-spec:"):
		ancestors = append(ancestors, "section:canonical-specs")
	}
	if outputAt := strings.Index(identity, ":output:"); outputAt >= 0 {
		identity = identity[:outputAt]
		ancestors = append(ancestors, identity)
	}
	if artifactAt := strings.Index(identity, ":artifact:"); artifactAt >= 0 {
		ancestors = append(ancestors, identity[:artifactAt])
	}
	if requirementAt := strings.Index(identity, ":requirement:"); requirementAt >= 0 {
		ancestors = append(ancestors, identity[:requirementAt])
	}
	return ancestors
}

func specSuffix(name string) string {
	if i := strings.LastIndex(name, "-"); i >= 0 {
		return name[i+1:]
	}
	return name
}

func (m *Model) buildSpecOrder() {
	n := len(m.projectSpecs)
	m.index.Order = make([]int, n)
	for i := range m.index.Order {
		m.index.Order[i] = i
	}
	if m.index.SortBySuffix {
		sort.SliceStable(m.index.Order, func(a, b int) bool {
			return specSuffix(m.projectSpecs[m.index.Order[a]].Name) < specSuffix(m.projectSpecs[m.index.Order[b]].Name)
		})
	}
}

func (m *Model) buildIndexItems() {
	m.buildSpecOrder()
	if m.index.ExpandedSpecs == nil {
		m.index.ExpandedSpecs = make(map[int]bool)
	}
	if m.index.ExpandedChanges == nil {
		m.index.ExpandedChanges = make(map[string]bool)
	}
	if m.index.ExpandedArtifacts == nil {
		m.index.ExpandedArtifacts = make(map[string]bool)
	}
	m.index.Items = nil

	m.index.Items = append(m.index.Items, indexItem{
		kind: indexKindSection, idx: sectionActive, identity: "section:active-work",
	})
	if (!m.index.CollapsedSections[sectionActive] || m.index.FilterText != "") && m.project != nil {
		for i := range m.project.Changes {
			m.appendChangeHierarchy(m.project.Changes[i], i, false)
		}
	}

	m.index.Items = append(m.index.Items, indexItem{
		kind: indexKindSection, idx: sectionSpecs, identity: "section:canonical-specs",
	})
	if !m.index.CollapsedSections[sectionSpecs] || m.index.FilterText != "" {
		for _, i := range m.index.Order {
			ps := m.projectSpecs[i]
			m.index.Items = append(m.index.Items, indexItem{
				kind: indexKindSpec, idx: i, identity: "canonical-spec:" + ps.Name, depth: 1,
			})
			if m.index.ExpandedSpecs[i] || m.index.FilterText != "" {
				for r, requirement := range ps.RequirementNames {
					m.index.Items = append(m.index.Items, indexItem{
						kind: indexKindRequirement, idx: i, reqIdx: r,
						identity: "canonical-spec:" + ps.Name + ":requirement:" + requirement,
						depth:    2,
					})
				}
			}
		}
	}

	m.index.Items = append(m.index.Items, indexItem{
		kind: indexKindSection, idx: sectionArchived, identity: "section:history",
	})
	if !m.index.CollapsedSections[sectionArchived] || m.index.FilterText != "" {
		for i := range m.index.ArchiveChanges {
			m.appendChangeHierarchy(m.index.ArchiveChanges[i], i, true)
		}
	}
}

func (m *Model) appendChangeHierarchy(change openspec.Change, changeIdx int, archived bool) {
	prefix := "active:"
	kind := indexKindActive
	if archived {
		prefix = "archive:"
		kind = indexKindArchived
	}
	changeIdentity := prefix + change.Name
	m.index.Items = append(m.index.Items, indexItem{
		kind: kind, idx: changeIdx, archived: archived, identity: changeIdentity, depth: 1,
	})
	if !m.index.ExpandedChanges[changeIdentity] && m.index.FilterText == "" {
		return
	}
	for artifactIdx, artifact := range change.Artifacts {
		artifactIdentity := changeIdentity + ":artifact:" + artifact.ID
		m.index.Items = append(m.index.Items, indexItem{
			kind: indexKindArtifact, idx: changeIdx, artifactIdx: artifactIdx,
			archived: archived, identity: artifactIdentity, depth: 2,
		})
		if len(artifact.Outputs) < 2 || (!m.index.ExpandedArtifacts[artifactIdentity] && m.index.FilterText == "") {
			continue
		}
		for outputIdx, output := range artifact.Outputs {
			m.index.Items = append(m.index.Items, indexItem{
				kind: indexKindArtifactOutput, idx: changeIdx, artifactIdx: artifactIdx, outputIdx: outputIdx,
				archived: archived, identity: artifactIdentity + ":output:" + output.RelativePath, depth: 3,
			})
		}
	}
}

func indexItemByIdentity(items []indexItem, identity string) int {
	for index := range items {
		if items[index].identity == identity {
			return index
		}
	}
	return -1
}

func (m *Model) selectedIndexIdentity() string {
	if m.index.Cursor < 0 || m.index.Cursor >= m.visibleItemCount() {
		return ""
	}
	rawIndex := m.index.Cursor
	if m.index.FilterIndices != nil {
		rawIndex = m.index.FilterIndices[m.index.Cursor]
	}
	if rawIndex < 0 || rawIndex >= len(m.index.Items) {
		return ""
	}
	return m.index.Items[rawIndex].identity
}

func (m *Model) rebuildIndexPreservingCursor() {
	rawCursor := m.rawCursorIndex()
	identity := ""
	var anchor cursorAnchor
	if rawCursor >= 0 && rawCursor < len(m.index.Items) {
		identity = m.index.Items[rawCursor].identity
		anchor = m.cursorAnchor(rawCursor)
	}
	m.buildIndexItems()
	m.applyFilter()
	if identity != "" {
		rawIndex := indexItemByIdentity(m.index.Items, identity)
		if rawIndex < 0 {
			// Enrichment swaps discovery-derived artifact IDs (filename
			// stems) for schema-declared ones: re-anchor the cursor by
			// (change, output path) instead of dropping it to the clamp.
			rawIndex = m.indexItemByAnchor(anchor)
		}
		if rawIndex >= 0 {
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
	}
	m.index.Cursor = min(m.index.Cursor, max(0, m.visibleItemCount()-1))
}

// cursorAnchor is a rebuild-stable pointer to a hierarchy row: the change
// identity plus (when known) a concrete output file path, which survives
// artifact-ID swaps between discovery and status enrichment.
type cursorAnchor struct {
	changeIdentity string
	outputPath     string
}

func (m *Model) rawCursorIndex() int {
	if m.index.Cursor < 0 || m.index.Cursor >= m.visibleItemCount() {
		return -1
	}
	if m.index.FilterIndices == nil {
		return m.index.Cursor
	}
	if m.index.Cursor >= len(m.index.FilterIndices) {
		return -1
	}
	return m.index.FilterIndices[m.index.Cursor]
}

func (m *Model) cursorAnchor(rawIndex int) cursorAnchor {
	if rawIndex < 0 || rawIndex >= len(m.index.Items) {
		return cursorAnchor{}
	}
	item := m.index.Items[rawIndex]
	split := strings.Index(item.identity, ":artifact:")
	if split < 0 {
		return cursorAnchor{}
	}
	anchor := cursorAnchor{changeIdentity: item.identity[:split]}
	change, ok := m.indexItemChange(item)
	if !ok || item.artifactIdx < 0 || item.artifactIdx >= len(change.Artifacts) {
		return anchor
	}
	artifact := change.Artifacts[item.artifactIdx]
	switch item.kind {
	case indexKindArtifactOutput:
		if item.outputIdx >= 0 && item.outputIdx < len(artifact.Outputs) {
			anchor.outputPath = artifact.Outputs[item.outputIdx].RelativePath
		}
	case indexKindArtifact:
		if len(artifact.Outputs) > 0 {
			anchor.outputPath = artifact.Outputs[0].RelativePath
		}
	}
	return anchor
}

func (m *Model) indexItemByAnchor(anchor cursorAnchor) int {
	if anchor.changeIdentity == "" {
		return -1
	}
	prefix := anchor.changeIdentity + ":artifact:"
	firstArtifact := -1
	for i := range m.index.Items {
		item := m.index.Items[i]
		if !strings.HasPrefix(item.identity, prefix) {
			continue
		}
		if item.kind == indexKindArtifact && firstArtifact < 0 {
			firstArtifact = i
		}
		if anchor.outputPath == "" {
			continue
		}
		switch item.kind {
		case indexKindArtifactOutput:
			if strings.HasSuffix(item.identity, ":output:"+anchor.outputPath) {
				return i
			}
		case indexKindArtifact:
			if change, ok := m.indexItemChange(item); ok && item.artifactIdx >= 0 && item.artifactIdx < len(change.Artifacts) {
				for _, output := range change.Artifacts[item.artifactIdx].Outputs {
					if output.RelativePath == anchor.outputPath {
						return i
					}
				}
			}
		}
	}
	return firstArtifact
}

func (m *Model) refreshIndexViewport() {
	content, cursorLine := m.renderIndexContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset() {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset()+m.vp.Height() {
		m.vp.SetYOffset(cursorLine - m.vp.Height() + 1)
	}
}

func (m *Model) isCursorAt(rawIdx int) bool {
	if m.index.FilterIndices != nil {
		return m.index.Cursor < len(m.index.FilterIndices) && m.index.FilterIndices[m.index.Cursor] == rawIdx
	}
	return m.index.Cursor == rawIdx
}

func (m *Model) renderIndexContent() (string, int) {
	contentWidth := m.width - 2
	activeLayout := m.buildActiveRowLayout(contentWidth)
	artifactIDWidth := m.buildArtifactIDWidth()
	artifactSecondaryWidth := m.buildArtifactSecondaryWidth()
	var sb strings.Builder
	line := 0
	cursorLine := 0

	// Precompute spec/archive display metrics
	maxSpecName := 0
	maxReqCount := 0
	for _, ps := range m.projectSpecs {
		if len(ps.Name) > maxSpecName {
			maxSpecName = len(ps.Name)
		}
		if ps.RequirementCount > maxReqCount {
			maxReqCount = ps.RequirementCount
		}
	}
	maxReqDigits := len(strconv.Itoa(maxReqCount))
	maxArchName := 0
	for _, ch := range m.index.ArchiveChanges {
		displayName := ch.Name
		if ch.Schema != "" {
			displayName += " [" + ch.Schema + "]"
		}
		if len(displayName) > maxArchName {
			maxArchName = len(displayName)
		}
	}

	sb.WriteString("\n")
	line++

	for i := 0; i < len(m.index.Items); {
		item := m.index.Items[i]

		if item.kind != indexKindSection {
			i++
			continue
		}

		sectionIdx := item.idx
		isCollapsed := m.index.CollapsedSections[sectionIdx]
		totalCount := sectionItemCount(sectionIdx, m)

		sectionName := sectionNames[sectionIdx]
		header := fmt.Sprintf("%s (%d)", sectionName, totalCount)
		if isCollapsed {
			header += " " + m.theme.Styles.Help.Render("…")
		}

		cursor := m.isCursorAt(i)
		if cursor {
			cursorLine = line
		}
		cursorMark := "  "
		if cursor {
			cursorMark = m.theme.Styles.ProgressDone.Render("▶") + " "
		}
		sb.WriteString(cursorMark + m.theme.Styles.Section.Render(header) + "\n")
		line++

		sb.WriteString("\n")
		line++

		i++

		if isCollapsed {
			for ; i < len(m.index.Items); i++ {
				if m.index.Items[i].kind == indexKindSection {
					break
				}
			}
			if sectionIdx != sectionArchived {
				sb.WriteString("\n")
				line++
			}
			continue
		}

		if totalCount == 0 {
			sb.WriteString(m.theme.Styles.Help.Render("  "+sectionEmptyMsg(sectionIdx)) + "\n")
			line++
			i = m.skipToNextSection(i)
			if sectionIdx != sectionArchived {
				sb.WriteString("\n")
				line++
			}
			continue
		}

		anyVisible := false
		for ; i < len(m.index.Items); i++ {
			if m.index.Items[i].kind == indexKindSection {
				break
			}
			if !m.isItemVisible(i) {
				continue
			}
			anyVisible = true
			childItem := m.index.Items[i]
			childCursor := m.isCursorAt(i)
			if childCursor {
				cursorLine = line
			}

			switch childItem.kind {
			case indexKindActive:
				ch := m.project.Changes[childItem.idx]
				sb.WriteString(m.renderActiveItem(ch, childCursor, activeLayout) + "\n")
			case indexKindArtifact:
				if ch, ok := m.indexItemChange(childItem); ok && childItem.artifactIdx < len(ch.Artifacts) {
					sb.WriteString(m.renderArtifactItem(ch.Artifacts[childItem.artifactIdx], childItem, childCursor, contentWidth, artifactIDWidth, artifactSecondaryWidth) + "\n")
				}
			case indexKindArtifactOutput:
				if ch, ok := m.indexItemChange(childItem); ok && childItem.artifactIdx < len(ch.Artifacts) {
					artifact := ch.Artifacts[childItem.artifactIdx]
					if childItem.outputIdx < len(artifact.Outputs) {
						sb.WriteString(m.renderArtifactOutputItem(artifact.Outputs[childItem.outputIdx], childCursor) + "\n")
					}
				}
			case indexKindSpec:
				ps := m.projectSpecs[childItem.idx]
				pad := strings.Repeat(" ", maxSpecName-len(ps.Name))
				label := m.theme.Styles.Help.Render(fmt.Sprintf("%*d requirements", maxReqDigits, ps.RequirementCount))
				cursorMark := "  "
				name := m.theme.Styles.BaseText.Render(ps.Name)
				if childCursor {
					cursorMark = m.theme.Styles.ProgressDone.Render("▶") + " "
					name = m.theme.Styles.IndexActive.Render(ps.Name)
				}
				sb.WriteString(cursorMark + name + pad + "  " + label + "\n")
			case indexKindRequirement:
				reqMark := "    "
				rName := m.theme.Styles.TaskPending.Render(m.projectSpecs[childItem.idx].RequirementNames[childItem.reqIdx])
				if childCursor {
					reqMark = "  " + m.theme.Styles.ProgressDone.Render("▶") + " "
					rName = m.theme.Styles.IndexActive.Render(m.projectSpecs[childItem.idx].RequirementNames[childItem.reqIdx])
				}
				sb.WriteString(reqMark + rName + "\n")
			case indexKindArchived:
				ch := m.index.ArchiveChanges[childItem.idx]
				sb.WriteString(m.renderArchivedItem(ch, childCursor, maxArchName) + "\n")
			}
			line++
		}

		if !anyVisible {
			sb.WriteString(m.theme.Styles.Help.Render("  No items match '"+m.index.FilterText+"'") + "\n")
			line++
		}

		if sectionIdx != sectionArchived {
			sb.WriteString("\n")
			line++
		}
	}

	return sb.String(), cursorLine
}

func sectionItemCount(sectionIdx int, m *Model) int {
	switch sectionIdx {
	case sectionActive:
		return len(m.project.Changes)
	case sectionSpecs:
		return len(m.projectSpecs)
	case sectionArchived:
		return len(m.index.ArchiveChanges)
	}
	return 0
}

func sectionEmptyMsg(sectionIdx int) string {
	switch sectionIdx {
	case sectionActive:
		return "No active work"
	case sectionSpecs:
		return "No canonical specifications available"
	case sectionArchived:
		return "No history"
	}
	return ""
}

func (m *Model) skipToNextSection(start int) int {
	for i := start; i < len(m.index.Items); i++ {
		if m.index.Items[i].kind == indexKindSection {
			return i
		}
	}
	return len(m.index.Items)
}

type activeRowLayout struct {
	nameWidth       int
	barWidth        int
	schemaWidth     int
	planningWidth   int
	tasksWidth      int
	diagnosticWidth int
}

func (m *Model) buildActiveRowLayout(contentWidth int) activeRowLayout {
	layout := activeRowLayout{nameWidth: 8}
	if m.project == nil {
		return layout
	}
	for _, change := range m.project.Changes {
		layout.nameWidth = max(layout.nameWidth, lipgloss.Width(change.Name))
		layout.schemaWidth = max(layout.schemaWidth, lipgloss.Width(activeSchemaText(change)))
		layout.planningWidth = max(layout.planningWidth, lipgloss.Width(activePlanningText(change)))
		layout.tasksWidth = max(layout.tasksWidth, lipgloss.Width(activeTasksText(change)))
		layout.diagnosticWidth = max(layout.diagnosticWidth, lipgloss.Width(activeDiagnosticText(change)))
	}
	layout.nameWidth = min(layout.nameWidth, 32)
	layout.diagnosticWidth = min(layout.diagnosticWidth, 24)

	metadataWidth := layout.schemaWidth
	if layout.planningWidth > 0 {
		metadataWidth += layout.planningWidth
		if layout.schemaWidth > 0 {
			metadataWidth++
		}
	}
	if layout.tasksWidth > 0 {
		metadataWidth += layout.tasksWidth
		if layout.schemaWidth > 0 || layout.planningWidth > 0 {
			metadataWidth += 3
		}
	}
	if layout.diagnosticWidth > 0 {
		metadataWidth += layout.diagnosticWidth
		if metadataWidth > layout.diagnosticWidth {
			metadataWidth++
		}
	}

	available := contentWidth - 2
	if metadataWidth > 0 {
		available -= metadataWidth + 1
	}
	if layout.tasksWidth > 0 {
		layout.barWidth = available - layout.nameWidth - 1
		if layout.barWidth < 4 {
			layout.nameWidth = max(8, layout.nameWidth-(4-layout.barWidth))
			layout.barWidth = max(4, available-layout.nameWidth-1)
		}
	} else {
		layout.nameWidth = min(layout.nameWidth, available)
	}
	return layout
}

func activeSchemaText(change openspec.Change) string {
	if change.Schema == "" {
		return ""
	}
	return "[" + change.Schema + "]"
}

func activePlanningText(change openspec.Change) string {
	if len(change.Artifacts) == 0 {
		return ""
	}
	if !change.StatusEnriched && change.Diagnostic != "" {
		count := 0
		for range change.Artifacts {
			count++
		}
		return fmt.Sprintf("files %d · status unavailable", count)
	}
	done, total, strays := 0, 0, 0
	for _, artifact := range change.Artifacts {
		if artifact.Source == openspec.ArtifactSourceStatus {
			total++
			if artifact.Status == openspec.ArtifactStatusDone {
				done++
			}
			continue
		}
		if change.StatusEnriched {
			strays++
			continue
		}
		total++
		if artifact.Status == openspec.ArtifactStatusDone {
			done++
		}
	}
	if total == 0 {
		return ""
	}
	text := fmt.Sprintf("planning %d/%d", done, total)
	if strays > 0 {
		text += " · " + strconv.Itoa(strays) + " strays"
	}
	return text
}

func activeTasksText(change openspec.Change) string {
	done, total := taskCounts(change)
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("tasks %d/%d", done, total)
}

func activeDiagnosticText(change openspec.Change) string {
	if change.Diagnostic == "" {
		return ""
	}
	return "! " + change.Diagnostic
}

func (m *Model) renderActiveItem(ch openspec.Change, cursor bool, layout activeRowLayout) string {
	doneTasks, totalTasks := taskCounts(ch)
	cursorMark := "  "
	if cursor {
		cursorMark = m.theme.Styles.ProgressDone.Render("▶") + " "
	}

	name := truncateIndexText(ch.Name, layout.nameWidth)
	name += strings.Repeat(" ", max(0, layout.nameWidth-lipgloss.Width(name)))
	renderedName := m.theme.Styles.BaseText.Render(name)
	if cursor {
		renderedName = m.theme.Styles.IndexActive.Render(name)
	}
	row := cursorMark + renderedName

	if layout.barWidth > 0 {
		bar := strings.Repeat(" ", layout.barWidth)
		if totalTasks > 0 {
			innerWidth := layout.barWidth - 2
			filled := (doneTasks * innerWidth) / totalTasks
			filledStyle := m.theme.Styles.ProgressDone
			if doneTasks == totalTasks {
				filled = innerWidth
				filledStyle = m.theme.Styles.ProgressComplete
			}
			bar = "[" + filledStyle.Render(strings.Repeat("█", filled)) +
				m.theme.Styles.ProgressEmpty.Render(strings.Repeat("░", innerWidth-filled)) + "]"
		}
		row += " " + bar
	}

	var metadata []string
	if layout.schemaWidth > 0 {
		text := activeSchemaText(ch)
		metadata = append(metadata, m.theme.Styles.ProgressDone.Render(text)+strings.Repeat(" ", layout.schemaWidth-lipgloss.Width(text)))
	}
	if layout.planningWidth > 0 {
		text := activePlanningText(ch)
		metadata = append(metadata, m.theme.Styles.Help.Render(text)+strings.Repeat(" ", layout.planningWidth-lipgloss.Width(text)))
	}
	if layout.tasksWidth > 0 {
		text := activeTasksText(ch)
		padding := strings.Repeat(" ", layout.tasksWidth-lipgloss.Width(text))
		if len(metadata) > 0 {
			metadata[len(metadata)-1] += " ·"
		}
		metadata = append(metadata, padding+m.theme.Styles.Help.Render(text))
	}
	if layout.diagnosticWidth > 0 {
		text := elideDiagnosticText(activeDiagnosticText(ch), layout.diagnosticWidth)
		metadata = append(metadata, m.theme.Styles.Error.Render(text)+strings.Repeat(" ", layout.diagnosticWidth-lipgloss.Width(text)))
	}
	if len(metadata) > 0 {
		row += " " + strings.Join(metadata, " ")
	}
	return row
}

func truncateIndexText(text string, width int) string {
	if lipgloss.Width(text) <= width {
		return text
	}
	if width <= 1 {
		return "."
	}
	runes := []rune(text)
	for len(runes) > 0 && lipgloss.Width(string(runes))+1 > width {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "."
}

// elideDiagnosticText keeps both ends of a diagnostic — the failing command
// at the head and the error at the tail — eliding only the middle, so the
// reason a change is degraded stays readable in narrow panes.
func elideDiagnosticText(text string, width int) string {
	if lipgloss.Width(text) <= width {
		return text
	}
	if width <= 1 {
		return "."
	}
	runes := []rune(text)
	head := width * 2 / 3
	tail := width - head - 1
	if head < 1 || tail < 1 || head+tail >= len(runes) {
		return truncateIndexText(text, width)
	}
	return string(runes[:head]) + "…" + string(runes[len(runes)-tail:])
}

func rightAlignIndexRow(left, right string, width int) string {
	padding := width - lipgloss.Width(left) - lipgloss.Width(right)
	if right != "" && padding < 1 {
		padding = 1
	}
	if padding < 0 {
		padding = 0
	}
	return left + strings.Repeat(" ", padding) + right
}

func (m *Model) indexItemChange(item indexItem) (openspec.Change, bool) {
	if item.archived {
		if item.idx >= 0 && item.idx < len(m.index.ArchiveChanges) {
			return m.index.ArchiveChanges[item.idx], true
		}
		return openspec.Change{}, false
	}
	if m.project != nil && item.idx >= 0 && item.idx < len(m.project.Changes) {
		return m.project.Changes[item.idx], true
	}
	return openspec.Change{}, false
}

func (m *Model) buildArtifactIDWidth() int {
	width := 0
	if m.project != nil {
		for _, change := range m.project.Changes {
			for _, artifact := range change.Artifacts {
				width = max(width, lipgloss.Width(artifactDisplayLabel(artifact)))
			}
		}
	}
	for _, change := range m.index.ArchiveChanges {
		for _, artifact := range change.Artifacts {
			width = max(width, lipgloss.Width(artifactDisplayLabel(artifact)))
		}
	}
	return min(width, 24)
}

func (m *Model) buildArtifactSecondaryWidth() int {
	width := 0
	secondaryWidth := func(artifacts []openspec.ChangeArtifact) {
		for _, artifact := range artifacts {
			_, secondary := artifactLabels(artifact)
			width = max(width, lipgloss.Width(secondary))
		}
	}
	if m.project != nil {
		for _, change := range m.project.Changes {
			secondaryWidth(change.Artifacts)
		}
	}
	for _, change := range m.index.ArchiveChanges {
		secondaryWidth(change.Artifacts)
	}
	return min(width, 24)
}

// artifactLabels splits an artifact's name into a stable primary label (the
// schema artifact ID — independent of disk state) and an optional secondary
// label: the display name of the single present output file (e.g. "1-idea"
// for a schema that numbers its templates), suppressed when it matches the
// ID so classic proposal/specs/design/tasks naming stays clean.
func artifactLabels(artifact openspec.ChangeArtifact) (primary, secondary string) {
	primary = artifact.ID
	label := ""
	present := 0
	for _, output := range artifact.Outputs {
		if !output.Present {
			continue
		}
		present++
		if label == "" {
			label = output.DisplayName
		}
	}
	if present == 1 && label != "" && label != artifact.ID {
		secondary = label
	}
	return primary, secondary
}

// artifactDisplayLabel is the single-name view of artifactLabels: the schema
// artifact ID — used by artifact tabs and wherever one stable name is needed.
func artifactDisplayLabel(artifact openspec.ChangeArtifact) string {
	primary, _ := artifactLabels(artifact)
	return primary
}

func (m *Model) renderArtifactItem(artifact openspec.ChangeArtifact, item indexItem, cursor bool, contentWidth, idWidth, secondaryWidth int) string {
	cursorMark := "    "
	if cursor {
		cursorMark = "  " + m.theme.Styles.ProgressDone.Render("▶") + " "
	}
	primary, secondary := artifactLabels(artifact)
	id := truncateIndexText(primary, idWidth)
	name := m.theme.Styles.BaseText.Render(id)
	if cursor {
		name = m.theme.Styles.IndexActive.Render(id)
	}
	left := cursorMark + name + strings.Repeat(" ", max(0, idWidth-lipgloss.Width(id)))
	if secondaryWidth > 0 {
		if secondary != "" {
			rendered := truncateIndexText(secondary, secondaryWidth)
			left += " " + m.theme.Styles.Help.Render(rendered) + strings.Repeat(" ", max(0, secondaryWidth-lipgloss.Width(secondary)))
		} else {
			left += strings.Repeat(" ", 1+secondaryWidth)
		}
	}
	if len(artifact.Requires) > 0 {
		left += "  " + m.theme.Styles.TaskCodeCyan.Render("requires") + " " +
			m.theme.Styles.Help.Render(strings.Join(artifact.Requires, ", "))
	}
	if len(artifact.Outputs) > 1 && !m.index.ExpandedArtifacts[item.identity] {
		left += m.theme.Styles.Help.Render(" …")
	}
	change, _ := m.indexItemChange(item)
	return rightAlignIndexRow(left, m.renderArtifactStatusBadge(change, artifact, item.archived), contentWidth)
}

func (m *Model) renderArtifactStatusBadge(change openspec.Change, artifact openspec.ChangeArtifact, archived bool) string {
	badge := func(style lipgloss.Style, label string) string {
		return style.Bold(true).Reverse(true).Render("[" + label + "]")
	}
	// Archived changes shipped: their discovered artifacts keep rendering as
	// authored history rather than stray or degraded rows.
	if artifact.Source == openspec.ArtifactSourceDiscovered && !archived {
		if change.StatusEnriched {
			return badge(m.theme.Styles.Section, "stray")
		}
		return badge(m.theme.Styles.ProgressDone, "discovered")
	}
	switch artifact.Status {
	case openspec.ArtifactStatusDone:
		return badge(m.theme.Styles.ProgressComplete, "authored")
	case openspec.ArtifactStatusReady:
		return badge(m.theme.Styles.Section, "ready to author")
	case openspec.ArtifactStatusBlocked:
		return badge(m.theme.Styles.Error, "blocked")
	}
	if artifact.Source == openspec.ArtifactSourceDiscovered {
		return badge(m.theme.Styles.ProgressDone, "discovered")
	}
	return badge(m.theme.Styles.Help, "unknown")
}

func (m *Model) renderArtifactOutputItem(output openspec.ArtifactOutput, cursor bool) string {
	cursorMark := "      "
	if cursor {
		cursorMark = "    " + m.theme.Styles.ProgressDone.Render("▶") + " "
	}
	name := output.DisplayName
	if name == "" {
		name = output.RelativePath
	}
	if cursor {
		return cursorMark + m.theme.Styles.IndexActive.Render(name)
	}
	return cursorMark + m.theme.Styles.TaskPending.Render(name)
}

func (m *Model) renderArchivedItem(ch openspec.Change, cursor bool, maxName int) string {
	cursorMark := "  "
	if cursor {
		cursorMark = m.theme.Styles.ProgressDone.Render("▶") + " "
	}

	displayName := ch.Name
	if ch.Schema != "" {
		displayName += " [" + ch.Schema + "]"
	}
	padWidth := maxName - len(displayName)
	if padWidth < 0 {
		padWidth = 0
	}
	pad := strings.Repeat(" ", padWidth)
	date := m.theme.Styles.Help.Render(ch.DisplayDate)
	name := m.theme.Styles.BaseText.Render(displayName) + pad
	if cursor {
		name = m.theme.Styles.IndexActive.Render(displayName) + pad
	}

	text := cursorMark + name + "  " + date
	if ch.Diagnostic != "" {
		text += m.theme.Styles.Error.Render(" ! " + ch.Diagnostic)
	}
	return text
}

const indexViewportContentStart = 3

func (m *Model) indexItemAtContentLine(contentLine int) (int, bool) {
	line := 0

	line++ // leading blank

	for i := 0; i < len(m.index.Items); {
		item := m.index.Items[i]
		if item.kind != indexKindSection {
			i++
			continue
		}

		sectionIdx := item.idx
		isCollapsed := m.index.CollapsedSections[sectionIdx]
		totalCount := sectionItemCount(sectionIdx, m)

		if line == contentLine {
			return i, true
		}
		line++ // section header

		line++ // blank after header

		i++

		if isCollapsed {
			i = m.skipToNextSection(i)
			if sectionIdx != sectionArchived {
				line++ // blank between sections
			}
			continue
		}

		if totalCount == 0 {
			line++ // "No X" message (no item)
			i = m.skipToNextSection(i)
			if sectionIdx != sectionArchived {
				line++ // blank between sections
			}
			continue
		}

		anyVisible := false
		for ; i < len(m.index.Items); i++ {
			if m.index.Items[i].kind == indexKindSection {
				break
			}
			if !m.isItemVisible(i) {
				continue
			}
			anyVisible = true
			if line == contentLine {
				return i, true
			}
			line++
		}

		if !anyVisible {
			line++ // "No items match" message (no item)
		}

		if sectionIdx != sectionArchived {
			line++ // blank between sections
		}
	}

	return 0, false
}

func taskCounts(ch openspec.Change) (int, int) {
	_, output, ok := artifactOutputByID(&ch, "tasks")
	if !ok || output == nil {
		return 0, 0
	}
	done, total := 0, 0
	for _, item := range openspec.ParseTasks(output.Content) {
		if item.Kind == openspec.KindTask {
			total++
			if item.Done {
				done++
			}
		}
	}
	return done, total
}

func sameNames(changes []openspec.Change, diskNames []string) bool {
	if len(changes) != len(diskNames) {
		return false
	}
	diskSet := make(map[string]struct{}, len(diskNames))
	for _, n := range diskNames {
		diskSet[n] = struct{}{}
	}
	for _, ch := range changes {
		if _, ok := diskSet[ch.Name]; !ok {
			return false
		}
	}
	return true
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (m Model) selectedIndexItem() (indexItem, bool) {
	if m.visibleItemCount() == 0 || m.index.Cursor < 0 || m.index.Cursor >= m.visibleItemCount() {
		return indexItem{}, false
	}
	return m.index.Items[m.visibleItemIdx(m.index.Cursor)], true
}

func (m Model) indexItemExpandable(item indexItem) bool {
	switch item.kind {
	case indexKindSection:
		return true
	case indexKindActive:
		return m.project != nil && item.idx < len(m.project.Changes) && len(m.project.Changes[item.idx].Artifacts) > 0
	case indexKindArchived:
		return item.idx < len(m.index.ArchiveChanges) && len(m.index.ArchiveChanges[item.idx].Artifacts) > 0
	case indexKindArtifact:
		if change, ok := m.indexItemChange(item); ok && item.artifactIdx < len(change.Artifacts) {
			return len(change.Artifacts[item.artifactIdx].Outputs) > 1
		}
	case indexKindSpec:
		return item.idx < len(m.projectSpecs) && len(m.projectSpecs[item.idx].RequirementNames) > 0
	}
	return false
}

func (m Model) primaryIndexItem(item indexItem) (tea.Model, tea.Cmd) {
	if m.indexItemExpandable(item) {
		return m.toggleIndexItem(item)
	}
	return m.inspectIndexItem(item)
}

func (m Model) toggleIndexItem(item indexItem) (tea.Model, tea.Cmd) {
	switch item.kind {
	case indexKindSection:
		m.index.CollapsedSections[item.idx] = !m.index.CollapsedSections[item.idx]
	case indexKindActive, indexKindArchived:
		m.index.ExpandedChanges[item.identity] = !m.index.ExpandedChanges[item.identity]
	case indexKindArtifact:
		m.index.ExpandedArtifacts[item.identity] = !m.index.ExpandedArtifacts[item.identity]
	case indexKindSpec:
		m.index.ExpandedSpecs[item.idx] = !m.index.ExpandedSpecs[item.idx]
	default:
		return m, nil
	}
	m.rebuildIndexPreservingCursor()
	m.refreshIndexViewport()
	return m, nil
}

func (m Model) inspectIndexItem(item indexItem) (tea.Model, tea.Cmd) {
	m.returnIndexIdentity = item.identity
	if rawIndex := indexItemByIdentity(m.index.Items, item.identity); rawIndex >= 0 {
		m.returnIndexAnchor = m.cursorAnchor(rawIndex)
	}
	m.artifactRenderCache = make(map[artifactOutputIdentity]string)
	switch item.kind {
	case indexKindActive:
		m.changeIdx = item.idx
		m.mode = ModeNormal
		m.selectDefaultArtifact()
		m.loadTaskItems()
		return m.commitStateChange()
	case indexKindArchived:
		m.index.ArchiveCursor = item.idx
		m.mode = ModeViewingArchive
		m.selectDefaultArtifact()
		return m.commitStateChange()
	case indexKindArtifact, indexKindArtifactOutput:
		if item.archived {
			m.index.ArchiveCursor = item.idx
			m.mode = ModeViewingArchive
		} else {
			m.changeIdx = item.idx
			m.mode = ModeNormal
			m.loadTaskItems()
		}
		if change, ok := m.indexItemChange(item); ok && item.artifactIdx < len(change.Artifacts) {
			artifact := change.Artifacts[item.artifactIdx]
			outputPath := firstArtifactOutputPath(artifact)
			if item.kind == indexKindArtifactOutput && item.outputIdx < len(artifact.Outputs) {
				outputPath = artifact.Outputs[item.outputIdx].RelativePath
			}
			m.selectArtifactOutput(artifact.ID, outputPath)
		}
		return m.commitStateChange()
	case indexKindSpec:
		m.specViewer.Cursor = item.idx
		m.specViewer.JumpTarget = ""
		m.specViewer.FocusMode = false
		m.specViewer.ReqCursor = 0
		m.mode = ModeViewingSpec
		return m.commitStateChange()
	case indexKindRequirement:
		m.specViewer.Cursor = item.idx
		m.specViewer.JumpTarget = m.projectSpecs[item.idx].RequirementNames[item.reqIdx]
		m.specViewer.FocusMode = true
		m.specViewer.ReqCursor = item.reqIdx
		m.mode = ModeViewingSpec
		return m.commitStateChange()
	}
	return m, nil
}

func (m Model) updateIndex(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	keyMap := m.effectiveKeyMap()
	if updated, cmd, handled := m.handleIndexActionInput(msg); handled {
		return updated, cmd
	} else {
		m = updated
	}
	if m.index.FilterActive {
		switch {
		case matchesKey(msg, keyMap.Filter.Cancel):
			m.index.FilterText = m.index.PrevFilterText
			m.index.FilterActive = false
			m.applyFilter()
			m.refreshIndexViewport()
			return m, nil

		case matchesKey(msg, keyMap.Filter.Accept):
			m.index.FilterActive = false
			m.refreshIndexViewport()
			return m, nil

		case matchesKey(msg, keyMap.Filter.Backspace):
			if len(m.index.FilterText) > 0 {
				m.index.FilterText = m.index.FilterText[:len(m.index.FilterText)-1]
				m.applyFilter()
				m.refreshIndexViewport()
			}
			return m, nil

		default:
			if len(msg.String()) == 1 {
				m.index.FilterText += msg.String()
				m.applyFilter()
				m.refreshIndexViewport()
			}
			return m, nil
		}
	}

	switch {

	case matchesKey(msg, keyMap.Index.New):
		m = m.startNewChange()

	case matchesKey(msg, keyMap.Index.Lifecycle):
		if item, ok := m.selectedIndexItem(); ok {
			m = m.startLifecycle(item)
		}

	case matchesKey(msg, keyMap.Index.Edit):
		if item, ok := m.selectedIndexItem(); ok {
			return m.editIndexItem(item)
		}

	case matchesKey(msg, keyMap.Index.Validate):
		if item, ok := m.selectedIndexItem(); ok {
			return m.startValidation(item)
		}

	case matchesKey(msg, keyMap.Index.Undo):
		if !m.readOnly && m.lifecycleUndo != nil {
			m.index.Action = indexActionState{Mode: indexActionConfirmUndo, Operation: indexOperationUndo, Name: m.lifecycleUndo.ChangeName}
		}

	case matchesKey(msg, keyMap.Index.Filter):
		m.index.PrevFilterText = m.index.FilterText
		m.index.FilterText = ""
		m.index.FilterActive = true
		m.index.FilterIndices = nil
		m.refreshIndexViewport()

	case matchesKey(msg, keyMap.Index.Info):
		m.prevMode = m.mode
		m.mode = ModeViewingConfig
		return m.commitStateChange()

	case matchesKey(msg, keyMap.Index.Back):
		if m.index.FilterText != "" {
			m.index.FilterText = ""
			m.index.FilterIndices = nil
			m.refreshIndexViewport()
			return m, nil
		}
		if m.lifecycle != nil {
			m.lifecycle.Cleanup(m.lifecycleUndo)
		}
		return m, tea.Quit

	case matchesKey(msg, keyMap.Index.Down):
		if m.index.Cursor < m.visibleItemCount()-1 {
			m.index.Cursor++
		}
		m.refreshIndexViewport()

	case matchesKey(msg, keyMap.Index.Up):
		if m.index.Cursor > 0 {
			m.index.Cursor--
		}
		m.refreshIndexViewport()

	case matchesKey(msg, keyMap.Index.Inspect):
		if item, ok := m.selectedIndexItem(); ok {
			return m.inspectIndexItem(item)
		}

	case matchesKey(msg, keyMap.Index.Open):
		if item, ok := m.selectedIndexItem(); ok {
			return m.primaryIndexItem(item)
		}

	case matchesKey(msg, keyMap.Index.Toggle):
		if item, ok := m.selectedIndexItem(); ok && m.indexItemExpandable(item) {
			return m.toggleIndexItem(item)
		}

	case matchesKey(msg, keyMap.Index.Sort):
		savedKind := indexKindActive
		savedIdx := -1
		savedReqIdx := 0
		if m.visibleItemCount() > 0 {
			item := m.index.Items[m.visibleItemIdx(m.index.Cursor)]
			savedKind = item.kind
			savedIdx = item.idx
			savedReqIdx = item.reqIdx
		}
		m.index.SortBySuffix = !m.index.SortBySuffix
		m.buildIndexItems()
		m.applyFilter()
		if savedIdx >= 0 {
			if m.index.FilterIndices != nil {
				for i, idx := range m.index.FilterIndices {
					if m.index.Items[idx].kind == savedKind && m.index.Items[idx].idx == savedIdx && m.index.Items[idx].reqIdx == savedReqIdx {
						m.index.Cursor = i
						break
					}
				}
			} else {
				for i, it := range m.index.Items {
					if it.kind == savedKind && it.idx == savedIdx && it.reqIdx == savedReqIdx {
						m.index.Cursor = i
						break
					}
				}
			}
		}
		m.refreshIndexViewport()
	}
	return m, nil
}
