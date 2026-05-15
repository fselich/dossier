## MODIFIED Requirements

### Requirement: Specs section in the index
In `ModeIndex`, the TUI SHALL display a "Specifications" section below the "Archived Changes" section. The section SHALL list each spec on a line with its name on the left and `N requirements` in secondary style on the right, aligned in two columns. The width of the name column SHALL adjust to the longest name. If there are no specs, the section SHALL display a message indicating that no specs are available.

#### Scenario: Specs present
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns at least one spec
- **THEN** the screen shows a "Specifications" section with each spec in the format `name  N requirements`, with the count column aligned

#### Scenario: Column alignment with names of different lengths
- **WHEN** there are specs with names of different lengths
- **THEN** all `N requirements` counts appear aligned in the same column

#### Scenario: No specs available
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns an empty list
- **THEN** the "Specifications" section shows the message "No specifications available"

#### Scenario: Specs loaded when entering the index
- **WHEN** the user enters `ModeIndex`
- **THEN** the list of specs is loaded from disk at that moment (same as archived changes)
