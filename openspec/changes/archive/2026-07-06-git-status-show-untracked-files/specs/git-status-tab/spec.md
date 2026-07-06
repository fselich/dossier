## MODIFIED Requirements

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
