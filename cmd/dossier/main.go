package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fselich/dossier/internal/openspec"
	"github.com/fselich/dossier/internal/ui"
)

func main() {
	var (
		project *openspec.Project
		err     error
		model   ui.Model
	)

	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("Usage: dossier [path]")
		fmt.Println()
		fmt.Println("A keyboard-driven TUI for navigating OpenSpec project artifacts.")
		fmt.Println()
		fmt.Println("  path  Optional path to a change directory (single-change mode)")
		os.Exit(0)
	}

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
