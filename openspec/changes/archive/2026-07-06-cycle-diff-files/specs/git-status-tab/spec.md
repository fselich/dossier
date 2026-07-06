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

#### Scenario: Press `[` cycles to previous file diff
- **GIVEN** the diff view is showing for file at cursor position N
- **WHEN** the user presses `[`
- **THEN** the cursor moves to the previous non-deleted file (N-1), the diff for that file is loaded, and the diff view remains active

#### Scenario: Press `]` cycles to next file diff
- **GIVEN** the diff view is showing for file at cursor position N
- **WHEN** the user presses `]`
- **THEN** the cursor moves to the next non-deleted file (N+1), the diff for that file is loaded, and the diff view remains active

#### Scenario: Cycling wraps around at first file
- **GIVEN** the diff view is showing for the first file (cursor = 0)
- **WHEN** the user presses `[`
- **THEN** the cursor wraps to the last non-deleted file, and its diff is loaded

#### Scenario: Cycling wraps around at last file
- **GIVEN** the diff view is showing for the last file
- **WHEN** the user presses `]`
- **THEN** the cursor wraps to the first non-deleted file, and its diff is loaded

#### Scenario: Cycling skips deleted files
- **GIVEN** the diff view is showing for file N, and file N+2 is the next non-deleted file (N+1 is deleted)
- **WHEN** the user presses `]`
- **THEN** the cursor skips N+1 and lands on N+2, and its diff is loaded

#### Scenario: Scroll resets when cycling to a new file
- **GIVEN** the diff view for file A has been scrolled horizontally (`ScrollX > 0`)
- **WHEN** the user presses `]` to cycle to file B
- **THEN** the horizontal scroll resets to 0 for the new diff

#### Scenario: Cycling is only available in diff view
- **GIVEN** the file list is showing (not in diff view)
- **WHEN** the user presses `[` or `]`
- **THEN** nothing happens

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

#### Scenario: Diff cleared when viewed file itself changes
- **GIVEN** the diff view is showing for `src/a.go`
- **WHEN** `src/a.go` itself is modified on disk (detected by polling)
- **THEN** the diff content is cleared, scroll resets to zero, and the file list is shown

#### Scenario: j/k scroll within diff view
- **GIVEN** the diff view is showing
- **WHEN** the user presses `j` or `k`
- **THEN** the diff content scrolls vertically via the viewport

#### Scenario: d/Enter/e does nothing on clean working tree
- **GIVEN** the working tree is clean (showing "working tree clean" message)
- **WHEN** the user presses `d`, `Enter`, or `e`
- **THEN** nothing happens
