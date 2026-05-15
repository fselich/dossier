## MODIFIED Requirements

### Requirement: Keyboard help bar
The TUI SHALL show a fixed help line at the bottom with the active shortcuts in the current context.

#### Scenario: Tasks tab selected
- **WHEN** the active tab is `tasks`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  q: quit`

#### Scenario: Proposal/design/specs tab selected
- **WHEN** the active tab is `proposal`, `design` or `specs`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: scroll  e: edit  q: quit`
