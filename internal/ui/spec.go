package ui

import tea "charm.land/bubbletea/v2"

func (m Model) updateSpec(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	keys := m.effectiveKeyMap().Spec
	switch {

	case matchesKey(msg, keys.Quit):
		if m.lifecycle != nil {
			m.lifecycle.Cleanup(m.lifecycleUndo)
		}
		return m, tea.Quit

	case matchesKey(msg, keys.Back):
		specIdx := m.specViewer.Cursor
		jumpTarget := m.specViewer.JumpTarget
		wasFocusMode := m.specViewer.FocusMode
		m.enterIndex()
		if wasFocusMode && jumpTarget != "" {
			m.index.ExpandedSpecs[specIdx] = true
			m.buildIndexItems()
			for j, it := range m.index.Items {
				if it.kind == indexKindRequirement && it.idx == specIdx &&
					it.reqIdx < len(m.projectSpecs[specIdx].RequirementNames) &&
					m.projectSpecs[specIdx].RequirementNames[it.reqIdx] == jumpTarget {
					m.index.Cursor = j
					break
				}
			}
		} else {
			for i, item := range m.index.Items {
				if item.kind == indexKindSpec && item.idx == specIdx {
					m.index.Cursor = i
					break
				}
			}
		}
		m.refreshIndexViewport()

	case matchesKey(msg, keys.Down):
		m.vp.ScrollDown(1)

	case matchesKey(msg, keys.Up):
		m.vp.ScrollUp(1)

	case matchesKey(msg, keys.PageDown):
		m.vp.PageDown()

	case matchesKey(msg, keys.PageUp):
		m.vp.PageUp()

	case matchesKey(msg, keys.PreviousRequirement):
		if m.specViewer.FocusMode {
			ps := m.projectSpecs[m.specViewer.Cursor]
			if len(ps.RequirementNames) > 0 {
				m.specViewer.ReqCursor = (m.specViewer.ReqCursor - 1 + len(ps.RequirementNames)) % len(ps.RequirementNames)
				m.specViewer.JumpTarget = ps.RequirementNames[m.specViewer.ReqCursor]
				return m, m.loadViewport()
			}
		}

	case matchesKey(msg, keys.NextRequirement):
		if m.specViewer.FocusMode {
			ps := m.projectSpecs[m.specViewer.Cursor]
			if len(ps.RequirementNames) > 0 {
				m.specViewer.ReqCursor = (m.specViewer.ReqCursor + 1) % len(ps.RequirementNames)
				m.specViewer.JumpTarget = ps.RequirementNames[m.specViewer.ReqCursor]
				return m, m.loadViewport()
			}
		}
	}
	return m, nil
}
