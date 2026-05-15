## ADDED Requirements

### Requirement: Open the archived changes modal with 'a'
From normal mode, pressing `a` SHALL open an overlay modal listing all changes in `openspec/changes/archive/`. If there are no archived changes, the modal SHALL display a message indicating this. Pressing `a` from `ViewingArchive` mode SHALL have the same effect as `Esc` (return to the picker).

#### Scenario: Open picker with available archived changes
- **WHEN** the mode is `Normal` and there is at least one archived change, and the user presses `a`
- **THEN** the modal is shown with the list of archived changes and the cursor on the first one

#### Scenario: Open picker with no archived changes
- **WHEN** the mode is `Normal` and there are no archived changes, and the user presses `a`
- **THEN** the modal is shown with the message "No hay changes archivados"

#### Scenario: 'a' from ViewingArchive returns to the picker
- **WHEN** the mode is `ViewingArchive` and the user presses `a`
- **THEN** the mode transitions to `ArchivePicker` and the modal is shown

### Requirement: Navigation within the modal
Within the modal, `j` / down arrow SHALL move the cursor to the next item and `k` / up arrow SHALL move it to the previous one. The cursor SHALL NOT go past the first or last element.

#### Scenario: Navigate downward
- **WHEN** the mode is `ArchivePicker`, there are multiple archived changes and the cursor is not on the last one, and the user presses `j`
- **THEN** the cursor moves to the next archived change

#### Scenario: Navigate upward
- **WHEN** the mode is `ArchivePicker` and the cursor is not on the first one, and the user presses `k`
- **THEN** the cursor moves to the previous archived change

#### Scenario: No overflow at the boundaries
- **WHEN** the cursor is on the last archived change and the user presses `j`
- **THEN** the cursor does not change

### Requirement: Name format in the modal
Each list item SHALL display the change name without the date prefix (`YYYY-MM-DD-`) and the date in `DD Mon` format right-aligned within the item, in secondary style.

#### Scenario: Clean name with date
- **WHEN** the archived change has the directory name `2026-05-02-specs-subnav`
- **THEN** the modal shows `specs-subnav` as the name and `02 May` as the date

#### Scenario: Directory name without standard date prefix
- **WHEN** the directory name does not follow the `YYYY-MM-DD-<name>` format
- **THEN** the modal shows the full directory name without a date

### Requirement: Close the modal with Esc
Pressing `Esc` from `ArchivePicker` SHALL close the modal and restore `Normal` mode.

#### Scenario: Close picker
- **WHEN** the mode is `ArchivePicker` and the user presses `Esc`
- **THEN** the mode returns to `Normal` and the modal disappears

### Requirement: Select an archived change with Enter
Pressing `Enter` on a modal item SHALL close the modal, load the selected archived change and transition to `ViewingArchive` mode.

#### Scenario: Select archived change
- **WHEN** the mode is `ArchivePicker`, there is an item under the cursor, and the user presses `Enter`
- **THEN** the mode transitions to `ViewingArchive` and the viewer shows the artifacts of the selected change
