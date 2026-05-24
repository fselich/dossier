# editor-launch Delta Specification

## MODIFIED Requirements

### Requirement: Abrir artefacto activo en editor externo
The TUI SHALL allow the user to open the file of the active tab's artifact in the system editor by pressing `e`. The editor SHALL be the value of the `$EDITOR` environment variable; if it is not defined, `vi` SHALL be used as a fallback. The TUI SHALL correctly suspend its control of the terminal before launching the editor using `tea.Exec` and resume it upon exit. After returning from the editor, mouse tracking SHALL still be functional because mouse mode is declared in `View()` and re-applied on every render frame.

#### Scenario: Abrir proposal en editor
- **WHEN** the active tab is `proposal` and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR proposal.md`; when the editor is closed the TUI resumes with functional mouse tracking

#### Scenario: Abrir tasks en editor
- **WHEN** the active tab is `tasks` and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR tasks.md`; when the editor is closed the TUI resumes with functional mouse tracking

#### Scenario: Fallback a vi cuando $EDITOR no está definido
- **WHEN** `$EDITOR` is not defined in the environment and the user presses `e`
- **THEN** the TUI launches `vi` with the path of the active artifact

#### Scenario: Tecla e en tab deshabilitada
- **WHEN** the user presses `e` and the active tab has no available artifact (`Present == false`)
- **THEN** nothing happens

#### Scenario: Mouse wheel works after editor return
- **WHEN** the user returns from the external editor
- **THEN** the user can immediately scroll the viewport with the mouse wheel without needing to restart the TUI
