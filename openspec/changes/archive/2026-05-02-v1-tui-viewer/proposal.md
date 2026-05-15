## Why

Browsing and editing OpenSpec artifacts (proposal, design, tasks, specs) requires opening multiple files in a text editor, with no overview or contextual interaction. A dedicated TUI eliminates that friction and turns daily change management into something fluid and fast.

## What Changes

- New CLI tool (`spec-viewer`) written in Go
- Reads `./openspec/` from the current working directory
- Shows active changes (not archive) with keyboard navigation
- Renders artifacts in markdown with ANSI formatting
- Allows marking tasks as completed directly from the TUI, without opening any editor

## Capabilities

### New Capabilities

- `openspec-loader`: Discovers and parses the `./openspec/` structure — active changes, artifacts and metadata
- `tui-viewer`: Bubble Tea TUI that shows active changes, navigates between artifacts (proposal, design, tasks, specs) and renders them with glamour
- `tasks-toggle`: Parses `tasks.md`, exposes checkboxes as interactive items and writes the `[ ]` ↔ `[x]` toggle directly to the file

### Modified Capabilities

## Impact

- New standalone Go tool (`spec-viewer`), installable with `go install`
- Dependencies: `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, `github.com/charmbracelet/glamour`, `github.com/charmbracelet/bubbles`, `gopkg.in/yaml.v3`
- Does not modify any OpenSpec file except `tasks.md` when toggling
