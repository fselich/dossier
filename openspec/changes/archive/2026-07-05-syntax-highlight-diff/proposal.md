## Why

The current diff view uses `git diff --color=always` which only colors diff markers (`+`/`-`/`@@`), leaving the code itself unhighlighted. For a developer reading diffs in the TUI, syntax highlighting of the code inside the diff makes a significant difference in readability — keywords, strings, comments, and types are visually distinct instead of being uniform monochrome text.

The `differ` project (github.com/JanSmrcka/differ) demonstrates the technique: parse the raw git diff into structured lines, then tokenize each line individually with chroma's language lexer, preserving only the foreground colors (syntax) while applying a background tint based on the diff type (added/removed/context).

## What Changes

- Replace `git diff --color=always` with raw `git diff` + custom parser.
- Parse diff output into structured `DiffLine` structs (type: added, removed, context, hunk header).
- For code lines, strip the `+`/`-`/` ` prefix and tokenize with chroma's language lexer (detected via `lexers.Match(filename)`).
- Apply chroma foreground colors (syntax highlighting) + line-type background tint (green for additions, red for removals).
- Cache chroma lexer per file extension to avoid repeated `lexers.Match` calls.
- Keep the untracked file viewer (already uses chroma).

## Capabilities

### Modified Capabilities

- `git-status-tab`: The diff view now displays syntax-highlighted code within diffs, matching the file's language. The diff parser replaces the raw `git diff --color=always` output.

## Impact

- Refactor `internal/ui/gitdiff.go`: replace `computeDiff` with diff parser + chroma-based renderer.
- Add diff parsing helpers (`parseDiff`, `parseDiffLine`) and chroma rendering (`highlightLine`).
- No new dependencies — chroma already a direct dependency.
- No changes to `model.go`, `view.go`, `viewer.go`, `git.go`, or OpenSpec domain logic.
