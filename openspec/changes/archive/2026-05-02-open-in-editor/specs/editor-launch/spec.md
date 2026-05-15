## ADDED Requirements

### Requirement: Open active artifact in external editor
The TUI SHALL allow the user to open the artifact file of the active tab in the system editor by pressing `e`. The editor SHALL be the value of the environment variable `$EDITOR`; if it is not defined, `vi` SHALL be used as a fallback. The TUI SHALL correctly suspend its terminal control before launching the editor and resume it on exit.

#### Scenario: Open proposal in editor
- **WHEN** the active tab is `proposal` and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR proposal.md`; when the editor is closed the TUI resumes

#### Scenario: Open tasks in editor
- **WHEN** the active tab is `tasks` and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR tasks.md`; when the editor is closed the TUI resumes

#### Scenario: Fallback to vi when $EDITOR is not defined
- **WHEN** `$EDITOR` is not defined in the environment and the user presses `e`
- **THEN** the TUI launches `vi` with the path of the active artifact

#### Scenario: e key on disabled tab
- **WHEN** the user presses `e` and the active tab has no available artifact (`Present == false`)
- **THEN** nothing happens

### Requirement: Immediate reload after editor is closed
The TUI SHALL reload the content of the edited artifact immediately upon returning from the editor, without waiting for the next polling cycle.

#### Scenario: Reload of tasks after editing
- **WHEN** the user edits `tasks.md` in the editor and closes the editor
- **THEN** the TUI shows the updated tasks content instantly, with the cursor restored by text

#### Scenario: Reload of markdown artifact after editing
- **WHEN** the user edits `proposal.md`, `design.md` or a `spec.md` and closes the editor
- **THEN** the TUI invalidates the render cache for that tab and re-renders with the new content
