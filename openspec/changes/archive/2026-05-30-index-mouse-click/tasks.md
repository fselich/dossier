## 1. Helper: map content line to index item

- [x] 1.1 Add `indexItemAtContentLine(line int) (int, bool)` method to Model that iterates `m.indexItems` while mirroring `renderIndexContent()` line counting, handling active items, spec items with expanded requirements, and archived items in order
- [x] 1.2 Add `indexViewportContentStart()` constant or helper returning the screen row where viewport content begins (row 3: boxTop + header + innerSep)

## 2. Extend mouse click handler

- [x] 2.1 Add `ModeIndex` case in `handleMouseClick` before the existing mode guard: validate click is within viewport content area (row 3 to 3 + vp.Height())
- [x] 2.2 Calculate `contentLine = msg.Y - 3 + m.vp.YOffset()` and find item via `indexItemAtContentLine`
- [x] 2.3 If click maps to an item: if it differs from `m.indexCursor`, move cursor there; if same, perform action (Enter for active/archived/requirement, Space for spec toggle; reuse existing logic from update.go Enter/Space handlers)

## 3. Tests

- [x] 3.1 Add test cases for `indexItemAtContentLine` with various index states (with/without specs, with/without expanded specs)
- [x] 3.2 Add test cases for click behavior in `ModeIndex`: click on unselected item, click on selected item, click on blank line, click outside viewport

## 4. Helpbar

- [x] 4.1 Update helpbar in `ModeIndex` to include mouse hint (e.g., "click: select")
