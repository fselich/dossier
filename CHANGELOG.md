# Changelog

## v0.10.0

### Added
- Tab bar now shows a distinct color (cyan) for progress bars that reach 100% completion.
- New project info view: press `i` to see `openspec/config.yaml` rendered as markdown.
- Mouse support: click on tabs to switch between them, scroll wheel works on viewports.
- `Tab` / `Shift+Tab` cycle forward and backward through available tabs.
- `--version` / `-v` flag to print the current version.

### Changed
- Progress bar at 100% completion now renders in cyan instead of green.
- Goreleaser releases are now created as drafts with changelog auto-generation disabled.
- Help bar updated to include `Tab` and mouse shortcuts.

### Internal
- Split `internal/ui/model.go` into six focused files (`model.go`, `update.go`, `viewport.go`, `index.go`, `tasks.go`, `view.go`).
