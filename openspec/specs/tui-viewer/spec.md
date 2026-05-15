# tui-viewer Specification

## Purpose
Defines the layout and main behavior of the TUI: screen structure with borders and fixed zones, navigation between changes and tabs, markdown rendering with glamour, periodic polling for disk changes, and a welcome screen when there are no active changes.

## Requirements


### Requirement: Layout del TUI
The TUI SHALL divide the screen into fixed zones separated by horizontal lines: header (1 line), separator (1 line), tab bar (1 line), separator (1 line), content area (remainder), separator (1 line), help bar (1 line). In the `tasks` tab, a global progress bar is also added between the content area and the bottom separator. The header SHALL show `<project> · <change-name> [N/M]` where N is the position of the current change and M is the total number of active changes.

#### Scenario: Separadores visibles entre zonas
- **WHEN** the TUI is rendered in any tab
- **THEN** a full-width horizontal line appears between the tab bar and the content, and another between the content and the help bar

#### Scenario: Un solo change activo
- **WHEN** there is a single active change
- **THEN** the header shows `my-project · feat-a [1/1]`

#### Scenario: Varios changes activos
- **WHEN** there are three active changes and the second is selected
- **THEN** the header shows `my-project · feat-b [2/3]`

### Requirement: Navegación entre changes
The TUI SHALL allow navigating between active changes with `h` (previous) and `l` (next). Changing the change SHALL reset the selected tab to `proposal` if available, or to the first available artifact otherwise. Pressing `a` or `Esc` from `ModeNormal` SHALL open `ModeIndex`. Pressing `q` or `Ctrl+C` SHALL exit the application from any mode.

#### Scenario: Avanzar al siguiente change
- **WHEN** the user presses `l` while on change N
- **THEN** the TUI shows change N+1 (wrapping to the first if on the last)

#### Scenario: Retroceder al change anterior
- **WHEN** the user presses `h` while on change N
- **THEN** the TUI shows change N-1 (wrapping to the last if on the first)

#### Scenario: 'a' desde ModeNormal abre el índice
- **WHEN** the mode is `ModeNormal` and the user presses `a`
- **THEN** the mode switches to `ModeIndex`

#### Scenario: 'Esc' desde ModeNormal abre el índice
- **WHEN** the mode is `ModeNormal` and the user presses `Esc`
- **THEN** the mode switches to `ModeIndex`

#### Scenario: Salir con q desde cualquier modo
- **WHEN** the user presses `q` from any mode
- **THEN** the TUI exits

### Requirement: Tabs de artifact
The TUI SHALL show a tab bar with tabs `proposal`, `design`, `tasks`, `specs`. Tabs for absent artifacts SHALL be shown visually disabled and not selectable. The user SHALL be able to change tabs with keys `1`, `2`, `3`, `4`. The `3` key SHALL have dual behavior: if the active tab is not `specs`, it switches to it; if it is already `specs`, it cycles to the next available spec. If an absent artifact appears on disk during the session, the corresponding tab SHALL be enabled without needing to restart the TUI.

#### Scenario: Seleccionar tab disponible con tecla numérica
- **WHEN** the user presses `2` and `design.md` exists
- **THEN** the content area shows the rendered design

#### Scenario: Intentar seleccionar tab deshabilitada
- **WHEN** the user presses `2` and `design.md` does not exist
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

### Requirement: Render de markdown con glamour
The TUI SHALL render `proposal`, `design`, and `specs` artifacts using glamour with the width of the content area. The content area SHALL be scrollable with `j`/`k` or the arrow keys.

#### Scenario: Scroll en contenido largo
- **WHEN** the artifact has more content than the screen height and the user presses `j`
- **THEN** the content scrolls down one line

#### Scenario: Wrap de glamour ajustado al ancho
- **WHEN** the terminal is 80 columns wide
- **THEN** glamour renders the markdown without exceeding those 80 columns

### Requirement: Pantalla de bienvenida sin changes activos
The TUI SHALL show an informational message when there are no active changes, instead of an empty state or an error. It SHALL also show a help line with the available actions: `a/Esc: index  q: quit`.

#### Scenario: Sin changes activos
- **WHEN** `openspec/changes/` exists but contains no active subdirectories
- **THEN** the TUI shows `"No active changes. Create one with /opsx:propose"` and the help line `a/Esc: index  q: quit`

### Requirement: Salir del TUI
The user SHALL be able to exit the TUI at any time with `q` or `Ctrl+C`.

#### Scenario: Salir con q
- **WHEN** the user presses `q`
- **THEN** the TUI exits and the terminal is left in a clean state

### Requirement: Barra de ayuda de teclado
The TUI SHALL show a fixed help line at the bottom with the active shortcuts in the current context.

#### Scenario: Tab de tasks seleccionada
- **WHEN** the active tab is `tasks` and the mode is `ModeNormal`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit`

#### Scenario: Tab de proposal/design/specs seleccionada
- **WHEN** the active tab is `proposal`, `design`, or `specs` and the mode is `ModeNormal`
- **THEN** the help line shows `h/l: change  1-4: artifact  j/k: scroll  e: edit  Esc: index  q: quit`

### Requirement: Polling periódico de artifacts
The TUI SHALL start a polling cycle every 500 ms on startup. On each tick it SHALL compare the on-disk content of the artifacts of the currently visible change with the in-memory content, AND detect changes in artifact presence (absent → present). If at the moment of the tick `len(m.project.Changes) == 0`, the tick SHALL attempt to reload the change list from disk and adopt the new state if at least one change is available. The cycle SHALL continue while the TUI is active.

#### Scenario: Tick sin cambios
- **WHEN** no file in the change has changed on disk
- **THEN** the TUI does not update any state or re-render anything

#### Scenario: Tick detecta cambio en tasks.md
- **WHEN** the content of `tasks.md` on disk differs from the in-memory content
- **THEN** the TUI re-parses the tasks, restores the cursor, and refreshes the view if the active tab is `tasks`

#### Scenario: Tick detecta cambio en artifact de markdown
- **WHEN** the content of `proposal.md`, `design.md`, or a `spec.md` on disk differs from the in-memory content
- **THEN** the TUI invalidates the corresponding entry in the render cache; the next time the user accesses that tab it is re-rendered with the new content

#### Scenario: Tick detecta aparición de artifact ausente
- **WHEN** an artifact that did not exist in the previous tick now exists on disk
- **THEN** the TUI updates the artifact presence state and enables the corresponding tab

#### Scenario: TUI arranca sin changes activos y se crea uno
- **WHEN** the TUI starts with `len(m.project.Changes) == 0` and during the session a change is created on disk
- **THEN** within a maximum of 500 ms the TUI reloads the change list and shows the new change

### Requirement: Actualización de tasks visible en tiempo real
When the TUI detects a change in `tasks.md` and the active tab is `tasks`, it SHALL refresh the view immediately without user intervention.

#### Scenario: Agente marca tarea como completada
- **WHEN** an external process changes `- [ ] tarea` to `- [x] tarea` in `tasks.md`
- **THEN** within a maximum of 500 ms the TUI shows the task as completed with the updated progress bar

### Requirement: Word wrap en todos los items de tarea
In the `tasks` tab, all task items SHALL word-wrap to the width of the content area (`m.width - 2`), regardless of whether the item is under the cursor or not.

#### Scenario: Item largo sin cursor
- **WHEN** a task item has more characters than the terminal width and the cursor is not on it
- **THEN** the text word-wraps and is shown in full across multiple lines

#### Scenario: Item largo con cursor
- **WHEN** a task item has more characters than the terminal width and the cursor is on it
- **THEN** the text word-wraps and is shown in full across multiple lines with the cursor style

### Requirement: Barra de progreso global en la vista de tasks
The TUI SHALL show a global progress bar as the first content line of the `tasks` tab, before any section. The bar SHALL reflect the total completed tasks over the total tasks in the change.

#### Scenario: Change con tareas parcialmente completadas
- **WHEN** a change has 3 completed tasks out of 8 total
- **THEN** the first line of the tasks view shows a progress bar with `3/8` and a proportional fraction of filled blocks

#### Scenario: Change con todas las tareas completadas
- **WHEN** all tasks in the change are marked as completed
- **THEN** the progress bar appears completely filled and shows the total as `N/N`

#### Scenario: Change sin tareas
- **WHEN** the change has no task items
- **THEN** no global progress bar is shown

### Requirement: Actualización inmediata del contador de progreso tras toggle
When the user toggles a task with `Space`, the progress counter in the tab bar SHALL update in the same frame, without waiting for the next polling cycle.

#### Scenario: Marcar tarea completa actualiza tab bar
- **WHEN** the user presses `Space` on a pending task and the disk write succeeds
- **THEN** the `N/M` counter and the progress bar in the tab bar update immediately in the same render

#### Scenario: Desmarcar tarea actualiza tab bar
- **WHEN** the user presses `Space` on a completed task and the disk write succeeds
- **THEN** the `N/M` counter and the progress bar in the tab bar decrement immediately in the same render
