package ui

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/fselich/dossier/internal/openspec"
)

func TestRenderTabBarCapsProgressWidth(t *testing.T) {
	m := &Model{
		tab:   TabProposal,
		width: 160,
		project: &openspec.Project{
			Changes: []openspec.Change{{
				Name:     "test",
				Proposal: openspec.Artifact{Present: true},
				Design:   openspec.Artifact{Present: true},
				Specs:    openspec.Artifact{Present: true},
				Tasks:    openspec.Artifact{Present: true},
			}},
		},
		tasks: taskState{
			Items: []openspec.TaskItem{
				{Kind: openspec.KindTask, Done: true},
				{Kind: openspec.KindTask, Done: true},
				{Kind: openspec.KindTask, Done: true},
				{Kind: openspec.KindTask, Done: true},
				{Kind: openspec.KindTask},
			},
		},
		mode: ModeNormal,
	}

	got := m.renderTabBar()

	if lipgloss.Width(got) > 80 {
		t.Errorf("expected compact tab bar, got width %d", lipgloss.Width(got))
	}
}
