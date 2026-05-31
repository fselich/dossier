# mouse-navigation Specification (Delta)

## ADDED Requirements

### Requirement: Index item selection via left-click

In `ModeIndex`, the TUI SHALL handle left-click on items rendered inside the viewport. Clicking SHALL move the cursor to the clicked item. If the clicked item is already under the cursor, the TUI SHALL perform the primary action for that item kind (Enter for active changes, archived changes, and requirements; Space for expanding/collapsing specs). Clicks on section headers, blank lines, or empty areas SHALL be ignored. The coordinate mapping SHALL account for the viewport's YOffset to correctly translate screen coordinates to content lines.

#### Scenario: Click on unselected index item moves cursor

- **WHEN** the cursor is on item A and the user left-clicks on item B in the index
- **THEN** the cursor moves to item B and no navigation (opening/expand) occurs

#### Scenario: Click on selected active change opens it

- **WHEN** the cursor is on an active change item and the user left-clicks on it
- **THEN** the mode switches to `ModeNormal` showing that change (same as pressing Enter)

#### Scenario: Click on selected archived change opens it

- **WHEN** the cursor is on an archived change item and the user left-clicks on it
- **THEN** the mode switches to `ModeViewingArchive` showing that change's artifacts

#### Scenario: Click on selected spec toggles expansion

- **WHEN** the cursor is on a spec item and the user left-clicks on it
- **THEN** the spec expands showing its requirements if collapsed, or collapses them if expanded (same as pressing Space)

#### Scenario: Click on selected requirement opens spec viewer

- **WHEN** the cursor is on a requirement item and the user left-clicks on it
- **THEN** the mode switches to `ModeViewingSpec` focused on that requirement (same as pressing Enter)

#### Scenario: Click on section header does nothing

- **WHEN** the user left-clicks on the "Active Changes", "Specifications", or "Archived Changes" header line
- **THEN** the cursor does not move and no action occurs

#### Scenario: Click on blank line does nothing

- **WHEN** the user left-clicks on an empty line between index sections
- **THEN** the cursor does not move and no action occurs

#### Scenario: Click outside viewport area does nothing

- **WHEN** the user left-clicks on a screen row outside the viewport content area (e.g., header, border, help bar)
- **THEN** the cursor does not move and no action occurs

#### Scenario: Click works correctly with scrolled viewport

- **WHEN** the viewport is scrolled down (YOffset > 0) and the user left-clicks on a visible item
- **THEN** the correct item is selected regardless of the scroll offset

#### Scenario: Click in normal mode does not trigger index selection

- **WHEN** the mode is `ModeNormal` and the user left-clicks at any Y coordinate in the viewport
- **THEN** no index item selection occurs (existing tab click behavior unchanged)
