## Why

Dossier is a tool for navigating OpenSpec project artifacts, but developers also need quick context on what files have changed in the repo. Switching to a terminal to run `git status` breaks the flow. A built-in git status tab keeps everything in one TUI.

## What Changes

- New `TabGit` (key `5`, label `changes`) shows a selectable list of modified/added/untracked/renamed/deleted files from `git status --porcelain`, excluding paths under `openspec/` (already tracked by other tabs).
- The tab only appears when inside a git worktree (detected via `git rev-parse --is-inside-work-tree`).
- `j/k` navigates cursor; `Enter` or `e` opens the selected file in `$EDITOR`.
- Deleted files are shown but the cursor skips them.
- Renamed files show `old → new`; `Enter` opens the new path.
- File count shown in tab bar (e.g. `changes (5)`).
- Git status is polled on the existing 500ms tick.

## Capabilities

### New Capabilities

- `git-status-tab`: A fifth tab in the change viewer that displays working-tree file changes from `git status --porcelain`, with cursor navigation and editor integration.

## Impact

- New package `internal/git/` with `IsInsideWorkTree()` and `Status()`.
- New file `internal/ui/git.go` with rendering and viewport logic.
- Modified files: `model.go` (Tab enum, state, tabAvailable), `view.go` (tab bar, help bar), `viewer.go` (key dispatch), `viewport.go` (load dispatch), `mouse.go` (click/scroll), `styles.go` (status colors).
- No changes to OpenSpec domain logic.
- No new dependencies.
- `go.mod` unchanged.
