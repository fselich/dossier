## 1. Model: add section kind and collapse state

- [x] 1.1 Add `indexKindSection` to the `indexItemKind` enum in `model.go`
- [x] 1.2 Add `CollapsedSections [3]bool` field to `indexState` in `model.go` (three sections: Active, Specs, Archived)
- [x] 1.3 Add a section name constant or helper to map section index to display name

## 2. buildIndexItems: insert section items

- [x] 2.1 In `buildIndexItems` (index.go), prepend an `indexKindSection` item before each section's children
- [x] 2.2 Guard child item insertion: only add children if `!m.index.CollapsedSections[sectionIdx]`
- [x] 2.3 Update `buildSpecOrder` / spec ordering to work with sections as items in the list (spec children still appear after the spec section header)

## 3. renderIndexContent: render section items

- [x] 3.1 Update `renderIndexContent` to handle `indexKindSection`: render header with `▼`/`▶` indicator and styled section name
- [x] 3.2 Ensure the cursor `▶` marker can appear on section header lines
- [x] 3.3 Remove the now-redundant hardcoded section header rendering that was done outside the item loop
- [x] 3.4 Verify the section header line count is consistent for `indexItemAtContentLine` mouse mapping

## 4. updateIndex: Space toggles sections

- [x] 4.1 In `updateIndex`, add a `Space` handler for `indexKindSection`: toggle `CollapsedSections[item.idx]`, rebuild items, apply filter, clamp cursor, refresh viewport
- [x] 4.2 Ensure existing `Space` on `indexKindSpec` still toggles `ExpandedSpecs` (no regression)
- [x] 4.3 Add section constants (e.g., `sectionActive = 0`, `sectionSpecs = 1`, `sectionArchived = 2`) to avoid magic numbers

## 5. Help bar + styles

- [x] 5.1 Update help bar text in `view.go` to show `Space: toggle section` (or context-sensitive hint)
- [x] 5.2 Add or adjust styles in `styles.go` for section header rendering if needed

## 6. Tests

- [x] 6.1 Update existing index tests to account for section items in the item list (cursor positions, item counts)
- [x] 6.2 Add test: collapsing a section hides its children from the item list
- [x] 6.3 Add test: Space on section toggles collapse state
- [x] 6.4 Add test: Space on spec still expands requirements (regression)
- [x] 6.5 Add test: filter respects collapsed state
- [x] 6.6 Add test: cursor navigation through sections and items works correctly
