package ui

import (
	"regexp"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"github.com/aemonge/dossier/internal/openspec"
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
	case m.mode == ModeNormal && m.isTasksView() && m.artifactSelection.OutputPath != "":
		return m.loadViewportForTasks()
	case m.viewingCode && m.mode == ModeNormal:
		return m.loadViewportForGit()
	default:
		return m.loadViewportForArtifact()
	}
}

func (m *Model) ensureRenderer(width int) {
	if m.glamourRenderer != nil && m.lastRenderWidth == width {
		return
	}
	options := []glamour.TermRendererOption{glamour.WithWordWrap(width)}
	if m.theme.GlamourStyle == "" {
		options = append(options, glamour.WithStandardStyle("dark"))
	} else {
		options = append(options, glamour.WithStyles(m.theme.GlamourStyles))
	}
	r, err := glamour.NewTermRenderer(options...)
	if err != nil {
		m.glamourRenderer = nil
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
	m.vp.SetContent("")
	width := m.renderWidth()
	m.ensureRenderer(width)
	if m.glamourRenderer == nil {
		return func() tea.Msg {
			return renderedConfigMsg{content: raw}
		}
	}
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
	m.vp.SetContent("")
	width := m.renderWidth()
	m.ensureRenderer(width)
	if m.glamourRenderer == nil {
		return func() tea.Msg {
			return specRenderedMsg{content: raw}
		}
	}

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

func (m *Model) loadViewportForGit() tea.Cmd {
	m.refreshGitViewport()
	return nil
}

func (m *Model) loadViewportForArtifact() tea.Cmd {
	if key, ok := m.selectedArtifactKey(); ok {
		if cached, cachedOK := m.artifactRenderCache[key]; cachedOK {
			m.vp.SetContent(cached)
			return nil
		}
		artifact, output, _ := m.selectedArtifactOutput()
		raw := artifactStatusView(artifact)
		if output != nil {
			raw = output.Content
		}
		if raw == "" {
			m.vp.SetContent("  (artifact not available)")
			return nil
		}
		m.loading = true
		m.vp.SetContent("")
		width := m.renderWidth()
		m.ensureRenderer(width)
		if m.glamourRenderer == nil {
			return func() tea.Msg { return artifactRenderedMsg{key: key, content: raw} }
		}
		return func() tea.Msg {
			out, err := m.glamourRenderer.Render(raw)
			if err != nil {
				return artifactRenderedMsg{key: key, content: raw}
			}
			return artifactRenderedMsg{key: key, content: out}
		}
	}

	m.vp.SetContent("  (artifact not available)")
	return nil
}

func (m *Model) renderWidth() int {
	width := m.width - 2
	if width < 20 {
		width = 80
	}
	return width
}
