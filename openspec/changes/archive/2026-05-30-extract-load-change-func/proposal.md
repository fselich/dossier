## Why

`LoadFrom`, `LoadFromPath`, and `ListArchiveChangesFrom` each duplicate the same ~10 lines that build a `Change` struct from a directory path and entry name. Extracting this into a single helper removes the duplication and ensures future loading logic changes only need one edit.

## What Changes

- Extract `loadChangeFromPath(dirPath, entryName string, isArchived bool) (Change, error)` helper in `loader.go`
- Refactor `LoadFrom` to call the helper (replaces lines 86–97)
- Refactor `LoadFromPath` to call the helper (replaces lines 324–335)
- Refactor `ListArchiveChangesFrom` to call the helper (replaces lines 206–214)

## Capabilities

### Modified Capabilities
- **openspec-loader**: Extracted `loadChangeFromPath` helper — no behavior change, same Change struct construction consolidated in one place

## Impact

- `internal/openspec/loader.go`: new `loadChangeFromPath` function; `LoadFrom`, `LoadFromPath`, `ListArchiveChangesFrom` refactored
- No API changes, no user-facing impact
