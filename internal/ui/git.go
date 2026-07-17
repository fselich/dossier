package ui

import (
	"strings"

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

	m.gitState.ErrMsg = ""

	preserveDiff := m.gitState.ShowingDiff && m.gitState.DiffFile != "" &&
		diffViewPreservable(m.gitState.DiffFile, files, m.gitState.Files)

	m.gitState.Files = files

	if preserveDiff {
		if m.gitState.Cursor >= len(files) {
			m.gitState.Cursor = 0
		}
	} else {
		m.gitState.ShowingDiff = false
		m.gitState.DiffLines = nil
		m.gitState.DiffFile = ""
		m.gitState.ScrollX = 0
		if m.gitState.Cursor >= len(files) {
			m.gitState.Cursor = 0
		} else if m.gitState.Cursor > 0 {
			m.gitState.Cursor = clampGitCursor(m.gitState.Cursor, files)
		}
	}
	if m.tab == TabGit {
		m.refreshGitViewport()
	}
}

func diffViewPreservable(diffFile string, newFiles, oldFiles []git.FileStatus) bool {
	var newStatus, oldStatus *git.FileStatus
	for i := range newFiles {
		if newFiles[i].Path == diffFile {
			newStatus = &newFiles[i]
			break
		}
	}
	for i := range oldFiles {
		if oldFiles[i].Path == diffFile {
			oldStatus = &oldFiles[i]
			break
		}
	}
	if newStatus == nil || oldStatus == nil {
		return false
	}
	return newStatus.X == oldStatus.X && newStatus.Y == oldStatus.Y
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
	return cursor
}

func (m *Model) moveGitCursorDown() {
	n := len(m.gitState.Files)
	if n == 0 {
		return
	}
	m.gitState.Cursor = (m.gitState.Cursor + 1) % n
}

func (m *Model) moveGitCursorUp() {
	n := len(m.gitState.Files)
	if n == 0 {
		return
	}
	m.gitState.Cursor = (m.gitState.Cursor - 1 + n) % n
}

func (m *Model) moveGitDiffCursorDown() {
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

func (m *Model) moveGitDiffCursorUp() {
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

func (m *Model) restoreGitCursor(path string) {
	for i, f := range m.gitState.Files {
		if f.Path == path {
			m.gitState.Cursor = i
			return
		}
	}
	if m.gitState.Cursor >= len(m.gitState.Files) {
		m.gitState.Cursor = 0
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
	if m.gitState.ShowingDiff {
		return m.renderDiffContent()
	}

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

func (m *Model) renderDiffContent() (string, int) {
	var sb strings.Builder
	header := m.gitState.DiffFile
	sb.WriteString("\n")
	sb.WriteString("  " + sectionStyle.Render("diff: "+header) + "\n")
	sb.WriteString("\n")
	if m.gitState.DiffLines == nil {
		sb.WriteString("  (no diff available)\n")
	} else {
		content := renderDiff(m.gitState.DiffLines, header, m.width-2, m.gitState.ScrollX)
		sb.WriteString(content)
	}
	return sb.String(), 0
}
