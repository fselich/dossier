# add-core-tests Specification

## Purpose

Comprehensive test coverage for core UI functions: Update keypress dispatch, doToggle task toggling, loadViewport content loading, handleTick polling, and renderTabBar rendering.

## Requirements

### Requirement: Update keypress smoke tests

Tests SHALL verify that each keybinding in `Update()` produces expected state transitions when invoked in its primary operating mode. Each test SHALL construct a `Model` in the relevant mode, send a `tea.KeyPressMsg`, and assert the resulting mode, tab, or cursor position.

#### Scenario: Enter on index active item opens normal mode
- **WHEN** mode is ModeIndex, cursor is on an active change item, and Enter is pressed
- **THEN** mode changes to ModeNormal and the change is loaded

#### Scenario: h/l navigate between changes in normal mode
- **WHEN** mode is ModeNormal with 2+ changes, and h or l is pressed
- **THEN** changeIdx increments or decrements (wrapping)

#### Scenario: j/k move cursor in index mode
- **WHEN** mode is ModeIndex with multiple items, and j or k is pressed
- **THEN** indexCursor increments or decrements within bounds

#### Scenario: Space toggles task in normal tasks tab
- **WHEN** mode is ModeNormal, tab is TabTasks, cursor is on an undone task, and space is pressed
- **THEN** doToggle is called (verified via disk write to the task file)

#### Scenario: 1-4 switch tabs in normal mode
- **WHEN** mode is ModeNormal and keys 1, 2, 3, or 4 are pressed
- **THEN** tab changes to the corresponding tab if available

#### Scenario: Tab/shift+tab cycle through available tabs
- **WHEN** mode is ModeNormal with multiple available tabs, and tab or shift+tab is pressed
- **THEN** tab cycles forward or backward

#### Scenario: Esc returns to index from normal mode
- **WHEN** mode is ModeNormal and Esc is pressed
- **THEN** mode changes to ModeIndex

#### Scenario: a enters index from normal mode
- **WHEN** mode is ModeNormal and a is pressed
- **THEN** mode changes to ModeIndex

#### Scenario: s toggles spec order in index mode
- **WHEN** mode is ModeIndex and s is pressed
- **THEN** specSortBySuffix toggles

#### Scenario: e opens editor in normal mode
- **WHEN** mode is ModeNormal on an available tab and e is pressed
- **THEN** an exec command is returned (EDITOR env or "vi")

#### Scenario: i opens config in index mode
- **WHEN** mode is ModeIndex and i is pressed
- **THEN** mode changes to ModeViewingConfig

#### Scenario: q quits in normal mode
- **WHEN** mode is ModeNormal and q is pressed
- **THEN** a Quit command is returned

### Requirement: doToggle tests

Tests SHALL verify that `doToggle()` writes task completion to disk correctly and returns nil for invalid states (no tasks, no current change, cursor on section).

#### Scenario: Toggling an undone task marks it done and writes to disk
- **GIVEN** a tasks.md file with an unchecked checkbox, and Model in ModeNormal with TabTasks active
- **WHEN** doToggle() is called
- **THEN** the task item's Done field is true and the file on disk contains `[x]`

#### Scenario: Toggling a done task marks it undone and writes to disk
- **GIVEN** a tasks.md file with a checked checkbox
- **WHEN** doToggle() is called
- **THEN** the task item's Done field is false and the file on disk contains `[ ]`

#### Scenario: doToggle on empty task list returns nil
- **WHEN** doToggle() is called with no task items
- **THEN** it returns nil (no error)

#### Scenario: doToggle with nil current change returns nil
- **WHEN** doToggle() is called and current() returns nil
- **THEN** it returns nil

### Requirement: loadViewport tests

Tests SHALL verify that `loadViewport()` returns the correct type of command and sets viewport content for each mode.

#### Scenario: loadViewport in ModeIndex sets index content and returns nil
- **WHEN** mode is ModeIndex and loadViewport() is called
- **THEN** viewport content is set (non-empty) and the returned command is nil

#### Scenario: loadViewport in ModeViewingConfig returns a glamour render command
- **WHEN** mode is ModeViewingConfig and loadViewport() is called
- **THEN** the returned command is non-nil (async glamour render)

#### Scenario: loadViewport in ModeViewingSpec returns a glamour render command
- **WHEN** mode is ModeViewingSpec and loadViewport() is called
- **THEN** the returned command is non-nil (async glamour render)

#### Scenario: loadViewport for TabTasks returns nil
- **WHEN** mode is ModeNormal, tab is TabTasks, and loadViewport() is called
- **THEN** viewport content is set and the returned command is nil

#### Scenario: loadViewport uses render cache
- **WHEN** mode is ModeNormal on a non-tasks tab with content already in renderCache
- **THEN** viewport content is set from cache and the returned command is nil

#### Scenario: loadViewport returns nil when vpReady is false
- **WHEN** vpReady is false and loadViewport() is called
- **THEN** the returned command is nil

### Requirement: handleTick tests

Tests SHALL verify that `handleTick()` detects changes on disk and updates the model.

#### Scenario: handleTick detects new change on disk
- **GIVEN** a root directory with one change, and a second change is added to disk
- **WHEN** handleTick() is called
- **THEN** project.Changes reflects the new change count

#### Scenario: handleTick returns nil in ModeViewingSpec
- **WHEN** mode is ModeViewingSpec and handleTick() is called
- **THEN** the returned command is nil (no polling in spec viewer)

#### Scenario: handleTick detects task content change
- **GIVEN** a tasks.md file whose content changes on disk
- **WHEN** handleTick() is called
- **THEN** the cached task content is invalidated and updated

### Requirement: renderTabBar tests

Tests SHALL verify that `renderTabBar()` produces correct output showing active, inactive, and disabled tabs, and the progress bar when tasks are present.

#### Scenario: renderTabBar shows active tab highlighted
- **WHEN** renderTabBar() is called with tab set to TabProposal
- **THEN** the output contains the active style markup for "proposal"

#### Scenario: renderTabBar shows disabled tabs dimmed
- **WHEN** a change has no design artifact and renderTabBar() is called
- **THEN** the output contains disabled style markup for "design"

#### Scenario: renderTabBar shows progress bar when tasks present
- **WHEN** task items contain at least one task and renderTabBar() is called
- **THEN** the output contains a progress bar with `N/M` format

#### Scenario: renderTabBar shows no progress bar when no tasks
- **WHEN** task items are nil or contain no KindTask items and renderTabBar() is called
- **THEN** the output does not contain a progress bar
