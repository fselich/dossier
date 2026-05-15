## Context

The tasks view is rendered in `renderTasksContent()` in `model.go`. Each item produces a line of text that is passed to a lipgloss style. The item under the cursor receives `.Width(m.width - 2)`, which causes lipgloss to wrap the content. The remaining items have no Width defined and therefore do not wrap.

The text of each task arrives as a plain string from `TaskItem.Text` — the `tasks.md` parser extracts the text following `- [ ]` or `- [x]` without processing markdown.

## Goals / Non-Goals

**Goals:**
- Consistent word wrap on all items
- Rendering of `**bold**` and `` `code` `` inline in each task's text
- Global progress bar visible at the top of the tasks view

**Non-Goals:**
- Full markdown support (lists, headings, links) — inline only
- Using glamour for individual item rendering

## Decisions

### 1. Word wrap: Width on all items
Apply `.Width(contentWidth)` to all item styles (cursor, done, pending). `contentWidth = m.width - 2` (lateral margin). For items without a background, padding spaces are invisible. Lipgloss automatically wraps content when Width is defined.

Alternative considered: truncate with `...`. Discarded because it loses information.

### 2. Inline markdown: regex over the text, before lipgloss rendering
A transformation pass converts plain text to a string with embedded ANSI before applying the item's lipgloss style. Sequence:

```
item.Text  →  inlineMarkdown(text)  →  lipgloss.Render(styled)
```

Patterns supported in v1:
- `` `code` `` → `codeStyle.Render(match)`
- `**bold**` → `boldStyle.Render(match)`

Processed left-to-right with regexp. Order matters: code first to prevent `**` inside backticks from being processed as bold.

Risk: strings with embedded ANSI + lipgloss `.Width()`. Lipgloss uses `runewidth` and `uniseg` to measure widths — it is aware of ANSI sequences and does not count them in the width. Word wrap should work correctly.

### 3. Global progress bar: first line of the tasks view
Calculated before the main items loop: iterate all `KindTask` items of the change and count `done`. The bar uses the existing `progressBar()` function but with a larger width (10 blocks) and is placed on the first content line, before any section.

## Risks / Trade-offs

- **ANSI + Width**: in extreme cases (many nested styles on a short line) wrap may be imperfect. Acceptable for real openspec usage.
- **Greedy regex**: if a line has multiple `` `code` `` or `**bold**`, the non-greedy regex (`*?`) captures them individually. `*?` must be used, not `*`.
