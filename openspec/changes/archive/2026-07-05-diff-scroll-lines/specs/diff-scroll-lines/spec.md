## MODIFIED Requirements

### Requirement: Diff view toggle within git changes tab

The TUI SHALL allow the user to view the diff of a file in the git changes tab by pressing `d`. For tracked files, the TUI SHALL parse the raw `git diff` output into structured lines with line numbers (`OldNum`/`NewNum`), render them with chroma syntax highlighting, display line numbers (4-char columns), and support horizontal scrolling via `→`/`←` keys (10 runes per step) by truncating raw content before chroma tokenization to avoid ANSI corruption. For untracked files, the TUI SHALL read the file contents and render them with chroma syntax highlighting. Pressing `d` or `Esc` in the diff view SHALL reset scroll to zero and return to the file list. The diff content SHALL be invalidated when git status changes.

#### Scenario: Horizontal scroll reveals overflow content
- **GIVEN** the diff view shows a line wider than the viewport
- **WHEN** the user presses `→`
- **THEN** the content shifts left by 10 runes, revealing previously hidden code

#### Scenario: Scroll resets on return to list
- **GIVEN** the diff view is scrolled horizontally by 20 runes
- **WHEN** the user presses `d` or `Esc`
- **THEN** the view returns to the file list and scroll offset resets to 0

#### Scenario: Line numbers in diff view
- **GIVEN** the diff view is showing
- **WHEN** the user views the diff
- **THEN** each code line shows its line number (4-char) from the file

#### Scenario: Line numbers increment correctly
- **GIVEN** a diff with added and removed lines
- **WHEN** the diff view is showing
- **THEN** added lines show only the new number, removed lines show only the old number, context lines show both
