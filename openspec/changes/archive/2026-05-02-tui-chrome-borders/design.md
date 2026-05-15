## Context

The current layout in `View()` concatenates four zones with `\n` without any visual separator:

```
header
tabBar
viewport
[globalProgressBar]   ← only in tasks tab
helpBar
```

The zones are distinguished only by the colour/style of their content, but without explicit delimitation.

## Goals / Non-Goals

**Goals:**
- Horizontal separator between tabBar and viewport
- Horizontal separator between viewport (or globalProgressBar) and helpBar
- Separators must span the full width of the terminal

**Non-Goals:**
- A box with side borders around the content (complicates text wrapping)
- Changing the layout or the order of zones
- Adding extra horizontal padding

## Decisions

### 1. Separator lines with `─` at full width

Add a line of `strings.Repeat("─", m.width)` styled with a subtle colour (dark grey, same as `helpStyle`) between zones. This is the cleanest option: a single line, does not occupy visible height intrusively, and is compatible with any terminal width.

Considered alternative: `╔═══╗`-style border with sides. Discarded because the sides add a margin character on every content line and complicate the calculation of `contentWidth`.

### 2. Impact on `contentHeight()`

Each separator occupies 1 line. Two separators are added (one above the viewport, one below), so `contentHeight()` goes from `m.height - 3` to `m.height - 5` (or `m.height - 6` in the tasks tab with the global bar). The separators are rendered outside the viewport, just like the header and helpBar.

### 3. Separator style

`separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))` — black/dark grey. Subtle, does not compete with the content. Color "0" is sufficiently visible on a dark terminal background without being distracting.

## Risks / Trade-offs

- **Lost height**: −2 lines of visible content. On very small terminals (< 10 lines) the content becomes very compressed. Acceptable for real-world use.
- **`contentHeight()` more complex**: now depends on the active tab and the number of separators. Must be kept in sync with `View()`.
