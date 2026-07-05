## ADDED Requirements

### Requirement: Git status tab visible in change viewer

The TUI SHALL show a fifth tab labeled `changes` in the tab bar when running inside a git worktree. The tab SHALL NOT appear when the project is not inside a git repo. The tab label SHALL include a count of changed files (e.g. `changes (5)`) when there are changes, and SHALL show `changes` (without a count) when the working tree is clean.

#### Scenario: Tab visible inside git repo
- **GIVEN** the project is inside a git worktree
- **WHEN** the TUI starts with at least one active change
- **THEN** the tab bar includes a fifth tab labeled `changes (N)` where N is the number of changed files

#### Scenario: Tab hidden outside git repo
- **GIVEN** the project is NOT inside a git worktree
- **WHEN** the TUI starts
- **THEN** the tab bar does NOT include the `changes` tab

#### Scenario: Working tree clean shows no count
- **GIVEN** the project is inside a git worktree and the working tree is clean
- **WHEN** the user navigates to or views the tab bar
- **THEN** the tab shows `changes` without a numeric suffix

### Requirement: List of changed files with status indicators

The TUI SHALL display files from `git status --porcelain` as a selectable list in the `changes` tab. Each file SHALL show a two-character status code that indicates index (staged) and worktree status. Files SHALL include: modified, added, untracked, renamed, and deleted. Files under the `openspec/` directory SHALL be excluded from the list. When the working tree is clean, the view SHALL show `(working tree clean)`.

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
