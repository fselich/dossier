## Context

`LoadFrom`, `LoadFromPath`, and `ListArchiveChangesFrom` in `loader.go` each contain ~10 identical lines that:
1. Resolve a directory path from the entry name
2. Read `.openspec.yaml` for the `created` date
3. Scan artifact subdirectories
4. Build and return a `Change` struct

## Goals / Non-Goals

**Goals:**
- Extract the duplicated block into `loadChangeFromPath(dirPath, entryName string, isArchived bool) (Change, error)`
- All three callers use the same helper

**Non-Goals:**
- Changing the `Change` struct or artifact loading logic
- Adding new functionality

## Decisions

**Decision: `isArchived bool` parameter**
- Archive changes need this set to `true` on the returned `Change`
- Active changes pass `false`
- Cleaner than passing a separate `archived` boolean after construction

**Decision: Function remains unexported (lowercase)**
- It's purely an internal helper — no external package needs it
- Matches existing convention of `loadArtifacts`, `loadProjectSpec` etc.

## Risks / Trade-offs

- [Risk: accidental behavior difference] → Mitigation: existing tests must pass unchanged; the refactor is purely mechanical
