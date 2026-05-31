## Why

The codebase has only 521 lines of tests covering UI logic — primarily index rendering, mouse click dispatch, and `extractRequirement`. Critical functions like `Update()`, `doToggle()`, `loadViewport()`, `handleTick()`, and `renderTabBar()` have zero test coverage. This presents a risk for regressions as the codebase grows.

## What Changes

- Add `Update()` keypress smoke tests: verify each mode's main keybindings (enter, h/l, j/k, space, 1-4, tab, esc, a, s, e, i, q) produce expected state transitions
- Add `doToggle()` tests: verify task completion writes to disk correctly and invalid cases return nil
- Add `loadViewport()` tests: verify correct content loading for each mode (ModeIndex, ModeViewingConfig, ModeViewingSpec, TabTasks, artifact tabs)
- Add `handleTick()` tests: verify polling detects disk changes and reloads
- Add `renderTabBar()` tests: verify tab rendering with available/disabled tabs and progress bar calculations

## Capabilities

### New Capabilities

- `add-core-tests`: comprehensive test coverage for core UI functions (Update keypress, doToggle, loadViewport, handleTick, renderTabBar)

### Modified Capabilities

*(none)*

## Impact

- **Affected code**: new tests in `internal/ui/view_test.go` (or new `internal/ui/update_test.go`)
- **No production code changes**
- **No new dependencies** (uses existing `testing` + `os` + `filepath` packages from stdlib)
