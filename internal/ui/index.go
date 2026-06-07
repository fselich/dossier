package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func (m *Model) handleTick() tea.Cmd {
	if m.mode == ModeViewingArchive || m.mode == ModeViewingSpec {
		return nil
	}
	if m.mode == ModeIndex {
		return m.pollIndexMode()
	}
	if !m.singlePath {
		if cmd := m.pollNormalModeChanges(); cmd != nil {
			return cmd
		}
	}
	return m.pollNormalModeContent()
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
			ch := &m.project.Changes[i]
			fresh := m.loader.ReloadChange(*ch)
			if fresh.Tasks.Present != ch.Tasks.Present || fresh.Tasks.Content != ch.Tasks.Content {
				ch.Tasks = fresh.Tasks
				needsRefresh = true
			}
		}
		if needsRefresh {
			m.buildIndexItems()
			m.applyFilter()
			if m.index.Cursor >= m.visibleItemCount() {
				m.index.Cursor = max(0, m.visibleItemCount()-1)
			}
			m.refreshIndexViewport()
		}
		return nil
	}

	if p, err := m.loader.LoadFrom(m.root); err == nil {
		m.project = p
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
	m.index.ExpandedSpecs = make(map[int]bool)
	m.buildIndexItems()
	m.applyFilter()
	if m.index.Cursor >= m.visibleItemCount() {
		m.index.Cursor = max(0, m.visibleItemCount()-1)
	}
	m.refreshIndexViewport()
	return nil
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
			m.project = p
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
			m.renderCache = make(map[Tab]string)
			m.tab = m.defaultTab()
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

	if tasksChanged {
		if cursorText != "" {
			m.tasks.Cursor = openspec.FindCursorByText(m.tasks.Items, cursorText)
		}
		if m.tab == TabTasks {
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
	switch item.kind {
	case indexKindActive:
		if item.idx < len(m.project.Changes) {
			return strings.Contains(strings.ToLower(m.project.Changes[item.idx].Name), lowerQuery)
		}
	case indexKindArchived:
		if item.idx < len(m.index.ArchiveChanges) {
			return strings.Contains(strings.ToLower(m.index.ArchiveChanges[item.idx].Name), lowerQuery)
		}
	case indexKindSpec:
		if item.idx < len(m.projectSpecs) {
			return strings.Contains(strings.ToLower(m.projectSpecs[item.idx].Name), lowerQuery)
		}
	case indexKindRequirement:
		if item.idx < len(m.projectSpecs) && item.reqIdx < len(m.projectSpecs[item.idx].RequirementNames) {
			return strings.Contains(strings.ToLower(m.projectSpecs[item.idx].RequirementNames[item.reqIdx]), lowerQuery)
		}
	}
	return false
}

func (m *Model) isItemVisible(idx int) bool {
	if m.index.FilterText == "" {
		return true
	}
	return m.matchesFilter(m.index.Items[idx], strings.ToLower(m.index.FilterText))
}

func (m *Model) applyFilter() {
	if m.index.FilterText == "" {
		m.index.FilterIndices = nil
		return
	}
	lower := strings.ToLower(m.index.FilterText)
	m.index.FilterIndices = nil
	for i := range m.index.Items {
		if m.matchesFilter(m.index.Items[i], lower) {
			m.index.FilterIndices = append(m.index.FilterIndices, i)
		}
	}
	if m.index.Cursor >= len(m.index.FilterIndices) {
		m.index.Cursor = 0
	}
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
	m.index.Items = nil
	for i := range m.project.Changes {
		m.index.Items = append(m.index.Items, indexItem{kind: indexKindActive, idx: i})
	}
	for _, i := range m.index.Order {
		ps := m.projectSpecs[i]
		m.index.Items = append(m.index.Items, indexItem{kind: indexKindSpec, idx: i})
		if m.index.ExpandedSpecs[i] {
			for r := range ps.RequirementNames {
				m.index.Items = append(m.index.Items, indexItem{kind: indexKindRequirement, idx: i, reqIdx: r})
			}
		}
	}
	for i := range m.index.ArchiveChanges {
		m.index.Items = append(m.index.Items, indexItem{kind: indexKindArchived, idx: i})
	}
}

func (m *Model) refreshIndexViewport() {
	m.vp.SetHeight(m.contentHeight())
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
	var sb strings.Builder
	line := 0
	cursorLine := 0

	activeEnd := 0
	for activeEnd < len(m.index.Items) && m.index.Items[activeEnd].kind == indexKindActive {
		activeEnd++
	}
	specEnd := activeEnd
	for specEnd < len(m.index.Items) && (m.index.Items[specEnd].kind == indexKindSpec || m.index.Items[specEnd].kind == indexKindRequirement) {
		specEnd++
	}

	sb.WriteString("\n")
	line++
	sb.WriteString("  " + sectionStyle.Render("Active Changes") + "\n\n")
	line += 2

	if len(m.project.Changes) == 0 {
		sb.WriteString(helpStyle.Render("  No active changes") + "\n")
		line++
	} else {
		anyVisible := false
		for i := 0; i < activeEnd; i++ {
			if !m.isItemVisible(i) {
				continue
			}
			anyVisible = true
			ch := m.project.Changes[m.index.Items[i].idx]
			cursor := m.isCursorAt(i)
			if cursor {
				cursorLine = line
			}
			sb.WriteString(m.renderActiveItem(ch, cursor, contentWidth) + "\n")
			line++
		}
		if !anyVisible {
			sb.WriteString(helpStyle.Render("  No items match '"+m.index.FilterText+"'") + "\n")
			line++
		}
	}

	sb.WriteString("\n")
	line++

	sb.WriteString("  " + sectionStyle.Render("Specifications") + "\n\n")
	line += 2

	if len(m.projectSpecs) == 0 {
		sb.WriteString(helpStyle.Render("  No specifications available") + "\n")
		line++
	} else {
		maxName := 0
		maxReqCount := 0
		for _, ps := range m.projectSpecs {
			if len(ps.Name) > maxName {
				maxName = len(ps.Name)
			}
			if ps.RequirementCount > maxReqCount {
				maxReqCount = ps.RequirementCount
			}
		}
		maxReqDigits := len(strconv.Itoa(maxReqCount))
		anyVisible := false
		for i := activeEnd; i < specEnd; i++ {
			if !m.isItemVisible(i) {
				continue
			}
			anyVisible = true
			item := m.index.Items[i]
			cursor := m.isCursorAt(i)
			if cursor {
				cursorLine = line
			}

			if item.kind == indexKindSpec {
				ps := m.projectSpecs[item.idx]
				pad := strings.Repeat(" ", maxName-len(ps.Name))
				label := helpStyle.Render(fmt.Sprintf("%*d requirements", maxReqDigits, ps.RequirementCount))
				cursorMark := "  "
				name := ps.Name
				if cursor {
					cursorMark = progressDoneStyle.Render("▶") + " "
					name = indexActiveStyle.Render(ps.Name)
				}
				sb.WriteString(cursorMark + name + pad + "  " + label + "\n")
				line++
			} else {
				reqMark := "    "
				rName := taskPendingStyle.Render(m.projectSpecs[item.idx].RequirementNames[item.reqIdx])
				if cursor {
					reqMark = "  " + progressDoneStyle.Render("▶") + " "
					rName = indexActiveStyle.Render(m.projectSpecs[item.idx].RequirementNames[item.reqIdx])
				}
				sb.WriteString(reqMark + rName + "\n")
				line++
			}
		}
		if !anyVisible {
			sb.WriteString(helpStyle.Render("  No items match '"+m.index.FilterText+"'") + "\n")
			line++
		}
	}

	sb.WriteString("\n")
	line++

	sb.WriteString("  " + sectionStyle.Render("Archived Changes") + "\n\n")
	line += 2

	if len(m.index.ArchiveChanges) == 0 {
		sb.WriteString(helpStyle.Render("  No archived changes") + "\n")
	} else {
		maxName := 0
		for _, ch := range m.index.ArchiveChanges {
			if len(ch.Name) > maxName {
				maxName = len(ch.Name)
			}
		}
		anyVisible := false
		for i := specEnd; i < len(m.index.Items); i++ {
			if !m.isItemVisible(i) {
				continue
			}
			anyVisible = true
			ch := m.index.ArchiveChanges[m.index.Items[i].idx]
			cursor := m.isCursorAt(i)
			if cursor {
				cursorLine = line
			}
			sb.WriteString(m.renderArchivedItem(ch, cursor, maxName) + "\n")
			line++
		}
		if !anyVisible {
			sb.WriteString(helpStyle.Render("  No items match '"+m.index.FilterText+"'") + "\n")
			line++
		}
	}

	return sb.String(), cursorLine
}

func (m *Model) renderActiveItem(ch openspec.Change, cursor bool, contentWidth int) string {
	done, total := taskCounts(ch)

	cursorMark := "  "
	if cursor {
		cursorMark = progressDoneStyle.Render("▶") + " "
	}

	const nameColWidth = 32
	name := ch.Name
	if len(name) > nameColWidth {
		name = name[:nameColWidth-1] + "."
	}
	paddedName := name + strings.Repeat(" ", nameColWidth-len(name))

	var renderedName string
	if cursor {
		renderedName = indexActiveStyle.Render(paddedName)
	} else {
		renderedName = paddedName
	}

	if total == 0 {
		return cursorMark + renderedName
	}

	countStr := fmt.Sprintf(" %d/%d", done, total)
	barSpace := contentWidth - 2 - nameColWidth - 2 - len(countStr)
	if barSpace < 4 {
		barSpace = 4
	}
	filled := (done * barSpace) / total
	filledStyle := progressDoneStyle
	if done == total {
		filled = barSpace
		filledStyle = progressCompleteStyle
	}
	bar := "[" + filledStyle.Render(strings.Repeat("█", filled)) +
		progressEmptyStyle.Render(strings.Repeat("░", barSpace-filled)) + "]" +
		helpStyle.Render(countStr)

	return cursorMark + renderedName + bar
}

func (m *Model) renderArchivedItem(ch openspec.Change, cursor bool, maxName int) string {
	cursorMark := "  "
	if cursor {
		cursorMark = progressDoneStyle.Render("▶") + " "
	}

	pad := strings.Repeat(" ", maxName-len(ch.Name))
	date := helpStyle.Render(ch.DisplayDate)
	name := ch.Name + pad
	if cursor {
		name = indexActiveStyle.Render(ch.Name) + pad
	}

	return cursorMark + name + "  " + date
}

const indexViewportContentStart = 3

func (m *Model) indexItemAtContentLine(contentLine int) (int, bool) {
	line := 0

	line += 3

	activeEnd := 0
	for activeEnd < len(m.index.Items) && m.index.Items[activeEnd].kind == indexKindActive {
		activeEnd++
	}

	if activeEnd > 0 {
		for itemIdx := range activeEnd {
			if !m.isItemVisible(itemIdx) {
				continue
			}
			if line == contentLine {
				return itemIdx, true
			}
			line++
		}
	}
	if activeEnd == 0 {
		line++
	}

	line++

	line += 2

	specEnd := activeEnd
	for specEnd < len(m.index.Items) && (m.index.Items[specEnd].kind == indexKindSpec || m.index.Items[specEnd].kind == indexKindRequirement) {
		specEnd++
	}

	if specEnd > activeEnd {
		for itemIdx := activeEnd; itemIdx < specEnd; itemIdx++ {
			if !m.isItemVisible(itemIdx) {
				continue
			}
			if line == contentLine {
				return itemIdx, true
			}
			line++
		}
	}
	if specEnd == activeEnd {
		line++
	}

	line++

	line += 2

	if specEnd >= len(m.index.Items) {
		line++
	}

	for itemIdx := specEnd; itemIdx < len(m.index.Items); itemIdx++ {
		if !m.isItemVisible(itemIdx) {
			continue
		}
		if line == contentLine {
			return itemIdx, true
		}
		line++
	}

	return 0, false
}

func taskCounts(ch openspec.Change) (int, int) {
	if !ch.Tasks.Present {
		return 0, 0
	}
	done, total := 0, 0
	for _, item := range openspec.ParseTasks(ch.Tasks.Content) {
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

func (m Model) updateIndex(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.index.FilterActive {
		switch msg.String() {
		case "esc":
			m.index.FilterText = m.index.PrevFilterText
			m.index.FilterActive = false
			m.applyFilter()
			m.refreshIndexViewport()
			return m, nil

		case "enter":
			m.index.FilterActive = false
			m.refreshIndexViewport()
			return m, nil

		case "backspace":
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

	switch msg.String() {

	case "/":
		m.index.PrevFilterText = m.index.FilterText
		m.index.FilterText = ""
		m.index.FilterActive = true
		m.index.FilterIndices = nil
		m.refreshIndexViewport()

	case "i":
		m.prevMode = m.mode
		m.mode = ModeViewingConfig
		return m.commitStateChange()

	case "esc":
		if m.index.FilterText != "" {
			m.index.FilterText = ""
			m.index.FilterIndices = nil
			m.refreshIndexViewport()
			return m, nil
		}
		return m, tea.Quit

	case "j", "down":
		if m.index.Cursor < m.visibleItemCount()-1 {
			m.index.Cursor++
		}
		m.refreshIndexViewport()

	case "k", "up":
		if m.index.Cursor > 0 {
			m.index.Cursor--
		}
		m.refreshIndexViewport()

	case "enter":
		if m.visibleItemCount() > 0 {
			item := m.index.Items[m.visibleItemIdx(m.index.Cursor)]
			m.renderCache = make(map[Tab]string)
			if item.kind == indexKindActive {
				m.changeIdx = item.idx
				m.mode = ModeNormal
				m.tab = m.defaultTab()
				m.loadTaskItems()
				return m.commitStateChange()
			}
			if item.kind == indexKindSpec {
				m.specViewer.Cursor = item.idx
				m.specViewer.JumpTarget = ""
				m.specViewer.FocusMode = false
				m.specViewer.ReqCursor = 0
				m.mode = ModeViewingSpec
				return m.commitStateChange()
			}
			if item.kind == indexKindRequirement {
				m.specViewer.Cursor = item.idx
				m.specViewer.JumpTarget = m.projectSpecs[item.idx].RequirementNames[item.reqIdx]
				m.specViewer.FocusMode = true
				m.specViewer.ReqCursor = item.reqIdx
				m.mode = ModeViewingSpec
				return m.commitStateChange()
			}
			m.index.ArchiveCursor = item.idx
			m.tab = firstAvailableTab(m.index.ArchiveChanges[item.idx])
			m.mode = ModeViewingArchive
			return m.commitStateChange()
		}

	case "space":
		if m.visibleItemCount() > 0 {
			item := m.index.Items[m.visibleItemIdx(m.index.Cursor)]
			if item.kind == indexKindSpec {
				specIdx := item.idx
				m.index.ExpandedSpecs[specIdx] = !m.index.ExpandedSpecs[specIdx]
				m.buildIndexItems()
				m.applyFilter()
				m.index.Cursor = 0
				if m.index.FilterIndices != nil {
					for i, idx := range m.index.FilterIndices {
						if m.index.Items[idx].kind == indexKindSpec && m.index.Items[idx].idx == specIdx {
							m.index.Cursor = i
							break
						}
					}
				} else {
					for i, it := range m.index.Items {
						if it.kind == indexKindSpec && it.idx == specIdx {
							m.index.Cursor = i
							break
						}
					}
				}
				if m.index.Cursor >= m.visibleItemCount() {
					m.index.Cursor = max(0, m.visibleItemCount()-1)
				}
				m.refreshIndexViewport()
			}
		}

	case "s":
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
