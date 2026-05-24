package ui

import (
	"image/color"
	"os"
	"path/filepath"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/openspec"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeIndex
	ModeViewingArchive
	ModeViewingSpec
	ModeViewingConfig
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

// renderedConfigMsg carries async glamour output for ModeViewingConfig.
type renderedConfigMsg struct {
	content string
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

type Theme struct {
	ViewBg color.Color
}

type Model struct {
	root string

	project   *openspec.Project
	changeIdx int
	tab       Tab

	vp      viewport.Model
	vpReady bool

	taskItems  []openspec.TaskItem
	taskCursor int

	specIdx int

	errMsg     string
	loading    bool
	singlePath bool

	width, height int

	renderCache map[Tab]string

	mode           Mode
	prevMode       Mode
	archiveChanges []openspec.Change
	archiveCursor  int
	indexItems     []indexItem
	indexCursor    int
	expandedSpecs  map[int]bool
	projectSpecs     []openspec.ProjectSpec
	specSortBySuffix bool
	specOrder        []int
	specViewerCursor int
	specJumpTarget   string
	specFocusMode    bool
	specReqCursor    int
	projectConfig    openspec.ProjectConfig
	theme            Theme
}

func New(project *openspec.Project, cfg openspec.ProjectConfig, root string) Model {
	m := Model{
		root:          root,
		project:       project,
		renderCache:   make(map[Tab]string),
		projectConfig: cfg,
		theme:         Theme{},
	}
	if len(project.Changes) > 0 {
		m.tab = m.defaultTab()
		m.loadTaskItems()
	}
	return m
}

func NewSinglePath(project *openspec.Project, cfg openspec.ProjectConfig, root string) Model {
	m := New(project, cfg, root)
	m.singlePath = true
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) View() tea.View {
	if !m.vpReady {
		return tea.NewView("")
	}

	var content string
	if m.mode == ModeViewingConfig {
		content = m.viewConfigContent()
	} else if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		content = m.viewIndexContent()
	} else if len(m.project.Changes) == 0 && m.mode == ModeNormal {
		content = m.emptyViewContent()
	} else {
		content = m.mainViewContent()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.BackgroundColor = m.theme.ViewBg
	return v
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

func (m *Model) nextAvailableTab(current Tab, delta int) Tab {
	next := current
	for i := 0; i < int(tabCount); i++ {
		next = Tab((int(next) + delta + int(tabCount)) % int(tabCount))
		if m.tabAvailable(next) {
			return next
		}
	}
	return current
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
	if m.mode == ModeIndex || m.mode == ModeViewingSpec || m.mode == ModeViewingConfig {
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

