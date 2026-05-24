## MODIFIED Requirements

### Requirement: Tabs de artifact

The TUI SHALL show a tab bar with tabs `proposal`, `design`, `tasks`, `specs`. Tabs for absent artifacts SHALL be shown visually disabled and not selectable. The user SHALL be able to change tabs with keys `1`, `2`, `3`, `4` or by left-clicking on the tab label with the mouse. The `3` key SHALL have dual behavior: if the active tab is not `specs`, it switches to it; if it is already `specs`, it cycles to the next available spec. If an absent artifact appears on disk during the session, the corresponding tab SHALL be enabled without needing to restart the TUI.

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
