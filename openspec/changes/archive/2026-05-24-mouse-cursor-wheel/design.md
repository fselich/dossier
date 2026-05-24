## Context

`mouse.go` currently handles all wheel events by scrolling the viewport 3 lines, regardless of mode. In `update.go`, `j`/`k` and `up`/`down` keys already have mode-aware behavior: they move the cursor in index and tasks modes, and scroll the viewport everywhere else. The wheel should mirror this behavior.

## Goals / Non-Goals

**Goals:**
- In `ModeIndex`: wheel moves the index cursor, viewport auto-follows
- In `ModeNormal` + `TabTasks`: wheel moves the task cursor, viewport auto-follows
- All other modes/views: wheel scrolls viewport as before

**Non-Goals:**
- Changing wheel behavior in any other mode
- Changing scroll speed or direction
- Adding wheel-driven navigation to archive mode tabs or spec subnav

## Decisions

### 1. Branch wheel handling on mode, not a new function

The existing `j`/`k` handler in `update.go` already does this exact dispatch. The wheel handler follows the same pattern inline within `handleMouse()`, keeping the logic collocated in `mouse.go`.

```go
case tea.MouseButtonWheelDown:
    if m.mode == ModeIndex {
        if m.indexCursor < len(m.indexItems)-1 { m.indexCursor++ }
        m.refreshIndexViewport()
    } else if m.tab == TabTasks && m.mode == ModeNormal {
        m.moveCursorDown()
        m.refreshTasksViewport()
    } else {
        m.vp.LineDown(3)
    }
```

**Alternative considered:** Extract a shared `dispatchVertical(direction)` used by both keyboard and mouse. Rejected because:
- Keyboard handles `down`/`up` via `tea.KeyMsg` string matching in `Update()`, not a method call; refactoring adds risk for a 15-line change
- Wheel scrolls 3 lines, keyboard scrolls 1 — they intentionally differ

### 2. Viewport auto-follows via existing refresh methods

`refreshIndexViewport()` and `refreshTasksViewport()` already handle scroll-to-cursor logic. No new viewport logic needed — cursor movement triggers the existing refresh.

### 3. One wheel tick = one cursor step

Unlike content scrolling (3 lines per tick), cursor movement is always one item per tick. This matches keyboard behavior (`j`/`k` = one item) and prevents wheel from skipping items.

## Risks / Trade-offs

- **Index with many expanded specs**: rapid wheel could trigger many `refreshIndexViewport()` calls → Minimal risk, the render is O(n) on index items (typically < 50)
- **Task view with wrapped items**: `refreshTasksViewport()` already handles word-wrap and cursor positioning → No new risk
