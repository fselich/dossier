package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/openspec"
	"github.com/fselich/dossier/internal/ui"
)

var version string

func main() {
	var (
		project *openspec.Project
		err     error
		model   ui.Model
	)

	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("dossier", version)
		os.Exit(0)
	}

	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("Usage: dossier [path]")
		fmt.Println()
		fmt.Println("A keyboard-driven TUI for navigating OpenSpec project artifacts.")
		fmt.Println()
		fmt.Println("  path  Optional path to a change directory (single-change mode)")
		os.Exit(0)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot determine working directory:", err)
		os.Exit(1)
	}

	cfg, err := openspec.LoadConfigFrom(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: loading config:", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		project, err = openspec.LoadFromPath(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		model = ui.NewSinglePath(project, cfg, os.Args[1])
	} else {
		project, err = openspec.LoadFrom(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		model = ui.New(project, cfg, cwd)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
