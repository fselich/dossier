## Why

When a change has multiple specs (one per capability), the `specs` tab shows them all concatenated into a single long document. With two or more specs the result is a wall of text that is difficult to read and navigate.

## What Changes

- Specs are stored individually in the model as `[]NamedSpec` instead of a flattened string
- The `specs` tab shows a navigation sub-bar with the name of each spec as selectable chips
- The `3` key cycles through specs when already on the `specs` tab; if on another tab, it switches to `specs` (showing the last selected spec)
- `contentHeight` is reduced by 1 when the sub-nav is visible so the viewport does not overflow
- The help bar shows `3: spec` instead of (or in addition to) `1-4: artifact` when on the specs tab

## Capabilities

### New Capabilities

- `specs-subnav`: Sub-navigation within the specs tab with chips per capability and cycling with the `3` key

### Modified Capabilities

- `tui-viewer`: The tab navigation requirement changes — the `3` key now has dual behaviour (go to specs / cycle specs)

## Impact

- `internal/openspec/loader.go`: add `NamedSpec` and `[]NamedSpec` to `Change`; `loadSpecs` populates both
- `internal/ui/model.go`: add `specIdx int`, sub-nav logic, `contentHeight` adjustment, updated `3` key handler
- No new dependencies
