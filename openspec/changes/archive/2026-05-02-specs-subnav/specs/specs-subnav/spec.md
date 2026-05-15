## ADDED Requirements

### Requirement: Specs sub-navigation
When the active tab is `specs` and at least one spec is available, the TUI SHALL display a row of chips below the tab bar with the name of each spec. The chip of the currently visible spec SHALL be shown with active style (same as an active tab). The other chips SHALL be shown with inactive style.

#### Scenario: Single spec
- **WHEN** the change has a single spec and the active tab is `specs`
- **THEN** a sub-bar is shown with one chip representing that spec, marked as active

#### Scenario: Multiple specs
- **WHEN** the change has two or more specs and the active tab is `specs`
- **THEN** a sub-bar is shown with one chip per spec; the chip of the visible spec is active and the others are inactive

#### Scenario: Sub-bar absent on other tabs
- **WHEN** the active tab is not `specs`
- **THEN** no specs sub-bar is shown

### Requirement: Spec cycling with key 3
The `3` key SHALL have dual behaviour: if the active tab is not `specs`, switch to it showing the last selected spec (or the first if there is no prior selection). If the active tab is already `specs`, SHALL advance to the next spec in the list, wrapping back to the first after the last.

#### Scenario: Enter specs from another tab
- **WHEN** the active tab is `proposal` and the user presses `3`
- **THEN** the active tab becomes `specs` and the previously selected spec is shown (or the first one)

#### Scenario: Cycle to the next spec
- **WHEN** the active tab is `specs`, there are 3 specs and the active spec is the second, and the user presses `3`
- **THEN** the active spec becomes the third and the viewport shows its content

#### Scenario: Cycling from the last spec wraps to the first
- **WHEN** the active tab is `specs`, the active spec is the last, and the user presses `3`
- **THEN** the active spec becomes the first

#### Scenario: Single spec, pressing 3 changes nothing visible
- **WHEN** the active tab is `specs` and there is only one spec, and the user presses `3`
- **THEN** the active spec remains the same and the content does not change

### Requirement: Content height adjustment with visible sub-nav
When the specs sub-nav is visible, the content area SHALL reduce its height by 1 line to accommodate the extra row, preventing the viewport from overflowing outside the box.

#### Scenario: Reduced height in specs tab
- **WHEN** the active tab is `specs` and specs are available
- **THEN** the viewport has 1 fewer line of height than in the other tabs

#### Scenario: Normal height in other tabs
- **WHEN** the active tab is `proposal`, `design` or `tasks`
- **THEN** the viewport has the standard height
