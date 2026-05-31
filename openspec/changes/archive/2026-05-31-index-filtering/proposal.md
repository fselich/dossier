## Why

The index view (`ModeIndex`) shows all active changes, specs, and archived changes in a flat list. As the project grows, users need to quickly find a specific item by name. Currently the only navigation is linear (j/k), which is slow when there are many items.

## What Changes

- Pressing `/` in `ModeIndex` enters a filter prompt at the bottom of the screen
- Typing filters the index items in real-time using case-insensitive substring matching
- `Enter` confirms the filter (stays applied), `Esc` clears it or cancels typing
- Filter matches against change names (active and archived), spec names, and requirement names
- When no items match, a "No items match" message is shown in the relevant section
- Esc behavior becomes contextual: clears filter first, then quits

## Capabilities

### New Capabilities
- `index-filtering`: Real-time filtering of the index view by pressing `/` and typing a substring query

### Modified Capabilities
- `change-index`: Add filtering requirements — `/` key to filter, real-time matching, contextual Esc, no-match message

## Impact

- `internal/ui/index.go`: New filter state fields, filter input handling in `updateIndex`, filtered rendering in `renderIndexContent`, filtered mouse mapping in `indexItemAtContentLine`
- `internal/ui/model.go`: New fields in `indexState`
- `internal/ui/view.go`: Filter prompt and indicator in `renderHelpBar`
- `internal/ui/mouse.go`: Mouse click mapping through filter indices
