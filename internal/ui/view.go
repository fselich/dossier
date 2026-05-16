package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fselich/dossier/internal/openspec"
)

func (m Model) View() string {
	if !m.vpReady {
		return ""
	}
	if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		return m.viewIndex()
	}
	if len(m.project.Changes) == 0 && m.mode == ModeNormal {
		return m.emptyView()
	}
	rows := []string{
		m.boxTop(),
		m.addBorderSides(m.renderHeader()),
		m.addBorderSides(m.renderTabBar()),
		m.boxInnerSep(),
	}
	if m.hasSpecSubnav() {
		rows = append(rows, m.addBorderSides(m.renderSpecSubnav()))
	}
	rows = append(rows,
		m.addBorderSides(m.vp.View()),
		m.boxInnerSep(),
		m.addBorderSides(m.renderHelpBar()),
		m.boxBottom(),
	)
	return strings.Join(rows, "\n")
}

func (m *Model) viewIndex() string {
	rows := []string{
		m.boxTop(),
		m.addBorderSides(m.renderHeader()),
		m.boxInnerSep(),
		m.addBorderSides(m.vp.View()),
		m.boxInnerSep(),
		m.addBorderSides(m.renderHelpBar()),
		m.boxBottom(),
	}
	return strings.Join(rows, "\n")
}

func (m *Model) renderHeader() string {
	if m.mode == ModeIndex {
		return headerStyle.Width(m.width - 2).Render(m.project.Name + "  ·  index")
	}
	if m.mode == ModeViewingSpec {
		specName := ""
		if m.specViewerCursor < len(m.projectSpecs) {
			specName = m.projectSpecs[m.specViewerCursor].Name
		}
		if m.specFocusMode && m.specViewerCursor < len(m.projectSpecs) {
			ps := m.projectSpecs[m.specViewerCursor]
			return headerStyle.Width(m.width - 2).Render(
				fmt.Sprintf("%s  ·  %s  ·  Req %d/%d", m.project.Name, specName, m.specReqCursor+1, len(ps.RequirementNames)),
			)
		}
		return headerStyle.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s  [spec]", m.project.Name, specName),
		)
	}
	ch := m.current()
	if ch == nil {
		return headerStyle.Render(m.project.Name)
	}
	if m.mode == ModeViewingArchive {
		return headerStyle.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s  [archive]", m.project.Name, ch.Name),
		)
	}
	nav := fmt.Sprintf("[%d/%d]", m.changeIdx+1, len(m.project.Changes))
	return headerStyle.Width(m.width - 2).Render(
		fmt.Sprintf("%s  ·  %s  %s", m.project.Name, ch.Name, nav),
	)
}

func (m *Model) renderTabBar() string {
	var parts []string
	for t := Tab(0); t < tabCount; t++ {
		label := tabLabels[t]
		switch {
		case t == m.tab:
			parts = append(parts, tabActiveStyle.Render(label))
		case !m.tabAvailable(t):
			parts = append(parts, tabDisabledStyle.Render(label))
		default:
			parts = append(parts, tabInactiveStyle.Render(label))
		}
	}
	tabs := strings.Join(parts, " ")

	// Progress bar right-aligned on the same line as the tabs
	taskItems := m.taskItems
	if m.mode == ModeViewingArchive {
		if ch := m.currentArchive(); ch != nil && ch.Tasks.Present {
			taskItems = openspec.ParseTasks(ch.Tasks.Content)
		} else {
			taskItems = nil
		}
	}
	total, done := 0, 0
	for _, item := range taskItems {
		if item.Kind == openspec.KindTask {
			total++
			if item.Done {
				done++
			}
		}
	}
	if total > 0 {
		label := fmt.Sprintf(" %d/%d", done, total)
		barSpace := (m.width-2) - lipgloss.Width(tabs) - 3 - len(label)
		if barSpace >= 3 {
			filled := (done * barSpace) / total
			if done == total {
				filled = barSpace
			}
			bar := "[" + progressDoneStyle.Render(strings.Repeat("█", filled)) +
				progressEmptyStyle.Render(strings.Repeat("░", barSpace-filled)) + "]"
			tabs = tabs + " " + bar + helpStyle.Render(label)
		}
	}
	return tabs
}

func (m *Model) renderSpecSubnav() string {
	ch := m.current()
	if ch == nil {
		return ""
	}
	var parts []string
	for i, s := range ch.SpecFiles {
		if i == m.specIdx {
			parts = append(parts, tabActiveStyle.Render(s.Name))
		} else {
			parts = append(parts, tabInactiveStyle.Render(s.Name))
		}
	}
	return strings.Join(parts, " ")
}

func (m *Model) hasSpecSubnav() bool {
	ch := m.current()
	return m.tab == TabSpecs && ch != nil && len(ch.SpecFiles) > 0
}

func (m *Model) boxTop() string {
	return separatorStyle.Render("┌" + strings.Repeat("─", m.width-2) + "┐")
}

func (m *Model) boxBottom() string {
	return separatorStyle.Render("└" + strings.Repeat("─", m.width-2) + "┘")
}

func (m *Model) boxInnerSep() string {
	return separatorStyle.Render("├" + strings.Repeat("─", m.width-2) + "┤")
}

func (m *Model) addBorderSides(content string) string {
	lines := strings.Split(content, "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	inner := m.width - 2
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		pad := inner - lipgloss.Width(line)
		if pad < 0 {
			pad = 0
		}
		result = append(result, separatorStyle.Render("│")+line+strings.Repeat(" ", pad)+separatorStyle.Render("│"))
	}
	return strings.Join(result, "\n")
}

func (m *Model) renderHelpBar() string {
	if m.errMsg != "" {
		return errStyle.Render(m.errMsg)
	}
	if m.mode == ModeIndex {
		sortHint := "s: sort by suffix"
		if m.specSortBySuffix {
			sortHint = "s: sort by name"
		}
		return helpStyle.Render("j/k: navigate  Enter: open  Space: expand  " + sortHint + "  Esc: quit")
	}
	if m.mode == ModeViewingSpec {
		if m.specFocusMode {
			return helpStyle.Render("h/l: req anterior/siguiente  j/k: scroll  Esc: índice  q: quit")
		}
		return helpStyle.Render("j/k: scroll  Esc: index  q: quit")
	}
	if m.mode == ModeViewingArchive {
		return helpStyle.Render("1-4: artifact  j/k: scroll  a/Esc: index  q: quit")
	}
	if m.tab == TabTasks {
		return helpStyle.Render("h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit")
	}
	return helpStyle.Render("h/l: change  1-4: artifact  j/k: scroll  e: edit  Esc: index  q: quit")
}

func (m *Model) emptyView() string {
	return headerStyle.Render(m.project.Name) +
		"\n\n\n  No active changes. Create one with /opsx:propose\n" +
		helpStyle.Render("\n  a/Esc: index  q: quit")
}
