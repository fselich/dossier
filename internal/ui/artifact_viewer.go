package ui

import (
	"fmt"
	"path/filepath"

	"github.com/aemonge/dossier/internal/openspec"
)

type artifactSelection struct {
	ArtifactID string
	OutputPath string
}

type artifactOutputIdentity struct {
	ChangeName string
	Archived   bool
	ArtifactID string
	OutputPath string
}

type artifactRenderedMsg struct {
	key     artifactOutputIdentity
	content string
}

func (m *Model) selectDefaultArtifact() bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		m.artifactSelection = artifactSelection{}
		return false
	}
	return m.selectArtifactOutput(ch.Artifacts[0].ID, firstArtifactOutputPath(ch.Artifacts[0]))
}

func firstArtifactOutputPath(artifact openspec.ChangeArtifact) string {
	for _, output := range artifact.Outputs {
		if output.Present {
			return output.RelativePath
		}
	}
	return ""
}

func (m *Model) selectArtifactOutput(artifactID, outputPath string) bool {
	ch := m.current()
	if ch == nil {
		return false
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != artifactID {
			continue
		}
		if outputPath != "" {
			found := false
			for _, output := range artifact.Outputs {
				if output.RelativePath == outputPath && output.Present {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		} else {
			outputPath = firstArtifactOutputPath(artifact)
		}
		m.artifactSelection = artifactSelection{ArtifactID: artifactID, OutputPath: outputPath}
		m.viewingCode = false
		return true
	}
	return false
}

func artifactOutputByID(change *openspec.Change, artifactID string) (openspec.ChangeArtifact, *openspec.ArtifactOutput, bool) {
	if change == nil {
		return openspec.ChangeArtifact{}, nil, false
	}
	for _, artifact := range change.Artifacts {
		if artifact.ID != artifactID {
			continue
		}
		for outputIndex := range artifact.Outputs {
			if artifact.Outputs[outputIndex].Present {
				return artifact, &artifact.Outputs[outputIndex], true
			}
		}
		return artifact, nil, true
	}
	return openspec.ChangeArtifact{}, nil, false
}

func (m *Model) selectedArtifactOutput() (openspec.ChangeArtifact, *openspec.ArtifactOutput, bool) {
	ch := m.current()
	if ch == nil {
		return openspec.ChangeArtifact{}, nil, false
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != m.artifactSelection.ArtifactID {
			continue
		}
		if m.artifactSelection.OutputPath == "" {
			return artifact, nil, true
		}
		for outputIndex := range artifact.Outputs {
			if artifact.Outputs[outputIndex].RelativePath == m.artifactSelection.OutputPath && artifact.Outputs[outputIndex].Present {
				return artifact, &artifact.Outputs[outputIndex], true
			}
		}
		return artifact, nil, true
	}
	return openspec.ChangeArtifact{}, nil, false
}

func (m *Model) selectedArtifactFilePath() (string, bool) {
	if m.mode == ModeViewingArchive || m.viewingCode {
		return "", false
	}
	_, output, ok := m.selectedArtifactOutput()
	change := m.current()
	if !ok || output == nil || change == nil {
		return "", false
	}
	return safeIndexFile(change.Path, filepath.Join(change.Path, filepath.FromSlash(output.RelativePath)))
}

func (m *Model) selectedTaskOutput() (*openspec.ArtifactOutput, bool) {
	_, output, ok := artifactOutputByID(m.current(), "tasks")
	return output, ok && output != nil
}

func (m *Model) isTasksView() bool {
	return !m.viewingCode && m.artifactSelection.ArtifactID == "tasks"
}

func (m *Model) selectedArtifactKey() (artifactOutputIdentity, bool) {
	artifact, output, ok := m.selectedArtifactOutput()
	if !ok {
		return artifactOutputIdentity{}, false
	}
	ch := m.current()
	key := artifactOutputIdentity{
		ChangeName: ch.Name,
		Archived:   m.mode == ModeViewingArchive,
		ArtifactID: artifact.ID,
	}
	if output != nil {
		key.OutputPath = output.RelativePath
	}
	return key, true
}

func (m *Model) invalidateArtifactRenderCache(changeName string) {
	for key := range m.artifactRenderCache {
		if key.ChangeName == changeName {
			delete(m.artifactRenderCache, key)
		}
	}
}

func artifactSelectionPosition(change *openspec.Change, selection artifactSelection) (int, int) {
	if change == nil {
		return 0, 0
	}
	for artifactIndex, artifact := range change.Artifacts {
		if artifact.ID != selection.ArtifactID {
			continue
		}
		for outputIndex, output := range artifact.Outputs {
			if output.RelativePath == selection.OutputPath {
				return artifactIndex, outputIndex
			}
		}
		return artifactIndex, 0
	}
	return 0, 0
}

func (m *Model) reconcileArtifactSelectionNear(artifactIndex, outputIndex int) bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		m.artifactSelection = artifactSelection{}
		return false
	}
	if m.selectArtifactOutput(m.artifactSelection.ArtifactID, m.artifactSelection.OutputPath) {
		return true
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != m.artifactSelection.ArtifactID {
			continue
		}
		var present []openspec.ArtifactOutput
		for _, output := range artifact.Outputs {
			if output.Present {
				present = append(present, output)
			}
		}
		if len(present) == 0 {
			return m.selectArtifactOutput(artifact.ID, "")
		}
		outputIndex = min(max(0, outputIndex), len(present)-1)
		return m.selectArtifactOutput(artifact.ID, present[outputIndex].RelativePath)
	}
	artifactIndex = min(max(0, artifactIndex), len(ch.Artifacts)-1)
	artifact := ch.Artifacts[artifactIndex]
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func (m *Model) reconcileArtifactSelection() bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		m.artifactSelection = artifactSelection{}
		return false
	}
	if m.artifactSelection.ArtifactID == "" {
		return m.selectDefaultArtifact()
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != m.artifactSelection.ArtifactID {
			continue
		}
		if m.artifactSelection.OutputPath == "" {
			return true
		}
		for _, output := range artifact.Outputs {
			if output.RelativePath == m.artifactSelection.OutputPath && output.Present {
				return true
			}
		}
		return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
	}
	return m.selectDefaultArtifact()
}

func (m *Model) moveArtifactSelection(delta int) bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		return false
	}
	current := 0
	for index := range ch.Artifacts {
		if ch.Artifacts[index].ID == m.artifactSelection.ArtifactID {
			current = index
			break
		}
	}
	next := (current + delta%len(ch.Artifacts) + len(ch.Artifacts)) % len(ch.Artifacts)
	artifact := ch.Artifacts[next]
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func (m *Model) moveViewerDestination(delta int) bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		return false
	}
	includeCode := m.isGitRepo && m.mode == ModeNormal && len(m.gitState.Files) > 0
	count := len(ch.Artifacts)
	if includeCode {
		count++
	}
	current := 0
	if m.viewingCode {
		current = len(ch.Artifacts)
	} else {
		for index := range ch.Artifacts {
			if ch.Artifacts[index].ID == m.artifactSelection.ArtifactID {
				current = index
				break
			}
		}
	}
	next := (current + delta%count + count) % count
	if includeCode && next == len(ch.Artifacts) {
		m.viewingCode = true
		return true
	}
	artifact := ch.Artifacts[next]
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func (m *Model) selectArtifactPosition(position int) bool {
	ch := m.current()
	if ch == nil || position < 0 || position >= len(ch.Artifacts) {
		return false
	}
	artifact := ch.Artifacts[position]
	if artifact.ID == m.artifactSelection.ArtifactID && len(artifact.Outputs) > 1 {
		current := 0
		for index := range artifact.Outputs {
			if artifact.Outputs[index].RelativePath == m.artifactSelection.OutputPath {
				current = index
				break
			}
		}
		next := (current + 1) % len(artifact.Outputs)
		return m.selectArtifactOutput(artifact.ID, artifact.Outputs[next].RelativePath)
	}
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func artifactStatusView(artifact openspec.ChangeArtifact) string {
	status := "unknown"
	switch artifact.Status {
	case openspec.ArtifactStatusDone:
		status = "authored"
	case openspec.ArtifactStatusReady:
		status = "ready to author"
	case openspec.ArtifactStatusBlocked:
		status = "blocked"
	default:
		if artifact.Source == openspec.ArtifactSourceDiscovered {
			status = "discovered"
		}
	}
	text := fmt.Sprintf("# %s\n\nStatus: %s", artifact.ID, status)
	if len(artifact.Requires) > 0 {
		text += "\n\nRequires: "
		for index, requirement := range artifact.Requires {
			if index > 0 {
				text += ", "
			}
			text += requirement
		}
	}
	return text + "\n"
}
