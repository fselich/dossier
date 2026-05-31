**English** | **[Español](CHANGELOG.es.md)**

# Changelog

## v0.14.1

### Fixed
- Done-task code spans in the task list no longer show the first letter in a different color. Lipgloss renders underlined text character by character, resetting the foreground between them. The fix combines underline with the foreground color so each character inherits both.

## v0.14.0

### Added
- Press `/` in the index view to filter changes, specs, and archived items by name in real-time. Type to narrow down, `Enter` to lock the filter, `Esc` to clear it. A search box, basically.

## v0.13.0

### Internal
- Split the monolithic `handleKeyPress()` into per-mode update functions, each in its own file: `viewer.go`, `index.go`, `spec.go`, `config.go`. `update.go` is now a thin dispatcher.
- Introduced a `fileSystem` interface and `Loader` struct in `openspec`, so the package no longer depends on `os` directly. All public functions preserved via backward-compatible wrappers.
- Added `.golangci.yml` with errcheck, staticcheck, govet, unused, gofmt, goimports, and a `Makefile` with `test`, `lint`, and `fmt` targets.
- Eliminated all silent `log.Printf` error calls. Archive and spec load errors are now displayed in the help bar for 3 seconds via `m.errMsg`, exactly like toggle errors.

### Changed
- Tab bar `parts` slice is now preallocated to exactly 4 entries, and the tasks `items` slice preallocates to the line count. Everything is now 3 nanoseconds faster. Totally worth the token spend.
- The `taskCounts` function no longer uses naked returns (which were confusing to anyone who scrolled past line 491 of index.go).
- Layout constants (`chromeTop`, `chromeHeader`, etc.) replace magic numbers in `contentHeight()`. Now you know why it was subtracting 6.
- The reload-merge logic that was copy-pasted in two places is now a single `mergeReloadedChange()` method. DRY*2.

## v0.12.0

### Fixed
- Starting dossier with no pending changes now shows the index view with specs and archived changes instead of a blank screen.
- Task content updates inside existing changes now trigger a live refresh of the index list instead of silently ignoring them.
- The loading placeholder (`"Loading..."` / `"Cargando..."`) was removed. Raw markdown is shown immediately while the styled version renders in the background. Goodbye to the involuntary epilepsy mode.

### Changed
- Change list in the index view is now sorted by `created` date (descending). Before, they were sorted by whatever the filesystem felt like.

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
