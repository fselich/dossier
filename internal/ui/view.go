package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/aemonge/dossier/internal/openspec"
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
	keys := m.effectiveKeyMap().Viewer
	hint := combinedKeyLabel(keys.Back) + ": index  " + combinedKeyLabel(keys.Quit) + ": quit"
	return m.boxTop() + "\n" +
		m.addBorderSides(m.theme.Styles.Header.Render(m.project.Name)+
			"\n\n\n  No active changes. Create one with /opsx:propose\n"+
			m.theme.Styles.Help.Render("\n  "+hint)) + "\n" +
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
	schema := ""
	if ch.Schema != "" {
		schema = " [" + ch.Schema + "]"
	}
	if m.mode == ModeViewingArchive {
		return m.theme.Styles.Header.Width(m.width - 2).Render(
			fmt.Sprintf("%s  ·  %s%s  [archive]", m.project.Name, ch.Name, schema),
		)
	}
	nav := fmt.Sprintf("[%d/%d]", m.changeIdx+1, len(m.project.Changes))
	return m.theme.Styles.Header.Width(m.width - 2).Render(
		fmt.Sprintf("%s  ·  %s%s  %s", m.project.Name, ch.Name, schema, nav),
	)
}

func (m *Model) renderTabBar() string {
	var parts []string
	selectedPart := 0
	ch := m.current()
	if ch != nil {
		parts = make([]string, 0, len(ch.Artifacts)+1)
		for artifactIndex, artifact := range ch.Artifacts {
			label := artifactDisplayLabel(artifact) + artifactTabStatusSuffix(artifact)
			if !m.viewingCode && artifact.ID == m.artifactSelection.ArtifactID {
				selectedPart = artifactIndex
				parts = append(parts, m.theme.Styles.TabActive.Render(label))
			} else {
				parts = append(parts, m.theme.Styles.TabInactive.Render(label))
			}
		}
		if m.isGitRepo && m.mode == ModeNormal {
			label := "code"
			if len(m.gitState.Files) > 0 {
				label += " (" + fmt.Sprintf("%d", len(m.gitState.Files)) + ")"
			}
			if m.viewingCode {
				selectedPart = len(parts)
				parts = append(parts, m.theme.Styles.TabActive.Render(label))
			} else if len(m.gitState.Files) == 0 {
				parts = append(parts, m.theme.Styles.TabDisabled.Render(label))
			} else {
				parts = append(parts, m.theme.Styles.TabInactive.Render(label))
			}
		}
	}
	tabs := fitTabParts(parts, selectedPart, m.width-2, m.theme.Styles.Help)

	var taskItems []openspec.TaskItem
	if output, ok := m.selectedTaskOutput(); ok {
		taskItems = openspec.ParseTasks(output.Content)
	} else {
		taskItems = nil
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

func fitTabParts(parts []string, selected, width int, markerStyle lipgloss.Style) string {
	if len(parts) == 0 || width <= 0 {
		return ""
	}
	all := strings.Join(parts, " ")
	if lipgloss.Width(all) <= width {
		return all
	}
	selected = min(max(0, selected), len(parts)-1)
	start, end := selected, selected+1
	render := func(from, to int) string {
		visible := strings.Join(parts[from:to], " ")
		if from > 0 {
			visible = markerStyle.Render("‹") + " " + visible
		}
		if to < len(parts) {
			visible += " " + markerStyle.Render("›")
		}
		return visible
	}
	for {
		expanded := false
		if end < len(parts) && lipgloss.Width(render(start, end+1)) <= width {
			end++
			expanded = true
		}
		if start > 0 && lipgloss.Width(render(start-1, end)) <= width {
			start--
			expanded = true
		}
		if !expanded {
			break
		}
	}
	result := render(start, end)
	if lipgloss.Width(result) > width {
		return lipgloss.NewStyle().MaxWidth(width).Render(result)
	}
	return result
}

func artifactTabStatusSuffix(artifact openspec.ChangeArtifact) string {
	switch artifact.Status {
	case openspec.ArtifactStatusDone:
		return " ✓"
	case openspec.ArtifactStatusReady:
		return " ○"
	case openspec.ArtifactStatusBlocked:
		return " ×"
	}
	if artifact.Source == openspec.ArtifactSourceDiscovered {
		return " ~"
	}
	return ""
}

func (m *Model) renderSpecSubnav() string {
	artifact, _, ok := m.selectedArtifactOutput()
	if !ok || artifact.ID != "specs" {
		return ""
	}
	var parts []string
	for _, output := range artifact.Outputs {
		if !output.Present {
			continue
		}
		name := output.DisplayName
		if name == "" {
			name = output.RelativePath
		}
		if output.RelativePath == m.artifactSelection.OutputPath {
			parts = append(parts, m.theme.Styles.TabActive.Render(name))
		} else {
			parts = append(parts, m.theme.Styles.TabInactive.Render(name))
		}
	}
	return strings.Join(parts, " ")
}

func (m *Model) hasSpecSubnav() bool {
	artifact, _, ok := m.selectedArtifactOutput()
	if !ok || artifact.ID != "specs" {
		return false
	}
	present := 0
	for _, output := range artifact.Outputs {
		if output.Present {
			present++
		}
	}
	return present > 1
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

func (m *Model) helpText(text string) string {
	if m.readOnly {
		return "[read-only]  " + text
	}
	return text
}

func (m *Model) renderHelpBar() string {
	if m.errMsg != "" {
		return m.theme.Styles.Error.Render(m.helpText(m.errMsg))
	}
	keyMap := m.effectiveKeyMap()
	if m.mode == ModeIndex {
		if m.index.FilterActive {
			return m.theme.Styles.Help.Render(m.helpText(primaryKeyLabel(keyMap.Index.Filter) + m.index.FilterText + "█"))
		}
		action := m.index.Action
		switch action.Mode {
		case indexActionSchema:
			schemas := m.filteredSchemas()
			choice := "no matching schemas"
			if len(schemas) > 0 {
				cursor := min(action.SchemaCursor, len(schemas)-1)
				schema := schemas[cursor]
				choice = schema.Name
				if schema.Description != "" {
					choice += " — " + schema.Description
				}
				if len(schema.Artifacts) > 0 {
					choice += " [" + strings.Join(schema.Artifacts, ", ") + "]"
				}
				if schema.Source != "" {
					choice += " (" + schema.Source + ")"
				}
			}
			return m.theme.Styles.Help.Render(m.helpText("schema /" + action.SchemaFilter + "  " + choice + "  " + primaryKeyLabel(keyMap.Index.Open) + ": select  Esc: cancel"))
		case indexActionName:
			text := "new " + action.Schema + " change: " + action.Name + "█  " + primaryKeyLabel(keyMap.Index.Open) + ": create  Esc: cancel"
			if action.Status != "" {
				text += "  " + action.Status
			}
			return m.theme.Styles.Help.Render(m.helpText(text))
		case indexActionConfirmArchive:
			return m.theme.Styles.Help.Render(m.helpText("archive " + m.archiveConfirmationSummary(action.TargetIdentity) + "?  " + combinedKeyLabel(keyMap.Index.Lifecycle, keyMap.Index.Open) + ": confirm  Esc: cancel"))
		case indexActionConfirmReactivate:
			return m.theme.Styles.Help.Render(m.helpText("make " + strings.TrimPrefix(action.TargetIdentity, "archive:") + " active?  " + combinedKeyLabel(keyMap.Index.Lifecycle, keyMap.Index.Open) + ": confirm  Esc: cancel"))
		case indexActionConfirmUndo:
			return m.theme.Styles.Help.Render(m.helpText("undo latest lifecycle action for " + action.Name + "?  " + combinedKeyLabel(keyMap.Index.Undo, keyMap.Index.Open) + ": confirm  Esc: cancel"))
		case indexActionPending:
			return m.theme.Styles.Help.Render(m.helpText(indexOperationLabel(action.Operation) + " in progress…"))
		}
		if action.Status != "" {
			if action.StatusError {
				return m.theme.Styles.Error.Render(m.helpText(action.Status))
			}
			return m.theme.Styles.Help.Render(m.helpText(action.Status))
		}
		sortAction := "sort by suffix"
		if m.index.SortBySuffix {
			sortAction = "sort by name"
		}
		actionHelp := ""
		if item, ok := m.selectedIndexItem(); ok {
			if m.indexItemExpandable(item) {
				actionHelp = combinedKeyLabel(keyMap.Index.Open, keyMap.Index.Toggle) + ": toggle"
				if item.kind != indexKindSection {
					actionHelp += "  " + primaryKeyLabel(keyMap.Index.Inspect) + ": inspect"
				}
			} else if item.kind != indexKindSection {
				actionHelp = combinedKeyLabel(keyMap.Index.Inspect, keyMap.Index.Open) + ": inspect"
			}
		}
		text := pairedKeyLabel(keyMap.Index.Down, keyMap.Index.Up) + ": navigate"
		if actionHelp != "" {
			text += "  " + actionHelp
		}
		if item, ok := m.selectedIndexItem(); ok {
			if !m.readOnly {
				if item.kind == indexKindActive || item.kind == indexKindSpec || item.kind == indexKindRequirement {
					text += "  " + primaryKeyLabel(keyMap.Index.New) + ": new"
				}
				switch item.kind {
				case indexKindActive:
					text += "  " + primaryKeyLabel(keyMap.Index.Lifecycle) + ": archive"
				case indexKindArchived:
					text += "  " + primaryKeyLabel(keyMap.Index.Lifecycle) + ": make active"
				}
				if _, editable := m.indexEditPath(item); editable {
					text += "  " + primaryKeyLabel(keyMap.Index.Edit) + ": edit"
				}
				if m.lifecycleUndo != nil && item.kind != indexKindSection {
					text += "  " + primaryKeyLabel(keyMap.Index.Undo) + ": undo"
				}
			}
			if _, _, validatable := m.indexValidationTarget(item); validatable {
				text += "  " + primaryKeyLabel(keyMap.Index.Validate) + ": validate"
			}
		}
		text += "  click: select  " + primaryKeyLabel(keyMap.Index.Sort) + ": " + sortAction + "  " +
			primaryKeyLabel(keyMap.Index.Info) + ": info  " + primaryKeyLabel(keyMap.Index.Back) + ": quit"
		if m.index.FilterText != "" {
			text += "  [" + primaryKeyLabel(keyMap.Index.Filter) + m.index.FilterText + "]"
		}
		return m.theme.Styles.Help.Render(m.helpText(text))
	}
	if m.mode == ModeViewingConfig {
		keys := keyMap.Config
		return m.theme.Styles.Help.Render(m.helpText(
			pairedKeyLabel(keys.Down, keys.Up) + ": scroll  " + combinedKeyLabel(keys.Back) + ": back",
		))
	}
	if m.mode == ModeViewingSpec {
		keys := keyMap.Spec
		text := ""
		if m.specViewer.FocusMode {
			text = pairedKeyLabel(keys.PreviousRequirement, keys.NextRequirement) + ": previous/next requirement  "
		}
		text += pairedKeyLabel(keys.Down, keys.Up) + ": scroll  " + primaryKeyLabel(keys.Back) + ": index  " +
			combinedKeyLabel(keys.Quit) + ": quit"
		return m.theme.Styles.Help.Render(m.helpText(text))
	}
	keys := keyMap.Viewer
	tabKeys := combinedKeyLabel(keys.ProposalTab, keys.DesignTab, keys.SpecsTab, keys.TasksTab)
	if m.isGitRepo && m.mode == ModeNormal {
		tabKeys = combinedKeyLabel(keys.ProposalTab, keys.DesignTab, keys.SpecsTab, keys.TasksTab, keys.GitTab)
	}
	tabKeys += "/" + pairedKeyLabel(keys.PreviousTab, keys.NextTab)
	if m.mode == ModeViewingArchive {
		return m.theme.Styles.Help.Render(m.helpText(
			tabKeys + ": artifact  " + pairedKeyLabel(keys.Down, keys.Up) + ": scroll  " +
				combinedKeyLabel(keys.Back) + ": index  " + combinedKeyLabel(keys.Quit) + ": quit",
		))
	}
	if m.viewingCode {
		if m.gitState.ErrMsg != "" {
			return m.theme.Styles.Error.Render(m.helpText(m.gitState.ErrMsg))
		}
		if m.gitState.ShowingDiff {
			text := combinedKeyLabel(keys.ViewDiff, keys.Back) + ": back  " + pairedKeyLabel(keys.PreviousDiff, keys.NextDiff) +
				": previous/next  " + pairedKeyLabel(keys.Down, keys.Up) + ": vertical  " +
				pairedKeyLabel(keys.Left, keys.Right) + ": horizontal  " + combinedKeyLabel(keys.Quit) + ": quit"
			return m.theme.Styles.Help.Render(m.helpText(text))
		}
		text := pairedKeyLabel(keys.Previous, keys.Next) + ": change  " + tabKeys + ": artifact  " +
			pairedKeyLabel(keys.Down, keys.Up) + ": navigate  " + combinedKeyLabel(keys.ViewDiff, keys.Open) + ": view diff"
		if !m.readOnly {
			text += "  " + primaryKeyLabel(keys.Stage) + ": stage/unstage"
		}
		return m.theme.Styles.Help.Render(m.helpText(text + "  " + combinedKeyLabel(keys.Back) + ": index  " + combinedKeyLabel(keys.Quit) + ": quit"))
	}
	if m.isTasksView() {
		text := pairedKeyLabel(keys.Previous, keys.Next) + ": change  " + tabKeys + ": artifact  " +
			pairedKeyLabel(keys.Down, keys.Up) + ": navigate"
		if !m.readOnly {
			text += "  " + primaryKeyLabel(keys.ToggleTask) + ": toggle  " + primaryKeyLabel(keys.Open) + ": edit"
		}
		return m.theme.Styles.Help.Render(m.helpText(text + "  " + primaryKeyLabel(keys.Info) + ": info  " +
			combinedKeyLabel(keys.Back) + ": index  " + combinedKeyLabel(keys.Quit) + ": quit"))
	}
	text := pairedKeyLabel(keys.Previous, keys.Next) + ": change  " + tabKeys + ": artifact  " +
		pairedKeyLabel(keys.Down, keys.Up) + ": scroll"
	if !m.readOnly {
		text += "  " + primaryKeyLabel(keys.Open) + ": edit"
	}
	return m.theme.Styles.Help.Render(m.helpText(text + "  " + primaryKeyLabel(keys.Info) + ": info  " +
		combinedKeyLabel(keys.Back) + ": index  " + combinedKeyLabel(keys.Quit) + ": quit"))
}
