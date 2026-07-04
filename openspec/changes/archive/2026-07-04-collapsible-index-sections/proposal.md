## Why

The index view currently renders three sections (Active Changes, Specifications, Archived Changes) always fully expanded. Users with many changes or specs must scroll past content they don't need. Adding collapsible sections lets users focus on what matters, fold away distractions, and navigate more efficiently.

## What Changes

- Introduce sections as navigable items in the index item list so the cursor can land on section headers
- `Space` toggles collapse/expand on a section (context-sensitive: `Space` on a spec still toggles requirement expansion as today)
- Collapsed sections hide their child items from the list and rendering
- Visual indicator (`▼` / `▶`) on section headers shows collapsed state
- Section collapse state is preserved across index rebuilds (tick polls, filter toggles)
- Help bar updated to show the new section toggle action

## Capabilities

### New Capabilities

- `collapsible-sections`: Navigable, foldable section headers in the index view that hide/show their child items

### Modified Capabilities

- (none — no existing specs are affected)

## Impact

- `internal/ui/index.go`: Add section item kind, section items to buildIndexItems, collapse/expand rendering, Space key handler for sections
- `internal/ui/model.go`: Add `CollapsedSections` field to `indexState`
- `internal/ui/view.go`: Update help bar text
- `internal/ui/styles.go`: May need a style for collapsed section indicator
- No changes to domain types (`internal/openspec/`), model update pattern, or filesystem layer
