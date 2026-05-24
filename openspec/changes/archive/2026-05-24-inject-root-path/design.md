## Context

`internal/openspec/loader.go` contains 7 functions that call `os.Getwd()` internally. The codebase has no tests for this package, partly because testing requires `os.Chdir` which is incompatible with `t.Parallel()`.

The `*From(path string)` pattern already exists once: `LoadFromPath(path string)` is the only function that accepts an explicit path. This change extends that pattern to all loader functions.

The project has no external consumers of the `openspec` package (it's `internal/`), so breaking signature changes are acceptable. The only callers are `cmd/dossier/main.go` and `internal/ui/index.go`.

## Goals / Non-Goals

**Goals:**
- All openspec loader functions accept an explicit root path, making them testable with `t.TempDir()`
- Existing zero-argument functions remain as convenience wrappers, preserving ergonomics for the CLI entry point
- Silently-swallowed errors are surfaced: functions that currently return only a value now return `(value, error)`
- `LoadFrom(root)` becomes the canonical entry point, replacing the body of `Load()`

**Non-Goals:**
- Switch to `fs.FS` abstraction (that's a separate change, item 2.2 in mejoras.md)
- Refactor `LoadFromPath` (it already accepts an explicit path and doesn't need changes)
- Remove the zero-argument wrappers (they serve as a convenience for `main.go` one-liners)

## Decisions

### Decision 1: `From(root string)` suffix vs `WithPath(path string)` suffix

Chose `From` because it's shorter and already matches the convention set by `LoadFromPath`. Bubble Tea's standard library uses `From` (e.g., `viewport.NewFrom(...)`), reinforcing this pattern.

Alternatives considered: `WithRoot`, `At`. `At` reads poorly (`LoadAt(root)` — "load at root").

### Decision 2: Zero-arg wrappers call `os.Getwd()` AND return errors

Currently, functions like `LoadConfig()` swallow `os.Getwd()` errors and, in some cases (`yaml.Unmarshal`), silently return zero values. The wrappers will now:

```go
func LoadConfig() (ProjectConfig, error) {
    cwd, err := os.Getwd()
    if err != nil {
        return ProjectConfig{}, err
    }
    return LoadConfigFrom(cwd)
}
```

This is a **BREAKING** change to the public API of `LoadConfig`, `LoadProjectSpecs`, `ListArchiveChanges`, `ListArchiveNames`, `ListSpecNames`, `ListChangeNames`. All of them gain an `error` return value.

`Load()` already returns `(*Project, error)` — unchanged.

### Decision 3: `LoadFrom(root)` DOES call `os.Stat(root + "/openspec")` — not `fs.Stat`

The `fs.FS` migration is a separate concern (item 2.2). For now, `*From` variants use the same `os.*` calls they currently use, just relative to `root` instead of `os.Getwd()`. This minimizes the diff and keeps the change focused.

## Risks / Trade-offs

- **Breaking API change**: Callers of `LoadConfig()`, `LoadProjectSpecs()`, etc. must handle the new `error` return. Mitigation: only two files call these functions (`main.go` and `index.go`), and both benefit from explicit error handling.
- **Polling code uses more allocations**: `internal/ui/index.go` calls `ListChangeNames()` and friends on every 500ms tick. Switching to `ListChangeNamesFrom(root)` means `root` must be stored somewhere on the model. Mitigation: store `root string` on `ui.Model` — a single field, negligible memory impact.
