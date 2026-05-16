package ui

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func (m *Model) loadViewport() tea.Cmd {
	if !m.vpReady {
		return nil
	}
	if m.mode == ModeIndex {
		m.refreshIndexViewport()
		return nil
	}
	if m.mode == ModeViewingSpec {
		if m.specViewerCursor >= len(m.projectSpecs) {
			m.vp.SetContent("  (spec not available)")
			return nil
		}
		raw := m.projectSpecs[m.specViewerCursor].Content
		if raw == "" {
			m.vp.SetContent("  (spec not available)")
			return nil
		}
		m.loading = true
		m.vp.SetContent("\n  Cargando...")
		width := m.width - 2
		if width < 20 {
			width = 80
		}
		if m.specFocusMode {
			jumpTarget := m.specJumpTarget
			return func() tea.Msg {
				block := extractRequirement(raw, jumpTarget)
				if block == "" {
					return specRenderedMsg{content: "  (spec not available)"}
				}
				r, err := glamour.NewTermRenderer(
					glamour.WithStandardStyle("dark"),
					glamour.WithWordWrap(width),
				)
				if err != nil {
					return specRenderedMsg{content: block}
				}
				out, err := r.Render(block)
				if err != nil {
					return specRenderedMsg{content: block}
				}
				return specRenderedMsg{content: out}
			}
		}
		jumpTarget := m.specJumpTarget
		ansiRe := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		return func() tea.Msg {
			r, err := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(width),
			)
			if err != nil {
				return specRenderedMsg{content: raw}
			}
			out, err := r.Render(raw)
			if err != nil {
				return specRenderedMsg{content: raw}
			}
			jumpLine := 0
			if jumpTarget != "" {
				for i, l := range strings.Split(out, "\n") {
					if strings.Contains(ansiRe.ReplaceAllString(l, ""), jumpTarget) {
						jumpLine = i
						break
					}
				}
			}
			return specRenderedMsg{content: out, jumpLine: jumpLine}
		}
	}
	if m.tab == TabTasks && m.mode == ModeNormal {
		m.refreshTasksViewport()
		return nil
	}

	// Cache hit — instant.
	if cached, ok := m.renderCache[m.tab]; ok {
		m.vp.SetContent(cached)
		return nil
	}

	ch := m.current()
	if ch == nil {
		m.vp.SetContent("")
		return nil
	}
	var raw string
	switch m.tab {
	case TabProposal:
		raw = ch.Proposal.Content
	case TabDesign:
		raw = ch.Design.Content
	case TabSpecs:
		if m.specIdx < len(ch.SpecFiles) {
			raw = ch.SpecFiles[m.specIdx].Content
		}
	case TabTasks:
		raw = ch.Tasks.Content
	}
	if raw == "" {
		m.vp.SetContent("  (artifact not available)")
		return nil
	}

	// Show placeholder immediately, render in background.
	m.loading = true
	m.vp.SetContent("\n  Cargando...")

	tab := m.tab
	width := m.width - 2
	if width < 20 {
		width = 80
	}
	return func() tea.Msg {
		r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return renderedMsg{tab: tab, content: raw}
		}
		out, err := r.Render(raw)
		if err != nil {
			return renderedMsg{tab: tab, content: raw}
		}
		return renderedMsg{tab: tab, content: out}
	}
}
