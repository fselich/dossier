package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fselich/dossier/internal/openspec"
)

func (m *Model) refreshTasksViewport() {
	content, cursorLine := m.renderTasksContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset+m.vp.Height {
		m.vp.SetYOffset(cursorLine - m.vp.Height + 1)
	}
}

func (m *Model) loadTaskItems() {
	ch := m.current()
	if ch == nil || !ch.Tasks.Present {
		m.taskItems = nil
		m.taskCursor = 0
		return
	}
	m.taskItems = openspec.ParseTasks(ch.Tasks.Content)
	m.taskCursor = m.firstTaskIdx()
}

func (m *Model) firstTaskIdx() int {
	for i, item := range m.taskItems {
		if item.Kind == openspec.KindTask {
			return i
		}
	}
	return 0
}

func (m *Model) moveCursorDown() {
	for i := m.taskCursor + 1; i < len(m.taskItems); i++ {
		if m.taskItems[i].Kind == openspec.KindTask {
			m.taskCursor = i
			return
		}
	}
}

func (m *Model) moveCursorUp() {
	for i := m.taskCursor - 1; i >= 0; i-- {
		if m.taskItems[i].Kind == openspec.KindTask {
			m.taskCursor = i
			return
		}
	}
}

func (m *Model) doToggle() tea.Cmd {
	if len(m.taskItems) == 0 || m.taskCursor >= len(m.taskItems) {
		return nil
	}
	if m.taskItems[m.taskCursor].Kind != openspec.KindTask {
		return nil
	}
	ch := m.current()
	if ch == nil {
		return nil
	}
	if err := openspec.ToggleTask(ch.Path+"/tasks.md", m.taskItems, m.taskCursor); err != nil {
		m.errMsg = "error: " + err.Error()
		return tea.Tick(3*time.Second, func(time.Time) tea.Msg { return errClearMsg{} })
	}
	m.refreshTasksViewport()
	return nil
}

var (
	rxCode = regexp.MustCompile("`(.+?)`")
	rxBold = regexp.MustCompile(`\*\*(.+?)\*\*`)
)

func extractOpeningEscape(style lipgloss.Style) string {
	const marker = "\x00"
	rendered := style.Render(marker)
	if idx := strings.Index(rendered, marker); idx > 0 {
		return rendered[:idx]
	}
	return ""
}

func inlineMarkdown(s, restore string, done bool) string {
	if done {
		s = rxCode.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[4m" + rxCode.FindStringSubmatch(m)[1] + "\033[24m"
		})
		s = rxBold.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[1m" + rxBold.FindStringSubmatch(m)[1] + "\033[22m"
		})
	} else {
		s = rxCode.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[36m" + rxCode.FindStringSubmatch(m)[1] + "\033[0m" + restore
		})
		s = rxBold.ReplaceAllStringFunc(s, func(m string) string {
			return "\033[1m" + rxBold.FindStringSubmatch(m)[1] + "\033[0m" + restore
		})
	}
	return s
}

func (m *Model) renderTasksContent() (string, int) {
	var sb strings.Builder
	line, cursorLine := 0, 0
	contentWidth := m.width - 2

	pendingRestore := extractOpeningEscape(taskPendingStyle)
	doneRestore := extractOpeningEscape(taskDoneStyle)

	for i, item := range m.taskItems {
		switch item.Kind {
		case openspec.KindSection:
			if i > 0 {
				sb.WriteString("\n")
				line++
			}
			done, total := sectionProgress(m.taskItems, i)
			sectionLine := sectionStyle.Render("  "+item.Text) + "  " + progressBar(done, total, 5)
			sb.WriteString(sectionLine + "\n")
			line += lipgloss.Height(sectionLine)
			sb.WriteString("\n")
			line++
		case openspec.KindTask:
			if i == m.taskCursor {
				cursorLine = line
			}
			checkbox := "[ ]"
			if item.Done {
				checkbox = "[x]"
			}
			restore := pendingRestore
			if item.Done {
				restore = doneRestore
			}
			var prefix string
			if i == m.taskCursor {
				prefix = taskCursorMarkStyle.Render("▶") + restore + " "
				checkbox = taskCursorMarkStyle.Render(checkbox) + restore
			} else {
				prefix = "  "
			}
			text := prefix + checkbox + " " + inlineMarkdown(item.Text, restore, item.Done)
			var rendered string
			switch {
			case item.Done:
				rendered = taskDoneStyle.Width(contentWidth).Render(text)
			default:
				rendered = taskPendingStyle.Width(contentWidth).Render(text)
			}
			sb.WriteString(rendered + "\n")
			line += lipgloss.Height(rendered)
		}
	}
	return sb.String(), cursorLine
}

func sectionProgress(items []openspec.TaskItem, sectionIdx int) (done, total int) {
	for i := sectionIdx + 1; i < len(items); i++ {
		if items[i].Kind == openspec.KindSection {
			break
		}
		total++
		if items[i].Done {
			done++
		}
	}
	return
}

func progressBar(done, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := (done * width) / total
	filledStyle := progressDoneStyle
	if done == total {
		filled = width
		filledStyle = progressCompleteStyle
	}
	bar := filledStyle.Render(strings.Repeat("─", filled)) +
		progressEmptyStyle.Render(strings.Repeat("─", width-filled))
	return bar + helpStyle.Render(fmt.Sprintf(" %d/%d", done, total))
}
