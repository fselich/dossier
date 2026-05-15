## MODIFIED Requirements

### Requirement: Helpbar del índice
The helpbar in `ModeIndex` SHALL show navigation hints and SHALL reflect the current sort mode via the `s` binding label:
- When sort mode is **name**: `j/k: navigate  Enter: open  Space: expand  s: sort by suffix  Esc: quit`
- When sort mode is **suffix**: `j/k: navigate  Enter: open  Space: expand  s: sort by name  Esc: quit`

#### Scenario: Helpbar en modo sort normal
- **WHEN** the mode is `ModeIndex` and the sort order is **name**
- **THEN** the helpbar shows `j/k: navigate  Enter: open  Space: expand  s: sort by suffix  Esc: quit`

#### Scenario: Helpbar en modo sort por sufijo
- **WHEN** the mode is `ModeIndex` and the sort order is **suffix**
- **THEN** the helpbar shows `j/k: navigate  Enter: open  Space: expand  s: sort by name  Esc: quit`
