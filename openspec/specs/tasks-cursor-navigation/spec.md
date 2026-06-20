# Tasks Cursor Navigation

## Purpose

Define how the tasks tab cursor interacts with section headers during keyboard navigation, ensuring section headers are reachable so users can see context above the task list.

## Requirements

### Requirement: Tasks cursor can navigate to section headers
The system SHALL allow the tasks tab cursor to land on section header items (`KindSection`) during keyboard navigation.

#### Scenario: Cursor moves onto section header from above
- **WHEN** the user presses `<up>` and the item above the current cursor position is a section header
- **THEN** the cursor SHALL move to that section header

#### Scenario: Cursor moves onto section header from below
- **WHEN** the user presses `<down>` and the next item after the current cursor position is a section header
- **THEN** the cursor SHALL move to that section header

#### Scenario: Toggle is rejected on section header
- **WHEN** the cursor is on a section header and the user presses `<Enter>` to toggle
- **THEN** the system SHALL NOT toggle anything (no change to task state)
