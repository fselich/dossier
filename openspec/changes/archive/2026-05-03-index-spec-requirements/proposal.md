## Why

The index shows specs with a requirement count but no names, so you can't tell at a glance which spec covers what you need. Expanding a spec to show its requirement names — and jumping directly to a requirement in the viewer — turns the index into a useful navigation surface instead of just a list.

## What Changes

- `ProjectSpec` gains a `RequirementNames []string` field; `LoadProjectSpecs()` populates it alongside the existing `RequirementCount`
- Each spec item in the index becomes expandable: pressing `Space` toggles a list of its requirement names below the spec name
- Requirements in the expanded list are navigable with `j`/`k` as flat items alongside specs and changes
- Pressing `Enter` on a requirement item opens `ModeViewingSpec` and scrolls the viewport to that requirement
- `buildIndexItems()` inserts `indexKindRequirement` items after each expanded spec item; expand state is tracked per spec in a `map[int]bool`

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `index-specs-section`: expand/collapse of spec items to show requirement names; requirement items as navigable entries in the index cursor model; display format changes
- `spec-detail-viewer`: scroll-to-requirement on entry when opened from a requirement item

## Impact

- `internal/openspec/loader.go`: add `RequirementNames []string` to `ProjectSpec`, populate in `LoadProjectSpecs()`
- `internal/ui/model.go`: new `indexKindRequirement` kind; `expandedSpecs map[int]bool` field on `Model`; `buildIndexItems()` flattens requirements when expanded; `handleTick()` and `enterIndex()` reset/preserve expand state; `renderIndexContent()` renders indented requirement lines; `loadViewport()` computes scroll offset for target requirement; `specRenderedMsg` carries optional `jumpToLine int`
