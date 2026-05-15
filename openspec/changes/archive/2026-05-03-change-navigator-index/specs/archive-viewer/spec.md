## MODIFIED Requirements

### Requirement: Return to the index with Esc or 'a'
In `ViewingArchive` mode, pressing `Esc` or `a` SHALL close the archived viewer and return to `ModeIndex`.

#### Scenario: Esc returns to the index
- **WHEN** the mode is `ViewingArchive` and the user presses `Esc`
- **THEN** the mode transitions to `ModeIndex`

#### Scenario: 'a' returns to the index
- **WHEN** the mode is `ViewingArchive` and the user presses `a`
- **THEN** the mode transitions to `ModeIndex`

### Requirement: Adapted helpbar in archive mode
In `ViewingArchive` mode, the helpbar SHALL display the actual available keys, omitting `e` and `Space`, and including `Esc: index`.

#### Scenario: Read-only helpbar
- **WHEN** the mode is `ViewingArchive`
- **THEN** the helpbar shows `1-4: artifact   j/k: scroll   a/Esc: index   q: quit`
