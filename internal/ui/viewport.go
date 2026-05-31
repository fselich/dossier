package ui

import (
	"regexp"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func (m *Model) loadViewport() tea.Cmd {
	if !m.vpReady {
		return nil
	}
	switch {
	case m.mode == ModeIndex:
		return m.loadViewportForIndex()
	case m.mode == ModeViewingConfig:
		return m.loadViewportForConfig()
	case m.mode == ModeViewingSpec:
		return m.loadViewportForSpec()
	case m.tab == TabTasks && m.mode == ModeNormal:
		return m.loadViewportForTasks()
	default:
		return m.loadViewportForArtifact()
	}
}

func (m *Model) ensureRenderer(width int) {
	if m.glamourRenderer != nil && m.lastRenderWidth == width {
		return
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return
	}
	m.glamourRenderer = r
	m.lastRenderWidth = width
}

func (m *Model) loadViewportForIndex() tea.Cmd {
	m.refreshIndexViewport()
	return nil
}

func (m *Model) loadViewportForConfig() tea.Cmd {
	raw := openspec.ConfigToMarkdown(m.projectConfig)
	if raw == "" {
		m.vp.SetContent("  (no project config found)")
		return nil
	}
	m.loading = true
	m.vp.SetContent(raw)
	width := m.renderWidth()
	m.ensureRenderer(width)
	return func() tea.Msg {
		out, err := m.glamourRenderer.Render(raw)
		if err != nil {
			return renderedConfigMsg{content: raw}
		}
		return renderedConfigMsg{content: out}
	}
}

func (m *Model) loadViewportForSpec() tea.Cmd {
	if m.specViewer.Cursor >= len(m.projectSpecs) {
		m.vp.SetContent("  (spec not available)")
		return nil
	}
	raw := m.projectSpecs[m.specViewer.Cursor].Content
	if raw == "" {
		m.vp.SetContent("  (spec not available)")
		return nil
	}
	m.loading = true
	m.vp.SetContent(raw)
	width := m.renderWidth()
	m.ensureRenderer(width)

	if m.specViewer.FocusMode {
		jumpTarget := m.specViewer.JumpTarget
		return func() tea.Msg {
			block := openspec.ExtractRequirement(raw, jumpTarget)
			if block == "" {
				return specRenderedMsg{content: "  (spec not available)"}
			}
			out, err := m.glamourRenderer.Render(block)
			if err != nil {
				return specRenderedMsg{content: block}
			}
			return specRenderedMsg{content: out}
		}
	}

	jumpTarget := m.specViewer.JumpTarget
	ansiRe := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return func() tea.Msg {
		out, err := m.glamourRenderer.Render(raw)
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

func (m *Model) loadViewportForTasks() tea.Cmd {
	m.refreshTasksViewport()
	return nil
}

func (m *Model) loadViewportForArtifact() tea.Cmd {
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

	m.loading = true
	m.vp.SetContent(raw)

	tab := m.tab
	width := m.renderWidth()
	m.ensureRenderer(width)
	return func() tea.Msg {
		out, err := m.glamourRenderer.Render(raw)
		if err != nil {
			return renderedMsg{tab: tab, content: raw}
		}
		return renderedMsg{tab: tab, content: out}
	}
}

func (m *Model) renderWidth() int {
	width := m.width - 2
	if width < 20 {
		width = 80
	}
	return width
}
