## Why

In the index view (`ModeIndex`), users can navigate items with the keyboard (j/k) and scroll the cursor with the mouse wheel, but clicking an item does nothing. This is inconsistent with the tab bar (where left-click switches tabs) and creates a UX gap: mouse users expect clickable items.

## What Changes

- Left-click on an index item (`ModeIndex`) moves the cursor to that item
- Left-click on the currently selected item opens it (same action as Enter)
- Spec items toggle expansion on click (same action as Space)
- All click behaviors handle viewport scrolling offsets correctly

## Capabilities

### New Capabilities

*(none — this is an interaction improvement to existing capabilities)*

### Modified Capabilities

- `mouse-navigation`: Add requirement for index item selection via left-click
- `change-index`: Add scenario noting that click is now a selection trigger

## Impact

- `internal/ui/mouse.go`: extend `handleMouseClick` to handle `ModeIndex`
- No new dependencies
- No API changes
