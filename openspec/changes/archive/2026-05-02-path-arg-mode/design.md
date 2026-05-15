## Context

`main.go` calls `openspec.Load()` without arguments, which scans `openspec/changes/` in the working directory. There is no mechanism to point to a specific change by path. The loader only knows the concept of "full project".

## Goals / Non-Goals

**Goals:**
- `./spec-viewer <path>` loads a single change from `<path>` and shows it in the TUI
- If `<path>` does not exist or does not contain a valid change, the binary exits with a descriptive error before opening the TUI
- Without argument: identical current behavior

**Non-Goals:**
- Support for multiple paths as arguments (`./spec-viewer path1 path2`)
- Flags (no `flag.Parse`, only `os.Args[1]`)
- Changing the polling logic in path mode (the tick keeps working — `ReloadChange` uses `ch.Path` directly, which already points to the correct path)

## Decisions

**D1: `openspec.LoadFromPath(path string)` as a separate function**

Instead of modifying `Load()` to accept an optional argument, `LoadFromPath(path string) (*Project, error)` is added. `main.go` decides which one to call based on `len(os.Args)`. This separation keeps `Load()` intact and makes the new mode explicit and independently testable.

Discarded alternative: passing `path string` to `Load()` with `""` as "normal mode". This pollutes the signature without real benefit and breaks compatibility if there are other callers.

**D2: Minimal validation — `.openspec.yaml` must exist**

`LoadFromPath` verifies that the directory exists and that it contains `.openspec.yaml`. If not, it returns an error. The yaml content and the presence of artifacts are not validated (that is the TUI's responsibility, same as in normal mode).

**D3: `Project.Name` in path mode is taken from the name of the change's parent directory**

`Load()` uses `filepath.Base(cwd)` as the project name. In path mode, the change could be in an archive of another project. `filepath.Base(filepath.Dir(path))` is used to infer the project name from the path — if the path is `.../my-project/openspec/changes/archive/2026-05-02-feat`, the "project" would be `archive`. As a more useful alternative, the change directory name is used directly as the sole visible identifier, and `Project.Name` is set to `""` or to the basename of the change's parent directory. Since the header shows `<project> · <change-name>`, the basename of the directory containing the change (e.g., `archive`) is used as the project name — honest and clear.

## Risks / Trade-offs

[Risk: Relative vs absolute path] → `main.go` passes the argument as-is to `LoadFromPath`; if it is relative, `os.Stat` resolves it correctly from the process cwd. No risk.

[Risk: Polling in path mode tries `ListChangeNames()` which scans the cwd] → The tick handler calls `ListChangeNames()` to detect new changes. In path mode, the cwd may not have `openspec/changes/`. The result will be an empty list or an ignored error — `sameNames` will return false and the tick will attempt `Load()` (not `LoadFromPath`), loading the normal project if it exists. This is undesirable. Mitigation: in path mode, the tick must not try to reload the list of changes. A `m.singlePath string` flag is introduced in the Model; if set, the tick skips the `ListChangeNames` check and only does `ReloadChange`.

## Migration Plan

No migration. The binary continues working the same without an argument.

## Open Questions

None.
