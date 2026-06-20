## Why

In the Tasks tab, section headers (e.g., `## 1. UI Changes`) anchor subsections but the cursor only ever lands on task items. When a section has many tasks and the user scrolls down past the viewport height, then scrolls back up, the section header is hidden above the visible area. The cursor cannot land on the header, so there is no way to scroll the viewport up enough to reveal it.

## What Changes

- Allow the cursor to land on `KindSection` items in addition to `KindTask` items during up/down navigation
- Preserve existing behavior: `doToggle` still skips section headers, `firstTaskIdx` still starts on the first task
- The cursor renders on section headers (no checkbox toggle indicator, just highlight/position)
- `refreshTasksViewport` ensures the cursor line (now possibly a section header) stays visible

## Capabilities

### New Capabilities

- `<none>`: no new capabilities — this is a pure UI fix within the existing tasks view

### Modified Capabilities

- `<none>`: no spec-level requirement changes — purely an implementation/UX fix

## Impact

- `internal/ui/tasks.go` — `moveCursorUp`, `moveCursorDown`, `refreshTasksViewport`, `renderTasksContent`, `firstTaskIdx` all need adjustment
- `internal/ui/tasks.go` - `doToggle` already guards against `KindSection` — no change needed there
- No API, dependency, or filesystem changes
