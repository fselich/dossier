package ui

import (
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/git"
)

func (m Model) updateViewer(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.tab == TabGit && m.gitState.ErrMsg != "" {
		m.gitState.ErrMsg = ""
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "s":
		if m.tab != TabGit || m.gitState.ShowingDiff || len(m.gitState.Files) == 0 {
			return m, nil
		}
		f := m.gitState.Files[m.gitState.Cursor]

		var err error
		if f.Y != ' ' || f.X == '?' && f.Y == '?' {
			paths := []string{f.Path}
			if f.OldPath != "" {
				paths = []string{f.OldPath, f.Path}
			}
			err = git.Stage(m.gitRoot, paths...)
		} else {
			paths := []string{f.Path}
			if f.OldPath != "" {
				paths = []string{f.OldPath, f.Path}
			}
			err = git.Unstage(m.gitRoot, paths...)
		}
		if err != nil {
			m.gitState.ErrMsg = strings.TrimSpace(err.Error())
		}
		files, statusErr := git.Status(m.gitRoot)
		if statusErr == nil {
			m.gitState.Files = files
		}
		m.restoreGitCursor(f.Path)
		m.refreshGitViewport()
		return m, nil

	case "i":
		m.prevMode = m.mode
		m.mode = ModeViewingConfig
		return m.commitStateChange()

	case "a", "esc":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.toggleGitDiff()
			return m, nil
		}
		m.enterIndex()
		return m, nil

	case "h":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.gitState.ScrollX -= 10
			if m.gitState.ScrollX < 0 {
				m.gitState.ScrollX = 0
			}
			m.refreshGitViewport()
		} else if len(m.project.Changes) > 0 {
			m.changeIdx = (m.changeIdx - 1 + len(m.project.Changes)) % len(m.project.Changes)
			m.renderCache = make(map[Tab]string)
			m.loadTaskItems()
			m.tab = m.defaultTab()
			m.specIdx = 0
			return m.commitStateChange()
		}

	case "l":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.gitState.ScrollX += 10
			m.refreshGitViewport()
		} else if len(m.project.Changes) > 0 {
			m.changeIdx = (m.changeIdx + 1) % len(m.project.Changes)
			m.renderCache = make(map[Tab]string)
			m.loadTaskItems()
			m.tab = m.defaultTab()
			m.specIdx = 0
			return m.commitStateChange()
		}

	case "right":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.gitState.ScrollX += 10
			m.refreshGitViewport()
		}

	case "left":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.gitState.ScrollX -= 10
			if m.gitState.ScrollX < 0 {
				m.gitState.ScrollX = 0
			}
			m.refreshGitViewport()
		}

	case "]":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.moveGitDiffCursorDown()
			m.loadDiffForFile(m.gitState.Cursor)
			return m, nil
		}

	case "[":
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.moveGitDiffCursorUp()
			m.loadDiffForFile(m.gitState.Cursor)
			return m, nil
		}

	case "1":
		if m.tabAvailable(TabProposal) {
			m.tab = TabProposal
			return m.commitStateChange()
		}
	case "2":
		if m.tabAvailable(TabDesign) {
			m.tab = TabDesign
			return m.commitStateChange()
		}
	case "3":
		if m.tabAvailable(TabSpecs) {
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
			return m.commitStateChange()
		}
	case "4":
		if m.tabAvailable(TabTasks) {
			m.tab = TabTasks
			return m.commitStateChange()
		}
	case "5":
		if m.tabAvailable(TabGit) {
			m.tab = TabGit
			return m.commitStateChange()
		}

	case "tab":
		nxt := m.nextAvailableTab(m.tab, 1)
		if nxt != m.tab {
			m.tab = nxt
			return m.commitStateChange()
		}
	case "shift+tab":
		prv := m.nextAvailableTab(m.tab, -1)
		if prv != m.tab {
			m.tab = prv
			return m.commitStateChange()
		}

	case "d":
		if m.tab == TabGit {
			m.toggleGitDiff()
		}

	case "j", "down":
		switch m.tab {
		case TabTasks:
			if m.mode == ModeViewingArchive {
				m.vp.ScrollDown(1)
			} else {
				m.moveCursorDown()
				m.refreshTasksViewport()
			}
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.ScrollDown(1)
			} else {
				m.moveGitCursorDown()
				m.refreshGitViewport()
			}
		default:
			m.vp.ScrollDown(1)
		}

	case "pgdown":
		switch m.tab {
		case TabTasks:
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.PageDown()
			}
		default:
			m.vp.PageDown()
		}

	case "k", "up":
		switch m.tab {
		case TabTasks:
			if m.mode == ModeViewingArchive {
				m.vp.ScrollUp(1)
			} else {
				m.moveCursorUp()
				m.refreshTasksViewport()
			}
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.ScrollUp(1)
			} else {
				m.moveGitCursorUp()
				m.refreshGitViewport()
			}
		default:
			m.vp.ScrollUp(1)
		}

	case "pgup":
		switch m.tab {
		case TabTasks:
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.PageUp()
			}
		default:
			m.vp.PageUp()
		}

	case "space":
		if m.mode == ModeViewingArchive {
			return m, nil
		}
		if m.tab == TabTasks {
			return m, m.doToggle()
		}

	case "enter":
		if m.tab == TabGit {
			m.toggleGitDiff()
			return m, nil
		}

	case "e":
		if m.tab == TabGit {
			m.toggleGitDiff()
			return m, nil
		}
		if m.mode == ModeViewingArchive {
			return m, nil
		}
		if m.tabAvailable(m.tab) {
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
	return m, nil
}
