# change-index Specification (Delta)

## MODIFIED Requirements

### Requirement: Actualización en tiempo real del índice

While the mode is `ModeIndex`, the TUI SHALL detect on each tick (≤ 500 ms) whether the list of active changes, the list of archived changes, or the list of project specs has changed on disk. If any structural change is detected, the index SHALL reload all three lists, rebuild the navigable items, and refresh the viewport without the user having to leave and re-enter `ModeIndex`. Additionally, when no structural change is detected, the TUI SHALL reload the task content of each active change from disk and, if any task content has changed, SHALL rebuild the index items and refresh the viewport so that progress bars reflect the latest task completion state. The cursor SHALL be preserved if the resulting index has at least as many items as the current position; otherwise it SHALL move to the last available item.

*Note: This requirement is re-listed without behavior changes. The tick handler implementation was restructured into `pollIndexMode()`, `pollNormalModeChanges()`, and `pollNormalModeContent()` methods extracted from `handleTick()`. All scenarios remain identical.*

#### Scenario: Nuevo spec aparece en disco mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and a new directory is created in `openspec/specs/`
- **THEN** within a maximum of 500 ms the index shows the new spec in the "Specifications" section without user intervention

#### Scenario: Spec archivado desaparece de specs mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and a directory is deleted from `openspec/specs/`
- **THEN** within a maximum of 500 ms the spec disappears from the "Specifications" section

#### Scenario: Nuevo change archivado mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and a change is moved to `openspec/changes/archive/`
- **THEN** within a maximum of 500 ms the change appears in the "Archived Changes" section

#### Scenario: Nuevo change activo mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and a new change is created in `openspec/changes/`
- **THEN** within a maximum of 500 ms the change appears in the "Active Changes" section

#### Scenario: Cursor preservado cuando el ítem sigue existiendo

- **WHEN** the index reloads and the number of items does not decrease below the cursor position
- **THEN** the cursor stays at the same numeric position

#### Scenario: Cursor reajustado cuando el ítem desaparece

- **WHEN** the index reloads and the number of items is less than the current cursor position
- **THEN** the cursor moves to the last available item

#### Scenario: Tareas actualizadas en disco mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and the `tasks.md` file of an active change is externally modified (e.g., a checkbox is toggled)
- **THEN** within a maximum of 500 ms the progress bar for that change in the index reflects the updated `done/total` count without user intervention
