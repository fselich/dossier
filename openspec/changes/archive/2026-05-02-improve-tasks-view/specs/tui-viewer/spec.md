## MODIFIED Requirements

### Requirement: TUI layout
The TUI SHALL divide the screen into three fixed zones: header (1 line), tab bar (1 line) and content area (the rest). The header SHALL show `<project> · <change-name> [N/M]` where N is the current change's position and M the total number of active changes.

#### Scenario: Single active change
- **WHEN** only one active change exists
- **THEN** the header shows `my-project · feat-a [1/1]`

#### Scenario: Multiple active changes
- **WHEN** three active changes exist and the second one is selected
- **THEN** the header shows `my-project · feat-b [2/3]`

### Requirement: Word wrap on all task items
In the `tasks` tab, all task items SHALL word-wrap to the width of the content area (`m.width - 2`), regardless of whether the item is under the cursor or not.

#### Scenario: Long item without cursor
- **WHEN** a task item has more characters than the terminal width and the cursor is not on it
- **THEN** the text word-wraps and is shown in full across multiple lines

#### Scenario: Long item with cursor
- **WHEN** a task item has more characters than the terminal width and the cursor is on it
- **THEN** the text word-wraps and is shown in full across multiple lines with the cursor style

### Requirement: Global progress bar in the tasks view
The TUI SHALL show a global progress bar as the first content line of the `tasks` tab, before any section. The bar SHALL reflect the total completed tasks over the total tasks in the change.

#### Scenario: Change with partially completed tasks
- **WHEN** a change has 3 completed tasks out of 8 total
- **THEN** the first line of the tasks view shows a progress bar with `3/8` and a proportional fraction of filled blocks

#### Scenario: Change with all tasks completed
- **WHEN** all tasks in the change are marked as completed
- **THEN** the progress bar appears fully filled and shows the total as `N/N`

#### Scenario: Change with no tasks
- **WHEN** the change has no task-type items
- **THEN** no global progress bar is shown
