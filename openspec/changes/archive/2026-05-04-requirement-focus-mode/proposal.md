## Why

When navigating from the index to a specific requirement, the TUI shows the full spec scrolled to that requirement — but in specs with many requirements it is hard to visually locate the target. Showing only the selected requirement removes all noise and makes the intent of the navigation obvious.

## What Changes

- When the user presses `Enter` on a requirement item in the index, the TUI enters a new **focus mode** that renders only that requirement's content instead of the full spec.
- `h` / `l` navigate to the previous / next requirement within the same spec while staying in focus mode.
- The header shows a `Req N/M` counter so the user always knows where they are.
- The HelpBar updates to reflect focus-mode controls.
- Entering `ModeViewingSpec` from a spec item (not a requirement item) continues to show the full spec, unchanged.

## Capabilities

### New Capabilities

- `requirement-focus-view`: Focused single-requirement rendering mode within the spec viewer, with intra-spec navigation.

### Modified Capabilities

- `spec-detail-viewer`: The "open spec from a requirement item" requirement changes behaviour — instead of scrolling the full spec to the requirement, it now renders only that requirement. The scroll-to-line mechanism is replaced by content extraction.

## Impact

- `internal/ui/model.go`: new fields `specFocusMode bool` and `specReqCursor int` on `Model`; changes to `loadViewport()`, header render, helpbar render, and key handling for `ModeViewingSpec`.
- `internal/openspec/loader.go` or a new helper: `extractRequirement(raw, name string) string` utility function.
- No changes to on-disk data formats, loader structs (beyond what is already there), or external dependencies.
