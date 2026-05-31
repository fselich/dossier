package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func (m *Model) refreshTasksViewport() {
	content, cursorLine := m.renderTasksContent()
	m.vp.SetContent(content)
	if cursorLine < m.vp.YOffset() {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine >= m.vp.YOffset()+m.vp.Height() {
		m.vp.SetYOffset(cursorLine - m.vp.Height() + 1)
	}
}

func (m *Model) loadTaskItems() {
	ch := m.current()
	if ch == nil || !ch.Tasks.Present {
		m.tasks.Items = nil
		m.tasks.Cursor = 0
		return
	}
	m.tasks.Items = openspec.ParseTasks(ch.Tasks.Content)
	m.tasks.Cursor = m.firstTaskIdx()
}

func (m *Model) firstTaskIdx() int {
	for i, item := range m.tasks.Items {
		if item.Kind == openspec.KindTask {
			return i
		}
	}
	return 0
}

func (m *Model) moveCursorDown() {
	for i := m.tasks.Cursor + 1; i < len(m.tasks.Items); i++ {
		if m.tasks.Items[i].Kind == openspec.KindTask {
			m.tasks.Cursor = i
			return
		}
	}
}

func (m *Model) moveCursorUp() {
	for i := m.tasks.Cursor - 1; i >= 0; i-- {
		if m.tasks.Items[i].Kind == openspec.KindTask {
			m.tasks.Cursor = i
			return
		}
	}
}

func (m *Model) doToggle() tea.Cmd {
	if len(m.tasks.Items) == 0 || m.tasks.Cursor >= len(m.tasks.Items) {
		return nil
	}
	if m.tasks.Items[m.tasks.Cursor].Kind != openspec.KindTask {
		return nil
	}
	ch := m.current()
	if ch == nil {
		return nil
	}
	if err := m.loader.ToggleTask(ch.Path+"/tasks.md", m.tasks.Items, m.tasks.Cursor); err != nil {
		m.errMsg = "error: " + err.Error()
		return tea.Tick(3*time.Second, func(time.Time) tea.Msg { return errClearMsg{} })
	}
	m.refreshTasksViewport()
	return nil
}

var (
	rxCode = regexp.MustCompile("`(.+?)`")
	rxBold = regexp.MustCompile(`\*\*(.+?)\*\*`)

	underlineStyle = lipgloss.NewStyle().Underline(true)
	boldStyle      = lipgloss.NewStyle().Bold(true)
	cyanStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
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
			return underlineStyle.Render(rxCode.FindStringSubmatch(m)[1]) + restore
		})
		s = rxBold.ReplaceAllStringFunc(s, func(m string) string {
			return boldStyle.Render(rxBold.FindStringSubmatch(m)[1]) + restore
		})
	} else {
		s = rxCode.ReplaceAllStringFunc(s, func(m string) string {
			return cyanStyle.Render(rxCode.FindStringSubmatch(m)[1]) + restore
		})
		s = rxBold.ReplaceAllStringFunc(s, func(m string) string {
			return boldStyle.Render(rxBold.FindStringSubmatch(m)[1]) + restore
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

	for i, item := range m.tasks.Items {
		switch item.Kind {
		case openspec.KindSection:
			if i > 0 {
				sb.WriteString("\n")
				line++
			}
			done, total := sectionProgress(m.tasks.Items, i)
			sectionLine := sectionStyle.Render("  "+item.Text) + "  " + progressBar(done, total, 5)
			sb.WriteString(sectionLine + "\n")
			line += lipgloss.Height(sectionLine)
			sb.WriteString("\n")
			line++
		case openspec.KindTask:
			if i == m.tasks.Cursor {
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
			if i == m.tasks.Cursor {
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

func renderProgressBar(done, total, width int, filledChar, emptyChar string) string {
	if total == 0 || width <= 0 {
		return ""
	}
	filled := (done * width) / total
	filledStyle := progressDoneStyle
	if done == total {
		filled = width
		filledStyle = progressCompleteStyle
	}
	return filledStyle.Render(strings.Repeat(filledChar, filled)) +
		progressEmptyStyle.Render(strings.Repeat(emptyChar, width-filled))
}

func progressBar(done, total, width int) string {
	bar := renderProgressBar(done, total, width, "─", "─")
	return bar + helpStyle.Render(fmt.Sprintf(" %d/%d", done, total))
}
