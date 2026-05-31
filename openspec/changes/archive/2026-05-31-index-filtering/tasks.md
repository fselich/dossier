## 1. Model: add filter state

- [x] 1.1 Add `FilterText string`, `FilterActive bool`, `FilterIndices []int` fields to `indexState` in `model.go`
- [x] 1.2 Initialize `FilterIndices` as nil (no filter) in `enterIndex()` and `New()`/`NewSinglePath()` (zero values satisfy this)

## 2. Core: filter logic

- [x] 2.1 Implement `matchesFilter(item indexItem, lowerQuery string) bool` that checks the item's name (change name, spec name, or requirement name) against the query
- [x] 2.2 Implement `isItemVisible(idx int) bool` that returns true when no filter is active or the item passes `matchesFilter`
- [x] 2.3 Implement `applyFilter()` that rebuilds `FilterIndices` from `Items`, preserves cursor on the same logical item if still matching, or moves cursor to 0

## 3. Navigation: updateIndex key handling

- [x] 3.1 Add `/` case to the `updateIndex` switch that saves `PrevFilterText`, clears `FilterText`, and sets `FilterActive = true`
- [x] 3.2 Add filter-input block at the top of `updateIndex` that intercepts all keys when `FilterActive`: printable chars append to `FilterText` and call `applyFilter()`, Backspace removes last char and calls `applyFilter()`, Enter confirms (`FilterActive = false`), Esc cancels (`FilterActive = false`, revert `FilterText` to `PrevFilterText`)
- [x] 3.3 Modify existing `Esc` case: if `FilterText != ""`, clear filter (set `FilterText = ""`, `FilterIndices = nil`); otherwise quit
- [x] 3.4 Wrap all item access in navigation (j/k, Enter, Space, s) to use `visibleItemIdx()` indirection when filter is active
- [x] 3.5 Bound cursor navigation (j/k) to `visibleItemCount()` instead of `len(Items)`

## 4. Render: filtered index content

- [x] 4.1 Update `renderIndexContent()` to skip items that fail `isItemVisible()`: iterate `Items` via section boundaries but only render matching ones
- [x] 4.2 After each section (Active/Specs/Archives), if the section has items in the full list but zero matching items, render "No items match '<query>'" instead of the section items

## 5. Render: filter prompt in help bar

- [x] 5.1 In `renderHelpBar()`, when `m.mode == ModeIndex && m.index.FilterActive`, render `/<FilterText>█` instead of normal help text
- [x] 5.2 When `m.mode == ModeIndex && m.index.FilterText != "" && !m.index.FilterActive`, append `  [/<query>]` to the normal help text

## 6. Mouse: click through filter

- [x] 6.1 Update `indexItemAtContentLine()` to skip items that fail `isItemVisible()` (same skipping logic as `renderIndexContent`)
- [x] 6.2 Update `handleMouseClick()` for `ModeIndex` to map cursor comparison through `FilterIndices` when filter is active

## 7. Rebuild: re-apply filter after structural changes

- [x] 7.1 Call `applyFilter()` after `buildIndexItems()` in `pollIndexMode()` and in the tick handler's task-refresh path
- [x] 7.2 Verify that tick-based reload (`pollIndexMode`) and spec expansion (Space) both re-apply the filter and clamp the cursor (Space in updateIndex calls applyFilter; pollIndexMode updated in 7.1)

## 8. Tests

- [x] 8.1 Add test for `matchesFilter` with various name/query combinations (case-insensitive, substring, non-matching)
- [x] 8.2 Add test for `applyFilter` that verifies `FilterIndices` is correctly built and cursor is clamped
- [x] 8.3 Add integration test for `/` key in `ModeIndex`: press `/`, type "d", verify filter updates, press Enter, verify filter stays, press Esc, verify filter clears
- [x] 8.4 Add test for no-match message rendering when filter matches zero items in a section
