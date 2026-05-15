## Why

The specs sub-bar currently appears between the tab bar and the horizontal separator, within the same top navigation visual block. This looks visually untidy: the spec chips appear to be an extension of the main tabs rather than belonging to the content area of the specs section.

## What Changes

- The spec chip row is no longer rendered between the tab bar and the horizontal separator.
- The chip row is now rendered as the first line inside the content block, immediately after the separator, when the active tab is `specs`.
- The viewport height reduction is maintained (the subnav still occupies 1 line), but the line is deducted from the content area, not the navigation area.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `specs-subnav`: The subnav is rendered inside the content block (after the `├───┤` separator), not in the top navigation block (before the separator). The position requirement changes from "below the tab bar" to "first line of the content area".

## Impact

- `internal/ui/model.go`: `View()` function — move the subnav render from before `boxInnerSep()` to after it; adjust `contentHeight()` if necessary.
- No changes to business logic, navigation or keys.
