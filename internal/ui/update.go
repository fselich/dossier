package ui

import (
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fselich/dossier/internal/openspec"
)

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

		case "s":
			if m.mode == ModeIndex {
				savedKind := indexKindActive
				savedIdx := -1
				savedReqIdx := 0
				if m.indexCursor < len(m.indexItems) {
					item := m.indexItems[m.indexCursor]
					savedKind = item.kind
					savedIdx = item.idx
					savedReqIdx = item.reqIdx
				}
				m.specSortBySuffix = !m.specSortBySuffix
				m.buildIndexItems()
				if savedIdx >= 0 {
					for i, it := range m.indexItems {
						if it.kind == savedKind && it.idx == savedIdx && it.reqIdx == savedReqIdx {
							m.indexCursor = i
							break
						}
					}
				}
				m.refreshIndexViewport()
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
