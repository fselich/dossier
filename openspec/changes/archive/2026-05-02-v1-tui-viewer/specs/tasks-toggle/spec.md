## ADDED Requirements

### Requirement: Parse tasks.md into navigable items
The system SHALL parse `tasks.md` line by line and produce a flat list of items. Each line matching `- [ ] <text>` or `- [x] <text>` SHALL become a task item (pending or completed respectively). Lines matching a section heading (`## <text>`) SHALL become non-interactive section items. All other lines SHALL be ignored.

#### Scenario: Parse pending task
- **WHEN** the line is `- [ ] Initialize Go module`
- **THEN** a task item is produced with `done=false` and `text="Initialize Go module"`

#### Scenario: Parse completed task
- **WHEN** the line is `- [x] Initialize Go module`
- **THEN** a task item is produced with `done=true`

#### Scenario: Parse section
- **WHEN** the line is `## 1. Setup`
- **THEN** a section item is produced with `text="1. Setup"`, non-interactive

#### Scenario: Ignored line
- **WHEN** the line is a free-text paragraph
- **THEN** no item is produced

### Requirement: Navigate between tasks with j/k
In the `tasks` tab the cursor SHALL move between task items (not sections) with `j` (down) and `k` (up). The cursor SHALL skip sections automatically.

#### Scenario: Skip section when moving down
- **WHEN** the cursor is on the last task of section 1 and the user presses `j`
- **THEN** the cursor moves to the first task of section 2, skipping the section header

#### Scenario: Lower bound
- **WHEN** the cursor is on the last task and the user presses `j`
- **THEN** the cursor does not move

### Requirement: Checkbox toggle with Space
When pressing `Space` on a task item, the system SHALL:
1. Invert the `done` state of the item in memory
2. Modify only the corresponding line in `tasks.md` on disk, changing `[ ]` to `[x]` or vice versa
3. Update the render without reloading the entire file

#### Scenario: Mark task as completed
- **WHEN** the cursor is on `- [ ] Create structure` and the user presses `Space`
- **THEN** the line on disk becomes `- [x] Create structure` and the item shows the completed state

#### Scenario: Unmark completed task
- **WHEN** the cursor is on `- [x] Create structure` and the user presses `Space`
- **THEN** the line on disk becomes `- [ ] Create structure` and the item shows the pending state

#### Scenario: Write error
- **WHEN** `tasks.md` does not have write permissions and the user presses `Space`
- **THEN** the toggle is not applied and the TUI shows a temporary error message

### Requirement: Progress bar per section
The TUI SHALL show alongside each section header a progress bar and the `completed/total` count of tasks in that section.

#### Scenario: Section partially completed
- **WHEN** a section has 2 completed tasks out of 5
- **THEN** `██░░░ 2/5` is shown next to the section header

#### Scenario: Section complete
- **WHEN** all tasks in a section are completed
- **THEN** the bar appears completely filled

### Requirement: Visual cursor indicator
The task item under the cursor SHALL be shown with a distinct visual indicator (e.g. `▶`) to differentiate it from the rest.

#### Scenario: Cursor on task
- **WHEN** the cursor is on task N
- **THEN** that task shows the `▶` prefix and a differentiated visual style
