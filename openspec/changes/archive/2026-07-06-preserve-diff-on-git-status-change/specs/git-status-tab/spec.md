## MODIFIED Requirements

### Requirement: Diff view toggle within git changes tab

The TUI SHALL allow the user to view the diff of a file in the git changes tab by pressing `d`, `Enter`, or `e`. For tracked files, the TUI SHALL parse the raw `git diff` output into structured lines with line numbers (`OldNum`/`NewNum`), render them with chroma syntax highlighting, display line numbers (4-char columns), and support horizontal scrolling via `h`/`l` keys (10 runes per step) by truncating raw content before chroma tokenization to avoid ANSI corruption. For untracked files (`??`), the TUI SHALL read the file contents and render them with chroma syntax highlighting. Pressing `d` or `Esc` in the diff view SHALL reset scroll to zero and return to the file list. The TUI SHALL NOT open the system editor from the git changes tab. The diff content SHALL be cleared only when the viewed file itself changes on disk.

#### Scenario: Press d on modified tracked file shows syntax-highlighted diff
- **GIVEN** the git changes tab is open and the cursor is on a modified tracked file
- **WHEN** the user presses `d`
- **THEN** the viewport shows the diff with chroma syntax highlighting, line numbers, and background tints for additions (green) and removals (red)

#### Scenario: Enter on file shows diff
- **GIVEN** the git changes tab is open and the cursor is on a file
- **WHEN** the user presses `Enter`
- **THEN** the diff view is shown for that file (same as pressing `d`)

#### Scenario: e on file shows diff
- **GIVEN** the git changes tab is open and the cursor is on a file
- **WHEN** the user presses `e`
- **THEN** the diff view is shown for that file (same as pressing `d`)

#### Scenario: Press d on untracked file shows syntax-highlighted content
- **GIVEN** the git changes tab is open and the cursor is on an untracked file (`??`)
- **WHEN** the user presses `d`
- **THEN** the viewport shows the file contents with chroma syntax highlighting and line numbers

#### Scenario: Press d or Esc returns to file list
- **GIVEN** the diff view is showing
- **WHEN** the user presses `d` or `Esc`
- **THEN** the viewport returns to the file list and the cursor is on the same file

#### Scenario: Horizontal scroll reveals overflow content
- **GIVEN** the diff view shows a line wider than the viewport
- **WHEN** the user presses `h` or `l` (or arrow keys)
- **THEN** the content shifts left or right by 10 runes, revealing previously hidden code; scroll resets to 0 on return to file list

#### Scenario: Line numbers in diff view
- **GIVEN** the diff view is showing
- **WHEN** the user views the diff
- **THEN** each code line shows its line number (4-char) from the file; added lines show new number, removed lines show old number, context lines show both

#### Scenario: Diff preserved when unrelated files change
- **GIVEN** the diff view is showing for `src/a.go`
- **WHEN** an unrelated file `src/b.go` is modified on disk (detected by polling)
- **THEN** the diff view for `src/a.go` remains visible and the viewport does not reset

#### Scenario: Diff cleared when viewed file itself changes
- **GIVEN** the diff view is showing for `src/a.go`
- **WHEN** `src/a.go` itself is modified on disk (detected by polling)
- **THEN** the diff content is cleared, scroll resets to zero, and the file list is shown

#### Scenario: Diff cleared when viewed file is removed
- **GIVEN** the diff view is showing for `src/a.go`
- **WHEN** `src/a.go` is deleted on disk
- **THEN** the diff content is cleared, scroll resets to zero, and the file list is shown

#### Scenario: j/k scroll within diff view
- **GIVEN** the diff view is showing
- **WHEN** the user presses `j` or `k`
- **THEN** the diff content scrolls vertically via the viewport

#### Scenario: d/Enter/e does nothing on clean working tree
- **GIVEN** the working tree is clean (showing "working tree clean" message)
- **WHEN** the user presses `d`, `Enter`, or `e`
- **THEN** nothing happens
