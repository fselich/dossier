## Context

The Tasks tab renders a list of `TaskItem` elements (both `KindSection` headers and `KindTask` items) parsed from `tasks.md`. Cursor navigation (`moveCursorUp`/`moveCursorDown`) currently skips `KindSection` items, landing only on `KindTask` items. The viewport is anchored to keep the cursor line visible. This means when a section has many tasks and the user scrolls down past the viewport height, the section header above the first visible task is unreachable — the cursor cannot land on it, and `refreshTasksViewport` has no reason to scroll up past the first task.

## Goals / Non-Goals

**Goals:**
- Allow cursor to land on `KindSection` items during up/down navigation
- Section headers become visible when scrolling up — cursor can now rest on a header
- Toggle action (`doToggle`) still rejects `KindSection` items
- `refreshTasksViewport` naturally handles section line heights since cursor line tracking already works for any index

**Non-Goals:**
- No visual changes to how sections are rendered (style, progress bar, spacing)
- No changes to the filesystem, task parsing, or toggle logic
- No changes to spec or proposal tabs

## Decisions

- **Cursor lands on sections**: `moveCursorUp` and `moveCursorDown` will iterate over all items instead of only `KindTask` items. This is the simplest change — one condition removed per function.
- **`firstTaskIdx` unchanged**: On initial load, cursor still starts on the first task (not the first section header). This preserves the existing UX for the common case.
- **No viewport logic change**: `refreshTasksViewport` and `renderTasksContent` already handle section line heights correctly. The cursor line tracked by `renderTasksContent` is just the line position in the rendered output — it works for any item index. No changes needed.
- **`doToggle` unchanged**: It already guards `m.tasks.Items[m.tasks.Cursor].Kind != openspec.KindTask` — section headers are already skipped.

## Risks / Trade-offs

- **[Low] Cursor on section header might feel unfamiliar**: users with large sections who scroll up will stop on the section header. Mitigation: the header is rendered with a visible style (already styled via `sectionStyle`), and the cursor mark `▶` will appear on the header line, making it clear the cursor has landed there.
- **[Low] Section header adds one more `KindSection` check**: The `renderTasksContent` switch already handles `KindSection` and `KindTask` separately. Minimal risk.
