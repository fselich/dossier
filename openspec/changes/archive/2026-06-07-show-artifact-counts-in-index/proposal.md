## Why

The index view section titles ("Active Changes", "Specifications", "Archived Changes") don't show how many items each section contains, forcing the user to mentally count or scroll to assess the contents at a glance.

## What Changes

- Section titles in the index view show item counts: "Active Changes (3)", "Specifications (5)", "Archived Changes (12)"
- Counts reflect the total number of items in each section (not filtered items)
- "No active changes" / "No specifications available" / "No archived changes" messages remain unchanged (no count shown when empty)
- No changes to data loading, filtering, or other UI areas

## Capabilities

### New Capabilities

None — this is a purely cosmetic enhancement to the existing index view.

### Modified Capabilities

- `change-index`: Section titles now display item counts (e.g., "Active Changes (3)")

## Impact

- `internal/ui/index.go`: Update `renderIndexContent()` to include counts in section titles
