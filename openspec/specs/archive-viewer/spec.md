# archive-viewer Specification

## Purpose
Defines the `ModeViewingArchive` mode for viewing artifacts of an archived change in read-only mode, with the same visual structure as normal mode but without editing or task toggling.

## Requirements


### Requirement: View artifacts of an archived change
In `ViewingArchive` mode, the TUI SHALL display the artifacts of the selected archived change using the same visual structure as active changes: header, tab bar, separator, content. Keys `1`-`4`, `j`/`k` and `h`/`l` SHALL work the same as in normal mode to navigate between artifacts and scroll.

#### Scenario: Navigate artifacts of an archived change
- **WHEN** the mode is `ViewingArchive` and the user presses `2`
- **THEN** the active tab changes to `design` and the viewport shows the content of the archive's design

#### Scenario: Scroll in an archived change
- **WHEN** the mode is `ViewingArchive` and the user presses `j`
- **THEN** the viewport scrolls down one line

#### Scenario: h/l does not change the change in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `h` or `l`
- **THEN** nothing changes (there is no lateral navigation between archived items)

### Requirement: Visual indicator for archive mode
When the mode is `ViewingArchive`, the header SHALL display the text `[archive]` instead of the usual position indicator `[N/M]`.

#### Scenario: Header in archive mode
- **WHEN** the mode is `ViewingArchive`
- **THEN** the header shows `<project>  ·  <archive-name>  [archive]`

### Requirement: Read-only in archive mode
In `ViewingArchive` mode, keys `e` (open editor) and `Space` (task toggle) SHALL be silently ignored.

#### Scenario: 'e' ignored in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `e`
- **THEN** no editor is opened and the state does not change

#### Scenario: 'Space' ignored in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `Space`
- **THEN** no task changes its state

### Requirement: Helpbar adaptado en modo archivo
In `ViewingArchive` mode, the helpbar SHALL show the actual available keys, omitting `e` and `Space`, and including `Esc: index`.

#### Scenario: Helpbar de solo lectura
- **WHEN** the mode is `ViewingArchive`
- **THEN** the helpbar shows `1-4: artifact  j/k: scroll  a/Esc: index  q: quit`

### Requirement: Volver al índice con Esc o 'a'
In `ViewingArchive` mode, pressing `Esc` or `a` SHALL close the archive viewer and return to `ModeIndex`.

#### Scenario: Esc vuelve al índice
- **WHEN** the mode is `ViewingArchive` and the user presses `Esc`
- **THEN** the mode switches to `ModeIndex`

#### Scenario: 'a' vuelve al índice
- **WHEN** the mode is `ViewingArchive` and the user presses `a`
- **THEN** the mode switches to `ModeIndex`
