package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/felipesanchez/dossier/internal/openspec"
	"github.com/felipesanchez/dossier/internal/ui"
)

func main() {
	var (
		project *openspec.Project
		err     error
		model   ui.Model
	)

	if len(os.Args) > 1 {
		project, err = openspec.LoadFromPath(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		model = ui.NewSinglePath(project)
	} else {
		project, err = openspec.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		model = ui.New(project)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
