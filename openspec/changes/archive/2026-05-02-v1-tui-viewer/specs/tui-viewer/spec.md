## ADDED Requirements

### Requirement: TUI layout
The TUI SHALL divide the screen into three fixed zones: header (1 line), tab bar (1 line) and content area (remainder). The header SHALL show `<project> · <change-name> [N/M]` where N is the position of the current change and M is the total number of active changes.

#### Scenario: Single active change
- **WHEN** a single active change exists
- **THEN** the header shows `my-project · feat-a [1/1]`

#### Scenario: Multiple active changes
- **WHEN** three active changes exist and the second is selected
- **THEN** the header shows `my-project · feat-b [2/3]`

### Requirement: Navigation between changes
The TUI SHALL allow navigating between active changes with `h` (previous) and `l` (next). Switching changes SHALL reset the selected tab to `tasks` if tasks is available, or to the first available artifact otherwise.

#### Scenario: Advance to the next change
- **WHEN** the user presses `l` while on change N
- **THEN** the TUI shows change N+1 (wrapping to the first if on the last)

#### Scenario: Go back to the previous change
- **WHEN** the user presses `h` while on change N
- **THEN** the TUI shows change N-1 (wrapping to the last if on the first)

### Requirement: Artifact tabs
The TUI SHALL display a tab bar with the tabs `proposal`, `design`, `tasks`, `specs`. Tabs for absent artifacts SHALL be shown visually disabled and not be selectable. The user SHALL be able to switch tabs with the keys `1`, `2`, `3`, `4`.

#### Scenario: Select available tab
- **WHEN** the user presses `2` and `design.md` exists
- **THEN** the content area shows the rendered design

#### Scenario: Attempt to select disabled tab
- **WHEN** the user presses `2` and `design.md` does not exist
- **THEN** the tab does not change and no error occurs

### Requirement: Markdown rendering with glamour
The TUI SHALL render the `proposal`, `design` and `specs` artifacts using glamour with the width of the content area. The content area SHALL be scrollable with `j`/`k` or the arrow keys.

#### Scenario: Scroll in long content
- **WHEN** the artifact has more content than the screen height and the user presses `j`
- **THEN** the content scrolls down one line

#### Scenario: Glamour wrap adjusted to width
- **WHEN** the terminal is 80 columns wide
- **THEN** glamour renders the markdown without exceeding those 80 columns

### Requirement: Welcome screen with no active changes
The TUI SHALL show an informative message when there are no active changes, instead of an empty state or an error.

#### Scenario: No active changes
- **WHEN** `openspec/changes/` exists but contains no active subdirectories
- **THEN** the TUI shows `"No active changes. Create one with /opsx:propose"`

### Requirement: Exit the TUI
The user SHALL be able to exit the TUI at any time with `q` or `Ctrl+C`.

#### Scenario: Exit with q
- **WHEN** the user presses `q`
- **THEN** the TUI exits and the terminal is left in a clean state

### Requirement: Keyboard help bar
The TUI SHALL show a fixed help line at the bottom with the active shortcuts in the current context.

#### Scenario: Tasks tab selected
- **WHEN** the active tab is `tasks`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: navigate  Space: toggle  q: quit`

#### Scenario: Proposal/design/specs tab selected
- **WHEN** the active tab is `proposal`, `design` or `specs`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: scroll  q: quit`
