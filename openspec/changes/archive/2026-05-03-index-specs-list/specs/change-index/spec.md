## MODIFIED Requirements

### Requirement: Archived changes format in the index
Each archived change SHALL be displayed with the clean name (without date prefix) on the left and the date `DD/MM/YYYY` in secondary style on the right, aligned in two columns. The width of the name column SHALL adjust to the longest name in the archived list.

#### Scenario: Archived with standard date format
- **WHEN** the archived directory is named `2026-05-02-specs-subnav`
- **THEN** the item shows `specs-subnav  02/05/2026` with the date in grey aligned to the right of the name

#### Scenario: Multiple archived items with different name lengths
- **WHEN** there are archived items with names of different lengths
- **THEN** all dates appear aligned in the same column

### Requirement: Full-screen index view
The TUI SHALL implement a `ModeIndex` mode that occupies the full screen with the same TUI chrome (borders, header, helpbar). The index SHALL show three sections: "Active" with active changes, "Archived" with changes in `openspec/changes/archive/`, and "Specs" with project specs in `openspec/specs/`; the three separated by a section line.

#### Scenario: Index with active, archived, and specs
- **WHEN** the mode is `ModeIndex` and there are active changes, archived changes, and project specs
- **THEN** the screen shows an "Active" section, an "Archived" section, and a "Specs" section, within the TUI chrome

#### Scenario: Index without active changes
- **WHEN** the mode is `ModeIndex` and there are no active changes
- **THEN** the "Active" section shows a message indicating there are no active changes

#### Scenario: Index without archived changes
- **WHEN** the mode is `ModeIndex` and there are no archived changes
- **THEN** the "Archived" section shows a message indicating there are no archived items

#### Scenario: Index without specs
- **WHEN** the mode is `ModeIndex` and there are no specs in `openspec/specs/`
- **THEN** the "Specs" section shows a message indicating that no specs are available
