## MODIFIED Requirements

### Requirement: View artifacts of an archived change
In `ViewingArchive` mode, the TUI SHALL display the artifacts of the selected archived change using the same visual structure as active changes: header, tab bar, separator, content. Keys `1`-`4`, `j`/`k` and `h`/`l` SHALL work the same as in normal mode to navigate between artifacts and scroll the viewport. The tasks tab in archive mode SHALL render as read-only markdown (via glamour) with scroll-only navigation — no cursor, no progress bars, no toggle.

#### Scenario: Navigate artifacts of an archived change
- **WHEN** the mode is `ViewingArchive` and the user presses `2`
- **THEN** the active tab changes to `design` and the viewport shows the content of the archive's design

#### Scenario: j/k scrolls on tasks tab
- **WHEN** the mode is `ViewingArchive`, the current tab is `tasks`, and the user presses `j` or `k`
- **THEN** the viewport scrolls down or up one line (same behavior as other archive tabs)

#### Scenario: j/k scrolls on non-tasks tabs
- **WHEN** the mode is `ViewingArchive`, the current tab is `proposal`, `design`, or `specs`, and the user presses `j` or `k`
- **THEN** the viewport scrolls down or up one line

#### Scenario: h/l does not change the change in archive mode
- **WHEN** the mode is `ViewingArchive` and the user presses `h` or `l`
- **THEN** nothing changes (there is no lateral navigation between archived items)

### Requirement: Helpbar adaptado en modo archivo
In `ViewingArchive` mode, the helpbar SHALL show the actual available keys, omitting `e` and `Space`, and including `Esc: index`. All tabs SHALL show "j/k: scroll" since no tab has cursor navigation in archive mode.

#### Scenario: Helpbar on any tab in archive mode
- **WHEN** the mode is `ViewingArchive` and the current tab is `proposal`, `design`, `specs`, or `tasks`
- **THEN** the helpbar shows `1-4/Tab: artifact  j/k: scroll  Esc: index  q: quit`
