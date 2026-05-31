# mouse-navigation Specification (Delta)

## MODIFIED Requirements

### Requirement: Wheel scrolling

The TUI SHALL handle mouse wheel events (up and down). In `ModeIndex`, wheel events SHALL move the index cursor up or down (one item per tick) and the viewport SHALL auto-follow the cursor to keep it visible. In `ModeNormal` with `TabTasks` active, wheel events SHALL move the task cursor up or down (one task per tick) and the viewport SHALL auto-follow the cursor. In all other modes and views, wheel events SHALL scroll the viewport by 3 lines per tick.

*Implementation note: Index cursor now accessed via `m.index.Cursor` (formerly `m.indexCursor`). Task cursor now accessed via `m.tasks.Cursor` (formerly `m.taskCursor`). Task items now accessed via `m.tasks.Items` (formerly `m.taskItems`). Behavior unchanged.*

#### Scenario: Wheel in index mode moves cursor

- **WHEN** the mode is `ModeIndex` and the user scrolls the mouse wheel
- **THEN** the index cursor moves up or down by one item per wheel tick and the viewport auto-follows to keep the cursor visible

#### Scenario: Wheel in tasks tab moves cursor

- **WHEN** the `tasks` tab is active in `ModeNormal` and the user scrolls the mouse wheel
- **THEN** the task cursor moves up or down by one task per wheel tick and the viewport auto-follows to keep the cursor visible

### Requirement: Index item selection via left-click

In `ModeIndex`, the TUI SHALL handle left-click on items rendered inside the viewport. Clicking SHALL move the cursor to the clicked item. If the clicked item is already under the cursor, the TUI SHALL perform the primary action for that item kind (Enter for active changes, archived changes, and requirements; Space for expanding/collapsing specs). Clicks on section headers, blank lines, or empty areas SHALL be ignored. The coordinate mapping SHALL account for the viewport's YOffset to correctly translate screen coordinates to content lines.

*Implementation note: Index items and cursor now accessed via `m.index.Items` and `m.index.Cursor`. Expanded specs now accessed via `m.index.ExpandedSpecs`. Behavior unchanged.*

#### Scenario: Click on unselected index item moves cursor

- **WHEN** the cursor is on item A and the user left-clicks on item B in the index
- **THEN** the cursor moves to item B and no navigation (opening/expand) occurs

#### Scenario: Click on selected active change opens it

- **WHEN** the cursor is on an active change item and the user left-clicks on it
- **THEN** the mode switches to `ModeNormal` showing that change (same as pressing Enter)

#### Scenario: Click on selected requirement opens spec viewer

- **WHEN** the cursor is on a requirement item and the user left-clicks on it
- **THEN** the mode switches to `ModeViewingSpec` focused on that requirement (same as pressing Enter)

#### Scenario: Click works correctly with scrolled viewport

- **WHEN** the viewport is scrolled down (YOffset > 0) and the user left-clicks on a visible item
- **THEN** the correct item is selected regardless of the scroll offset

### Requirement: Header click navigates to index

In `ModeNormal` and `ModeViewingArchive`, the TUI SHALL enter `ModeIndex` when the user left-clicks on the header row (screen Y=1). In all other modes, clicking the header row SHALL be ignored. The behavior SHALL be the same as pressing the `a` or `Esc` key.

*Implementation note: Archive cursor access changed from `m.archiveCursor` to `m.index.ArchiveCursor`. Behavior unchanged.*

#### Scenario: Click header in normal mode enters index

- **WHEN** the mode is `ModeNormal` and the user left-clicks at Y=1
- **THEN** the TUI enters `ModeIndex`

#### Scenario: Click header in archive view enters index

- **WHEN** the mode is `ModeViewingArchive` and the user left-clicks at Y=1
- **THEN** the TUI enters `ModeIndex`
