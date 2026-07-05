# AGENTS.md — Dossier

A keyboard-driven TUI for navigating [OpenSpec](https://openspec.dev) project artifacts, built with [Bubble Tea v2](https://charm.land/bubbletea/v2).

## Essential Commands

```bash
make build              # Build binary to ./dossier
make install            # go install
make test               # go test -race -cover ./...
make lint               # golangci-lint run ./...
make fmt                # goimports -w .
make clean              # Remove binary
# CI also runs: go vet ./... && go tool cover -func=coverage.out
```

## Project Structure

```
cmd/dossier/main.go         # Entry point: flags (--version, --help, path arg), tea.NewProgram
internal/
  openspec/                 # Domain types + filesystem logic
    fs.go                   # fileSystem interface + OSFS implementation
    loader.go               # Loader: reads changes, specs, config from disk
    tasks.go                # Task parsing (markdown checkboxes), toggling
    tasks_test.go
    loader_test.go
  git/                      # Git porcelain output parsing (zero deps)
    git.go                  # IsInsideWorkTree, WorkTreeRoot, Status
  ui/                       # Bubble Tea model, views, handlers
    model.go                # Model struct, Mode/Tab enums, New(), Init(), View()
    update.go               # Update loop: dispatch by msg type and mode
    view.go                 # Rendering: headers, tab bar, help bar, borders
    index.go                # Index mode: building/rendering/navigating index items + tick polling
    viewer.go               # Key dispatch for ModeNormal / ModeViewingArchive
    viewport.go             # Viewport loading: glamour rendering, caching, async msgs
    git.go                  # Git tab: rendering, cursor, polling, editor open
    tasks.go                # Tasks tab: cursor navigation, toggle, progress bars
    config.go               # Config mode key handler (minimal: q/esc back, j/k scroll)
    spec.go                 # Spec viewer mode key handler (Esc index restore, h/l req nav)
    mouse.go                # Mouse wheel + click handlers
    styles.go               # All lipgloss style definitions (19 styles)
    view_test.go            # Tests for rendering, key dispatch, index, tasks, tick
openspec/                   # OpenSpec project artifacts (not Go code)
  config.yaml
  specs/                    # 22 project specification files
  changes/                  # Active changes (dirs)
    archive/                # 40+ archived changes (dirs)
```

## Architecture

- **Model** is the single Bubble Tea model with value-driven updates (no pointers). `Update` returns a new `Model`.
- **Modes** control what the UI shows: `ModeNormal` (change viewer), `ModeIndex` (nav index), `ModeViewingArchive`, `ModeViewingSpec`, `ModeViewingConfig`.
- **Tabs** (`TabProposal=0`, `TabDesign=1`, `TabSpecs=2`, `TabTasks=3`, `TabGit=4`, `tabCount=5`) switch content within a change. `TabGit` only appears when inside a git worktree and in `ModeNormal`.
- **Filesystem** is abstracted behind `fileSystem` interface in `openspec/fs.go` — `OSFS` wraps real OS calls, enabling testability.
- **Glamour** renders Markdown to ANSI async via `renderedMsg`/`specRenderedMsg`/`renderedConfigMsg` messages. Renderer is cached by width (`ensureRenderer(width)`).
- **Tick** at 500ms polls disk for OpenSpec changes and git status. Stops for `ModeViewingArchive`/`ModeViewingSpec`.
- **Index mode** displays active changes, specs (expandable to requirements), and archived changes. Supports filtering via `/` key and sort-by-suffix via `s`.

## Key Gotchas & Non-Obvious Patterns

1. **Go 1.25.x** — latest Go. No generics. No `context.Context`. No `errors` package (plain `fmt.Errorf`).
2. **Value vs pointer receivers**: `updateConfig`, `updateViewer`, `updateSpec`, `dispatchKey` all take `Model` (value). Mutating methods (`buildIndexItems`, `refreshIndexViewport`, `loadTaskItems`, `pollGitStatus`, `moveGitCursor*`) use `*Model`. `Update` returns a new `Model`.
3. **Two-tier openspec API**: Every `Loader` method has a package-level wrapper (e.g., `loader.LoadFrom` → `openspec.LoadFromFrom`). Zero-argument forms (`Load()`, `LoadConfig()`) exist but are unused internally — they call `os.Getwd()`.
4. **`artifactPath()` uses direct `os.ReadDir/Stat`** (not `fileSystem` interface) — a testability gap in specs tab path resolution.
5. **Task list cursor synced by text**: `FindCursorByText` restores cursor position after reload by matching task text.
6. **Only lowercase `[x]` is recognized as done** — `[X]` (uppercase) is NOT matched by `rxDone`.
7. **`.openspec.yaml` parse errors silently ignored** (optional metadata).
8. **`renderCache` cleared on three events**: change switch, window resize, mode switch. `editorReturnMsg` deletes only the current tab's cache.
9. **`commitStateChange()`** adjusts viewport height and calls `loadViewport()` — used after every mode/tab/change change.
10. **Git porcelain parsing**: `XY path` format — `X`/`Y` are index/worktree status, separator at `[2]`, path starts at `[3:]`. Renames/copies split on ` -> `. Files under `openspec/` are filtered out. Never use `strings.TrimSpace` on the raw output (it strips leading whitespace, corrupting first-line XY codes) — use `strings.TrimRight(raw, "\n\r")` instead.
11. **Git cursor skips deleted files**: `moveGitCursorDown/Up` wraps via modulo and skips `IsDeleted`. `clampGitCursor` scans forward then backward.
12. **Git status poll is always-on**: `pollGitStatus()` runs every tick (guarded by `isGitRepo`), refreshes viewport only when `TabGit` is active.
13. **Git tab label is dynamic**: `changes` when clean, `changes (N)` when files exist.
14. **Archived change names**: `YYYY-MM-DD-name` format. `parseArchiveName` extracts first 10 chars as date (`DD/MM/YYYY`), rest as name. Non-matching names use full dir name as name with no date.
15. **Mouse mode** uses `tea.MouseModeCellMotion`. Header click (Y=1) goes to index. Tab bar click (Y=2) switches tabs.
16. **Index mouse click is two-phase**: first click selects (moves cursor), second click on same item actions it.
17. **Index has two-phase Esc**: during filter editing → revert to `PrevFilterText`. after filter editing with filter set → clear filter. after with no filter → quit.
18. **`pollIndexMode` compares by name then reloads**: uses `sameNames` (set) for changes, `sameStrings` (ordered) for specs/archives. Reloads task content even when names match (keeps progress bars current).
19. **`singlePath` disables change-list polling**: when launched with a path arg, `pollNormalModeChanges` is skipped.
20. **`indexItemAtContentLine` mirrors `renderIndexContent`**: must stay in sync — any layout change in rendering requires identical change in click handler.
21. **`specRenderedMsg.jumpLine`** only set in full-spec mode (not focus mode). Jump target is found by stripping ANSI codes from glamour output and substring-matching.
22. **Config mode returns to `m.prevMode`** on Esc. Spec viewer always returns to index (restoring focus state if applicable).
23. **Task inline markdown uses `extractOpeningEscape`**: renders a marker char with lipgloss, extracts the ANSI prefix, uses it as a "restore" to preserve outer style around inline spans. Fragile if lipgloss internals change.
24. **Editor** defaults to `vi` if `$EDITOR` is unset. Launch async via `tea.ExecProcess`. Git tab uses `m.gitRoot` for absolute paths.
25. **`renderWidth()` minimum is 80**: when `m.width-2 < 20`, glamour renders at 80 columns.
26. **Specs structure**: inside a change, `specs/<name>/spec.md`. Project-level: `openspec/specs/<name>/spec.md`.

## Testing Patterns

- Tests use `t.TempDir()` + `os.MkdirAll` — no mock filesystem.
- Table-driven tests used sparingly; function-per-test-case is more common.
- `testLoader()` returns `openspec.NewLoader(openspec.OSFS{})`.
- Model construction in tests is explicit field-by-field (no constructor helpers).
- Tests in `internal/ui/` and `internal/openspec/` as `_test.go` in the same package.

## Linting & Formatting

- `.golangci.yml`: `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, `intrange`, `copyloopvar`.
- Formatters: `gofmt` + `goimports`.
- `errcheck` has `check-type-assertions: true`.
