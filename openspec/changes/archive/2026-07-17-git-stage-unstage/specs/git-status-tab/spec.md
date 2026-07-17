## MODIFIED Requirements

### Requirement: Cursor navigation skips deleted files

The TUI SHALL support cursor navigation with `j`/`k` in the `changes` tab. Deleted files SHALL be shown in the list (dimmed) and the cursor SHALL be able to land on them, so that their deletion can be staged. Pressing `Enter`, `e`, or `d` on a deleted file SHALL do nothing. Cycling between diffs with `[`/`]` inside the diff view SHALL continue to skip deleted files.

#### Scenario: j/k land on deleted files
- **GIVEN** the list has files at indices 0 (modified), 1 (deleted), and 2 (added)
- **WHEN** the user presses `j` once from index 0
- **THEN** the cursor lands on index 1 (the deleted file)

#### Scenario: Enter on deleted file does nothing
- **GIVEN** the cursor is on a deleted file
- **WHEN** the user presses `Enter`
- **THEN** no file is opened and no error occurs

#### Scenario: d on deleted file does nothing
- **GIVEN** the cursor is on a deleted file
- **WHEN** the user presses `d`
- **THEN** no diff view is opened and no error occurs

#### Scenario: Diff cycling still skips deleted files
- **GIVEN** the diff view is showing for file N and file N+1 is deleted
- **WHEN** the user presses `]`
- **THEN** the cursor skips N+1 and lands on the next non-deleted file
