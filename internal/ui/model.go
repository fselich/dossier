package ui

import (
	"image/color"
	"os"
	"path/filepath"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
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

type indexState struct {
	Items          []indexItem
	Cursor         int
	ExpandedSpecs  map[int]bool
	SortBySuffix   bool
	Order          []int
	ArchiveChanges []openspec.Change
	ArchiveCursor  int
}

type specViewerState struct {
	Cursor     int
	JumpTarget string
	FocusMode  bool
	ReqCursor  int
}

type taskState struct {
	Items  []openspec.TaskItem
	Cursor int
}

type indexItemKind int

const (
	indexKindActive indexItemKind = iota
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
	root   string
	loader *openspec.Loader

	project   *openspec.Project
	changeIdx int
	tab       Tab

	vp      viewport.Model
	vpReady bool

	tasks taskState

	specIdx int

	errMsg     string
	loading    bool
	singlePath bool

	width, height int

	renderCache     map[Tab]string
	glamourRenderer *glamour.TermRenderer
	lastRenderWidth int

	mode          Mode
	prevMode      Mode
	index         indexState
	projectSpecs  []openspec.ProjectSpec
	specViewer    specViewerState
	projectConfig openspec.ProjectConfig
	theme         Theme
}

func New(project *openspec.Project, cfg openspec.ProjectConfig, root string, loader *openspec.Loader) Model {
	m := Model{
		root:          root,
		loader:        loader,
		project:       project,
		renderCache:   make(map[Tab]string),
		projectConfig: cfg,
		theme:         Theme{},
	}
	if len(project.Changes) > 0 {
		m.tab = m.defaultTab()
		m.loadTaskItems()
	} else {
		var archiveErr error
		m.index.ArchiveChanges, archiveErr = loader.ListArchiveChangesFrom(root)
		if archiveErr != nil {
			m.errMsg = "error loading archive changes: " + archiveErr.Error()
		}
		var specErr error
		m.projectSpecs, specErr = loader.LoadProjectSpecsFrom(root)
		if specErr != nil {
			m.errMsg = "error loading project specs: " + specErr.Error()
		}
		m.index.ExpandedSpecs = make(map[int]bool)
		m.buildIndexItems()
		m.mode = ModeIndex
	}
	return m
}

func NewSinglePath(project *openspec.Project, cfg openspec.ProjectConfig, root string, loader *openspec.Loader) Model {
	m := New(project, cfg, root, loader)
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
		content = m.viewContentWithChrome()
	} else if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		content = m.viewContentWithChrome()
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
	for range int(tabCount) {
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
	if m.index.ArchiveCursor < len(m.index.ArchiveChanges) {
		return &m.index.ArchiveChanges[m.index.ArchiveCursor]
	}
	return nil
}

const (
	chromeTop        = 1
	chromeHeader     = 1
	chromeTabBar     = 1
	chromeInnerSep   = 1
	chromeSpecSubnav = 1
	chromeHelpBar    = 1
	chromeBottom     = 1
)

func (m *Model) contentHeight() int {
	if m.mode == ModeIndex || m.mode == ModeViewingSpec || m.mode == ModeViewingConfig {
		h := m.height - (chromeTop + chromeHeader + chromeInnerSep + chromeInnerSep + chromeHelpBar + chromeBottom)
		if h < 1 {
			h = 1
		}
		return h
	}
	h := m.height - (chromeTop + chromeHeader + chromeTabBar + chromeInnerSep + chromeInnerSep + chromeHelpBar + chromeBottom)
	if m.hasSpecSubnav() {
		h -= chromeSpecSubnav
	}
	if h < 1 {
		h = 1
	}
	return h
}

// mergeReloadedChange updates in-memory state from a freshly reloaded Change
// and returns which artifacts changed. It does not handle cursor preservation
// or viewport refresh — the caller handles those.
func (m *Model) mergeReloadedChange(fresh openspec.Change) (tasksChanged bool, viewportDirty bool) {
	ch := m.current()
	if ch == nil {
		return false, false
	}

	if fresh.Tasks.Present != ch.Tasks.Present || fresh.Tasks.Content != ch.Tasks.Content {
		m.project.Changes[m.changeIdx].Tasks = fresh.Tasks
		m.tasks.Items = openspec.ParseTasks(fresh.Tasks.Content)
		tasksChanged = true
	}
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
	return
}
