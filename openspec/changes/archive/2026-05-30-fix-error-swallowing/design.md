## Context

`LoadConfigFrom` reads `openspec/config.yaml`. The current implementation treats all `os.ReadFile` errors the same — returning nil error — which means a permission-denied or corrupt-filesystem error is silently swallowed. The spec (`openspec-loader`) already requires: "If the config file does not exist, return empty config and nil error. If the YAML is malformed, return an error." The non-IsNotExist case is missing error propagation.

Separately, `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` are called with errors discarded via `_` in the UI layer (`index.go`, `model.go`). These functions already return proper errors; the callers just ignore them.

## Goals / Non-Goals

**Goals:**
- `LoadConfigFrom` returns the real error on non-IsNotExist `os.ReadFile` failures.
- UI callers log or surface errors from `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` instead of discarding them.

**Non-Goals:**
- Adding retry logic or fallback paths for transient errors.
- Changing the signature of any loader function.

## Decisions

**Return the raw error from `LoadConfigFrom`**
The fix is a one-word change: `return ProjectConfig{}, nil` → `return ProjectConfig{}, err` on line 126. This is the minimum correct fix and aligns with the existing spec requirement.
_Alternative_: wrap the error with context — unnecessary since the caller can wrap if needed.

**Use `m.logger.Printf` or `log.Printf` for UI error handling**
Since the Model doesn't currently have a logger field, the simplest approach is to log the error to stderr via `log.Printf` (or `fmt.Fprintf(os.Stderr, ...)`) so failures are visible in the terminal. The archive/specs lists will be nil on error, which is the existing fallback behavior.
_Alternative_: store error in `m.errMsg` — this would show it in the TUI, but the error is non-critical (UI still works without archive/spec data), so logging is less disruptive.

## Risks / Trade-offs

- **`LoadConfigFrom` callers may not handle the new errors** → Mitigation: the `err != nil` check in `main.go:41` already exists and will correctly surface the error to the user.
- **Logging to stderr may clutter the TUI** → Mitigation: errors from background data loading (archive, specs) are rare; a brief log line is acceptable.
