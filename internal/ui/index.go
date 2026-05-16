package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fselich/dossier/internal/openspec"
)

func (m *Model) handleTick() tea.Cmd {
	if m.mode == ModeViewingArchive || m.mode == ModeViewingSpec {
		return nil
	}

	if m.mode == ModeIndex {
		diskChanges := openspec.ListChangeNames()
		diskArchives := openspec.ListArchiveNames()
		diskSpecs := openspec.ListSpecNames()

		archiveNames := make([]string, len(m.archiveChanges))
		for i, ch := range m.archiveChanges {
			archiveNames[i] = filepath.Base(ch.Path)
		}
		specNames := make([]string, len(m.projectSpecs))
		for i, ps := range m.projectSpecs {
			specNames[i] = ps.Name
		}

		if sameNames(m.project.Changes, diskChanges) &&
			sameStrings(archiveNames, diskArchives) &&
			sameStrings(specNames, diskSpecs) {
			return nil
		}

		if p, err := openspec.Load(); err == nil {
			m.project = p
		}
		m.archiveChanges = openspec.ListArchiveChanges()
		m.projectSpecs = openspec.LoadProjectSpecs()
		m.expandedSpecs = make(map[int]bool)
		m.buildIndexItems()
		if m.indexCursor >= len(m.indexItems) {
			m.indexCursor = max(0, len(m.indexItems)-1)
		}
		m.refreshIndexViewport()
		return nil
	}
	// Detect change list additions/removals. Skipped in single-path mode.
	if !m.singlePath {
		diskNames := openspec.ListChangeNames()
		if !sameNames(m.project.Changes, diskNames) {
			currentName := ""
			if ch := m.current(); ch != nil {
				currentName = ch.Name
			}
			if p, err := openspec.Load(); err == nil {
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
	}

	ch := m.current()
	if ch == nil {
		return nil
	}
	fresh := openspec.ReloadChange(*ch)

	// tasks.md: presence or content change → full re-parse with cursor restoration
	if fresh.Tasks.Present != ch.Tasks.Present || fresh.Tasks.Content != ch.Tasks.Content {
		var cursorText string
		if m.taskCursor < len(m.taskItems) && m.taskItems[m.taskCursor].Kind == openspec.KindTask {
			cursorText = m.taskItems[m.taskCursor].Text
		}
		m.project.Changes[m.changeIdx].Tasks = fresh.Tasks
		m.taskItems = openspec.ParseTasks(fresh.Tasks.Content)
		m.taskCursor = openspec.FindCursorByText(m.taskItems, cursorText)
		if m.tab == TabTasks {
			m.refreshTasksViewport()
		}
	}

	// markdown artifacts: presence or content change → update + invalidate render cache
	viewportDirty := false
	if fresh.Proposal.Present != ch.Proposal.Present || fresh.Proposal.Content != ch.Proposal.Content {
		m.project.Changes[m.changeIdx].Proposal = fresh.Proposal
		delete(m.renderCache, TabProposal)
		if m.tab == TabProposal {
			viewportDirty = true
		}
	}
	if fresh.Design.Present != ch.Design.Present || fresh.Design.Content != ch.Design.Content {
		m.project.Changes[m.changeIdx].Design = fresh.Design
		delete(m.renderCache, TabDesign)
		if m.tab == TabDesign {
			viewportDirty = true
		}
	}
	if fresh.Specs.Present != ch.Specs.Present || fresh.Specs.Content != ch.Specs.Content {
		m.project.Changes[m.changeIdx].Specs = fresh.Specs
		m.project.Changes[m.changeIdx].SpecFiles = fresh.SpecFiles
		if m.specIdx >= len(fresh.SpecFiles) {
			m.specIdx = 0
		}
		delete(m.renderCache, TabSpecs)
		if m.tab == TabSpecs {
			viewportDirty = true
		}
	}
	if viewportDirty {
		return m.loadViewport()
	}
	return nil
}

func (m *Model) enterIndex() {
	if len(m.archiveChanges) == 0 {
		m.archiveChanges = openspec.ListArchiveChanges()
	}
	m.projectSpecs = openspec.LoadProjectSpecs()
	m.expandedSpecs = make(map[int]bool)
	m.buildIndexItems()
	m.indexCursor = 0
	m.mode = ModeIndex
	m.vp.Height = m.contentHeight()
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
	m.specOrder = make([]int, n)
	for i := range m.specOrder {
		m.specOrder[i] = i
	}
	if m.specSortBySuffix {
		sort.SliceStable(m.specOrder, func(a, b int) bool {
			return specSuffix(m.projectSpecs[m.specOrder[a]].Name) < specSuffix(m.projectSpecs[m.specOrder[b]].Name)
		})
	}
}

func (m *Model) buildIndexItems() {
	m.buildSpecOrder()
	m.indexItems = nil
	for i := range m.project.Changes {
		m.indexItems = append(m.indexItems, indexItem{kind: indexKindActive, idx: i})
	}
	for _, i := range m.specOrder {
		ps := m.projectSpecs[i]
		m.indexItems = append(m.indexItems, indexItem{kind: indexKindSpec, idx: i})
		if m.expandedSpecs[i] {
			for r := range ps.RequirementNames {
				m.indexItems = append(m.indexItems, indexItem{kind: indexKindRequirement, idx: i, reqIdx: r})
			}
		}
	}
	for i := range m.archiveChanges {
		m.indexItems = append(m.indexItems, indexItem{kind: indexKindArchived, idx: i})
	}
}

func (m *Model) refreshIndexViewport() {
	content, cursorLine := m.renderIndexContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset+m.vp.Height {
		m.vp.SetYOffset(cursorLine - m.vp.Height + 1)
	}
}

func (m *Model) renderIndexContent() (string, int) {
	contentWidth := m.width - 2
	var sb strings.Builder
	line := 0
	cursorLine := 0

	// ── Activos ──────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	line++
	sb.WriteString("  " + sectionStyle.Render("Active Changes") + "\n\n")
	line += 2

	if len(m.project.Changes) == 0 {
		sb.WriteString(helpStyle.Render("  No active changes") + "\n")
		line++
	} else {
		for i, ch := range m.project.Changes {
			cursor := m.indexCursor < len(m.indexItems) &&
				m.indexItems[m.indexCursor].kind == indexKindActive &&
				m.indexItems[m.indexCursor].idx == i
			if cursor {
				cursorLine = line
			}
			sb.WriteString(m.renderActiveItem(ch, cursor, contentWidth) + "\n")
			line++
		}
	}

	sb.WriteString("\n")
	line++

	// ── Specs ─────────────────────────────────────────────────────────────────
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
		for _, i := range m.specOrder {
			ps := m.projectSpecs[i]
			cursor := m.indexCursor < len(m.indexItems) &&
				m.indexItems[m.indexCursor].kind == indexKindSpec &&
				m.indexItems[m.indexCursor].idx == i
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
			if m.expandedSpecs[i] {
				for r, reqName := range ps.RequirementNames {
					reqCursor := m.indexCursor < len(m.indexItems) &&
						m.indexItems[m.indexCursor].kind == indexKindRequirement &&
						m.indexItems[m.indexCursor].idx == i &&
						m.indexItems[m.indexCursor].reqIdx == r
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

	// ── Archivados ────────────────────────────────────────────────────────────
	sb.WriteString("  " + sectionStyle.Render("Archived Changes") + "\n\n")
	line += 2

	if len(m.archiveChanges) == 0 {
		sb.WriteString(helpStyle.Render("  No archived changes") + "\n")
	} else {
		maxName := 0
		for _, ch := range m.archiveChanges {
			if len(ch.Name) > maxName {
				maxName = len(ch.Name)
			}
		}
		for i, ch := range m.archiveChanges {
			cursor := m.indexCursor < len(m.indexItems) &&
				m.indexItems[m.indexCursor].kind == indexKindArchived &&
				m.indexItems[m.indexCursor].idx == i
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
	// layout: 2 (cursor) + nameColWidth + 1 ([) + barSpace + 1 (]) + len(countStr)
	barSpace := contentWidth - 2 - nameColWidth - 2 - len(countStr)
	if barSpace < 4 {
		barSpace = 4
	}
	filled := (done * barSpace) / total
	if done == total {
		filled = barSpace
	}
	bar := "[" + progressDoneStyle.Render(strings.Repeat("█", filled)) +
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

func taskCounts(ch openspec.Change) (done, total int) {
	if !ch.Tasks.Present {
		return 0, 0
	}
	for _, item := range openspec.ParseTasks(ch.Tasks.Content) {
		if item.Kind == openspec.KindTask {
			total++
			if item.Done {
				done++
			}
		}
	}
	return
}

func sameNames(changes []openspec.Change, diskNames []string) bool {
	if len(changes) != len(diskNames) {
		return false
	}
	for i, ch := range changes {
		if ch.Name != diskNames[i] {
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

func extractRequirement(raw, name string) string {
	target := "### Requirement: " + name
	lines := strings.Split(raw, "\n")
	start := -1
	for i, l := range lines {
		if l == target {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	block := []string{lines[start]}
	for _, l := range lines[start+1:] {
		if strings.HasPrefix(l, "### Requirement: ") {
			break
		}
		block = append(block, l)
	}
	return strings.Join(block, "\n")
}
