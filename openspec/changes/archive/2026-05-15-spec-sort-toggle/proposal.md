## Why

When a project has many specs, alphabetical order by full name buries related concepts across the list. Sorting by the last segment of the name (e.g. `viewer` from `archive-viewer`) groups conceptually related specs together, making it easier to find all specs that touch the same concern.

## What Changes

- Add a sort toggle in `ModeIndex` bound to the `s` key
- Specs section cycles between two orderings: alphabetical by full name (default) and alphabetical by name suffix (last `-`-delimited segment)
- Help bar reflects the available toggle action contextually
- Sort state is session-only; it resets to default on restart

## Capabilities

### New Capabilities

- `index-spec-sort`: Toggle sort order of the Specifications section in `ModeIndex` between full-name and suffix ordering, preserving cursor position and expand state across toggles.

### Modified Capabilities

- `change-index`: Add `s` key binding to `ModeIndex`; update help bar text to reflect current sort state.

## Impact

- `internal/ui/model.go`: new field `specSortBySuffix bool`, new field `specOrder []int`, new helper `buildSpecOrder()`, modified `buildIndexItems()` and `renderIndexContent()`, new key handler for `s`, updated `renderHelpBar()`
- No new dependencies
- No breaking changes

## Non-goals

- Persisting sort preference across sessions
- Sorting archived changes or active changes
- More than two sort modes
