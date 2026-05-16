# tasks-toggle Specification

## Purpose
Manages the TUI tasks tab: parsing `tasks.md`, cursor navigation between items, checkbox toggle with `Space`, per-section progress bars, and inline markdown rendering in task text.

## Requirements


### Requirement: Parsear tasks.md en items navegables
The system SHALL parse `tasks.md` line by line and produce a flat list of items. Each line matching `- [ ] <text>` or `- [x] <text>` SHALL become a task item (pending or completed respectively). Lines matching a section heading (`## <text>`) SHALL become non-interactive section items. All other lines SHALL be ignored.

#### Scenario: Parseo de tarea pendiente
- **WHEN** the line is `- [ ] Inicializar módulo Go`
- **THEN** a task item is produced with `done=false` and `text="Inicializar módulo Go"`

#### Scenario: Parseo de tarea completada
- **WHEN** the line is `- [x] Inicializar módulo Go`
- **THEN** a task item is produced with `done=true`

#### Scenario: Parseo de sección
- **WHEN** the line is `## 1. Setup`
- **THEN** a section item is produced with `text="1. Setup"`, non-interactive

#### Scenario: Línea ignorada
- **WHEN** the line is a free-text paragraph
- **THEN** no item is produced

### Requirement: Navegar entre tasks con j/k
In the `tasks` tab the cursor SHALL move between task items (not sections) with `j` (down) and `k` (up). The cursor SHALL skip sections automatically.

#### Scenario: Saltar sección al bajar
- **WHEN** the cursor is on the last task of section 1 and the user presses `j`
- **THEN** the cursor moves to the first task of section 2, skipping the section header

#### Scenario: Límite inferior
- **WHEN** the cursor is on the last task and the user presses `j`
- **THEN** the cursor does not move

### Requirement: Toggle de checkbox con Space
When pressing `Space` on a task item, the system SHALL:
1. Invert the `done` state of the item in memory
2. Modify only the corresponding line in `tasks.md` on disk, changing `[ ]` to `[x]` or vice versa
3. Update the render without reloading the entire file

#### Scenario: Marcar tarea como completada
- **WHEN** the cursor is on `- [ ] Crear estructura` and the user presses `Space`
- **THEN** the line on disk becomes `- [x] Crear estructura` and the item shows the completed state

#### Scenario: Desmarcar tarea completada
- **WHEN** the cursor is on `- [x] Crear estructura` and the user presses `Space`
- **THEN** the line on disk becomes `- [ ] Crear estructura` and the item shows the pending state

#### Scenario: Error de escritura
- **WHEN** `tasks.md` does not have write permissions and the user presses `Space`
- **THEN** the toggle is not applied and the TUI shows a temporary error message

### Requirement: Barra de progreso por sección
The TUI SHALL show a progress bar next to each section header along with the `completed/total` count of tasks in that section.

#### Scenario: Sección parcialmente completada
- **WHEN** a section has 2 completed tasks out of 5
- **THEN** `██░░░ 2/5` is shown next to the section header

#### Scenario: Sección completa
- **WHEN** all tasks in a section are completed
- **THEN** the bar appears completely filled

### Requirement: Indicador visual del cursor
The task item under the cursor SHALL be shown with a distinct visual indicator (e.g. `▶`) to differentiate it from the rest.

#### Scenario: Cursor sobre tarea
- **WHEN** the cursor is on task N
- **THEN** that task shows the `▶` prefix and a differentiated visual style

### Requirement: Restaurar cursor por texto tras reload
When `tasks.md` is reloaded from disk, the system SHALL attempt to restore the cursor to the task whose text matches the text of the task that had the cursor before the reload. If the text is not found in the new list, the cursor SHALL be positioned on the first available task item.

#### Scenario: Tarea bajo el cursor sigue existiendo tras reload
- **WHEN** the cursor was on the task with text `"1.3 Crear estructura"` and the reload does not remove that task
- **THEN** the cursor is positioned on the same task `"1.3 Crear estructura"`

#### Scenario: Tarea bajo el cursor eliminada en el reload
- **WHEN** the cursor was on a task that no longer exists in the new `tasks.md`
- **THEN** the cursor is positioned on the first available task item in the new list

### Requirement: Renderizado de markdown inline en items de tarea
The TUI SHALL convert inline markdown marks present in the text of each task to ANSI styles before rendering the item with lipgloss. The supported patterns are `` `code` `` (backtick) and `**bold**` (double asterisk).

#### Scenario: Tarea con fragmento de código
- **WHEN** the text of a task item contains `` `func main()` ``
- **THEN** the fragment is shown with the visual code style (distinct background or color) in the TUI

#### Scenario: Tarea con texto en negrita
- **WHEN** the text of a task item contains `**importante**`
- **THEN** the word is shown in bold in the TUI

#### Scenario: Múltiples fragmentos en la misma tarea
- **WHEN** the text of an item contains several `` `code` `` or `**bold**` fragments separated from each other
- **THEN** each fragment is rendered with its corresponding style independently

#### Scenario: Tarea sin markdown inline
- **WHEN** the text of an item contains no backticks or double asterisks
- **THEN** the text is shown unchanged, without visual artifacts

### Requirement: Viewport scroll follows cursor
When navigating the tasks view, the viewport SHALL always scroll to keep the cursor-selected task fully visible, correctly accounting for task items that wrap across multiple terminal lines.

#### Scenario: Cursor moves below visible area
- **WHEN** the user navigates down and the selected task is below the bottom of the visible viewport
- **THEN** the viewport SHALL scroll down so the selected task is visible

#### Scenario: Cursor moves above visible area
- **WHEN** the user navigates up and the selected task is above the top of the visible viewport
- **THEN** the viewport SHALL scroll up so the selected task is visible

#### Scenario: Task text wraps across multiple lines
- **WHEN** a task's text is long enough that lipgloss renders it across more than one terminal line
- **THEN** the line counter SHALL advance by the actual rendered height of that item, not by 1

#### Scenario: Tasks beyond the initial visible area are reachable
- **WHEN** the task list contains more items than fit in the visible viewport height
- **THEN** navigating down with `j` SHALL eventually reach every task in the list
