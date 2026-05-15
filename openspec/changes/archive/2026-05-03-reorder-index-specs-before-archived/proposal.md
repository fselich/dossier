## Why

In the index, Specifications are active, living content of the project, while Archived Changes are history. Showing specs before archived changes reflects their greater relevance to the user browsing the index.

## What Changes

- The order of sections in `ModeIndex` changes from `Active → Archived → Specifications` to `Active → Specifications → Archived`.
- Navigation with `j`/`k` will traverse the sections in the new order.

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- `change-index`: The requirement that defines the order of the three sections changes: "Specifications" is now displayed between "Active Changes" and "Archived Changes".
- `index-specs-section`: The "Specs section in the index" requirement currently states that the section appears *below* "Archived Changes"; it must be updated to reflect that it appears *above* "Archived Changes".

## Impact

- Index rendering code (`ModeIndex`): reorder the construction of `indexItems`.
- Affected specs: `change-index` and `index-specs-section`.
- No changes to APIs or external dependencies.
