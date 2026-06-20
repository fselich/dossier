## MODIFIED Requirements

### Requirement: Vista índice de pantalla completa
The TUI SHALL implement a `ModeIndex` mode that occupies the full screen with the same TUI chrome (borders, header, helpbar). The index SHALL show three sections: "Active Changes (N)", "Specifications (N)", and "Archived Changes (N)", where N is the count of items in that section; followed by the active changes, specifications, and archived changes respectively; the three separated by a section line. When a section has zero items, the empty-state message (e.g., "No active changes") SHALL be shown instead of a count. When a view background color is configured, the entire index view SHALL render with that background color filling the full terminal viewport, including all whitespace areas between elements and the empty area below the box frame.

#### Scenario: Active changes section title shows count
- **WHEN** the mode is `ModeIndex` and there are 3 active changes
- **THEN** the "Active Changes" section title reads "Active Changes (3)"

#### Scenario: Archived changes section title shows count
- **WHEN** the mode is `ModeIndex` and there are 5 archived changes
- **THEN** the "Archived Changes" section title reads "Archived Changes (5)"

#### Scenario: Specifications section title shows count
- **WHEN** the mode is `ModeIndex` and there are 2 project specs
- **THEN** the "Specifications" section title reads "Specifications (2)"

#### Scenario: Zero items still shows empty-state message
- **WHEN** the mode is `ModeIndex` and there are no active changes
- **THEN** the "Active Changes" section shows "No active changes" without a count
