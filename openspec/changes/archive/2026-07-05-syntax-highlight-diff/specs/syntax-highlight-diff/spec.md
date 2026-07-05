## MODIFIED Requirements

### Requirement: Diff view toggle within git changes tab

The TUI SHALL allow the user to view the diff of a file in the git changes tab by pressing `d`. For tracked files (modified, added, renamed, deleted), the TUI SHALL parse the raw `git diff` output into structured lines and render them with chroma syntax highlighting: each code line SHALL be tokenized by chroma's language lexer (detected via `lexers.Match(filename)`), and SHALL display the chroma foreground colors for syntax while applying a background tint to distinguish additions (green tint), removals (red tint), and context (no tint). Diff indicators (`+` in green, `-` in red) and hunk headers (`@@ ... @@` in cyan) SHALL be color-coded. For untracked files (`??`), the TUI SHALL read the file contents and render them with chroma syntax highlighting. Pressing `d` or `Esc` in the diff view SHALL return to the file list. The diff content SHALL be invalidated when git status changes.

#### Scenario: Diff shows syntax-highlighted code
- **GIVEN** the git changes tab is open and the cursor is on a modified Go file
- **WHEN** the user presses `d`
- **THEN** the diff view shows the code lines with Go syntax highlighting (keywords, strings, types in different colors) and diff indicators (+ in green, - in red)

#### Scenario: Diff highlights multiple languages correctly
- **GIVEN** the git changes tab has files of different types (Go, Python, YAML, Markdown)
- **WHEN** the user presses `d` on each file
- **THEN** each file's diff uses the correct language lexer for syntax highlighting

#### Scenario: Added lines have green background tint
- **GIVEN** the diff view is showing
- **WHEN** there are added lines (`+` prefix) in the diff
- **THEN** those lines appear with a subtle green background tint and chroma foreground colors

#### Scenario: Removed lines have red background tint
- **GIVEN** the diff view is showing
- **WHEN** there are removed lines (`-` prefix) in the diff
- **THEN** those lines appear with a subtle red background tint and chroma foreground colors
