package ui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.mode == ModeIndex {
			if m.indexCursor > 0 {
				m.indexCursor--
			}
			m.refreshIndexViewport()
			return m, nil
		}
		if m.tab == TabTasks && m.mode == ModeNormal {
			m.moveCursorUp()
			m.refreshTasksViewport()
			return m, nil
		}
		m.vp.LineUp(3)
		return m, nil
	case tea.MouseButtonWheelDown:
		if m.mode == ModeIndex {
			if m.indexCursor < len(m.indexItems)-1 {
				m.indexCursor++
			}
			m.refreshIndexViewport()
			return m, nil
		}
		if m.tab == TabTasks && m.mode == ModeNormal {
			m.moveCursorDown()
			m.refreshTasksViewport()
			return m, nil
		}
		m.vp.LineDown(3)
		return m, nil
	}

	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return m, nil
	}

	if m.mode != ModeNormal && m.mode != ModeViewingArchive {
		return m, nil
	}

	if msg.Y != 2 {
		return m, nil
	}

	x := 1
	for t := Tab(0); t < tabCount; t++ {
		w := len(tabLabels[t]) + 2
		if msg.X >= x && msg.X <= x+w-1 {
			if !m.tabAvailable(t) {
				return m, nil
			}
			if t == TabSpecs && m.tab == TabSpecs {
				ch := m.current()
				if ch != nil && len(ch.SpecFiles) > 1 {
					m.specIdx = (m.specIdx + 1) % len(ch.SpecFiles)
					delete(m.renderCache, TabSpecs)
				}
			} else {
				m.tab = t
				if t == TabSpecs {
					m.specIdx = 0
				}
			}
			m.vp.Height = m.contentHeight()
			return m, m.loadViewport()
		}
		x += w + 1
	}

	return m, nil
}
