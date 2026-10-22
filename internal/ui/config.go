package ui

import tea "charm.land/bubbletea/v2"

func (m Model) updateConfig(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	keys := m.effectiveKeyMap().Config
	switch {

	case matchesKey(msg, keys.Back):
		m.mode = m.prevMode
		return m.commitStateChange()

	case matchesKey(msg, keys.Down):
		m.vp.ScrollDown(1)

	case matchesKey(msg, keys.Up):
		m.vp.ScrollUp(1)

	case matchesKey(msg, keys.PageDown):
		m.vp.PageDown()

	case matchesKey(msg, keys.PageUp):
		m.vp.PageUp()
	}
	return m, nil
}
