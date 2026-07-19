## Context

`ModeViewingArchive` shares the key dispatcher and viewport with `ModeNormal` via the `updateViewer` method (`viewer.go:12`). For `TabTasks`, `j`/`k` navigation calls `moveCursorDown/Up` + `refreshTasksViewport`, which iterates over `m.tasks.Items`. This field is initialized via `loadTaskItems()`, which extracts `Change.Tasks` through `m.current()` (which redirects to `m.currentArchive()` when the mode is `ModeViewingArchive`).

The problem: when entering `ModeViewingArchive` from the index (keyboard or mouse), `loadTaskItems()` is never called. `m.tasks.Items` retains data from the last active change (or stays `nil`). Pressing `j`/`k` causes `refreshTasksViewport` to overwrite the viewport with empty content.

The initial viewport works because `loadViewport` falls through to the `default` case → `loadViewportForArtifact`, which renders raw `tasks.md` via glamour. When switching back to the tab (`4`), the render cache returns that glamour content without going through `m.tasks.Items`.

Additionally, the `archive-viewer` spec requires `Space` and `e` to be ignored in archive mode, but `updateViewer` has no such guards. `doToggle` would attempt to write to disk without checking the mode.

## Goals / Non-Goals

**Goals:**
- `j`/`k` on `TabTasks` of an archived change must navigate the task cursor correctly (same behavior as `ModeNormal`)
- `Space` and `e` must be silently ignored in `ModeViewingArchive`
- The initial viewport of the tasks tab in archive mode must show the interactive list (not raw glamour), same as `ModeNormal`

**Non-Goals:**
- `h`/`l` behavior in archive mode is not changed
- The `archive-viewer` spec baseline is not modified (only a delta spec is created to clarify `j`/`k` behavior on TabTasks)
- No new navigation features are added

## Decisions

### 1. Call `loadTaskItems()` when entering archive mode

Invoke `m.loadTaskItems()` in both `updateIndex` (keyboard, line 762) and `clickIndexItem` (mouse, line 152), immediately after setting `m.index.ArchiveCursor` and before `commitStateChange`. This mirrors exactly what is done for active changes (`indexKindActive`, line 737).

Alternative considered: defer loading to `updateViewer`. Rejected because it breaks symmetry with active changes and adds a conditional check in the hot path of every keypress.

### 2. Extend `loadViewport` for `ModeViewingArchive`

Change the condition in `viewport.go:23` from:
```go
case m.tab == TabTasks && m.mode == ModeNormal:
```
to:
```go
case m.tab == TabTasks && (m.mode == ModeNormal || m.mode == ModeViewingArchive):
```

This makes the initial viewport also use `loadViewportForTasks` (interactive rendering with cursor and progress bars), not raw glamour.

### 3. Guards for `Space` and `e` in `updateViewer`

Add `m.mode == ModeViewingArchive` check:
- `Space`: return `m, nil` if in archive mode (don't call `doToggle`)
- `e`: return `m, nil` if in archive mode and not git diff (don't open editor)

### 4. Mouse wheel in archive mode for TabTasks

The mouse wheel (`mouse.go:19,39`) already handles the archive case correctly because it only acts on `TabTasks && ModeNormal`; in archive mode it falls through to the `default` scroll case. No changes needed.

## Risks / Trade-offs

- **[Risk] Transition from glamour to interactive list**: In `ModeNormal`, the tasks viewport shows the interactive list with cursor, checkboxes, and progress bars. In archive mode it will now do the same (read-only). This is the correct behavior per the spec ("work the same as in normal mode"), but visually differs from the raw markdown that was previously shown. → **Mitigation**: The spec already defines this behavior. Users expect consistency between modes.
