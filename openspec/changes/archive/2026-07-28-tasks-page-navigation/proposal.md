## Why

`pgup`/`pgdown` scroll every document tab (Proposal, Design, Specs, Git diff) by a full viewport page, but the Tasks tab silently ignores them — both while viewing an archived (read-only) change's tasks and while navigating the live, editable task cursor. Users reviewing a change with a long task list have no fast way to jump a screen at a time.

## What Changes

- Tasks tab, `ModeViewingArchive` (read-only tasks): `pgdown`/`pgup` scroll the raw viewport by a page, matching the existing `j`/`k` → `vp.ScrollDown/Up(1)` behavior already used there.
- Tasks tab, `ModeNormal` (live cursor): `pgdown`/`pgup` move `m.tasks.Cursor` to the item whose rendered line is approximately one viewport-height away from the cursor's current line — a line-accurate jump, not a fixed item count. Cursor movement past either end clamps to the first/last item (same clamping already used by `moveCursorDown`/`moveCursorUp`).
- `renderTasksContent` additionally reports the rendered start line of every task item (not just the cursor's), so the page-jump can look up a target item by line without duplicating the render's height-counting logic.
- Git tab's non-diff file list keeps its current `pgdown`/`pgup` no-op — explicitly out of scope for this change.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `tasks-cursor-navigation`: adds page-jump (`pgdown`/`pgup`) as a keyboard navigation method for the tasks cursor, alongside the existing single-step `up`/`down`.

## Impact

- `internal/ui/tasks.go`: `renderTasksContent` signature grows an `itemLines []int` return value; `refreshTasksViewport` stores it; new `moveCursorPageDown`/`moveCursorPageUp` helpers.
- `internal/ui/model.go`: `taskState` gains an `ItemLines []int` field.
- `internal/ui/viewer.go`: `pgdown`/`pgup` cases gain a `TabTasks` branch (archive scroll vs. live cursor jump), mirroring the existing `j`/`k` branch structure.
- `internal/ui/view_test.go`: replace the `"pgdown is a no-op on task list"` assertion (behavior it tests is being removed); add coverage for live-cursor page jump (down and up) and archive-mode page scroll.
- No public API, storage format, or CLI flag changes.
</content>
