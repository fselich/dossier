## Why

Dossier is purely keyboard-driven. Users who prefer the mouse or switch between keyboard and mouse have no clickable targets. Tab switching is a natural first target: the tab bar is a fixed-position UI element with clear affordance for clicking.

## What Changes

- Enable `tea.WithMouseCellMotion()` to receive mouse events from the terminal
- Handle `tea.MouseMsg` wheel events (up/down) to maintain scroll functionality that currently works via terminal translation of wheel to arrow keys
- Handle left-click on the tab bar (`Y=2`) by mapping the X coordinate to the corresponding tab label, triggering the same behavior as pressing `1`/`2`/`3`/`4`
- Disabled tabs are not clickable; clicks on gaps between tabs or outside tab boundaries are no-ops

## Capabilities

### New Capabilities

- `mouse-navigation`: Mouse input handling for the TUI, starting with tab clicking and wheel scrolling. Covers `tea.WithMouseCellMotion` enablement, `tea.MouseMsg` dispatching, and coordinate-to-action mapping.

### Modified Capabilities

- `tui-viewer`: Tab navigation requirement gains a mouse-click alternative alongside the existing keyboard shortcuts. Help bar updated to include mouse hints where relevant.

## Impact

| Area | Detail |
|---|---|
| `cmd/dossier/main.go` | Add `tea.WithMouseCellMotion()` to program options |
| `internal/ui/update.go` | Add `case tea.MouseMsg` handler; delegate wheel to viewport, map tab bar clicks |
| `internal/ui/model.go` | New `handleMouse()` method or inline logic in update |
| `internal/ui/view.go` | Help bar text updated to mention mouse (optional, low priority) |
| Scroll behavior | When mouse mode is on, terminal no longer translates wheel to arrow keys; app must forward wheel events to viewport explicitly |
| Text selection | Holding Shift while selecting text bypasses mouse capture (standard terminal behavior) |

## Non-goals

- Clicking on index items (deferred to future change)
- Right-click for "back" navigation
- Clicking on spec subnav labels
- Drag-to-scroll or hover effects
- `tea.WithMouseAllMotion()` (unnecessary for click events)
