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
cmd/dossier/main.go         # Entry point: flags (-h/--help, --read-only, --version, path arg), tea.NewProgram
internal/
  openspec/                 # Domain types + filesystem logic
    fs.go                   # fileSystem interface + OSFS implementation
    loader.go               # Loader: reads changes, specs, config from disk
    tasks.go                # Task parsing (markdown checkboxes), toggling
    tasks_test.go
    loader_test.go
  git/                      # Git porcelain output parsing (zero deps)
    git.go                  # IsInsideWorkTree, WorkTreeRoot, Status
  settings/                 # XDG TOML config, embedded themes, keymap defaults/validation
  ui/                       # Bubble Tea model, views, handlers
    model.go                # Model struct, Mode enum, constructors, Init(), View()
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
- **Artifact destinations** use `{artifactID, outputPath}` identities from each change's dynamic `Artifacts` list. The synthetic code destination is active-only, appears when inside a Git worktree with changes, and is tracked by `viewingCode` rather than a schema artifact.
- **Filesystem** is abstracted behind `fileSystem` interface in `openspec/fs.go` — `OSFS` wraps real OS calls, enabling testability.
- **Glamour** renders Markdown to ANSI async via identity-keyed `artifactRenderedMsg`, `specRenderedMsg`, and `renderedConfigMsg` messages. Renderer is cached by width (`ensureRenderer(width)`).
- **Tick** at 500ms polls disk for OpenSpec changes and git status. Stops for `ModeViewingArchive`/`ModeViewingSpec`.
- **Index mode** displays active changes, specs (expandable to requirements), and archived changes. Supports filtering via `/` key and sort-by-suffix via `s`.

## Key Gotchas & Non-Obvious Patterns

1. **Go 1.25.x** — latest Go. No generics. `context.Context` used only locally for subprocess timeouts (never plumbed through APIs). No `errors` package (plain `fmt.Errorf`).
2. **Value vs pointer receivers**: `updateConfig`, `updateViewer`, `updateSpec`, `dispatchKey` all take `Model` (value). Mutating methods (`buildIndexItems`, `refreshIndexViewport`, `loadTaskItems`, `pollGitStatus`, `moveGitCursor*`) use `*Model`. `Update` returns a new `Model`.
3. **Two-tier openspec API**: Every `Loader` method has a package-level wrapper (e.g., `loader.LoadFrom` → `openspec.LoadFrom`). Zero-argument forms (`Load()`, `LoadConfig()`) exist but are unused internally — they call `os.Getwd()`.
4. **Dynamic output paths are containment-checked**: editor and task mutations resolve the selected output beneath the change root and reject missing files or symlink escapes.
5. **Task list cursor synced by text**: `FindCursorByText` restores cursor position after reload by matching task text.
6. **Only lowercase `[x]` is recognized as done** — `[X]` (uppercase) is NOT matched by `rxDone`.
7. **`.openspec.yaml` is optional but diagnostic-aware**: missing metadata degrades safely; malformed metadata is surfaced on the change without hiding discovered outputs.
8. **Artifact render caching is identity-based**: `artifactRenderCache` keys include change name, archive state, artifact ID, and output path; reload/editor return invalidates only the affected change.
9. **`commitStateChange()`** adjusts viewport height and calls `loadViewport()` — used after every mode/artifact/change transition.
10. **Git porcelain parsing**: uses `git status --porcelain=v1 -z -u`. Output is NUL-separated: each entry is `XY <path>\0` (3-byte header + path). Renames/copies: `R <old_path>\0<new_path>\0` — the entry token has the old path, the next NUL token has the new path (`Path` = new, `OldPath` = old). All git subprocess calls go through the exported `RunGit(dir, args...)` helper in `internal/git`, which enforces a 2s timeout via `context.WithTimeout` + `exec.CommandContext`. Files under `openspec/` are filtered out.
11. **Git cursor no longer skips deleted files**: `moveGitCursorDown/Up` wraps via simple modulo (no `IsDeleted` skip). `clampGitCursor` is a simple clamp. Diff cycling (`[`/`]`) uses `moveGitDiffCursorDown/Up` which still skips deleted files. `s` key toggles stage/unstage on the file under cursor.
12. **Git status poll is always-on**: `pollGitStatus()` runs every tick (guarded by `isGitRepo`), refreshing the viewport only when `viewingCode` is active.
13. **The code destination label is dynamic**: `code` when clean/disabled and `code (N)` when files exist.
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
24. **Editor** defaults to `vi` if `$EDITOR` is unset. Launch async via `tea.ExecProcess` on the exact safe selected artifact output; Git uses `m.gitRoot` for absolute paths.
25. **`renderWidth()` minimum is 80**: when `m.width-2 < 20`, glamour renders at 80 columns.
26. **Specs structure**: change-local `specs/<name>/spec.md` outputs are delta specs under the dynamic `specs` artifact; project-level `openspec/specs/<name>/spec.md` files are canonical specs.
27. **Read-only is explicit launch state**: `--read-only` is passed into `ui.New`/`NewSinglePath` and blocks change creation, lifecycle/undo, task toggles, external editor launch, and Git stage/unstage. Inspect, validate, Git status/diff, and navigation remain active; disabled mutation hints are omitted from the help bar.
28. **User settings are strict XDG TOML**: `internal/settings.Load` reads `${XDG_CONFIG_HOME:-os.UserConfigDir()}/dossier/config.toml`; a missing default is allowed, but explicit/malformed/unknown settings fail startup. Built-in palettes (`none`, `dark`, `light`, `dracula`, `gruvbox-light-soft`) are embedded from `internal/settings/themes/*.toml`.
29. **Keybindings are context-specific actions with a named base**: `nvim` is the default keystyle; `[keys].style` and `--keystyle` select the base while sparse action overrides remain applied. Handlers match `settings.KeyConfig` via `matchesKey`; zero-value test models use `effectiveKeyMap`. The nvim profile uses `q` back, `Q` quit, `h`/`l` for visible artifacts, `H`/`L` horizontal, Tab/Shift+Tab for changes, and page keys plus `Ctrl+U`/`Ctrl+D`. Help labels must use the key-label helpers rather than hardcoded key names.
30. **Runtime themes merge sparse overrides onto a base**: `ui.BuildTheme` constructs Lip Gloss styles, Glamour `ansi.StyleConfig`, and a Chroma style. Do not reintroduce fallback color literals in rendering paths.
31. **Enriched state is the persistent truth; discovery fingerprints are discovery-shaped**: `Model.lastStatusByName` caches each change's last `openspec status`; every disk reload (`adoptDiscoveredChange`, `adoptDiscoveredProject`, `reloadIndex`, `mergeReloadedChange`) re-applies it synchronously so schema-declared artifact IDs never flap back to filename stems, and `applyStatusEnrichment` skips rebuilds when the enriched state fingerprint is unchanged. `discoveryFingerprints` must always be computed from pre-enrichment (discovered) changes — hashing enriched changes would make every later adoption mismatch and rebuild forever; reload detection likewise compares stored discovery fingerprints (`mergeReloadedChange`), never the enriched change itself. `adoptDiscoveredChange` is the single writer for change reloads and reconciles the artifact selection with `reconcileArtifactSelectionNear` (position-preserving). `rebuildIndexPreservingCursor` and `enterIndex` re-anchor the cursor by (change identity, output path) when the exact row identity disappears (discovery↔schema ID swaps). `newModel` runs `bootstrapEnrichment` — a 750 ms bounded synchronous status pass — so the first paint is schema-aware; misses degrade to async enrichment. The pass prefers `CLI.StatusAll` (`status --all --json`, one spawn for the whole repo) via the optional `batchStatusClient` interface assertion, falling back to bounded per-change `Status` calls.
32. **Labels show what is read; identities stay schema-declared**: index artifact rows and viewer tabs use `artifactDisplayLabel` — the single present output's display name (file stem, e.g. `1-idea`) when an artifact resolves to exactly one file, else the schema artifact ID (multi-output like `specs`, or unwritten). Row identity strings and selection keep using `artifact.ID`; only the rendered label differs, and index filtering matches both.

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
