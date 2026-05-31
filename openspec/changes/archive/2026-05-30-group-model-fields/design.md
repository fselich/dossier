## Context

`Model` in `internal/ui/model.go` is a flat struct with 34 fields belonging to different features (index, spec viewer, tasks, archive, config). Finding which fields relate to which feature requires scanning the entire struct. The fields are private (lowercase), so grouping has no external impact.

## Goals / Non-Goals

**Goals:**
- Define three sub-struct types: `IndexState`, `SpecViewerState`, `TaskState`
- Embed instances in Model (not pointers) — zero values work without initialization
- Update all field references from `m.field` to `m.substruct.field` across the ui package
- Tests continue to pass without modification

**Non-Goals:**
- Grouping the remaining fields (project, changeIdx, tab, vp, specIdx, mode, etc.)
- Adding methods to the new sub-structs
- Using pointers (value types only, for struct embedding)

## Decisions

### Decision: Value-type embedding, not pointers

```go
type Model struct {
    // ...
    index        IndexState
    specViewer   SpecViewerState
    tasks        TaskState
    // ...
}
```

Using value types means no nil checks, no `NewIndexState()` constructors, and the zero value is valid. All current usage already accesses fields on `m Model` or `*Model` directly — this just adds a `.` in the path.

**Alternatives considered:**
- Pointer fields — rejected; adds nil-pointer risk and requires initialization in `New()`.
- Named wrapper methods — rejected; adds unnecessary indirection for what is straightforward field access.

### Decision: Field-to-substruct mapping

| Sub-struct | Fields |
|---|---|
| `IndexState` | `IndexItems []indexItem`, `IndexCursor int`, `ExpandedSpecs map[int]bool`, `SpecSortBySuffix bool`, `SpecOrder []int`, `ArchiveChanges []openspec.Change`, `ArchiveCursor int` |
| `SpecViewerState` | `SpecViewerCursor int`, `SpecJumpTarget string`, `SpecFocusMode bool`, `SpecReqCursor int` |
| `TaskState` | `TaskItems []openspec.TaskItem`, `TaskCursor int` |

Fields are exported (uppercase) within the sub-structs since they're aggregate types. The sub-struct types are unexported (lowercase) since they're implementation details.

### Decision: Sub-struct naming

`IndexState`, `SpecViewerState`, `TaskState` — suffixes with `State` to distinguish from other concepts already in the package (e.g., `indexItem` is a different type).

## Risks / Trade-offs

- **Find-replace mistakes**: A mechanical rename could miss dynamic field access or string-based references. Mitigation: compile-check after each rename; Go's type system catches all mismatches.
- **Merge conflicts**: Any concurrent change touching the same fields will conflict. Mitigation: this is a fast, mechanical change; do it when no other changes to model fields are in flight.
