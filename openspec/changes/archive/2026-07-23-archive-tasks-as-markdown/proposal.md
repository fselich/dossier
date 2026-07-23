## Why

Archived task lists use a special interactive renderer (cursor navigation, progress bars) inherited from `ModeNormal`, but they're read-only data. This caused a regression (cursor movement would blank the viewport) that was patched with brittle guards in multiple files. Rendering archived tasks as markdown via glamour — the same path used by every other tab in archive mode — eliminates the bug at its root and simplifies the codebase.

## What Changes

- **Tasks tab in `ModeViewingArchive` renders via glamour** instead of the interactive `renderTasksContent()` path
- **j/k scroll** the viewport on the tasks tab in archive mode (matching proposal/design/specs tabs), instead of moving a cursor
- **Remove archive-mode guards** in `viewer.go` that blocked `space` and `e` — these are no longer reachable
- **Remove `loadTaskItems()` calls** added by the previous fix in `index.go` and `mouse.go`
- **Unify help bar** for all tabs in archive mode (all show "j/k: scroll")
- **Remove `ModeViewingArchive` from `loadViewport()` tasks dispatch** — it falls through to `loadViewportForArtifact()` which already handles `TabTasks` via glamour

## Capabilities

### New Capabilities

None. This is a behavior change to an existing capability.

### Modified Capabilities

- `archive-viewer`: Tasks tab in archive mode changes from interactive (cursor navigation, progress bars, toggle) to read-only markdown rendered view (scroll only, glamour-rendered like proposal/design/specs)

## Impact

- `internal/ui/viewport.go` — one-line change in `loadViewport()` dispatch
- `internal/ui/viewer.go` — add mode guards to j/k handlers, remove dead archive guards for space/e
- `internal/ui/view.go` — simplify help bar (remove archive+tasks special case)
- `internal/ui/index.go` — remove `loadTaskItems()` call
- `internal/ui/mouse.go` — remove `loadTaskItems()` call
- `openspec/specs/archive-viewer/spec.md` — update requirement for tasks tab behavior
- Tests in `internal/ui/view_test.go` — update archive+tasks scenarios

## Non-goals

- Does not change `ModeNormal` tasks tab behavior at all
- Does not remove the interactive task renderer (`renderTasksContent`) — it remains used for active changes
- Does not change how archived changes are listed in the index
