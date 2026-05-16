## Why

When an active change has more tasks than fit in the visible terminal window, navigating with j/k does not scroll the viewport — tasks beyond the fold are unreachable. The root cause is a line-counting bug in `renderTasksContent()`: each task always increments the line counter by 1, but lipgloss may render a task as multiple lines when its text wraps to fit the content width. This desynchronizes `cursorLine` from the actual rendered position, causing `refreshTasksViewport()` to scroll to the wrong offset (or not at all).

## What Changes

- Fix `renderTasksContent()` to count rendered lines using `lipgloss.Height()` instead of a flat `+1` per item.
- Apply the same fix to section header lines, which can also wrap on narrow terminals.

## Non-goals

- Adding free/arbitrary scroll (independent of the cursor) to the tasks view.
- Changing navigation key bindings.
- Modifying the visual appearance of tasks.

## Capabilities

### New Capabilities

<!-- none -->

### Modified Capabilities

- `tasks-toggle`: The tasks view SHALL correctly scroll to keep the cursor-selected task visible, including when task text wraps across multiple lines.

## Impact

- `internal/ui/tasks.go` — `renderTasksContent()` line-counting logic.
