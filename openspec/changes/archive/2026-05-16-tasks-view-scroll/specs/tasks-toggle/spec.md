## ADDED Requirements

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
