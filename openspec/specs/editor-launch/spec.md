# editor-launch Specification

## Purpose
Allows editing the active artifact of a change in the system editor (`$EDITOR`) by pressing `e`, and automatically reloads the content when the editor is closed without needing to restart the TUI.

## Requirements


### Requirement: Abrir artefacto activo en editor externo
The TUI SHALL allow the user to open the file of the active tab's artifact in the system editor by pressing `e`. The editor SHALL be the value of the `$EDITOR` environment variable; if it is not defined, `vi` SHALL be used as a fallback. The TUI SHALL correctly suspend its control of the terminal before launching the editor and resume it upon exit.

#### Scenario: Abrir proposal en editor
- **WHEN** the active tab is `proposal` and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR proposal.md`; when the editor is closed the TUI resumes

#### Scenario: Abrir tasks en editor
- **WHEN** the active tab is `tasks` and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR tasks.md`; when the editor is closed the TUI resumes

#### Scenario: Fallback a vi cuando $EDITOR no está definido
- **WHEN** `$EDITOR` is not defined in the environment and the user presses `e`
- **THEN** the TUI launches `vi` with the path of the active artifact

#### Scenario: Tecla e en tab deshabilitada
- **WHEN** the user presses `e` and the active tab has no available artifact (`Present == false`)
- **THEN** nothing happens

### Requirement: Recarga inmediata tras cierre del editor
The TUI SHALL reload the content of the edited artifact immediately upon returning from the editor, without waiting for the next polling cycle.

#### Scenario: Recarga de tasks tras edición
- **WHEN** the user edits `tasks.md` in the editor and closes the editor
- **THEN** the TUI shows the updated tasks content instantly, with the cursor restored by text

#### Scenario: Recarga de artifact markdown tras edición
- **WHEN** the user edits `proposal.md`, `design.md`, or a `spec.md` and closes the editor
- **THEN** the TUI invalidates the render cache for that tab and re-renders with the new content
