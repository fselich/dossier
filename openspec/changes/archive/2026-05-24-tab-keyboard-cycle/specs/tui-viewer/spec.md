## MODIFIED Requirements

### Requirement: Tabs de artifact
The TUI SHALL show a tab bar with tabs `proposal`, `design`, `tasks`, `specs`. Tabs for absent artifacts SHALL be shown visually disabled and not selectable. The user SHALL be able to change tabs with keys `1`, `2`, `3`, `4`, with `Tab` (next available) and `Shift+Tab` (previous available), or by left-clicking on the tab label with the mouse. `Tab` and `Shift+Tab` SHALL skip disabled tabs and wrap around at the ends. The `3` key SHALL have dual behavior: if the active tab is not `specs`, it switches to it; if it is already `specs`, it cycles to the next available spec. If an absent artifact appears on disk during the session, the corresponding tab SHALL be enabled without needing to restart the TUI.

#### Scenario: Seleccionar tab disponible con tecla numérica
- **WHEN** the user presses `2` and `design.md` exists
- **THEN** the content area shows the rendered design

#### Scenario: Intentar seleccionar tab deshabilitada con tecla
- **WHEN** the user presses `2` and `design.md` does not exist
- **THEN** the tab does not change and no error occurs

#### Scenario: Seleccionar tab disponible con click del mouse
- **WHEN** the user left-clicks on the "design" tab label and `design.md` exists
- **THEN** the content area shows the rendered design

#### Scenario: Intentar seleccionar tab deshabilitada con click
- **WHEN** the user left-clicks on a disabled tab label and the artifact does not exist
- **THEN** the tab does not change and no error occurs

#### Scenario: Tab se habilita al aparecer artifact
- **WHEN** the TUI starts without `proposal.md` and an external process creates that file
- **THEN** within a maximum of 500 ms the `proposal` tab is shown as enabled and is selectable

#### Scenario: Tecla 3 desde otra tab va a specs
- **WHEN** the active tab is `proposal` and the user presses `3`
- **THEN** the active tab changes to `specs`

#### Scenario: Tecla 3 en specs cicla al siguiente spec
- **WHEN** the active tab is `specs` and the user presses `3`
- **THEN** the visible spec advances to the next one (wrapping to the first)

#### Scenario: Ciclar hacia adelante con Tab
- **WHEN** the active tab is `proposal`, `design` is disabled, and `specs` is available
- **THEN** the user pressing `Tab` changes the active tab to `specs` (skipping disabled `design`)

#### Scenario: Ciclar hacia atrás con Shift+Tab
- **WHEN** the active tab is `tasks`, `specs` is disabled, and `design` is available
- **THEN** the user pressing `Shift+Tab` changes the active tab to `design` (skipping disabled `specs`)

#### Scenario: Tab da la vuelta al final
- **WHEN** the active tab is the last available tab and the user presses `Tab`
- **THEN** the active tab wraps around to the first available tab

#### Scenario: Shift+Tab da la vuelta al principio
- **WHEN** the active tab is the first available tab and the user presses `Shift+Tab`
- **THEN** the active tab wraps around to the last available tab

#### Scenario: Tab no actúa en modo configuración
- **WHEN** the mode is `ModeViewingConfig` and the user presses `Tab`
- **THEN** the tab does not change and the key is handled by the text input instead
