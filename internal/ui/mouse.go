package ui

import tea "charm.land/bubbletea/v2"

func (m Model) handleMouseWheel(msg tea.MouseWheelMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseWheelUp:
		if m.mode == ModeIndex {
			if m.index.Cursor > 0 {
				m.index.Cursor--
			}
			m.refreshIndexViewport()
			return m, nil
		}
		if m.tab == TabTasks && m.mode == ModeNormal {
			m.moveCursorUp()
			m.refreshTasksViewport()
			return m, nil
		}
		m.vp.ScrollUp(3)
		return m, nil
	case tea.MouseWheelDown:
		if m.mode == ModeIndex {
			if m.index.Cursor < len(m.index.Items)-1 {
				m.index.Cursor++
			}
			m.refreshIndexViewport()
			return m, nil
		}
		if m.tab == TabTasks && m.mode == ModeNormal {
			m.moveCursorDown()
			m.refreshTasksViewport()
			return m, nil
		}
		m.vp.ScrollDown(3)
		return m, nil
	}
	return m, nil
}

func (m Model) handleMouseClick(msg tea.MouseClickMsg) (tea.Model, tea.Cmd) {
	if msg.Button != tea.MouseLeft {
		return m, nil
	}

	if m.mode == ModeIndex {
		if msg.Y < indexViewportContentStart || msg.Y >= indexViewportContentStart+m.vp.Height() {
			return m, nil
		}
		contentLine := msg.Y - indexViewportContentStart + m.vp.YOffset()
		idx, found := m.indexItemAtContentLine(contentLine)
		if !found {
			return m, nil
		}
		if m.index.Cursor != idx {
			m.index.Cursor = idx
			m.refreshIndexViewport()
			return m, nil
		}
		return m.clickIndexItem(idx)
	}

	if m.mode != ModeNormal && m.mode != ModeViewingArchive {
		return m, nil
	}

	if msg.Y == 1 {
		m.enterIndex()
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
			m.vp.SetHeight(m.contentHeight())
			return m, m.loadViewport()
		}
		x += w + 1
	}

	return m, nil
}

func (m Model) clickIndexItem(idx int) (tea.Model, tea.Cmd) {
	item := m.index.Items[idx]
	m.renderCache = make(map[Tab]string)
	switch item.kind {
	case indexKindActive:
		m.changeIdx = item.idx
		m.mode = ModeNormal
		m.tab = m.defaultTab()
		m.loadTaskItems()
		m.vp.SetHeight(m.contentHeight())
		return m, m.loadViewport()

	case indexKindArchived:
		m.index.ArchiveCursor = item.idx
		m.tab = firstAvailableTab(m.index.ArchiveChanges[item.idx])
		m.mode = ModeViewingArchive
		m.vp.SetHeight(m.contentHeight())
		return m, m.loadViewport()

	case indexKindSpec:
		m.index.ExpandedSpecs[item.idx] = !m.index.ExpandedSpecs[item.idx]
		m.buildIndexItems()
		m.index.Cursor = 0
		for i, it := range m.index.Items {
			if it.kind == indexKindSpec && it.idx == item.idx {
				m.index.Cursor = i
				break
			}
		}
		if m.index.Cursor >= len(m.index.Items) {
			m.index.Cursor = max(0, len(m.index.Items)-1)
		}
		m.refreshIndexViewport()
		return m, nil

	case indexKindRequirement:
		m.specViewer.Cursor = item.idx
		m.specViewer.JumpTarget = m.projectSpecs[item.idx].RequirementNames[item.reqIdx]
		m.specViewer.FocusMode = true
		m.specViewer.ReqCursor = item.reqIdx
		m.mode = ModeViewingSpec
		m.vp.SetHeight(m.contentHeight())
		return m, m.loadViewport()
	}
	return m, nil
}
