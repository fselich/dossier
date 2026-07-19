## Context

The index view (`renderActiveItem`, `renderArchivedItem`, spec rendering in `renderIndexContent`) renders non-selected item names as plain unstyled text. When the active theme sets a background color (e.g., light theme: `#ffffff`), unstyled text uses the terminal's default foreground, which may lack contrast against the themed background. Selected items use `IndexActive` (blue background + white text), so they are always visible.

The `ThemeColors` struct already defines `PrimaryFg` — the theme's primary text color (`"15"` for dark, `"0"` for light). But no style in `ThemeStyles` exposes `PrimaryFg` as a standalone base text style.

## Goals / Non-Goals

**Goals:**
- Add `BaseText` style to `ThemeStyles` applying `PrimaryFg` (no bold, no background)
- Apply it to non-selected item names in index view: active changes, archived changes, spec names

**Non-Goals:**
- Change rendering of requirements (already styled with `TaskPending`)
- Change rendering of section headers, cursor marks, or progress bars
- Add new color roles to `ThemeColors`

## Decisions

### 1. New `BaseText` style in `ThemeStyles`

One new field:

```go
BaseText lipgloss.Style  // PrimaryFg (no bold, no background)
```

Built with:
```go
BaseText: lipgloss.NewStyle().Foreground(c.PrimaryFg),
```

### 2. Three call sites in index.go

| Function | Line | Current | New |
|----------|------|---------|-----|
| `renderActiveItem` | 505 | `renderedName = paddedName` | `renderedName = m.theme.Styles.BaseText.Render(paddedName)` |
| `renderArchivedItem` | 539 | `name = ch.Name + pad` | `name = m.theme.Styles.BaseText.Render(ch.Name) + pad` |
| `renderIndexContent` (spec) | 418 | `name = ps.Name` | `name = m.theme.Styles.BaseText.Render(ps.Name)` |

## Risks / Trade-offs

- **[Zero] Dark theme regression**: `PrimaryFg` = `"15"` (white) — same as terminal default on most dark terminals. No visual change.
- **[Zero] Performance**: One extra lipgloss `Render()` call per index item. Lipgloss styles are value types and rendering is fast.
