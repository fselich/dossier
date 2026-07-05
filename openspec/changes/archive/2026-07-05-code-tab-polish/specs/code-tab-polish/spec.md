## MODIFIED Requirements

### Requirement: Git status tab visible in change viewer

The TUI SHALL show a tab labeled `code` in the tab bar when running inside a git worktree AND there are changed files (`len(gitState.Files) > 0`). The tab SHALL NOT appear when the project is not inside a git repo. The tab SHALL NOT appear in archive mode (`ModeViewingArchive`). The tab SHALL be shown as disabled when the working tree is clean.

#### Scenario: Tab labeled code in normal mode
- **GIVEN** the project is inside a git worktree with changed files
- **WHEN** the TUI is in normal mode
- **THEN** the tab bar shows a tab labeled `code (N)` where N is the number of changed files

#### Scenario: Tab hidden in archive mode
- **GIVEN** the user is viewing an archived change
- **WHEN** the tab bar is rendered
- **THEN** the `code` tab is not shown at all
