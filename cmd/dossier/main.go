package main

import (
	"flag"
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

	themeName := flag.String("theme", "dark", "Visual theme (dark, light, dracula)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	showHelp := flag.Bool("help", false, "Print help and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dossier [--theme <name>] [--version] [--help] [path]\n")
		fmt.Fprintf(os.Stderr, "\nA keyboard-driven TUI for navigating OpenSpec project artifacts.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("dossier", version)
		os.Exit(0)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	theme, ok := ui.LookupTheme(*themeName)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unknown theme %q. Available: dark, none, light, dracula\n", *themeName)
		os.Exit(1)
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

	loader := openspec.NewLoader(openspec.OSFS{})

	pathArg := flag.Arg(0)
	if pathArg != "" {
		project, err = openspec.LoadFromPath(pathArg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		model = ui.NewSinglePath(project, cfg, pathArg, loader, theme)
	} else {
		project, err = openspec.LoadFrom(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		model = ui.New(project, cfg, cwd, loader, theme)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
