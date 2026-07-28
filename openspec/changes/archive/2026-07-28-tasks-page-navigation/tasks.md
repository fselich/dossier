## 1. Render: expose per-item line offsets

- [x] 1.1 In `internal/ui/tasks.go`, change `renderTasksContent` to return `(string, int, []int)`: append the local `line` counter to a new `itemLines []int` slice for every item (both `KindSection` and `KindTask`) before its height is added, replacing the existing `cursorLine`-only capture with `itemLines[m.tasks.Cursor]`.
- [x] 1.2 Update `refreshTasksViewport` to receive the new `itemLines` return value and store it on `m.tasks.ItemLines`.
- [x] 1.3 Add `ItemLines []int` to the `taskState` struct in `internal/ui/model.go`.

## 2. Cursor: page-jump helpers

- [x] 2.1 In `internal/ui/tasks.go`, add `moveCursorPageDown(pageHeight int)`: using `m.tasks.ItemLines`, find the last item index whose line is `<= itemLines[Cursor] + pageHeight`; if that index is not past the current cursor, clamp to `len(Items)-1`; set `m.tasks.Cursor`.
- [x] 2.2 Add `moveCursorPageUp(pageHeight int)`: mirror of 2.1 — last item index whose line is `<= itemLines[Cursor] - pageHeight`, clamped to `0`.
- [x] 2.3 Guard both helpers against `len(m.tasks.Items) == 0` (no-op) and `len(m.tasks.ItemLines) != len(m.tasks.Items)` (stale/empty state — no-op) before indexing.

## 3. Key dispatch: wire pgdown/pgup for the Tasks tab

- [x] 3.1 In `internal/ui/viewer.go`, `pgdown` case: add a `TabTasks` branch mirroring the existing `j`/`k` structure — `ModeViewingArchive` calls `m.vp.PageDown()`; otherwise call `m.moveCursorPageDown(m.vp.Height())` followed by `m.refreshTasksViewport()`.
- [x] 3.2 In the same file, `pgup` case: add the matching `TabTasks` branch calling `m.vp.PageUp()` in archive mode, or `m.moveCursorPageUp(m.vp.Height())` + `m.refreshTasksViewport()` otherwise.

## 4. Tests

- [x] 4.1 In `internal/ui/view_test.go`, replace `"pgdown is a no-op on task list"` with a test asserting `pgdown` in `ModeNormal`/`TabTasks` moves `m.tasks.Cursor` forward by roughly a page (construct a task list longer than the viewport height, assert the resulting cursor index increased and lands near the expected target item).
- [x] 4.2 Add a symmetric `pgup` test: page down first to move off the top, then page up, and assert the cursor moves back toward (but not necessarily past) its earlier position.
- [x] 4.3 Add a test for `ModeViewingArchive` + `TabTasks`: assert `pgdown`/`pgup` scroll `m.vp`'s `YOffset()` (parity with the existing archive-mode `j`/`k` scroll test, if one exists — otherwise model it on the `pgdown scrolls a full page in proposal tab` test).
- [x] 4.4 Add a boundary test: cursor near the end of the list, `pgdown` clamps to the last item without going out of range; cursor near the start, `pgup` clamps to `0`.

## 5. Verification

- [x] 5.1 Run `go build ./...` and `make test` (`go test -race -cover ./...`); confirm all tests pass including the new/updated ones in `internal/ui`.
- [x] 5.2 Run `make lint`; confirm no new findings.
- [x] 5.3 Manually smoke-test: open a change with a task list longer than one screen, switch to the Tasks tab, press `<pgdown>`/`<pgup>` and confirm the cursor jumps a full screen and stays visible; repeat while viewing an archived change's tasks.
</content>
