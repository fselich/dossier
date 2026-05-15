## MODIFIED Requirements

### Requirement: Periodic artifact polling
The TUI SHALL start a polling cycle every 500 ms on startup. On each tick it SHALL compare the on-disk content of the artifacts of the currently visible change with the in-memory content, AND detect changes in artifact presence (absent → present). If at the time of the tick `len(m.project.Changes) == 0`, the tick SHALL attempt to reload the change list from disk and adopt the new state if at least one change is available. The cycle SHALL continue while the TUI is active.

#### Scenario: Tick with no changes
- **WHEN** no file of the change has changed on disk
- **THEN** the TUI does not update any state or re-render anything

#### Scenario: Tick detects change in tasks.md
- **WHEN** the content of `tasks.md` on disk differs from the in-memory content
- **THEN** the TUI re-parses the tasks, restores the cursor and refreshes the view if the active tab is `tasks`

#### Scenario: Tick detects change in a markdown artifact
- **WHEN** the content of `proposal.md`, `design.md` or a `spec.md` on disk differs from the in-memory content
- **THEN** the TUI invalidates the corresponding entry in the render cache; the next time the user accesses that tab it is re-rendered with the new content

#### Scenario: Tick detects appearance of an absent artifact
- **WHEN** an artifact that did not exist on the previous tick now exists on disk
- **THEN** the TUI updates the artifact's presence state and enables the corresponding tab

#### Scenario: TUI starts with no active changes and one is created
- **WHEN** the TUI starts with `len(m.project.Changes) == 0` and during the session a change is created on disk
- **THEN** within a maximum of 500 ms the TUI reloads the change list and displays the new change

### Requirement: Artifact tabs
The TUI SHALL show a tab bar with tabs `proposal`, `design`, `specs`, `tasks`. Tabs for absent artifacts SHALL be displayed as visually disabled and not be selectable. The user SHALL be able to switch tabs with keys `1`, `2`, `3`, `4`. If an absent artifact appears on disk during the session, the corresponding tab SHALL be enabled without needing to restart the TUI.

#### Scenario: Select an available tab
- **WHEN** the user presses `2` and `design.md` exists
- **THEN** the content area shows the rendered design

#### Scenario: Attempt to select a disabled tab
- **WHEN** the user presses `2` and `design.md` does not exist
- **THEN** the tab does not change and no error occurs

#### Scenario: Tab becomes enabled when artifact appears
- **WHEN** the TUI starts without `proposal.md` and an external process creates that file
- **THEN** within a maximum of 500 ms the `proposal` tab is shown as enabled and is selectable

### Requirement: Task updates visible in real time
When the TUI detects a change in `tasks.md` and the active tab is `tasks`, it SHALL refresh the view immediately without user intervention.

#### Scenario: Agent marks a task as completed
- **WHEN** an external process changes `- [ ] task` to `- [x] task` in `tasks.md`
- **THEN** within a maximum of 500 ms the TUI shows the task as completed with the updated progress bar

## ADDED Requirements

### Requirement: Immediate progress counter update after toggle
When the user toggles a task with `Space`, the progress counter in the tab bar SHALL update in the same frame, without waiting for the next polling cycle.

#### Scenario: Marking a task complete updates the tab bar
- **WHEN** the user presses `Space` on a pending task and the disk write succeeds
- **THEN** the `N/M` counter and the tab bar progress bar update immediately in the same render

#### Scenario: Unmarking a task updates the tab bar
- **WHEN** the user presses `Space` on a completed task and the disk write succeeds
- **THEN** the `N/M` counter and the tab bar progress bar decrement immediately in the same render
