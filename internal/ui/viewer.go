package ui

import (
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/aemonge/dossier/internal/git"
)

func (m Model) updateViewer(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.viewingCode && m.gitState.ErrMsg != "" {
		m.gitState.ErrMsg = ""
	}

	keys := m.effectiveKeyMap().Viewer
	switch {

	case matchesKey(msg, keys.Quit):
		if m.lifecycle != nil {
			m.lifecycle.Cleanup(m.lifecycleUndo)
		}
		return m, tea.Quit

	case matchesKey(msg, keys.Stage):
		if m.readOnly || !m.viewingCode || m.gitState.ShowingDiff || len(m.gitState.Files) == 0 {
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

	case matchesKey(msg, keys.Info):
		m.prevMode = m.mode
		m.mode = ModeViewingConfig
		return m.commitStateChange()

	case matchesKey(msg, keys.Back):
		if m.viewingCode && m.gitState.ShowingDiff {
			m.toggleGitDiff()
			return m, nil
		}
		m.enterIndex()
		return m, nil

	case matchesKey(msg, keys.Previous):
		if len(m.project.Changes) > 0 {
			m.changeIdx = (m.changeIdx - 1 + len(m.project.Changes)) % len(m.project.Changes)
			m.artifactRenderCache = make(map[artifactOutputIdentity]string)
			m.selectDefaultArtifact()
			m.loadTaskItems()
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.Next):
		if len(m.project.Changes) > 0 {
			m.changeIdx = (m.changeIdx + 1) % len(m.project.Changes)
			m.artifactRenderCache = make(map[artifactOutputIdentity]string)
			m.selectDefaultArtifact()
			m.loadTaskItems()
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.Right):
		if m.viewingCode && m.gitState.ShowingDiff {
			m.gitState.ScrollX += 10
			m.refreshGitViewport()
		}

	case matchesKey(msg, keys.Left):
		if m.viewingCode && m.gitState.ShowingDiff {
			m.gitState.ScrollX -= 10
			if m.gitState.ScrollX < 0 {
				m.gitState.ScrollX = 0
			}
			m.refreshGitViewport()
		}

	case matchesKey(msg, keys.NextDiff):
		if m.viewingCode && m.gitState.ShowingDiff {
			m.moveGitDiffCursorDown()
			m.loadDiffForFile(m.gitState.Cursor)
			return m, nil
		}

	case matchesKey(msg, keys.PreviousDiff):
		if m.viewingCode && m.gitState.ShowingDiff {
			m.moveGitDiffCursorUp()
			m.loadDiffForFile(m.gitState.Cursor)
			return m, nil
		}

	case matchesKey(msg, keys.ProposalTab):
		if m.selectArtifactPosition(0) {
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.DesignTab):
		if m.selectArtifactPosition(1) {
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.SpecsTab):
		if m.selectArtifactPosition(2) {
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.TasksTab):
		if m.selectArtifactPosition(3) {
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.GitTab):
		if m.isGitRepo && m.mode == ModeNormal && len(m.gitState.Files) > 0 {
			m.viewingCode = true
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.NextTab):
		if m.moveViewerDestination(1) {
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.PreviousTab):
		if m.moveViewerDestination(-1) {
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.ViewDiff):
		if m.viewingCode {
			m.toggleGitDiff()
		}

	case matchesKey(msg, keys.Down):
		switch {
		case m.viewingCode:
			if m.gitState.ShowingDiff {
				m.vp.ScrollDown(1)
			} else {
				m.moveGitCursorDown()
				m.refreshGitViewport()
			}
		case m.isTasksView() && m.mode == ModeNormal:
			m.moveCursorDown()
			m.refreshTasksViewport()
		default:
			m.vp.ScrollDown(1)
		}

	case matchesKey(msg, keys.PageDown):
		if !m.isTasksView() || m.mode == ModeViewingArchive {
			m.vp.PageDown()
		}

	case matchesKey(msg, keys.Up):
		switch {
		case m.viewingCode:
			if m.gitState.ShowingDiff {
				m.vp.ScrollUp(1)
			} else {
				m.moveGitCursorUp()
				m.refreshGitViewport()
			}
		case m.isTasksView() && m.mode == ModeNormal:
			m.moveCursorUp()
			m.refreshTasksViewport()
		default:
			m.vp.ScrollUp(1)
		}

	case matchesKey(msg, keys.PageUp):
		if !m.isTasksView() || m.mode == ModeViewingArchive {
			m.vp.PageUp()
		}

	case matchesKey(msg, keys.ToggleTask):
		if m.mode == ModeNormal && m.isTasksView() {
			return m, m.doToggle()
		}

	case matchesKey(msg, keys.Open):
		if m.viewingCode {
			m.toggleGitDiff()
			return m, nil
		}
		if m.readOnly || m.mode == ModeViewingArchive {
			return m, nil
		}
		if path, ok := m.selectedArtifactFilePath(); ok {
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			cmd := exec.Command(editor, path)
			identity, _ := m.selectedArtifactKey()
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return editorReturnMsg{artifactIdentity: identity}
			})
		}
	}
	return m, nil
}
