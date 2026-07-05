## MODIFIED Requirements

### Requirement: Git status tab visible in change viewer

The TUI SHALL show a fifth tab labeled `changes` in the tab bar when running inside a git worktree AND there are changed files (`len(gitState.Files) > 0`). The tab SHALL NOT appear when the project is not inside a git repo or when the working tree is clean. The tab label SHALL include a count of changed files (e.g. `changes (5)`) when there are changes.

#### Scenario: Tab disabled when working tree is clean
- **GIVEN** the project is inside a git worktree with no changed files
- **WHEN** the TUI is in normal mode
- **THEN** the `changes` tab is shown as disabled (grayed out)

#### Scenario: Tab becomes enabled when files appear
- **GIVEN** the changes tab is disabled with a clean working tree
- **WHEN** a file is modified on disk (detected by polling)
- **THEN** within a maximum of 500 ms the tab becomes enabled

### Requirement: Git status tab key actions show diff

The TUI SHALL show the diff view when the user presses `Enter` or `e` on a file in the git changes tab, matching the behavior of the `d` key. The TUI SHALL NOT open the system editor from the git changes tab.

#### Scenario: Enter on file shows diff
- **GIVEN** the git changes tab is open and the cursor is on a file
- **WHEN** the user presses `Enter`
- **THEN** the diff view is shown for that file (same as pressing `d`)

#### Scenario: e on file shows diff
- **GIVEN** the git changes tab is open and the cursor is on a file
- **WHEN** the user presses `e`
- **THEN** the diff view is shown for that file (same as pressing `d`)

#### Scenario: Enter on clean working tree does nothing
- **GIVEN** the git changes tab shows "working tree clean"
- **WHEN** the user presses `Enter` or `e`
- **THEN** nothing happens
