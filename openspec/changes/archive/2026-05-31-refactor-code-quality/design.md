## Context

Dossier is a terminal UI for navigating OpenSpec project artifacts. The current codebase is clean but has accumulated technical debt in three areas:

1. **Monolithic `handleKeyPress`** (`internal/ui/update.go:112-381`): a 270-line switch handling keybindings for all five UI modes interleaved. Adding a new mode requires editing this function and fighting with the existing nesting.

2. **No testability abstraction**: the `openspec` package calls `os.ReadFile`/`os.WriteFile` directly, forcing all tests to use `t.TempDir()` with real disk I/O. No in-memory testing possible.

3. **Missing project tooling**: no `.golangci.yml`, `Makefile` missing `test`/`lint`/`fmt` targets.

These issues surfaced during a structured code review against Go code style and pattern guidelines.

## Goals / Non-Goals

**Goals:**

- Split key handling by mode, each mode in its own file, with `update.go` becoming a thin dispatcher
- Define a `fileSystem` interface for testability, convert loader functions to methods on a struct
- Add `.golangci.yml` and complete `Makefile` targets
- Fix code style issues: naked return, missing preallocation, undocumented ignored error, magic layout numbers
- DRY the duplicated reload-merge logic

**Non-Goals:**

- Changing user-facing TUI behavior or keybindings
- Replacing Bubble Tea with another TUI framework
- Adding new features (this is pure refactor + tooling)
- Moving to a sub-model pattern where each mode implements `tea.Model` (would require significant restructuring of `View()`); we stay with `update*()` methods on the single `Model`
- Extracting `openspec` to a separate module

## Decisions

### Decision 1: Per-mode files with `update*()` methods vs. sub-models implementing `tea.Model`

**Chosen:** Per-mode update methods on the same `Model` struct.

**Rationale:** Sub-models (like Glow's `pagerModel`/`stashModel`) would require each mode to implement its own `View()`. Dossier's viewport is shared across all modes; only the content and keybindings differ. Extracting per-mode `update*()` methods preserves the shared viewport and `View()` logic while eliminating the monolithic switch. This follows the Bubble Tea "composable-views" pattern where sub-functions return `(tea.Model, tea.Cmd)`.

**Alternative considered:** Full sub-models with `tea.Model` per mode. Rejected because it would duplicate viewport/chrome rendering across modes without clear benefit.

### Decision 2: File structure for per-mode handlers

**Chosen:**
```
internal/ui/
├── update.go      → dispatcher (thin switch on mode)
├── viewer.go      → ModeNormal + ModeViewingArchive (shared artifact viewer logic)
├── index.go       → ModeIndex (already exists; gets updateIndex method)
├── spec.go        → ModeViewingSpec (new file)
├── config.go      → ModeViewingConfig (new file)
```

`ModeNormal` and `ModeViewingArchive` share a single `viewer.go` because they have identical keybindings (1-4/Tab for artifacts, j/k for scroll, e for edit). The only difference is which `Change` pointer they operate on (`current()` vs `currentArchive()`), handled internally.

### Decision 3: Keymap pattern with `key.Binding`

**Chosen:** Each mode's file defines a `keymap` struct using `key.Binding` with descriptive names instead of raw string comparisons.

```go
// viewer.go
var viewerKeys = struct {
    proposal, design, specs, tasks key.Binding
    nextChange, prevChange         key.Binding
    edit, quit, index              key.Binding
    scrollDown, scrollUp           key.Binding
    toggle                         key.Binding
}{...}
```

**Rationale:** `key.Binding` enables:
- Automatic help text generation (replacing the manual `renderHelpBar()` strings)
- Disabling bindings via `.SetEnabled()` (e.g., hide "edit" when no artifact is available)
- Future key customization without code changes
- Single source of truth for key → action mapping

Existing `renderHelpBar()` is replaced by `m.help.ShortHelpView()` fed from the active mode's key bindings.

### Decision 4: FileSystem interface location

**Chosen:** Define `fileSystem` interface in `internal/ui/` (consumer side), convert `openspec` functions to methods on a `*Loader` struct.

```go
// internal/ui/fs.go
type fileSystem interface {
    ReadFile(name string) ([]byte, error)
    WriteFile(name string, data []byte, perm os.FileMode) error
    ReadDir(name string) ([]os.DirEntry, error)
    Stat(name string) (os.FileInfo, error)
}

// internal/openspec/loader.go
type Loader struct {
    fs fileSystem
}
```

The `fileSystem` interface is **not** exported from `openspec` — it's defined in `ui` where it's consumed. The `openspec` package accepts it via constructor injection. Backward compatibility is preserved with package-level wrapper functions that use a default `osFS` implementation.

### Decision 5: Layout constant extraction

**Chosen:** Named constants replacing magic numbers in `contentHeight()`.

```go
const (
    chromeTop     = 1
    chromeHeader  = 1
    chromeInnerSep = 1
    chromeTabBar  = 1
    chromeHelpBar = 1
    chromeBottom  = 1
    chromeSpecSubnav = 1
)
```

### Decision 6: Error surfacing for polling failures

**Chosen:** Store errors in `m.errMsg` (existing mechanism) instead of `log.Printf`.

The error auto-clears after 3 seconds via the existing `errClearMsg` tick. This is consistent with how toggle errors already work.

## Risks / Trade-offs

- **Diff size** → The refactor touches most UI files. Mitigation: each commit is atomic per concern (per-mode split, interface, tooling, style fixes).
- **key.Binding integration** → Requires adding the `help` bubble dependency. Mitigation: Bubble's `help` package is already a transitive dependency via other bubbles.
- **Loader struct migration** → `cmd/dossier/main.go` needs wiring changes. Mitigation: wrapper functions preserve the existing API surface; main.go only adds dependency injection.
- **Test breakage** → Tests that construct `Model` directly may need field adjustments. Mitigation: run `go test -race ./...` after each commit to catch regressions immediately.
