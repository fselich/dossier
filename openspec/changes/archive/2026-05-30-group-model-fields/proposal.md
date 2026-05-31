## Why

The `Model` struct has 34 fields with no logical grouping, making it hard to understand which fields belong to which feature. Grouping related fields into sub-structs improves readability and makes intent clearer for future contributors.

## What Changes

- Define three sub-struct types:
  - `IndexState`: groups `indexItems`, `indexCursor`, `expandedSpecs`, `specSortBySuffix`, `specOrder`, `archiveChanges`, `archiveCursor`
  - `SpecViewerState`: groups `specViewerCursor`, `specJumpTarget`, `specFocusMode`, `specReqCursor`
  - `TaskState`: groups `taskItems`, `taskCursor`
- Replace flat fields in `Model` with embedded (not pointer) sub-struct instances
- Update all field references across the `ui` package
- **BREAKING**: field access paths change (e.g., `m.indexCursor` → `m.index.cursor`). Internal-only; no external API impact.

## Capabilities

### New Capabilities

*(none — pure internal refactor)*

### Modified Capabilities

- `change-index`: Field access paths for index state change; no requirement-level behavior changes
- `mouse-navigation`: Field access paths for index/state fields change; no requirement-level behavior changes

## Non-goals

- Adding new methods to the sub-structs
- Using pointers for sub-structs (value types only)
- Refactoring any logic beyond field access updates
- Changing any behavior visible to the user

## Impact

- `internal/ui/model.go`: new `IndexState`, `SpecViewerState`, `TaskState` types; updated `Model` struct
- All files referencing the affected fields (index.go, update.go, mouse.go, view.go, etc.)
- No new dependencies, no API changes
