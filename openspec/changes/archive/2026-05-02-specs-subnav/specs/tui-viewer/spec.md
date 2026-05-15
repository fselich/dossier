## MODIFIED Requirements

### Requirement: Artifact tabs
The TUI SHALL display a tab bar with the tabs `proposal`, `design`, `specs`, `tasks`. Tabs for absent artifacts SHALL be shown visually disabled and not be selectable. The user SHALL be able to switch tabs with the keys `1`, `2`, `3`, `4`. The `3` key SHALL have dual behaviour: if the active tab is not `specs`, switch to it; if it is already `specs`, cycle to the next available spec. If an absent artifact appears on disk during the session, the corresponding tab SHALL become enabled without needing to restart the TUI.

#### Scenario: Select available tab with numeric key
- **WHEN** the user presses `2` and `design.md` exists
- **THEN** the content area shows the rendered design

#### Scenario: Attempt to select disabled tab
- **WHEN** the user presses `2` and `design.md` does not exist
- **THEN** the tab does not change and no error occurs

#### Scenario: Tab becomes enabled when artifact appears
- **WHEN** the TUI starts without `proposal.md` and an external process creates that file
- **THEN** within a maximum of 500 ms the `proposal` tab is shown as enabled and is selectable

#### Scenario: Key 3 from another tab goes to specs
- **WHEN** the active tab is `proposal` and the user presses `3`
- **THEN** the active tab changes to `specs`

#### Scenario: Key 3 in specs cycles to the next spec
- **WHEN** the active tab is `specs` and the user presses `3`
- **THEN** the visible spec advances to the next one (wrapping to the first)
