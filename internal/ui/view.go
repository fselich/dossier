package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func (m *Model) mainViewContent() string {
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

func (m *Model) viewContentWithChrome() string {
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

func (m *Model) emptyViewContent() string {
	return m.boxTop() + "\n" +
		m.addBorderSides(m.theme.Styles.Header.Render(m.project.Name)+
			"\n\n\n  No active changes. Create one with /opsx:propose\n"+
			m.theme.Styles.Help.Render("\n  a/Esc: index  q: quit")) + "\n" +
		m.boxInnerSep() + "\n" +
		m.addBorderSides(m.renderHelpBar()) + "\n" +
		m.boxBottom()
}

func (m *Model) renderHeader() string {
	if m.mode == ModeViewingConfig {
		return m.theme.Styles.Header.Width(m.width - 2).Render(m.project.Name + "  ·  project config")
	}
	if m.mode == ModeIndex {
		return m.theme.Styles.Header.Width(m.width - 2).Render(m.project.Name + "  ·  index")
	}
	if m.mode == ModeViewingSpec {
		specName := ""
		if m.specViewer.Cursor < len(m.projectSpecs) {
			specName = m.projectSpecs[m.specViewer.Cursor].Name
		}
		if m.specViewer.FocusMode && m.specViewer.Cursor < len(m.projectSpecs) {
			ps := m.projectSpecs[m.specViewer.Cursor]
			return m.theme.Styles.Header.Width(m.width - 2).Render(
				fmt.Sprintf("%s  ·  %s  ·  Req %d/%d", m.project.Name, specName, m.specViewer.ReqCursor+1, len(ps.RequirementNames)),
			)
		}
		return m.theme.Styles.Header.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s  [spec]", m.project.Name, specName),
		)
	}
	ch := m.current()
	if ch == nil {
		return m.theme.Styles.Header.Render(m.project.Name)
	}
	if m.mode == ModeViewingArchive {
		return m.theme.Styles.Header.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s  [archive]", m.project.Name, ch.Name),
		)
	}
	nav := fmt.Sprintf("[%d/%d]", m.changeIdx+1, len(m.project.Changes))
	return m.theme.Styles.Header.Width(m.width - 2).Render(
		fmt.Sprintf("%s  ·  %s  %s", m.project.Name, ch.Name, nav),
	)
}

func (m *Model) renderTabBar() string {
	parts := make([]string, 0, tabCount)
	for t := Tab(0); t < tabCount; t++ {
		if t == TabGit && m.mode != ModeNormal {
			continue
		}
		label := tabLabels[t]
		if t == TabGit && len(m.gitState.Files) > 0 {
			label = "code (" + fmt.Sprintf("%d", len(m.gitState.Files)) + ")"
		}
		switch {
		case t == m.tab:
			parts = append(parts, m.theme.Styles.TabActive.Render(label))
		case !m.tabAvailable(t):
			parts = append(parts, m.theme.Styles.TabDisabled.Render(label))
		default:
			parts = append(parts, m.theme.Styles.TabInactive.Render(label))
		}
	}
	tabs := strings.Join(parts, " ")

	taskItems := m.tasks.Items
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
		barSpace := (m.width - 2) - lipgloss.Width(tabs) - 3 - len(label)
		if barSpace >= 3 {
			tabs = tabs + " [" + m.renderProgressBar(done, total, barSpace, "█", "░") + "]" + m.theme.Styles.Help.Render(label)
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
			parts = append(parts, m.theme.Styles.TabActive.Render(s.Name))
		} else {
			parts = append(parts, m.theme.Styles.TabInactive.Render(s.Name))
		}
	}
	return strings.Join(parts, " ")
}

func (m *Model) hasSpecSubnav() bool {
	ch := m.current()
	return m.tab == TabSpecs && ch != nil && len(ch.SpecFiles) > 0
}

func (m *Model) boxTop() string {
	return m.theme.Styles.Separator.Render("┌" + strings.Repeat("─", m.width-2) + "┐")
}

func (m *Model) boxBottom() string {
	return m.theme.Styles.Separator.Render("└" + strings.Repeat("─", m.width-2) + "┘")
}

func (m *Model) boxInnerSep() string {
	return m.theme.Styles.Separator.Render("├" + strings.Repeat("─", m.width-2) + "┤")
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
		result = append(result, m.theme.Styles.Separator.Render("│")+line+strings.Repeat(" ", pad)+m.theme.Styles.Separator.Render("│"))
	}
	return strings.Join(result, "\n")
}

func (m *Model) renderHelpBar() string {
	if m.errMsg != "" {
		return m.theme.Styles.Error.Render(m.errMsg)
	}
	if m.mode == ModeIndex {
		if m.index.FilterActive {
			return m.theme.Styles.Help.Render("/" + m.index.FilterText + "█")
		}
		sortHint := "s: sort by suffix"
		if m.index.SortBySuffix {
			sortHint = "s: sort by name"
		}
		text := "j/k: navigate  Enter: open  Space: toggle  click: select  " + sortHint + "  i: info  Esc: quit"
		if m.index.FilterText != "" {
			text += "  [/" + m.index.FilterText + "]"
		}
		return m.theme.Styles.Help.Render(text)
	}
	if m.mode == ModeViewingConfig {
		return m.theme.Styles.Help.Render("j/k: scroll  i/Esc: back  q: quit")
	}
	if m.mode == ModeViewingSpec {
		if m.specViewer.FocusMode {
			return m.theme.Styles.Help.Render("h/l: req anterior/siguiente  j/k: scroll  Esc: index  q: quit")
		}
		return m.theme.Styles.Help.Render("j/k: scroll  Esc: index  q: quit")
	}
	if m.mode == ModeViewingArchive {
		return m.theme.Styles.Help.Render("1-4/Tab: artifact  j/k: scroll  a/Esc: index  q: quit")
	}
	tabRange := "1-4"
	if m.isGitRepo {
		tabRange = "1-5"
	}
	if m.tab == TabGit {
		if m.gitState.ErrMsg != "" {
			return m.theme.Styles.Error.Render(m.gitState.ErrMsg)
		}
		if m.gitState.ShowingDiff {
			return m.theme.Styles.Help.Render("d/Esc: back  [/]: prev/next  j/k: vertical  h/l: ←→ horizontal  q: quit")
		}
		return m.theme.Styles.Help.Render("h/l: change  " + tabRange + "/Tab: artifact  j/k: navigate  Enter/e: open file  d: view diff  s: stage/unstage  Esc: index  q: quit")
	}
	if m.tab == TabTasks {
		return m.theme.Styles.Help.Render("h/l: change  " + tabRange + "/Tab: artifact  j/k: navigate  Space: toggle  e: edit  i: info  Esc: index  q: quit")
	}
	return m.theme.Styles.Help.Render("h/l: change  " + tabRange + "/Tab: artifact  j/k: scroll  e: edit  i: info  Esc: index  q: quit")
}
