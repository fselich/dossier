## ADDED Requirements

### Requirement: Display artifacts of an archived change
In `ViewingArchive` mode, the TUI SHALL show the artifacts of the selected archived change using the same visual structure as active changes: header, tab bar, separator, content. Keys `1`–`4`, `j`/`k` and `h`/`l` SHALL work the same as in normal mode for navigating between artifacts and scrolling.

#### Scenario: Navigate artifacts of an archived change
- **WHEN** the mode is `ViewingArchive` and the user presses `2`
- **THEN** the active tab changes to `design` and the viewport shows the archived change's design content

#### Scenario: Scroll within an archived change
- **WHEN** the mode is `ViewingArchive` and the user presses `j`
- **THEN** the viewport scrolls down one line

#### Scenario: h/l does not switch change in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `h` or `l`
- **THEN** nothing changes (there is no lateral navigation between archived changes)

### Requirement: Visual indicator for archive mode
When the mode is `ViewingArchive`, the header SHALL display the text `[archivo]` instead of the usual position indicator `[N/M]`.

#### Scenario: Header in archive mode
- **WHEN** the mode is `ViewingArchive`
- **THEN** the header shows `<project>  ·  <archived-name>  [archivo]`

### Requirement: Read-only in archive mode
In `ViewingArchive` mode, keys `e` (open editor) and `Space` (task toggle) SHALL be silently ignored.

#### Scenario: 'e' ignored in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `e`
- **THEN** no editor is opened and state does not change

#### Scenario: 'Space' ignored in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `Space`
- **THEN** no task changes state

### Requirement: Adapted helpbar in archive mode
In `ViewingArchive` mode, the helpbar SHALL show the actually available keys, omitting `e` and `Space`, and including `Esc: back`.

#### Scenario: Read-only helpbar
- **WHEN** the mode is `ViewingArchive`
- **THEN** the helpbar shows `1-4: artifact   j/k: scroll   a/Esc: volver`

### Requirement: Return to picker with Esc
In `ViewingArchive` mode, pressing `Esc` SHALL close the archived viewer and return to `ArchivePicker` mode, with the cursor on the item that had been selected.

#### Scenario: Esc returns to the picker
- **WHEN** the mode is `ViewingArchive` and the user presses `Esc`
- **THEN** the mode transitions to `ArchivePicker` and the modal reappears with the cursor on the archived change that was being viewed
