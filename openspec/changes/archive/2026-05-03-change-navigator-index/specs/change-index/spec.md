## ADDED Requirements

### Requirement: Full-screen index view
The TUI SHALL implement a `ModeIndex` mode that occupies the full screen with the same TUI chrome (borders, header, helpbar). The index SHALL display two sections: "Active" with the active changes and "Archived" with the changes in `openspec/changes/archive/`, both separated by a section separator line.

#### Scenario: Index with active and archived changes
- **WHEN** the mode is `ModeIndex` and both active and archived changes exist
- **THEN** the screen shows an "Active" section with the active changes followed by an "Archived" section with the archived changes, within the TUI chrome

#### Scenario: Index with no active changes
- **WHEN** the mode is `ModeIndex` and there are no active changes
- **THEN** the "Active" section shows a message indicating there are no active changes

#### Scenario: Index with no archived changes
- **WHEN** the mode is `ModeIndex` and there are no archived changes
- **THEN** the "Archived" section shows a message indicating there are no archived changes

### Requirement: Active change format in the index
Each active change SHALL be shown with its name on the left and a progress bar `[█░] N/M` on the right, using the same bar style as the tab bar. The item under the cursor SHALL be visually highlighted.

#### Scenario: Active change with partial progress
- **WHEN** an active change has 6 out of 10 tasks completed and is under the cursor
- **THEN** it is displayed as `▶ change-name  [██████░░░░] 6/10` with highlighted style

#### Scenario: Active change with no tasks
- **WHEN** an active change has no `tasks.md`
- **THEN** the name is shown without a progress bar

### Requirement: Archived change format in the index
Each archived change SHALL be shown with the date `DD/MM/YYYY` in secondary style on the left and the clean name (without the date prefix) following it.

#### Scenario: Archived change with standard date format
- **WHEN** the archived change directory is named `2026-05-02-specs-subnav`
- **THEN** the item shows `02/05/2026  specs-subnav`

### Requirement: Navigation in the index
The cursor SHALL be able to move through all items (active and archived) with `j` (down) and `k` (up). Section separators are not selectable items. The cursor SHALL NOT go past the first or last item.

#### Scenario: Navigate from active to archived
- **WHEN** the cursor is on the last active change and the user presses `j`
- **THEN** the cursor jumps to the first archived change

#### Scenario: No overflow at the extremes
- **WHEN** the cursor is on the last item and the user presses `j`
- **THEN** the cursor does not change

### Requirement: Select a change with Enter
Pressing `Enter` on an item SHALL close the index and open the selected change. If it is an active change, the mode transitions to `ModeNormal` with that active change. If it is an archived change, the mode transitions to `ModeViewingArchive` with that archived change.

#### Scenario: Select active change
- **WHEN** the cursor is on an active change and the user presses `Enter`
- **THEN** the mode transitions to `ModeNormal` showing that change

#### Scenario: Select archived change
- **WHEN** the cursor is on an archived change and the user presses `Enter`
- **THEN** the mode transitions to `ModeViewingArchive` showing the artifacts of that archived change

### Requirement: Index helpbar
The helpbar in `ModeIndex` SHALL show `j/k: navigate  Enter: open  Esc: quit`.

#### Scenario: Helpbar visible in the index
- **WHEN** the mode is `ModeIndex`
- **THEN** the helpbar shows `j/k: navigate  Enter: open  Esc: quit`
