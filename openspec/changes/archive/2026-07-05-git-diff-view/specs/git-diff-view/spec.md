## MODIFIED Requirements

### Requirement: Diff view toggle within git changes tab

The TUI SHALL allow the user to view the diff of a file in the git changes tab by pressing `d`. For tracked files (modified, added, renamed, deleted), the TUI SHALL show the output of `git diff HEAD --color=always -- <file>` using ANSI-colored diff output. For untracked files (`??`), the TUI SHALL read the file contents and render them with chroma syntax highlighting detected via `lexers.Match(filename)`, with a header indicating it is a new file. Pressing `d` or `Esc` in the diff view SHALL return to the file list. The diff content SHALL be invalidated when git status changes.

#### Scenario: Press d on modified tracked file
- **GIVEN** the git changes tab is open and the cursor is on a modified tracked file
- **WHEN** the user presses `d`
- **THEN** the viewport shows the colored `git diff HEAD --color=always` output for that file

#### Scenario: Press d on untracked file shows syntax-highlighted content
- **GIVEN** the git changes tab is open and the cursor is on an untracked file (`??`)
- **WHEN** the user presses `d`
- **THEN** the viewport shows the file contents with chroma syntax highlighting and a header indicating it is a new file

#### Scenario: Press d or Esc returns to file list
- **GIVEN** the diff view is showing
- **WHEN** the user presses `d` or `Esc`
- **THEN** the viewport returns to the file list and the cursor is on the same file

#### Scenario: Diff cache invalidated on status change
- **GIVEN** the diff view is showing for a file
- **WHEN** the git status changes (e.g., file is modified externally)
- **THEN** the diff content is cleared and the file list is shown again

#### Scenario: j/k scroll within diff view
- **GIVEN** the diff view is showing
- **WHEN** the user presses `j` or `k`
- **THEN** the diff content scrolls normally via the viewport

#### Scenario: d does nothing on clean working tree
- **GIVEN** the working tree is clean (showing "working tree clean" message)
- **WHEN** the user presses `d`
- **THEN** nothing happens
