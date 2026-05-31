## Why

The header row (project name + change name) is the most prominent clickable area in the TUI, but clicking it does nothing. Since `a` and `Esc` already navigate to the index, clicking the header should be a convenient mouse shortcut for the same action.

## What Changes

- Left-click on the header row (screen Y=1) in `ModeNormal` or `ModeViewingArchive` enters the index view (same as `a` or `Esc`)
- Click on the header row in other modes does nothing

## Capabilities

### New Capabilities

*(none)*

### Modified Capabilities

- `mouse-navigation`: Add requirement for header click to navigate to index

## Impact

- `internal/ui/mouse.go`: extend `handleMouseClick` with header click handling
