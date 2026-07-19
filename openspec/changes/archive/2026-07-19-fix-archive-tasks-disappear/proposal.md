## Why

In the tasks view of archived changes (`ModeViewingArchive`), pressing `j`/`k` or arrow up/down causes all tasks to disappear. Exiting and re-entering the tab restores them. This happens because `m.tasks.Items` is never initialized when entering archive mode, and pressing `j`/`k` overwrites the viewport with an empty task list. Additionally, `Space` (toggle) and `e` (editor) are not blocked in archive mode as the spec requires.

## What Changes

- Initialize `m.tasks.Items` when entering `ModeViewingArchive` from the index
- Make `loadViewport` use the interactive task list rendering (not raw glamour) for archive mode too
- Block `Space` (task toggle) in archive mode
- Block `e` (editor) in archive mode

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `archive-viewer`: fix `j`/`k` on TabTasks to keep cursor navigation instead of overwriting the viewport; block `Space` and `e` as already required by the spec

## Impact

- `internal/ui/index.go`: `updateIndex` Enter handler for archived items — add `loadTaskItems()`
- `internal/ui/viewport.go`: `loadViewport` — extend TabTasks condition to include `ModeViewingArchive`
- `internal/ui/viewer.go`: `updateViewer` — add archive mode guard for `Space` and `e`
- `internal/ui/mouse.go`: `clickIndexItem` for archived items — add `loadTaskItems()`

## Non-goals

- `h`/`l` behavior in archive mode is not changed (out of scope)
- No new navigation features beyond what already exists in normal mode
