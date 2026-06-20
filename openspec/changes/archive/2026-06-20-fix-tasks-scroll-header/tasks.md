## 1. Cursor Navigation Changes

- [x] 1.1 Modify `moveCursorUp` to iterate over all items instead of only `KindTask` items — cursor can now land on section headers when scrolling up
- [x] 1.2 Modify `moveCursorDown` to iterate over all items instead of only `KindTask` items — cursor can now land on section headers when scrolling down
- [x] 1.3 Verify `firstTaskIdx` still starts cursor on the first `KindTask` on initial load (no change needed, but confirm)

## 2. Rendering and Toggle

- [x] 2.1 Verify `renderTasksContent` renders the cursor mark `▶` on section header lines when cursor lands on them (check that `sectionStyle` and cursor markup compose correctly)
- [x] 2.2 Verify `doToggle` already rejects `KindSection` cursor (guard at line 66) — add test coverage if missing

## 3. Tests

- [x] 3.1 Add tests for `moveCursorUp`/`moveCursorDown` navigating to section headers
- [x] 3.2 Add tests for cursor rendering on section headers in `renderTasksContent`
