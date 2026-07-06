# git-status-tab Specification

## Purpose

Provides a `code` tab in the TUI that shows working-tree file changes from `git status --porcelain`, enabling the user to browse modified/added/untracked/renamed/deleted files and view diffs with syntax highlighting without leaving the TUI.

## Requirements

### Requirement: Git status tab visible in change viewer

The TUI SHALL show a tab labeled `code` in the tab bar when running inside a git worktree AND there are changed files (`len(gitState.Files) > 0`). The tab SHALL NOT appear when the project is not inside a git repo. The tab SHALL NOT appear in archive mode (`ModeViewingArchive`). The tab SHALL be shown as disabled (grayed out) when the working tree is clean. The tab label SHALL include a count of changed files (e.g. `code (5)`) when there are changes, and SHALL show `code` (without a count) when the working tree is clean.

#### Scenario: Tab visible inside git repo
- **GIVEN** the project is inside a git worktree
- **WHEN** the TUI starts with at least one active change
- **THEN** the tab bar includes a tab labeled `code (N)` where N is the number of changed files

#### Scenario: Tab hidden outside git repo
- **GIVEN** the project is NOT inside a git worktree
- **WHEN** the TUI starts
- **THEN** the tab bar does NOT include the `code` tab

#### Scenario: Tab disabled when working tree is clean
- **GIVEN** the project is inside a git worktree with no changed files
- **WHEN** the TUI is in normal mode
- **THEN** the `code` tab is shown as disabled (grayed out) and is not selectable

#### Scenario: Tab hidden in archive mode
- **GIVEN** the user is viewing an archived change
- **WHEN** the tab bar is rendered
- **THEN** the `code` tab is not shown at all

#### Scenario: Tab becomes enabled when files appear
- **GIVEN** the code tab is disabled with a clean working tree
- **WHEN** a file is modified on disk (detected by polling)
- **THEN** within a maximum of 500 ms the tab becomes enabled and shows the file count
- **THEN** within a maximum of 500 ms the tab becomes enabled and shows the file count

#### Scenario: Working tree clean shows no count
- **GIVEN** the project is inside a git worktree and the working tree is clean
- **WHEN** the user navigates to or views the tab bar
- **THEN** the tab shows `changes` without a numeric suffix

### Requirement: List of changed files with status indicators

The TUI SHALL display files from `git status --porcelain -u` as a selectable list in the `changes` tab. Each file SHALL show a two-character status code that indicates index (staged) and worktree status. Files SHALL include: modified, added, untracked, renamed, and deleted. Files under the `openspec/` directory SHALL be excluded from the list. When the working tree is clean, the view SHALL show `(working tree clean)`. Untracked files in new directories SHALL appear as individual file entries (e.g., `?? src/Domain/Cache/CacheInterface.php`) rather than being collapsed into a directory entry (`?? src/Domain/Cache/`).

#### Scenario: Shows modified, added, untracked, renamed, deleted
- **GIVEN** the working tree has modified (`M`), added (`A`), untracked (`??`), renamed (`R`), and deleted (`D`) files
- **WHEN** the user opens the `changes` tab
- **THEN** all five types appear in the list with their corresponding status codes

#### Scenario: Excludes files under openspec/ directory
- **GIVEN** there are changes inside `openspec/` and outside `openspec/`
- **WHEN** the user opens the `changes` tab
- **THEN** only files outside `openspec/` appear in the list

#### Scenario: Working tree clean
- **GIVEN** the working tree has no changed files
- **WHEN** the user opens the `changes` tab
- **THEN** the view shows `(working tree clean)`

#### Scenario: Untracked files in new directories show individually
- **GIVEN** the working tree has untracked files inside a new directory `src/Domain/Cache/`
- **WHEN** the user opens the `changes` tab
- **THEN** each file appears individually (e.g., `?? src/Domain/Cache/CacheInterface.php`) instead of a single `?? src/Domain/Cache/` entry

### Requirement: Cursor navigation skips deleted files

The TUI SHALL support cursor navigation with `j`/`k` in the `changes` tab. Deleted files SHALL be shown in the list (dimmed) but the cursor SHALL skip them when navigating. Pressing `Enter` or `e` on a deleted file SHALL do nothing.

#### Scenario: j/k navigate past deleted files
- **GIVEN** the list has files at indices 0 (modified), 1 (deleted), and 2 (added)
- **WHEN** the user presses `j` twice from index 0
- **THEN** the cursor lands on index 2, skipping the deleted file

#### Scenario: Enter on deleted file does nothing
- **GIVEN** the cursor is on a deleted file
- **WHEN** the user presses `Enter`
- **THEN** no file is opened and no error occurs

### Requirement: Open changed file in system editor

The TUI SHALL open the selected file in the system editor (`$EDITOR`, fallback `vi`) when the user presses `Enter` or `e` on a non-deleted file. For renamed files, the new path SHALL be opened.

#### Scenario: Enter on modified file opens editor
- **GIVEN** the cursor is on a modified file with path `internal/ui/model.go`
- **WHEN** the user presses `Enter`
- **THEN** `$EDITOR /abs/path/to/internal/ui/model.go` is launched via `tea.ExecProcess`

#### Scenario: e on modified file opens editor
- **GIVEN** the cursor is on a modified file
- **WHEN** the user presses `e`
- **THEN** the file is opened in the system editor (same as `Enter`)

#### Scenario: Renamed file opens new path
- **GIVEN** the cursor is on a renamed file with old path `old.go` and new path `new.go`
- **WHEN** the user presses `Enter`
- **THEN** the new path `/abs/path/to/new.go` is opened in the editor

#### Scenario: Fallback to vi when $EDITOR is unset
- **GIVEN** `$EDITOR` is not defined
- **WHEN** the user presses `Enter` on a file
- **THEN** `vi` is launched with the file path

### Requirement: Git status polling

The TUI SHALL poll `git status --porcelain` on the same 500ms tick used for artifact detection, but only when inside a git worktree and in `ModeNormal`. Changes discovered by polling SHALL update the file list immediately if the `changes` tab is visible. The file list SHALL also be available when switching to the tab (no extra delay).

#### Scenario: Polling detects new file
- **GIVEN** the `changes` tab is open and the list shows 3 files
- **WHEN** a new file is created on disk
- **THEN** within a maximum of 500 ms the list shows 4 files with the new entry

#### Scenario: Switching to changes tab shows current state
- **GIVEN** the user is on the `proposal` tab and a file was modified
- **WHEN** the user switches to the `changes` tab
- **THEN** the modified file appears in the list without waiting for the next tick

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
