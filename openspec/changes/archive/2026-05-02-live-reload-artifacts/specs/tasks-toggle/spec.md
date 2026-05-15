## ADDED Requirements

### Requirement: Restore cursor by text after reload
When `tasks.md` is reloaded from disk, the system SHALL attempt to restore the cursor to the task whose text matches the text of the task that had the cursor before the reload. If the text is not found in the new list, the cursor SHALL be positioned on the first available task item.

#### Scenario: Task under cursor still exists after reload
- **WHEN** the cursor was on the task with text `"1.3 Crear estructura"` and the reload does not remove that task
- **THEN** the cursor remains positioned on the same task `"1.3 Crear estructura"`

#### Scenario: Task under cursor deleted in the reload
- **WHEN** the cursor was on a task that no longer exists in the new `tasks.md`
- **THEN** the cursor is positioned on the first available task item in the new list
