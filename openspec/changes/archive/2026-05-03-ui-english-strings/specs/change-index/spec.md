## MODIFIED Requirements

### Requirement: Full-screen index view
The TUI SHALL implement a `ModeIndex` mode that occupies the entire screen with the same TUI chrome (borders, header, helpbar). The index SHALL display three sections: "Active Changes" with the active changes, "Archived Changes" with the changes in `openspec/changes/archive/`, and "Specifications" with the project specs in `openspec/specs/`; all three separated by a section divider line.

#### Scenario: Index with active, archived, and specs
- **WHEN** the mode is `ModeIndex` and there are active changes, archived changes, and project specs
- **THEN** the screen shows an "Active Changes" section, an "Archived Changes" section, and a "Specifications" section, within the TUI chrome

#### Scenario: Index with no active changes
- **WHEN** the mode is `ModeIndex` and there are no active changes
- **THEN** the "Active Changes" section shows a message indicating there are no active changes

#### Scenario: Index with no archived changes
- **WHEN** the mode is `ModeIndex` and there are no archived changes
- **THEN** the "Archived Changes" section shows a message indicating there are no archived changes

#### Scenario: Index with no specs
- **WHEN** the mode is `ModeIndex` and there are no specs in `openspec/specs/`
- **THEN** the "Specifications" section shows a message indicating there are no specs available

### Requirement: Index helpbar
The helpbar in `ModeIndex` SHALL show `j/k: navigate  Enter: open  Esc: quit`.

#### Scenario: Helpbar visible in the index
- **WHEN** the mode is `ModeIndex`
- **THEN** the helpbar shows `j/k: navigate  Enter: open  Esc: quit`
