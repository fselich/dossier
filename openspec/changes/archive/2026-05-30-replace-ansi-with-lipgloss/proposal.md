## Why

`inlineMarkdown()` in `tasks.go` uses raw ANSI escape codes (`\033[4m`, `\033[1m`, `\033[0m`, etc.) for underline, bold, and reset styling. The project already depends on Lipgloss, which provides a cleaner, more maintainable API for these styles. Raw ANSI strings are error-prone and harder to read.

## What Changes

- Replace raw ANSI sequences in `inlineMarkdown()` with Lipgloss styles.
- Define `underlineStyle` for code/underline rendering.
- Define `boldStyle` for bold rendering.
- Use Lipgloss reset via style inheritance instead of `\033[0m`.

## Capabilities

<!-- No spec changes — implementation detail only. -->

## Impact

- `internal/ui/tasks.go`: `inlineMarkdown()` rewritten to use Lipgloss styles.
- No rendering changes expected; visual output should be identical.
