**English** | **[Español](CHANGELOG.es.md)**

# Changelog

## v0.11.0

### Fixed
- Mouse stopped working after returning from the external editor (`e`). Turns out Bubble Tea v1 didn't save mouse state when suspending the terminal. It works now, but it doesn't matter because nobody should be using a mouse anyway.
- App would crash on startup if `archive/`, `specs/`, or `changes/` directories didn't exist. Now it returns empty lists as it should, without making a scene.
- The app background was black instead of the terminal's default color. `NoColor` means "no color," not "black." Who knew.
- `go.mod` had all dependencies marked as indirect. All of them. Including Bubble Tea, which is literally what the app is about.

### Changed
- Full migration to Bubble Tea v2, Bubbles v2, Lip Gloss v2, and Glamour v2. New imports, new declarative API for `View()`, key and mouse messages split into separate types. About 1300 lines touched. Don't ask for whom.
- `renderWithBackground()` and `bgSGRRestore()` removed. Bubble Tea v2 handles the background on its own. One less function to maintain.

### Added
- Unit tests. Yes, finally. ~30 tests across `loader_test.go`, `tasks_test.go`, and `view_test.go`. 74% coverage in `openspec`. UI tests are harder, don't judge me.
- CI via GitHub Actions: `go vet`, `go test -race`, and coverage on every push and PR to `main`. Failures are now caught before merging, not after.

### Internal
- The `openspec` package now accepts an explicit root path in all its functions (`LoadFrom`, `LoadConfigFrom`, etc.) instead of calling `os.Getwd()` internally. More testable, less coupled to global state.
- All loader functions now return `error` instead of silently swallowing failures. Malformed YAML errors are no longer swept under the rug.

## v0.10.0

### Added
- Tab bar now shows a distinct color (cyan) for progress bars that reach 100% completion. This change alone deserved a jump straight to v1.0, I know.
- New project info view: press `i` to see `openspec/config.yaml` rendered as markdown. Still can't edit it. I forgot to add that.
- Mouse support: click on tabs to switch between them, scroll wheel works on viewports. Still, don't use a mouse. It's for cowards.
- `Tab` / `Shift+Tab` cycle forward and backward through available tabs. Welcome to the world of keybinding incompatibilities between the app and the window system.
- `--version` / `-v` flag to print the current version. The AI did this on its own, without being asked.

### Changed
- Progress bar at 100% completion now renders in cyan instead of green. Cyan is like light blue, in case I forget.
- Goreleaser releases are now fully automated (no more drafts). Boring.
- Help bar updated to include `Tab` and mouse shortcuts.

### Internal
- Split `internal/ui/model.go` into six focused files (`model.go`, `update.go`, `viewport.go`, `index.go`, `tasks.go`, `view.go`). Super boring.
