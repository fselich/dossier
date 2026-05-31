# change-index Specification (Delta)

## MODIFIED Requirements

### Requirement: Seleccionar un change con Enter

Pressing `Enter` on an item or left-clicking an already-selected item SHALL close the index and open the selected change. If it is an active change, the mode switches to `ModeNormal` with that active change. If it is an archived change, the mode switches to `ModeViewingArchive` with that archived change.

#### Scenario: Seleccionar change activo

- **WHEN** the cursor is on an active change and the user presses `Enter` or left-clicks on it (when already selected)
- **THEN** the mode switches to `ModeNormal` showing that change

#### Scenario: Seleccionar change archivado

- **WHEN** the cursor is on an archived change and the user presses `Enter` or left-clicks on it (when already selected)
- **THEN** the mode switches to `ModeViewingArchive` showing the artifacts of that archived change
