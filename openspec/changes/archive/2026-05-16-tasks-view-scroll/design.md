## Context

`renderTasksContent()` in `tasks.go` builds the full task list as a string and returns `(content, cursorLine)`. `refreshTasksViewport()` uses `cursorLine` to set the viewport Y-offset so the selected task stays visible.

The bug: `renderTasksContent()` always does `line++` for each task item (line 162), regardless of how many terminal lines the rendered text occupies. When lipgloss wraps a task to fit `contentWidth`, the item takes more than one line — but `line` only advances by 1. Every task after a wrapped item has an underestimated `cursorLine`, so `refreshTasksViewport()` scrolls to the wrong position.

## Goals / Non-Goals

**Goals:**
- Fix `renderTasksContent()` to advance `line` by the actual rendered height of each item using `lipgloss.Height()`.
- Apply the same fix to section header rendering.

**Non-Goals:**
- Adding free/arbitrary scroll independent of the cursor.
- Any visual or behavioral changes to navigation, toggles, or progress bars.

## Decisions

**Use `lipgloss.Height()` on the already-rendered string**

`rendered` is computed before `line` is advanced. `lipgloss.Height(rendered)` counts the actual newlines in the rendered string, which is the exact number of terminal rows the item occupies. This is the minimal, zero-allocation fix — no re-rendering, no extra state.

Alternative considered: track wrap count manually by dividing text width by `contentWidth`. Rejected: fragile with ANSI escape sequences embedded in the string; `lipgloss.Height()` handles this correctly.

**Apply to section headers as well**

Section headers use a hardcoded `line++` (line 131). On very narrow terminals a long section name could wrap. Fixing them consistently prevents a latent version of the same bug.

## Risks / Trade-offs

- **Performance**: `lipgloss.Height()` does a string scan. For typical task lists (tens of items) this is negligible.
- **Edge case — zero-height render**: `lipgloss.Height()` returns at least 1 for non-empty strings. The `\n` appended after `rendered` in `sb.WriteString(rendered + "\n")` adds one extra newline not counted by `lipgloss.Height()`, but this is already accounted for by the existing +1 per item — so the replacement is a direct swap.

## Migration Plan

No data model or file format changes. Rollback is reverting the two `line +=` changes in `tasks.go`.
