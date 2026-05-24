**English** | **[Español](CHANGELOG.es.md)**

# Changelog

## v0.10.0

### Added
- Tab bar now shows a distinct color (cyan) for progress bars that reach 100% completion. This change alone deserved a jump straight to v1.0, I know.
- New project info view: press `i` to see `openspec/config.yaml` rendered as markdown. Still can't edit it. I forgot to add that.
- Mouse support: click on tabs to switch between them, scroll wheel works on viewports. Still, don't use a mouse. It's for cowards.
- `Tab` / `Shift+Tab` cycle forward and backward through available tabs. Welcome to the world of keybinding incompatibilities between the app and the window system.
- `--version` / `-v` flag to print the current version. The AI did this on its own, without being asked.

### Changed
- Progress bar at 100% completion now renders in cyan instead of green. Cyan is like light blue, in case I forget.
- Goreleaser releases are now created as drafts with changelog auto-generation disabled. Boring.
- Help bar updated to include `Tab` and mouse shortcuts.

### Internal
- Split `internal/ui/model.go` into six focused files (`model.go`, `update.go`, `viewport.go`, `index.go`, `tasks.go`, `view.go`). Super boring.
