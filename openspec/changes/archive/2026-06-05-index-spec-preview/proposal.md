## Why

When browsing specifications in the index view, the user sees only the spec name and requirement count — no indication of what the spec is actually about. They must press `Enter` to open the full spec viewer just to read its purpose. A persistent preview bar showing the spec's purpose text would let users quickly scan specs without leaving the index.

## What Changes

- Add a fixed 1-line preview bar between the viewport and the helpbar in `ModeIndex`
- The bar shows the currently selected spec's name and its purpose text (truncated with `…` if too long)
- The bar is always present (1 line reserved), showing content when cursor is on a `indexKindSpec` or `indexKindRequirement` item, empty otherwise
- Add `ExtractPurpose()` helper to the `openspec` package

## Capabilities

### New Capabilities
- `spec-preview-bar`: A persistent 1-line bar in the index view that displays the name and purpose of the currently selected specification

### Modified Capabilities

- `change-index`: The index chrome gains a new fixed row; `contentHeight` is reduced by 1 in `ModeIndex` to accommodate it

## Impact

- `internal/ui/view.go`: Insert `renderSpecPreview()` into `viewContentWithChrome()`
- `internal/ui/model.go`: Update `contentHeight()` for `ModeIndex`
- `internal/ui/index.go`: No changes needed (update keys in `updateIndex()` not required since bar is always visible)
- `internal/openspec/loader.go`: Add `ExtractPurpose()` function
- `internal/ui/styles.go`: Minor — reuses existing styles
