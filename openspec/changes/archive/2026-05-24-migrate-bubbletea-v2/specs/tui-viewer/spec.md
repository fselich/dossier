# tui-viewer Delta Specification

## MODIFIED Requirements

### Requirement: Layout del TUI
The TUI SHALL divide the screen into fixed zones separated by horizontal lines: header (1 line), separator (1 line), tab bar (1 line), separator (1 line), content area (remainder), separator (1 line), help bar (1 line). In the `tasks` tab, a global progress bar is also added between the content area and the bottom separator. The header SHALL show `<project> · <change-name> [N/M]` where N is the position of the current change and M is the total number of active changes. The `View()` method SHALL return a `tea.View` struct with `AltScreen = true` and `BackgroundColor` set to the configured theme background color, instead of manually filling the background with padding.

#### Scenario: Separadores visibles entre zonas
- **WHEN** the TUI is rendered in any tab
- **THEN** a full-width horizontal line appears between the tab bar and the content, and another between the content and the help bar

#### Scenario: Un solo change activo
- **WHEN** there is a single active change
- **THEN** the header shows `my-project · feat-a [1/1]`

#### Scenario: Varios changes activos
- **WHEN** there are three active changes and the second is selected
- **THEN** the header shows `my-project · feat-b [2/3]`

#### Scenario: Alt screen and background color managed by tea.View
- **WHEN** the TUI renders any view
- **THEN** `tea.View.AltScreen` is set to `true` and `tea.View.BackgroundColor` reflects the configured theme

### Requirement: Salir del TUI
The user SHALL be able to exit the TUI at any time with `q` or `Ctrl+C`.

#### Scenario: Salir con q
- **WHEN** the user presses `q`
- **THEN** the TUI exits and the terminal is left in a clean state

#### Scenario: Salir con Ctrl+C
- **WHEN** the user presses `Ctrl+C`
- **THEN** the TUI exits and the terminal is left in a clean state
