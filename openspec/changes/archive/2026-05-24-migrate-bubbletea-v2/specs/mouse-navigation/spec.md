# mouse-navigation Delta Specification

## MODIFIED Requirements

### Requirement: Mouse event capture

The TUI SHALL enable mouse event capture declaratively via `tea.View.MouseMode = tea.MouseModeCellMotion` instead of the imperative `tea.WithMouseCellMotion()` program option. Mouse events SHALL be delivered as `tea.MouseClickMsg`, `tea.MouseWheelMsg`, `tea.MouseMotionMsg`, and `tea.MouseReleaseMsg` instead of the unified `tea.MouseMsg`. Mouse tracking SHALL persist across external editor sessions because the mouse mode is re-declared on every frame.

#### Scenario: Mouse events are received after startup
- **WHEN** the TUI starts
- **THEN** `tea.MouseClickMsg` and `tea.MouseWheelMsg` messages are delivered to `Update`

#### Scenario: Mouse tracking persists after external editor
- **WHEN** the user opens an external editor and returns to the TUI
- **THEN** mouse events continue to be received and handled correctly

#### Scenario: Text selection bypasses mouse capture
- **WHEN** the user holds Shift and clicks / drags
- **THEN** the terminal performs native text selection instead of sending mouse events to the TUI
