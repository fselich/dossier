package ui

import (
	"testing"

	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

func TestRootStartupDefaultsToIndex(t *testing.T) {
	root := t.TempDir()
	project := &openspec.Project{Name: "project", Changes: []openspec.Change{{Name: "change", Path: root}}}
	m := New(project, openspec.ProjectConfig{}, root, openspec.NewLoader(openspec.OSFS{}), Theme{}, settings.DefaultKeys(), false)
	if m.mode != ModeIndex {
		t.Fatalf("default startup mode = %d, want ModeIndex", m.mode)
	}
}

func TestRootStartupCanOpenChange(t *testing.T) {
	root := t.TempDir()
	project := &openspec.Project{Name: "project", Changes: []openspec.Change{{Name: "change", Path: root}}}
	m := NewWithStartView(project, openspec.ProjectConfig{}, root, openspec.NewLoader(openspec.OSFS{}), Theme{}, settings.DefaultKeys(), false, "change")
	if m.mode != ModeNormal {
		t.Fatalf("configured startup mode = %d, want ModeNormal", m.mode)
	}
}

func TestChangeStartupWithoutChangesFallsBackToIndex(t *testing.T) {
	root := t.TempDir()
	m := NewWithStartView(&openspec.Project{Name: "project"}, openspec.ProjectConfig{}, root, openspec.NewLoader(openspec.OSFS{}), Theme{}, settings.DefaultKeys(), false, "change")
	if m.mode != ModeIndex {
		t.Fatalf("empty configured startup mode = %d, want ModeIndex", m.mode)
	}
}

func TestSinglePathStartupBypassesIndexPreference(t *testing.T) {
	root := t.TempDir()
	project := &openspec.Project{Name: "project", Changes: []openspec.Change{{Name: "change", Path: root}}}
	m := NewSinglePath(project, openspec.ProjectConfig{}, root, openspec.NewLoader(openspec.OSFS{}), Theme{}, settings.DefaultKeys(), false)
	if m.mode != ModeNormal || !m.singlePath {
		t.Fatalf("single path startup = mode %d single=%v, want normal direct target", m.mode, m.singlePath)
	}
}
