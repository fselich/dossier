## Context

`inlineMarkdown()` in `tasks.go:95-112` uses raw ANSI escape sequences to style task text:
- `\033[4m` / `\033[24m` for underline (code spans in done tasks)
- `\033[1m` / `\033[22m` for bold
- `\033[36m` for cyan (code spans in pending tasks)
- `\033[0m` for reset

The project already uses Lipgloss for all other TUI styling. Raw ANSI strings are opaque, error-prone, and inconsistent with the rest of the codebase.

## Goals / Non-Goals

**Goals:**
- Replace all raw ANSI codes in `inlineMarkdown()` with Lipgloss style objects.
- Visual rendering must remain unchanged.

**Non-Goals:**
- Refactoring `inlineMarkdown()` logic or regex patterns.
- Adding new styling capabilities.

## Decisions

**Define package-level `underlineStyle` and `boldStyle`**
```go
var (
    underlineStyle = lipgloss.NewStyle().Underline(true)
    boldStyle      = lipgloss.NewStyle().Bold(true)
    cyanStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)
```

For the code-span case in pending tasks, instead of `\033[36m...\033[0m` + restore string, use `cyanStyle.Render(...)` + restore.

**Inline style rendering**
Rather than calling `.Render()` on each substring (which returns a full styled string), use the style's start/end sequences: `underlineStyle.GetUnderline()` isn't directly exposed. Instead, we render the substring with the style and concatenate: `underlineStyle.Render(code)`.

**Reset handling**: When a Lipgloss style renders, it auto-resets at the end via `\033[0m`. This matches the current behavior. The `restore` parameter (which reapplies the parent task style after the inline markdown) is appended after the Lipgloss output, same as today.

## Risks / Trade-offs

- **Style output diff** → Mitigation: visual comparison before/after. Lipgloss `Underline(true).Render("x")` produces `\033[4mx\033[0m`, identical to the current raw ANSI. Same for Bold.
- **Cyan color code difference** → `\033[36m` (ANSI cyan) maps to Lipgloss `Color("6")`. Verify output matches.
