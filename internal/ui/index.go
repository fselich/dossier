package ui

import (
	"fmt"
	"path/filepath"
	"sort"
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
			if m.index.Cursor >= len(m.index.Items) {
				m.index.Cursor = max(0, len(m.index.Items)-1)
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
	if m.index.Cursor >= len(m.index.Items) {
		m.index.Cursor = max(0, len(m.index.Items)-1)
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
	content, cursorLine := m.renderIndexContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset() {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset()+m.vp.Height() {
		m.vp.SetYOffset(cursorLine - m.vp.Height() + 1)
	}
}

func (m *Model) renderIndexContent() (string, int) {
	contentWidth := m.width - 2
	var sb strings.Builder
	line := 0
	cursorLine := 0

	sb.WriteString("\n")
	line++
	sb.WriteString("  " + sectionStyle.Render("Active Changes") + "\n\n")
	line += 2

	if len(m.project.Changes) == 0 {
		sb.WriteString(helpStyle.Render("  No active changes") + "\n")
		line++
	} else {
		for i, ch := range m.project.Changes {
			cursor := m.index.Cursor < len(m.index.Items) &&
				m.index.Items[m.index.Cursor].kind == indexKindActive &&
				m.index.Items[m.index.Cursor].idx == i
			if cursor {
				cursorLine = line
			}
			sb.WriteString(m.renderActiveItem(ch, cursor, contentWidth) + "\n")
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
		for _, ps := range m.projectSpecs {
			if len(ps.Name) > maxName {
				maxName = len(ps.Name)
			}
		}
		for _, i := range m.index.Order {
			ps := m.projectSpecs[i]
			cursor := m.index.Cursor < len(m.index.Items) &&
				m.index.Items[m.index.Cursor].kind == indexKindSpec &&
				m.index.Items[m.index.Cursor].idx == i
			if cursor {
				cursorLine = line
			}
			pad := strings.Repeat(" ", maxName-len(ps.Name))
			label := helpStyle.Render(fmt.Sprintf("%d requirements", ps.RequirementCount))
			cursorMark := "  "
			name := ps.Name
			if cursor {
				cursorMark = progressDoneStyle.Render("▶") + " "
				name = indexActiveStyle.Render(ps.Name)
			}
			sb.WriteString(cursorMark + name + pad + "  " + label + "\n")
			line++
			if m.index.ExpandedSpecs[i] {
				for r, reqName := range ps.RequirementNames {
					reqCursor := m.index.Cursor < len(m.index.Items) &&
						m.index.Items[m.index.Cursor].kind == indexKindRequirement &&
						m.index.Items[m.index.Cursor].idx == i &&
						m.index.Items[m.index.Cursor].reqIdx == r
					if reqCursor {
						cursorLine = line
					}
					reqMark := "    "
					rName := taskPendingStyle.Render(reqName)
					if reqCursor {
						reqMark = "  " + progressDoneStyle.Render("▶") + " "
						rName = indexActiveStyle.Render(reqName)
					}
					sb.WriteString(reqMark + rName + "\n")
					line++
				}
			}
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
		for i, ch := range m.index.ArchiveChanges {
			cursor := m.index.Cursor < len(m.index.Items) &&
				m.index.Items[m.index.Cursor].kind == indexKindArchived &&
				m.index.Items[m.index.Cursor].idx == i
			if cursor {
				cursorLine = line
			}
			sb.WriteString(m.renderArchivedItem(ch, cursor, maxName) + "\n")
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
	switch msg.String() {

	case "i":
		m.prevMode = m.mode
		m.mode = ModeViewingConfig
		return m.commitStateChange()

	case "esc":
		return m, tea.Quit

	case "j", "down":
		if m.index.Cursor < len(m.index.Items)-1 {
			m.index.Cursor++
		}
		m.refreshIndexViewport()

	case "k", "up":
		if m.index.Cursor > 0 {
			m.index.Cursor--
		}
		m.refreshIndexViewport()

	case "enter":
		if len(m.index.Items) > 0 {
			item := m.index.Items[m.index.Cursor]
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
		if len(m.index.Items) > 0 {
			item := m.index.Items[m.index.Cursor]
			if item.kind == indexKindSpec {
				specIdx := item.idx
				m.index.ExpandedSpecs[specIdx] = !m.index.ExpandedSpecs[specIdx]
				m.buildIndexItems()
				m.index.Cursor = 0
				for i, it := range m.index.Items {
					if it.kind == indexKindSpec && it.idx == specIdx {
						m.index.Cursor = i
						break
					}
				}
				if m.index.Cursor >= len(m.index.Items) {
					m.index.Cursor = max(0, len(m.index.Items)-1)
				}
				m.refreshIndexViewport()
			}
		}

	case "s":
		savedKind := indexKindActive
		savedIdx := -1
		savedReqIdx := 0
		if m.index.Cursor < len(m.index.Items) {
			item := m.index.Items[m.index.Cursor]
			savedKind = item.kind
			savedIdx = item.idx
			savedReqIdx = item.reqIdx
		}
		m.index.SortBySuffix = !m.index.SortBySuffix
		m.buildIndexItems()
		if savedIdx >= 0 {
			for i, it := range m.index.Items {
				if it.kind == savedKind && it.idx == savedIdx && it.reqIdx == savedReqIdx {
					m.index.Cursor = i
					break
				}
			}
		}
		m.refreshIndexViewport()
	}
	return m, nil
}
