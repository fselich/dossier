## MODIFIED Requirements

### Requirement: Navigation between changes
The TUI SHALL allow navigating between active changes with `h` (previous) and `l` (next). Switching changes SHALL reset the selected tab to `proposal` if available, or to the first available artifact otherwise. Pressing `a` or `Esc` from `ModeNormal` SHALL open `ModeIndex`. Pressing `q` or `Ctrl+C` SHALL quit the application from any mode.

#### Scenario: Advance to the next change
- **WHEN** the user presses `l` while on change N
- **THEN** the TUI shows change N+1 (wrapping to the first if on the last)

#### Scenario: Go back to the previous change
- **WHEN** the user presses `h` while on change N
- **THEN** the TUI shows change N-1 (wrapping to the last if on the first)

#### Scenario: 'a' from ModeNormal opens the index
- **WHEN** the mode is `ModeNormal` and the user presses `a`
- **THEN** the mode transitions to `ModeIndex`

#### Scenario: 'Esc' from ModeNormal opens the index
- **WHEN** the mode is `ModeNormal` and the user presses `Esc`
- **THEN** the mode transitions to `ModeIndex`

#### Scenario: Quit with q from any mode
- **WHEN** the user presses `q` while in any mode
- **THEN** the TUI exits

### Requirement: Keyboard help bar
The TUI SHALL show a fixed help line at the bottom with the active shortcuts for the current context.

#### Scenario: Tasks tab selected
- **WHEN** the active tab is `tasks` and the mode is `ModeNormal`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit`

#### Scenario: Proposal/design/specs tab selected
- **WHEN** the active tab is `proposal`, `design`, or `specs` and the mode is `ModeNormal`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: scroll  e: edit  Esc: index  q: quit`
