## MODIFIED Requirements

### Requirement: Wheel scrolling

The TUI SHALL handle `tea.MouseMsg` wheel events (up and down). In `ModeIndex`, wheel events SHALL move the index cursor up or down (one item per tick) and the viewport SHALL auto-follow the cursor to keep it visible. In `ModeNormal` with `TabTasks` active, wheel events SHALL move the task cursor up or down (one task per tick) and the viewport SHALL auto-follow the cursor. In all other modes and views, wheel events SHALL scroll the viewport by 3 lines per tick.

#### Scenario: Wheel down scrolls content
- **WHEN** the user scrolls the mouse wheel down while viewing a proposal, design, spec, config, or archive
- **THEN** the viewport scrolls down by 3 lines

#### Scenario: Wheel up scrolls content
- **WHEN** the user scrolls the mouse wheel up while viewing a proposal, design, spec, config, or archive
- **THEN** the viewport scrolls up by 3 lines

#### Scenario: Wheel at top of content does not crash
- **WHEN** the user scrolls the mouse wheel up while the viewport is already at the top
- **THEN** no error occurs and the viewport remains at the top

#### Scenario: Wheel in index mode moves cursor
- **WHEN** the mode is `ModeIndex` and the user scrolls the mouse wheel
- **THEN** the index cursor moves up or down by one item per wheel tick and the viewport auto-follows to keep the cursor visible

#### Scenario: Wheel in tasks tab moves cursor
- **WHEN** the `tasks` tab is active in `ModeNormal` and the user scrolls the mouse wheel
- **THEN** the task cursor moves up or down by one task per wheel tick and the viewport auto-follows to keep the cursor visible
