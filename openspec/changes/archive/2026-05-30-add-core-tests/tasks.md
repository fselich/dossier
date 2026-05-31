## 1. Core function tests

- [x] 1.1 Add `Update()` keypress smoke tests covering enter, h/l, j/k, space, 1-4, tab, esc, a, s, e, i, q in their primary modes
- [x] 1.2 Add `doToggle()` tests verifying disk writes for task completion and nil returns for edge cases
- [x] 1.3 Add `loadViewport()` tests verifying correct content loading for ModeIndex, ModeViewingConfig, ModeViewingSpec, TabTasks, cache hit, and vpReady=false
- [x] 1.4 Add `handleTick()` tests verifying disk change detection, ModeViewingSpec no-op, and task content reload
- [x] 1.5 Add `renderTabBar()` tests verifying active/inactive/disabled tab rendering and progress bar display

## 2. Verify

- [x] 1.6 Run `go test ./internal/...` and verify all tests pass
