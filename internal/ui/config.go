package ui

import tea "charm.land/bubbletea/v2"

func (m Model) updateConfig(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "q", "ctrl+c", "i", "esc":
		m.mode = m.prevMode
		return m.commitStateChange()

	case "j", "down":
		m.vp.ScrollDown(1)

	case "k", "up":
		m.vp.ScrollUp(1)
	}
	return m, nil
}
