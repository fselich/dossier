## ADDED Requirements

### Requirement: Real-time index refresh
While the mode is `ModeIndex`, the TUI SHALL detect on every tick (≤ 500 ms) whether the list of active changes, the list of archived changes, or the list of project specs has changed on disk. If any change is detected, the index SHALL reload all three lists, rebuild the navigable items, and refresh the viewport without the user having to exit and re-enter `ModeIndex`. The cursor SHALL be preserved if the resulting index has at least as many items as the current position; otherwise it SHALL move to the last available item.

#### Scenario: New spec appears on disk while the index is open
- **WHEN** the mode is `ModeIndex` and a new directory is created in `openspec/specs/`
- **THEN** within a maximum of 500 ms the index shows the new spec in the "Specifications" section without user intervention

#### Scenario: Spec disappears from specs while the index is open
- **WHEN** the mode is `ModeIndex` and a directory is deleted from `openspec/specs/`
- **THEN** within a maximum of 500 ms the spec disappears from the "Specifications" section

#### Scenario: New archived change while the index is open
- **WHEN** the mode is `ModeIndex` and a change is moved to `openspec/changes/archive/`
- **THEN** within a maximum of 500 ms the change appears in the "Archived Changes" section

#### Scenario: New active change while the index is open
- **WHEN** the mode is `ModeIndex` and a new change is created in `openspec/changes/`
- **THEN** within a maximum of 500 ms the change appears in the "Active Changes" section

#### Scenario: Cursor preserved when the item still exists
- **WHEN** the index is reloaded and the number of items does not drop below the cursor position
- **THEN** the cursor remains at the same numeric position

#### Scenario: Cursor adjusted when the item disappears
- **WHEN** the index is reloaded and the number of items is less than the current cursor position
- **THEN** the cursor moves to the last available item
