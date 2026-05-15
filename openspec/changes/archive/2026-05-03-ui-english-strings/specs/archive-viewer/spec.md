## MODIFIED Requirements

### Requirement: Visual archive mode indicator
When the mode is `ViewingArchive`, the header SHALL show the text `[archive]` instead of the usual position indicator `[N/M]`.

#### Scenario: Header in archive mode
- **WHEN** the mode is `ViewingArchive`
- **THEN** the header shows `<project>  ·  <archive-name>  [archive]`

### Requirement: Adapted helpbar in archive mode
In `ViewingArchive` mode, the helpbar SHALL show the actual available keys, omitting `e` and `Space`, and including `Esc: index`.

#### Scenario: Read-only helpbar
- **WHEN** the mode is `ViewingArchive`
- **THEN** the helpbar shows `1-4: artifact  j/k: scroll  a/Esc: index  q: quit`
