package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/fselich/dossier/internal/git"
)

func (m *Model) pollGitStatus() {
	if !m.isGitRepo {
		return
	}
	files, err := git.Status(m.root)
	if err != nil {
		return
	}
	if gitStatusesEqual(m.gitState.Files, files) {
		return
	}
	m.gitState.Files = files
	if m.gitState.Cursor >= len(files) {
		m.gitState.Cursor = 0
	} else if m.gitState.Cursor > 0 {
		m.gitState.Cursor = clampGitCursor(m.gitState.Cursor, files)
	}
	if m.tab == TabGit {
		m.refreshGitViewport()
	}
}

func gitStatusesEqual(a, b []git.FileStatus) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Path != b[i].Path || a[i].X != b[i].X || a[i].Y != b[i].Y {
			return false
		}
	}
	return true
}

func clampGitCursor(cursor int, files []git.FileStatus) int {
	if len(files) == 0 {
		return 0
	}
	if cursor >= len(files) {
		return len(files) - 1
	}
	if files[cursor].IsDeleted {
		for i := cursor; i < len(files); i++ {
			if !files[i].IsDeleted {
				return i
			}
		}
		for i := cursor; i >= 0; i-- {
			if !files[i].IsDeleted {
				return i
			}
		}
	}
	return cursor
}

func (m *Model) moveGitCursorDown() {
	n := len(m.gitState.Files)
	if n == 0 {
		return
	}
	start := m.gitState.Cursor
	for {
		m.gitState.Cursor = (m.gitState.Cursor + 1) % n
		if !m.gitState.Files[m.gitState.Cursor].IsDeleted {
			return
		}
		if m.gitState.Cursor == start {
			return
		}
	}
}

func (m *Model) moveGitCursorUp() {
	n := len(m.gitState.Files)
	if n == 0 {
		return
	}
	start := m.gitState.Cursor
	for {
		m.gitState.Cursor = (m.gitState.Cursor - 1 + n) % n
		if !m.gitState.Files[m.gitState.Cursor].IsDeleted {
			return
		}
		if m.gitState.Cursor == start {
			return
		}
	}
}

func (m *Model) refreshGitViewport() {
	content, cursorLine := m.renderGitContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset() {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset()+m.vp.Height() {
		m.vp.SetYOffset(cursorLine - m.vp.Height() + 1)
	}
}

func gitFileColorStyle(x, y byte, isDeleted bool) lipgloss.Style {
	if isDeleted {
		return gitStatusDeleted
	}
	switch {
	case x == '?' && y == '?':
		return gitStatusUntracked
	case x == 'A' || y == 'A':
		return gitStatusAdded
	case x == 'R' || y == 'R', x == 'C' || y == 'C':
		return gitStatusRenamed
	case x == 'D' || y == 'D':
		return gitStatusDeleted
	default:
		return gitStatusModified
	}
}

func (m *Model) renderGitContent() (string, int) {
	var sb strings.Builder
	line := 0
	cursorLine := 0

	sb.WriteString("\n")
	line++

	if len(m.gitState.Files) == 0 {
		sb.WriteString("  " + helpStyle.Render("(working tree clean)") + "\n")
		return sb.String(), cursorLine
	}

	for i, f := range m.gitState.Files {
		cursor := m.gitState.Cursor == i
		if cursor {
			cursorLine = line
		}

		cursorMark := "  "
		if cursor {
			cursorMark = taskCursorMarkStyle.Render("▶") + " "
		}

		statusStyle := gitFileColorStyle(f.X, f.Y, f.IsDeleted)

		xPart := gitStatusDot.Render("·")
		if f.X > ' ' {
			xPart = statusStyle.Render(string(f.X))
		}
		yPart := gitStatusDot.Render("·")
		if f.Y > ' ' {
			yPart = statusStyle.Render(string(f.Y))
		}
		styledStatus := xPart + yPart

		var pathStr string
		if f.OldPath != "" {
			pathStr = helpStyle.Render(f.OldPath+" → ") + statusStyle.Render(f.Path)
		} else if f.IsDeleted {
			pathStr = gitStatusDeleted.Render(f.Path)
		} else {
			pathStr = statusStyle.Render(f.Path)
		}

		lineText := cursorMark + styledStatus + " " + pathStr
		sb.WriteString(lineText + "\n")
		line++
	}

	return sb.String(), cursorLine
}

func (m *Model) openGitFile() tea.Cmd {
	if m.gitState.Cursor >= len(m.gitState.Files) {
		return nil
	}
	f := m.gitState.Files[m.gitState.Cursor]
	if f.IsDeleted {
		return nil
	}
	path := f.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(m.gitRoot, path)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, path)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorReturnMsg{}
	})
}
