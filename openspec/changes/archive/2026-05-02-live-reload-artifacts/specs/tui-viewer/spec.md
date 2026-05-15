## ADDED Requirements

### Requirement: Periodic artifact polling
The TUI SHALL start a polling cycle every 2 seconds at startup. On each tick it SHALL compare the on-disk content of the currently visible change's artifacts with the in-memory content. The cycle SHALL continue while the TUI is active.

#### Scenario: Tick with no changes
- **WHEN** no file of the change has changed on disk
- **THEN** the TUI does not update any state or re-render anything

#### Scenario: Tick detects change in tasks.md
- **WHEN** the content of `tasks.md` on disk differs from the in-memory content
- **THEN** the TUI re-parses the tasks, restores the cursor and refreshes the view if the active tab is `tasks`

#### Scenario: Tick detects change in a markdown artifact
- **WHEN** the content of `proposal.md`, `design.md` or a `spec.md` on disk differs from the in-memory content
- **THEN** the TUI invalidates the corresponding entry in the render cache; the next time the user accesses that tab it re-renders with the new content

### Requirement: Real-time tasks view update
When the TUI detects a change in `tasks.md` and the active tab is `tasks`, it SHALL refresh the view immediately without user intervention.

#### Scenario: Agent marks task as completed
- **WHEN** an external process changes `- [ ] tarea` to `- [x] tarea` in `tasks.md`
- **THEN** within a maximum of 2 seconds the TUI shows the task as completed with the updated progress bar
