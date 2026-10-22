package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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
		if m.isTasksView() && m.mode == ModeNormal {
			m.moveCursorUp()
			m.refreshTasksViewport()
			return m, nil
		}
		if m.viewingCode && m.mode == ModeNormal {
			m.moveGitCursorUp()
			m.refreshGitViewport()
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
		if m.isTasksView() && m.mode == ModeNormal {
			m.moveCursorDown()
			m.refreshTasksViewport()
			return m, nil
		}
		if m.viewingCode && m.mode == ModeNormal {
			m.moveGitCursorDown()
			m.refreshGitViewport()
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
		if m.index.FilterIndices != nil {
			cursorItem := m.index.FilterIndices[m.index.Cursor]
			if cursorItem != idx {
				for ci, ri := range m.index.FilterIndices {
					if ri == idx {
						m.index.Cursor = ci
						break
					}
				}
				m.refreshIndexViewport()
				return m, nil
			}
		} else {
			if m.index.Cursor != idx {
				m.index.Cursor = idx
				m.refreshIndexViewport()
				return m, nil
			}
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
	if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
		for index, artifact := range ch.Artifacts {
			label := artifact.ID + artifactTabStatusSuffix(artifact)
			w := lipgloss.Width(label) + 2
			if msg.X >= x && msg.X <= x+w-1 {
				if m.selectArtifactPosition(index) {
					m.vp.SetHeight(m.contentHeight())
					return m, m.loadViewport()
				}
				return m, nil
			}
			x += w + 1
		}
		if m.isGitRepo && m.mode == ModeNormal {
			label := "code"
			if len(m.gitState.Files) > 0 {
				label += " (" + fmt.Sprintf("%d", len(m.gitState.Files)) + ")"
			}
			w := lipgloss.Width(label) + 2
			if msg.X >= x && msg.X <= x+w-1 && len(m.gitState.Files) > 0 {
				m.viewingCode = true
				m.vp.SetHeight(m.contentHeight())
				return m, m.loadViewport()
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) clickIndexItem(idx int) (tea.Model, tea.Cmd) {
	if idx < 0 || idx >= len(m.index.Items) {
		return m, nil
	}
	return m.primaryIndexItem(m.index.Items[idx])
}
