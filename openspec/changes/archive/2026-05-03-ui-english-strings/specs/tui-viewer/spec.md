## MODIFIED Requirements

### Requirement: Welcome screen with no active changes
The TUI SHALL display an informative message when there are no active changes, instead of an empty state or an error.

#### Scenario: No active changes
- **WHEN** `openspec/changes/` exists but contains no active subdirectories
- **THEN** the TUI shows `"No active changes. Create one with /opsx:propose"`

### Requirement: Keyboard help bar
The TUI SHALL display a fixed help line at the bottom with the active shortcuts in the current context.

#### Scenario: Tasks tab selected
- **WHEN** the active tab is `tasks` and the mode is `ModeNormal`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit`

#### Scenario: Proposal/design/specs tab selected
- **WHEN** the active tab is `proposal`, `design`, or `specs` and the mode is `ModeNormal`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: scroll  e: edit  Esc: index  q: quit`
