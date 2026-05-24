## Why

In index and tasks views, the mouse wheel currently scrolls the viewport independently of cursor position — the cursor stays put while content moves under it. Since these views have selectable items, the wheel should move the cursor between items instead (like `j`/`k`), with the viewport auto-following the cursor. This matches how list navigation works in every GUI application.

## What Changes

- In `ModeIndex`: mouse wheel up moves cursor up, wheel down moves cursor down (same as `j`/`k` in index mode)
- In `ModeNormal` with `TabTasks`: mouse wheel up moves task cursor up, wheel down moves task cursor down (same as `j`/`k` in tasks tab)
- In all other views (proposal, design, specs, config, archive): wheel continues scrolling the viewport as before
- Viewport auto-scrolls to keep cursor visible after wheel-driven cursor movement

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `mouse-navigation`: The "Wheel scrolling" requirement changes — in index and tasks modes, wheel now moves the cursor instead of scrolling the viewport.

## Impact

| Area | Detail |
|---|---|
| `internal/ui/mouse.go` | Modify wheel handling in `handleMouse()` — dispatch to cursor movement in index/tasks modes instead of `vp.LineUp`/`vp.LineDown` |
| No other files | Single-function change |

## Non-goals

- Changing wheel scroll speed in content views
- Adding wheel-driven cursor to any other mode
- Inverting wheel direction
