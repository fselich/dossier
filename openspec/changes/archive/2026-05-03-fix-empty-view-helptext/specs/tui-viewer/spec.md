## MODIFIED Requirements

### Requirement: Welcome screen with no active changes
The TUI SHALL show an informational message when there are no active changes, instead of an empty state or an error. It SHALL also show a help line in English with the available actions: `a/Esc: index  q: quit`.

#### Scenario: No active changes
- **WHEN** `openspec/changes/` exists but contains no active subdirectories
- **THEN** the TUI shows `"No active changes. Create one with /opsx:propose"` and the help line `a/Esc: index  q: quit`
