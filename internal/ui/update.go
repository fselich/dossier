package ui

import (
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentH := m.contentHeight()
		if !m.vpReady {
			m.vp = viewport.New(viewport.WithWidth(m.width-2), viewport.WithHeight(contentH))
			m.vpReady = true
		} else {
			m.vp.SetWidth(m.width - 2)
			m.vp.SetHeight(contentH)
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

	case renderedConfigMsg:
		m.loading = false
		if m.mode == ModeViewingConfig {
			m.vp.SetContent(msg.content)
			m.vp.GotoTop()
		}
		return m, nil

	case tickMsg:
		cmd := m.handleTick()
		nextTick := tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
		return m, tea.Batch(nextTick, cmd)

	case editorReturnMsg:
		ch := m.current()
		if ch != nil {
			var cursorText string
			if m.tasks.Cursor < len(m.tasks.Items) && m.tasks.Items[m.tasks.Cursor].Kind == openspec.KindTask {
				cursorText = m.tasks.Items[m.tasks.Cursor].Text
			}
			fresh := m.loader.ReloadChange(*ch)
			tasksChanged, _ := m.mergeReloadedChange(fresh)
			if tasksChanged {
				m.tasks.Cursor = openspec.FindCursorByText(m.tasks.Items, cursorText)
			}
		}
		return m, m.loadViewport()

	case errClearMsg:
		m.errMsg = ""
		return m, nil

	case tea.MouseWheelMsg:
		return m.handleMouseWheel(msg)

	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)

	case tea.KeyPressMsg:
		return m.dispatchKey(msg)
	}
	return m, nil
}

func (m Model) dispatchKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeNormal, ModeViewingArchive:
		return m.updateViewer(msg)
	case ModeIndex:
		return m.updateIndex(msg)
	case ModeViewingSpec:
		return m.updateSpec(msg)
	case ModeViewingConfig:
		return m.updateConfig(msg)
	}
	return m, nil
}

func (m Model) commitStateChange() (tea.Model, tea.Cmd) {
	m.vp.SetHeight(m.contentHeight())
	return m, m.loadViewport()
}
