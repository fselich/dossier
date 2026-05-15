## ADDED Requirements

### Requirement: Project specs loading
The loader SHALL expose a `LoadProjectSpecs()` function that reads `openspec/specs/` in the current working directory and returns a list of `ProjectSpec`, each with its `Name` and the `RequirementCount` obtained by counting `### Requirement: ` lines in the corresponding `spec.md`. Specs SHALL be sorted alphabetically by name.

#### Scenario: Specs available
- **WHEN** `openspec/specs/` contains two or more subdirectories with `spec.md`
- **THEN** `LoadProjectSpecs()` returns one entry per spec with the correct requirement count

#### Scenario: Specs directory absent
- **WHEN** `openspec/specs/` does not exist or is empty
- **THEN** `LoadProjectSpecs()` returns an empty list without error

#### Scenario: spec.md with no requirements
- **WHEN** a `spec.md` contains no `### Requirement:` line
- **THEN** the corresponding `ProjectSpec` has `RequirementCount` 0 and still appears in the list

### Requirement: Specs section in the index
In `ModeIndex`, the TUI SHALL display a "Specs" section below the "Archived" section. The section SHALL list each spec on a line with its name on the left and `N requirements` in secondary style on the right, aligned in two columns. The width of the name column SHALL adjust to the longest name. If there are no specs, the section SHALL show a message indicating that no specs are available.

#### Scenario: Specs present
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns at least one spec
- **THEN** the screen shows a "Specs" section with each spec in the format `name  N requirements`, with the count column aligned

#### Scenario: Column alignment with different name lengths
- **WHEN** there are specs with names of different lengths
- **THEN** all `N requirements` counts appear aligned in the same column

#### Scenario: No specs available
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns an empty list
- **THEN** the "Specs" section shows the message "No specs available"

#### Scenario: Specs loaded when entering the index
- **WHEN** the user enters `ModeIndex`
- **THEN** the spec list is loaded from disk at that moment (same as archived items)

### Requirement: Specs not selectable in the index
Specs listed in the "Specs" section of `ModeIndex` SHALL be purely informational. They SHALL NOT be part of `indexItems`, the cursor SHALL NOT be able to position on them, and the `j`/`k`/`Enter` keys SHALL NOT interact with them.

#### Scenario: Cursor does not enter the Specs section
- **WHEN** the cursor is on the last navigable item (last archived) and the user presses `j`
- **THEN** the cursor does not move (existing boundary behavior, no overflow into specs)

#### Scenario: Enter on a navigable item does not activate specs
- **WHEN** the user presses `Enter` on an active or archived change
- **THEN** the behavior is the existing one; the Specs section does not interfere
