## MODIFIED Requirements

### Requirement: Specs sub-navigation
When the active tab is `specs` and at least one spec is available, the TUI SHALL display a row of chips as the first line inside the content block, immediately after the horizontal separator (`├───┤`), with the name of each spec. The chip of the currently visible spec SHALL be shown with active style (same as an active tab). The other chips SHALL be shown with inactive style. The row SHALL be static (it is not part of the scrollable viewport area).

#### Scenario: Single spec
- **WHEN** the change has a single spec and the active tab is `specs`
- **THEN** the chip row is shown as the first line of the content block, with one chip representing that spec marked as active

#### Scenario: Multiple specs
- **WHEN** the change has two or more specs and the active tab is `specs`
- **THEN** the chip row is shown as the first line of the content block; the chip of the visible spec is active and the others are inactive

#### Scenario: Row absent on other tabs
- **WHEN** the active tab is not `specs`
- **THEN** no specs chip row is shown in the content block

#### Scenario: The row does not disappear on scroll
- **WHEN** the active tab is `specs` and the user scrolls down
- **THEN** the chip row remains visible as the first line of the content block

#### Scenario: Visual separation between tab bar and spec chips
- **WHEN** the active tab is `specs`
- **THEN** the horizontal separator `├───┤` appears between the tab bar and the spec chip row
