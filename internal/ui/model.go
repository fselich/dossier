package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/fselich/dossier/internal/openspec"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeIndex
	ModeViewingArchive
	ModeViewingSpec
)

type Tab int

const (
	TabProposal Tab = iota
	TabDesign
	TabSpecs
	TabTasks
	tabCount
)

var tabLabels = [tabCount]string{"proposal", "design", "specs", "tasks"}

type errClearMsg struct{}
type editorReturnMsg struct{}

// renderedMsg carries async glamour output back to the event loop.
type renderedMsg struct {
	tab     Tab
	content string
}

// specRenderedMsg carries async glamour output for ModeViewingSpec.
type specRenderedMsg struct {
	content  string
	jumpLine int // line offset to scroll to after render; 0 = start of document
}

type tickMsg time.Time

type indexItemKind int

const (
	indexKindActive      indexItemKind = iota
	indexKindArchived
	indexKindSpec
	indexKindRequirement
)

type indexItem struct {
	kind   indexItemKind
	idx    int // into project.Changes (active), archiveChanges (archived), or projectSpecs (spec/requirement)
	reqIdx int // index into projectSpecs[idx].RequirementNames; only used for indexKindRequirement
}

type Model struct {
	project   *openspec.Project
	changeIdx int
	tab       Tab

	vp      viewport.Model
	vpReady bool

	taskItems  []openspec.TaskItem
	taskCursor int

	specIdx int

	errMsg     string
	loading    bool // true while glamour renders in background
	singlePath bool // true when launched with an explicit change path

	width, height int

	renderCache map[Tab]string

	mode           Mode
	archiveChanges []openspec.Change
	archiveCursor  int // which archived change is viewed in ModeViewingArchive
	indexItems     []indexItem
	indexCursor    int
	expandedSpecs  map[int]bool
	projectSpecs     []openspec.ProjectSpec
	specViewerCursor int    // which projectSpec is shown in ModeViewingSpec
	specJumpTarget   string // requirement name to scroll to when entering ModeViewingSpec; empty = top
	specFocusMode    bool   // true when ModeViewingSpec shows only the selected requirement
	specReqCursor    int    // index into projectSpecs[specViewerCursor].RequirementNames in focus mode
}

func New(project *openspec.Project) Model {
	m := Model{project: project, renderCache: make(map[Tab]string)}
	if len(project.Changes) > 0 {
		m.tab = m.defaultTab()
		m.loadTaskItems()
	}
	return m
}

func NewSinglePath(project *openspec.Project) Model {
	m := New(project)
	m.singlePath = true
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentH := m.contentHeight()
		if !m.vpReady {
			m.vp = viewport.New(m.width-2, contentH)
			m.vpReady = true
		} else {
			m.vp.Width = m.width - 2
			m.vp.Height = contentH
		}
		m.renderCache = make(map[Tab]string)
		return m, m.loadViewport()

	case renderedMsg:
		m.renderCache[msg.tab] = msg.content
		m.loading = false
		if m.tab == msg.tab {
			m.vp.SetContent(msg.content)
			m.vp.GotoTop()
		}
		return m, nil

	case specRenderedMsg:
		m.loading = false
		if m.mode == ModeViewingSpec {
			m.vp.SetContent(msg.content)
			if msg.jumpLine > 0 {
				m.vp.SetYOffset(msg.jumpLine)
			} else {
				m.vp.GotoTop()
			}
		}
		return m, nil

	case tickMsg:
		cmd := m.handleTick()
		nextTick := tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
		return m, tea.Batch(nextTick, cmd)

	case editorReturnMsg:
		ch := m.current()
		if ch != nil {
			fresh := openspec.ReloadChange(*ch)
			var cursorText string
			if m.taskCursor < len(m.taskItems) && m.taskItems[m.taskCursor].Kind == openspec.KindTask {
				cursorText = m.taskItems[m.taskCursor].Text
			}
			if fresh.Tasks.Present != ch.Tasks.Present || fresh.Tasks.Content != ch.Tasks.Content {
				m.project.Changes[m.changeIdx].Tasks = fresh.Tasks
				m.taskItems = openspec.ParseTasks(fresh.Tasks.Content)
				m.taskCursor = openspec.FindCursorByText(m.taskItems, cursorText)
			}
			if fresh.Proposal.Present != ch.Proposal.Present || fresh.Proposal.Content != ch.Proposal.Content {
				m.project.Changes[m.changeIdx].Proposal = fresh.Proposal
				delete(m.renderCache, TabProposal)
			}
			if fresh.Design.Present != ch.Design.Present || fresh.Design.Content != ch.Design.Content {
				m.project.Changes[m.changeIdx].Design = fresh.Design
				delete(m.renderCache, TabDesign)
			}
			if fresh.Specs.Present != ch.Specs.Present || fresh.Specs.Content != ch.Specs.Content {
				m.project.Changes[m.changeIdx].Specs = fresh.Specs
				m.project.Changes[m.changeIdx].SpecFiles = fresh.SpecFiles
				if m.specIdx >= len(fresh.SpecFiles) {
					m.specIdx = 0
				}
				delete(m.renderCache, TabSpecs)
			}
		}
		return m, m.loadViewport()

	case errClearMsg:
		m.errMsg = ""
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "a":
			if m.mode == ModeNormal || m.mode == ModeViewingArchive {
				m.enterIndex()
			}

		case "esc":
			switch m.mode {
			case ModeNormal, ModeViewingArchive:
				m.enterIndex()
			case ModeIndex:
				return m, tea.Quit
			case ModeViewingSpec:
				specIdx := m.specViewerCursor
				jumpTarget := m.specJumpTarget
				wasFocusMode := m.specFocusMode
				m.enterIndex()
				if wasFocusMode && jumpTarget != "" {
					// Expand the spec and place cursor on the requirement we were viewing.
					m.expandedSpecs[specIdx] = true
					m.buildIndexItems()
					for j, it := range m.indexItems {
						if it.kind == indexKindRequirement && it.idx == specIdx &&
							it.reqIdx < len(m.projectSpecs[specIdx].RequirementNames) &&
							m.projectSpecs[specIdx].RequirementNames[it.reqIdx] == jumpTarget {
							m.indexCursor = j
							break
						}
					}
				} else {
					for i, item := range m.indexItems {
						if item.kind == indexKindSpec && item.idx == specIdx {
							m.indexCursor = i
							break
						}
					}
				}
				m.refreshIndexViewport()
			}

		case "enter":
			if m.mode == ModeIndex && len(m.indexItems) > 0 {
				item := m.indexItems[m.indexCursor]
				m.renderCache = make(map[Tab]string)
				if item.kind == indexKindActive {
					m.changeIdx = item.idx
					m.mode = ModeNormal
					m.tab = m.defaultTab()
					m.loadTaskItems()
					m.vp.Height = m.contentHeight()
					return m, m.loadViewport()
				}
				if item.kind == indexKindSpec {
					m.specViewerCursor = item.idx
					m.specJumpTarget = ""
					m.specFocusMode = false
					m.specReqCursor = 0
					m.mode = ModeViewingSpec
					m.vp.Height = m.contentHeight()
					return m, m.loadViewport()
				}
				if item.kind == indexKindRequirement {
					m.specViewerCursor = item.idx
					m.specJumpTarget = m.projectSpecs[item.idx].RequirementNames[item.reqIdx]
					m.specFocusMode = true
					m.specReqCursor = item.reqIdx
					m.mode = ModeViewingSpec
					m.vp.Height = m.contentHeight()
					return m, m.loadViewport()
				}
				// archived
				m.archiveCursor = item.idx
				m.tab = firstAvailableTab(m.archiveChanges[item.idx])
				m.mode = ModeViewingArchive
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}

		case "h":
			if m.mode == ModeViewingSpec && m.specFocusMode {
				ps := m.projectSpecs[m.specViewerCursor]
				if len(ps.RequirementNames) > 0 {
					m.specReqCursor = (m.specReqCursor - 1 + len(ps.RequirementNames)) % len(ps.RequirementNames)
					m.specJumpTarget = ps.RequirementNames[m.specReqCursor]
					return m, m.loadViewport()
				}
			}
			if m.mode == ModeNormal && len(m.project.Changes) > 0 {
				m.changeIdx = (m.changeIdx - 1 + len(m.project.Changes)) % len(m.project.Changes)
				m.renderCache = make(map[Tab]string)
				m.loadTaskItems()
				m.tab = m.defaultTab()
				m.specIdx = 0
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}

		case "l":
			if m.mode == ModeViewingSpec && m.specFocusMode {
				ps := m.projectSpecs[m.specViewerCursor]
				if len(ps.RequirementNames) > 0 {
					m.specReqCursor = (m.specReqCursor + 1) % len(ps.RequirementNames)
					m.specJumpTarget = ps.RequirementNames[m.specReqCursor]
					return m, m.loadViewport()
				}
			}
			if m.mode == ModeNormal && len(m.project.Changes) > 0 {
				m.changeIdx = (m.changeIdx + 1) % len(m.project.Changes)
				m.renderCache = make(map[Tab]string)
				m.loadTaskItems()
				m.tab = m.defaultTab()
				m.specIdx = 0
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}

		case "1":
			if m.tabAvailable(TabProposal) {
				m.tab = TabProposal
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}
		case "2":
			if m.tabAvailable(TabDesign) {
				m.tab = TabDesign
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}
		case "3":
			if m.mode != ModeIndex && m.tabAvailable(TabSpecs) {
				if m.tab == TabSpecs {
					ch := m.current()
					if ch != nil && len(ch.SpecFiles) > 1 {
						m.specIdx = (m.specIdx + 1) % len(ch.SpecFiles)
						delete(m.renderCache, TabSpecs)
					}
				} else {
					m.tab = TabSpecs
					m.specIdx = 0
				}
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}
		case "4":
			if m.mode != ModeIndex && m.tabAvailable(TabTasks) {
				m.tab = TabTasks
				m.vp.Height = m.contentHeight()
				return m, m.loadViewport()
			}

		case "j", "down":
			switch m.mode {
			case ModeIndex:
				if m.indexCursor < len(m.indexItems)-1 {
					m.indexCursor++
				}
				m.refreshIndexViewport()
			default:
				if m.tab == TabTasks && m.mode == ModeNormal {
					m.moveCursorDown()
					m.refreshTasksViewport()
				} else {
					m.vp.LineDown(1)
				}
			}

		case "k", "up":
			switch m.mode {
			case ModeIndex:
				if m.indexCursor > 0 {
					m.indexCursor--
				}
				m.refreshIndexViewport()
			default:
				if m.tab == TabTasks && m.mode == ModeNormal {
					m.moveCursorUp()
					m.refreshTasksViewport()
				} else {
					m.vp.LineUp(1)
				}
			}

		case " ":
			if m.mode == ModeNormal && m.tab == TabTasks {
				return m, m.doToggle()
			}
			if m.mode == ModeIndex && len(m.indexItems) > 0 {
				item := m.indexItems[m.indexCursor]
				if item.kind == indexKindSpec {
					specIdx := item.idx
					m.expandedSpecs[specIdx] = !m.expandedSpecs[specIdx]
					m.buildIndexItems()
					// Restore cursor to the spec item (it may have shifted).
					m.indexCursor = 0
					for i, it := range m.indexItems {
						if it.kind == indexKindSpec && it.idx == specIdx {
							m.indexCursor = i
							break
						}
					}
					if m.indexCursor >= len(m.indexItems) {
						m.indexCursor = max(0, len(m.indexItems)-1)
					}
					m.refreshIndexViewport()
				}
				// Space on a requirement item: no-op.
			}

		case "e":
			if m.mode == ModeNormal && m.tabAvailable(m.tab) {
				path := m.artifactPath()
				if path != "" {
					editor := os.Getenv("EDITOR")
					if editor == "" {
						editor = "vi"
					}
					cmd := exec.Command(editor, path)
					return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
						return editorReturnMsg{}
					})
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if !m.vpReady {
		return ""
	}
	if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		return m.viewIndex()
	}
	if len(m.project.Changes) == 0 && m.mode == ModeNormal {
		return m.emptyView()
	}
	rows := []string{
		m.boxTop(),
		m.addBorderSides(m.renderHeader()),
		m.addBorderSides(m.renderTabBar()),
		m.boxInnerSep(),
	}
	if m.hasSpecSubnav() {
		rows = append(rows, m.addBorderSides(m.renderSpecSubnav()))
	}
	rows = append(rows,
		m.addBorderSides(m.vp.View()),
		m.boxInnerSep(),
		m.addBorderSides(m.renderHelpBar()),
		m.boxBottom(),
	)
	return strings.Join(rows, "\n")
}

func (m *Model) viewIndex() string {
	rows := []string{
		m.boxTop(),
		m.addBorderSides(m.renderHeader()),
		m.boxInnerSep(),
		m.addBorderSides(m.vp.View()),
		m.boxInnerSep(),
		m.addBorderSides(m.renderHelpBar()),
		m.boxBottom(),
	}
	return strings.Join(rows, "\n")
}

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

// ── index ─────────────────────────────────────────────────────────────────────

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

func (m *Model) buildIndexItems() {
	m.indexItems = nil
	for i := range m.project.Changes {
		m.indexItems = append(m.indexItems, indexItem{kind: indexKindActive, idx: i})
	}
	for i, ps := range m.projectSpecs {
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
		for i, ps := range m.projectSpecs {
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

// ── viewport loading ──────────────────────────────────────────────────────────

func (m *Model) loadViewport() tea.Cmd {
	if !m.vpReady {
		return nil
	}
	if m.mode == ModeIndex {
		m.refreshIndexViewport()
		return nil
	}
	if m.mode == ModeViewingSpec {
		if m.specViewerCursor >= len(m.projectSpecs) {
			m.vp.SetContent("  (spec not available)")
			return nil
		}
		raw := m.projectSpecs[m.specViewerCursor].Content
		if raw == "" {
			m.vp.SetContent("  (spec not available)")
			return nil
		}
		m.loading = true
		m.vp.SetContent("\n  Cargando...")
		width := m.width - 2
		if width < 20 {
			width = 80
		}
		if m.specFocusMode {
			jumpTarget := m.specJumpTarget
			return func() tea.Msg {
				block := extractRequirement(raw, jumpTarget)
				if block == "" {
					return specRenderedMsg{content: "  (spec not available)"}
				}
				r, err := glamour.NewTermRenderer(
					glamour.WithStandardStyle("dark"),
					glamour.WithWordWrap(width),
				)
				if err != nil {
					return specRenderedMsg{content: block}
				}
				out, err := r.Render(block)
				if err != nil {
					return specRenderedMsg{content: block}
				}
				return specRenderedMsg{content: out}
			}
		}
		jumpTarget := m.specJumpTarget
		ansiRe := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		return func() tea.Msg {
			r, err := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(width),
			)
			if err != nil {
				return specRenderedMsg{content: raw}
			}
			out, err := r.Render(raw)
			if err != nil {
				return specRenderedMsg{content: raw}
			}
			jumpLine := 0
			if jumpTarget != "" {
				for i, l := range strings.Split(out, "\n") {
					if strings.Contains(ansiRe.ReplaceAllString(l, ""), jumpTarget) {
						jumpLine = i
						break
					}
				}
			}
			return specRenderedMsg{content: out, jumpLine: jumpLine}
		}
	}
	if m.tab == TabTasks && m.mode == ModeNormal {
		m.refreshTasksViewport()
		return nil
	}

	// Cache hit — instant.
	if cached, ok := m.renderCache[m.tab]; ok {
		m.vp.SetContent(cached)
		return nil
	}

	ch := m.current()
	if ch == nil {
		m.vp.SetContent("")
		return nil
	}
	var raw string
	switch m.tab {
	case TabProposal:
		raw = ch.Proposal.Content
	case TabDesign:
		raw = ch.Design.Content
	case TabSpecs:
		if m.specIdx < len(ch.SpecFiles) {
			raw = ch.SpecFiles[m.specIdx].Content
		}
	case TabTasks:
		raw = ch.Tasks.Content
	}
	if raw == "" {
		m.vp.SetContent("  (artifact not available)")
		return nil
	}

	// Show placeholder immediately, render in background.
	m.loading = true
	m.vp.SetContent("\n  Cargando...")

	tab := m.tab
	width := m.width - 2
	if width < 20 {
		width = 80
	}
	return func() tea.Msg {
		r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return renderedMsg{tab: tab, content: raw}
		}
		out, err := r.Render(raw)
		if err != nil {
			return renderedMsg{tab: tab, content: raw}
		}
		return renderedMsg{tab: tab, content: out}
	}
}

func (m *Model) refreshTasksViewport() {
	content, cursorLine := m.renderTasksContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset+m.vp.Height {
		m.vp.SetYOffset(cursorLine - m.vp.Height + 1)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (m *Model) current() *openspec.Change {
	if m.mode == ModeViewingArchive {
		return m.currentArchive()
	}
	if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		return nil
	}
	if len(m.project.Changes) == 0 {
		return nil
	}
	return &m.project.Changes[m.changeIdx]
}

func firstAvailableTab(ch openspec.Change) Tab {
	if ch.Proposal.Present {
		return TabProposal
	}
	if ch.Design.Present {
		return TabDesign
	}
	if ch.Specs.Present {
		return TabSpecs
	}
	if ch.Tasks.Present {
		return TabTasks
	}
	return TabProposal
}

func (m *Model) tabAvailable(t Tab) bool {
	ch := m.current()
	if ch == nil {
		return false
	}
	switch t {
	case TabProposal:
		return ch.Proposal.Present
	case TabDesign:
		return ch.Design.Present
	case TabTasks:
		return ch.Tasks.Present
	case TabSpecs:
		return ch.Specs.Present
	}
	return false
}

func (m *Model) defaultTab() Tab {
	for t := Tab(0); t < tabCount; t++ {
		if m.tabAvailable(t) {
			return t
		}
	}
	return TabProposal
}

func (m *Model) artifactPath() string {
	ch := m.current()
	if ch == nil {
		return ""
	}
	switch m.tab {
	case TabProposal:
		return filepath.Join(ch.Path, "proposal.md")
	case TabDesign:
		return filepath.Join(ch.Path, "design.md")
	case TabTasks:
		return filepath.Join(ch.Path, "tasks.md")
	case TabSpecs:
		if m.specIdx < len(ch.SpecFiles) {
			specsDir := filepath.Join(ch.Path, "specs")
			entries, err := os.ReadDir(specsDir)
			if err != nil {
				return ""
			}
			dirIdx := 0
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				if dirIdx == m.specIdx {
					p := filepath.Join(specsDir, e.Name(), "spec.md")
					if _, err := os.Stat(p); err == nil {
						return p
					}
					return ""
				}
				dirIdx++
			}
		}
	}
	return ""
}

func (m *Model) currentArchive() *openspec.Change {
	if m.archiveCursor < len(m.archiveChanges) {
		return &m.archiveChanges[m.archiveCursor]
	}
	return nil
}

func (m *Model) contentHeight() int {
	if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		// top+bottom borders + header + 2 inner seps + helpBar (no tab bar)
		h := m.height - 6
		if h < 1 {
			h = 1
		}
		return h
	}
	h := m.height - 7 // top+bottom borders + header + tabBar + 2 inner seps + helpBar
	if m.hasSpecSubnav() {
		h--
	}
	if h < 1 {
		h = 1
	}
	return h
}

func (m *Model) loadTaskItems() {
	ch := m.current()
	if ch == nil || !ch.Tasks.Present {
		m.taskItems = nil
		m.taskCursor = 0
		return
	}
	m.taskItems = openspec.ParseTasks(ch.Tasks.Content)
	m.taskCursor = m.firstTaskIdx()
}

func (m *Model) firstTaskIdx() int {
	for i, item := range m.taskItems {
		if item.Kind == openspec.KindTask {
			return i
		}
	}
	return 0
}

func (m *Model) moveCursorDown() {
	for i := m.taskCursor + 1; i < len(m.taskItems); i++ {
		if m.taskItems[i].Kind == openspec.KindTask {
			m.taskCursor = i
			return
		}
	}
}

func (m *Model) moveCursorUp() {
	for i := m.taskCursor - 1; i >= 0; i-- {
		if m.taskItems[i].Kind == openspec.KindTask {
			m.taskCursor = i
			return
		}
	}
}

func (m *Model) doToggle() tea.Cmd {
	if len(m.taskItems) == 0 || m.taskCursor >= len(m.taskItems) {
		return nil
	}
	if m.taskItems[m.taskCursor].Kind != openspec.KindTask {
		return nil
	}
	ch := m.current()
	if ch == nil {
		return nil
	}
	if err := openspec.ToggleTask(ch.Path+"/tasks.md", m.taskItems, m.taskCursor); err != nil {
		m.errMsg = "error: " + err.Error()
		return tea.Tick(3*time.Second, func(time.Time) tea.Msg { return errClearMsg{} })
	}
	m.refreshTasksViewport()
	return nil
}

// ── inline markdown ───────────────────────────────────────────────────────────

var (
	rxCode = regexp.MustCompile("`(.+?)`")
	rxBold = regexp.MustCompile(`\*\*(.+?)\*\*`)
)

func extractOpeningEscape(style lipgloss.Style) string {
	const marker = "\x00"
	rendered := style.Render(marker)
	if idx := strings.Index(rendered, marker); idx > 0 {
		return rendered[:idx]
	}
	return ""
}

func inlineMarkdown(s, restore string, done bool) string {
	if done {
		s = rxCode.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[4m" + rxCode.FindStringSubmatch(m)[1] + "\033[24m"
		})
		s = rxBold.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[1m" + rxBold.FindStringSubmatch(m)[1] + "\033[22m"
		})
	} else {
		s = rxCode.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[36m" + rxCode.FindStringSubmatch(m)[1] + "\033[0m" + restore
		})
		s = rxBold.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[1m" + rxBold.FindStringSubmatch(m)[1] + "\033[0m" + restore
		})
	}
	return s
}

// ── tasks rendering ───────────────────────────────────────────────────────────

func (m *Model) renderTasksContent() (string, int) {
	var sb strings.Builder
	line, cursorLine := 0, 0
	contentWidth := m.width - 2

	pendingRestore := extractOpeningEscape(taskPendingStyle)
	doneRestore := extractOpeningEscape(taskDoneStyle)

	for i, item := range m.taskItems {
		switch item.Kind {
		case openspec.KindSection:
			if i > 0 {
				sb.WriteString("\n")
				line++
			}
			done, total := sectionProgress(m.taskItems, i)
			sb.WriteString(sectionStyle.Render("  "+item.Text) + "  " + progressBar(done, total, 5) + "\n")
			line++
			sb.WriteString("\n")
			line++
		case openspec.KindTask:
			if i == m.taskCursor {
				cursorLine = line
			}
			checkbox := "[ ]"
			if item.Done {
				checkbox = "[x]"
			}
			restore := pendingRestore
			if item.Done {
				restore = doneRestore
			}
			var prefix string
			if i == m.taskCursor {
				prefix = taskCursorMarkStyle.Render("▶") + restore + " "
				checkbox = taskCursorMarkStyle.Render(checkbox) + restore
			} else {
				prefix = "  "
			}
			text := prefix + checkbox + " " + inlineMarkdown(item.Text, restore, item.Done)
			var rendered string
			switch {
			case item.Done:
				rendered = taskDoneStyle.Width(contentWidth).Render(text)
			default:
				rendered = taskPendingStyle.Width(contentWidth).Render(text)
			}
			sb.WriteString(rendered + "\n")
			line++
		}
	}
	return sb.String(), cursorLine
}

func sectionProgress(items []openspec.TaskItem, sectionIdx int) (done, total int) {
	for i := sectionIdx + 1; i < len(items); i++ {
		if items[i].Kind == openspec.KindSection {
			break
		}
		total++
		if items[i].Done {
			done++
		}
	}
	return
}

func progressBar(done, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := (done * width) / total
	if done == total {
		filled = width
	}
	bar := progressDoneStyle.Render(strings.Repeat("─", filled)) +
		progressEmptyStyle.Render(strings.Repeat("─", width-filled))
	return bar + helpStyle.Render(fmt.Sprintf(" %d/%d", done, total))
}

// ── view ──────────────────────────────────────────────────────────────────────

func (m *Model) renderHeader() string {
	if m.mode == ModeIndex {
		return headerStyle.Width(m.width - 2).Render(m.project.Name + "  ·  index")
	}
	if m.mode == ModeViewingSpec {
		specName := ""
		if m.specViewerCursor < len(m.projectSpecs) {
			specName = m.projectSpecs[m.specViewerCursor].Name
		}
		if m.specFocusMode && m.specViewerCursor < len(m.projectSpecs) {
			ps := m.projectSpecs[m.specViewerCursor]
			return headerStyle.Width(m.width - 2).Render(
				fmt.Sprintf("%s  ·  %s  ·  Req %d/%d", m.project.Name, specName, m.specReqCursor+1, len(ps.RequirementNames)),
			)
		}
		return headerStyle.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s  [spec]", m.project.Name, specName),
		)
	}
	ch := m.current()
	if ch == nil {
		return headerStyle.Render(m.project.Name)
	}
	if m.mode == ModeViewingArchive {
		return headerStyle.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s  [archive]", m.project.Name, ch.Name),
		)
	}
	nav := fmt.Sprintf("[%d/%d]", m.changeIdx+1, len(m.project.Changes))
	return headerStyle.Width(m.width - 2).Render(
		fmt.Sprintf("%s  ·  %s  %s", m.project.Name, ch.Name, nav),
	)
}

func (m *Model) renderTabBar() string {
	var parts []string
	for t := Tab(0); t < tabCount; t++ {
		label := tabLabels[t]
		switch {
		case t == m.tab:
			parts = append(parts, tabActiveStyle.Render(label))
		case !m.tabAvailable(t):
			parts = append(parts, tabDisabledStyle.Render(label))
		default:
			parts = append(parts, tabInactiveStyle.Render(label))
		}
	}
	tabs := strings.Join(parts, " ")

	// Progress bar right-aligned on the same line as the tabs
	taskItems := m.taskItems
	if m.mode == ModeViewingArchive {
		if ch := m.currentArchive(); ch != nil && ch.Tasks.Present {
			taskItems = openspec.ParseTasks(ch.Tasks.Content)
		} else {
			taskItems = nil
		}
	}
	total, done := 0, 0
	for _, item := range taskItems {
		if item.Kind == openspec.KindTask {
			total++
			if item.Done {
				done++
			}
		}
	}
	if total > 0 {
		label := fmt.Sprintf(" %d/%d", done, total)
		barSpace := (m.width-2) - lipgloss.Width(tabs) - 3 - len(label)
		if barSpace >= 3 {
			filled := (done * barSpace) / total
			if done == total {
				filled = barSpace
			}
			bar := "[" + progressDoneStyle.Render(strings.Repeat("█", filled)) +
				progressEmptyStyle.Render(strings.Repeat("░", barSpace-filled)) + "]"
			tabs = tabs + " " + bar + helpStyle.Render(label)
		}
	}
	return tabs
}

func (m *Model) renderSpecSubnav() string {
	ch := m.current()
	if ch == nil {
		return ""
	}
	var parts []string
	for i, s := range ch.SpecFiles {
		if i == m.specIdx {
			parts = append(parts, tabActiveStyle.Render(s.Name))
		} else {
			parts = append(parts, tabInactiveStyle.Render(s.Name))
		}
	}
	return strings.Join(parts, " ")
}

func (m *Model) hasSpecSubnav() bool {
	ch := m.current()
	return m.tab == TabSpecs && ch != nil && len(ch.SpecFiles) > 0
}

func (m *Model) boxTop() string {
	return separatorStyle.Render("┌" + strings.Repeat("─", m.width-2) + "┐")
}

func (m *Model) boxBottom() string {
	return separatorStyle.Render("└" + strings.Repeat("─", m.width-2) + "┘")
}

func (m *Model) boxInnerSep() string {
	return separatorStyle.Render("├" + strings.Repeat("─", m.width-2) + "┤")
}

func (m *Model) addBorderSides(content string) string {
	lines := strings.Split(content, "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	inner := m.width - 2
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		pad := inner - lipgloss.Width(line)
		if pad < 0 {
			pad = 0
		}
		result = append(result, separatorStyle.Render("│")+line+strings.Repeat(" ", pad)+separatorStyle.Render("│"))
	}
	return strings.Join(result, "\n")
}

func (m *Model) renderHelpBar() string {
	if m.errMsg != "" {
		return errStyle.Render(m.errMsg)
	}
	if m.mode == ModeIndex {
		return helpStyle.Render("j/k: navigate  Enter: open  Space: expand  Esc: quit")
	}
	if m.mode == ModeViewingSpec {
		if m.specFocusMode {
			return helpStyle.Render("h/l: req anterior/siguiente  j/k: scroll  Esc: índice  q: quit")
		}
		return helpStyle.Render("j/k: scroll  Esc: index  q: quit")
	}
	if m.mode == ModeViewingArchive {
		return helpStyle.Render("1-4: artifact  j/k: scroll  a/Esc: index  q: quit")
	}
	if m.tab == TabTasks {
		return helpStyle.Render("h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit")
	}
	return helpStyle.Render("h/l: change  1-4: artifact  j/k: scroll  e: edit  Esc: index  q: quit")
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

func (m *Model) emptyView() string {
	return headerStyle.Render(m.project.Name) +
		"\n\n\n  No active changes. Create one with /opsx:propose\n" +
		helpStyle.Render("\n  a/Esc: index  q: quit")
}
