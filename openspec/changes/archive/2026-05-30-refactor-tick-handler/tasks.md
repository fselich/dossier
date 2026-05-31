## 1. Extract pollIndexMode

- [x] 1.1 Extract `pollIndexMode()` from `handleTick()` lines 18-74: disk polling, name comparison, reload logic, task refresh for index mode
- [x] 1.2 Verify the method returns `tea.Cmd` with same nil/early-return semantics

## 2. Extract pollNormalModeChanges

- [x] 2.1 Extract `pollNormalModeChanges()` from `handleTick()` lines 76-104: change-list detection, reload project, preserve current change
- [x] 2.2 Verify `m.singlePath` guard is preserved and early returns are exact

## 3. Extract pollNormalModeContent

- [x] 3.1 Extract `pollNormalModeContent()` from `handleTick()` lines 106-154: reload tasks/proposal/design/specs per-artifact, viewport dirty flag
- [x] 3.2 Verify `viewportDirty` logic and cache invalidation are preserved

## 4. Update handleTick dispatcher

- [x] 4.1 Rewrite `handleTick()` to guard and delegate to the three extracted methods
- [x] 4.2 Verify compilation and that all call sites (`Update()` tickMsg case) are unchanged

## 5. Verify tests pass

- [x] 5.1 Run `go test ./internal/ui/...` and confirm all tests pass
- [x] 5.2 Run `go build ./...` to confirm no compilation errors
