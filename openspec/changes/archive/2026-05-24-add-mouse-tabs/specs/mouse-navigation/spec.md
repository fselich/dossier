## ADDED Requirements

### Requirement: Mouse event capture

The TUI SHALL enable mouse event capture via `tea.WithMouseCellMotion()` so that `tea.MouseMsg` events are delivered to the `Update` function.

#### Scenario: Mouse events are received after startup
- **WHEN** the TUI starts
- **THEN** `tea.MouseMsg` messages for clicks, releases, and wheel events are delivered to `Update`

#### Scenario: Text selection bypasses mouse capture
- **WHEN** the user holds Shift and clicks / drags
- **THEN** the terminal performs native text selection instead of sending mouse events to the TUI

### Requirement: Wheel scrolling

The TUI SHALL handle `tea.MouseMsg` wheel events (up and down) by scrolling the viewport 3 lines per wheel tick. Wheel scrolling SHALL apply in all modes. In task and index modes, wheel scrolling SHALL scroll the viewport content and SHALL NOT move the task cursor or index cursor.

#### Scenario: Wheel down scrolls content
- **WHEN** the user scrolls the mouse wheel down while viewing any content
- **THEN** the viewport scrolls down by 3 lines

#### Scenario: Wheel up scrolls content
- **WHEN** the user scrolls the mouse wheel up while viewing any content
- **THEN** the viewport scrolls up by 3 lines

#### Scenario: Wheel at top of content does not crash
- **WHEN** the user scrolls the mouse wheel up while the viewport is already at the top
- **THEN** no error occurs and the viewport remains at the top

#### Scenario: Wheel in tasks tab scrolls viewport, not cursor
- **WHEN** the `tasks` tab is active and the user scrolls the mouse wheel
- **THEN** the viewport scrolls and the task cursor position does not change

#### Scenario: Wheel in index mode scrolls viewport, not cursor
- **WHEN** the mode is `ModeIndex` and the user scrolls the mouse wheel
- **THEN** the viewport scrolls and the index cursor position does not change

### Requirement: Tab selection via left-click

The TUI SHALL switch to the clicked tab on the tab bar when the user performs a left-click (press) on a tab label. The tab bar is present only in `ModeNormal` and `ModeViewingArchive`. The coordinate mapping SHALL use the tab label width including `Padding(0, 1)` plus one space between tabs, starting from X=1 (past the `│` border). Clicked disabled tabs (absent artifacts) SHALL be ignored.

#### Scenario: Click on an available tab switches to it
- **WHEN** the user left-clicks on the "design" label in the tab bar while on the "proposal" tab
- **THEN** the active tab changes to "design" and the content area shows the rendered design artifact

#### Scenario: Click on a disabled tab does nothing
- **WHEN** the user left-clicks on a tab label whose artifact does not exist on disk
- **THEN** the active tab does not change and no error occurs

#### Scenario: Click on the currently active tab reloads viewport
- **WHEN** the user left-clicks on the "proposal" tab while "proposal" is already the active tab
- **THEN** the viewport reloads (same behavior as pressing `1` when already on the proposal tab)

#### Scenario: Click between tabs does nothing
- **WHEN** the user left-clicks on the space between two tab labels
- **THEN** the active tab does not change

#### Scenario: Click outside tab bar area does nothing
- **WHEN** the user left-clicks outside the X range of any tab label
- **THEN** the active tab does not change

#### Scenario: Click in index mode does not trigger tab switch
- **WHEN** the mode is `ModeIndex` and the user left-clicks at any Y coordinate
- **THEN** no tab switching occurs (the tab bar is not present in index layout)
