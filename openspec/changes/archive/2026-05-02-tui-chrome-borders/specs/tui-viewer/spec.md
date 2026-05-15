## MODIFIED Requirements

### Requirement: TUI layout
The TUI SHALL divide the screen into fixed zones separated by horizontal lines: header (1 line), separator (1 line), tab bar (1 line), separator (1 line), content area (remainder), separator (1 line), help bar (1 line). In the `tasks` tab the global progress bar is also added between the content area and the bottom separator. The header SHALL show `<project> · <change-name> [N/M]` where N is the position of the current change and M is the total number of active changes.

#### Scenario: Separators visible between zones
- **WHEN** the TUI is rendered on any tab
- **THEN** a full-width horizontal line appears between the tab bar and the content, and another between the content and the help bar

#### Scenario: Single active change
- **WHEN** a single active change exists
- **THEN** the header shows `my-project · feat-a [1/1]`

#### Scenario: Multiple active changes
- **WHEN** three active changes exist and the second is selected
- **THEN** the header shows `my-project · feat-b [2/3]`
