## MODIFIED Requirements

### Requirement: Specs no seleccionables en el índice
Specs listed in the "Specs" section of `ModeIndex` SHALL be navigable. They SHALL be part of `indexItems` with kind `indexKindSpec`, the cursor SHALL be able to position on them with `j`/`k`, and `Enter` SHALL open `ModeViewingSpec` for the spec under the cursor.

#### Scenario: Cursor enters the Specs section
- **WHEN** the cursor is on the last previously navigable item (last archived) and the user presses `j`
- **THEN** the cursor advances to the first spec in the "Specs" section

#### Scenario: Enter on a spec activates the viewer
- **WHEN** the index cursor is on a spec and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` showing the content of that spec

#### Scenario: Cursor does not go past the last spec
- **WHEN** the cursor is on the last spec and the user presses `j`
- **THEN** the cursor does not move (existing boundary behavior)
