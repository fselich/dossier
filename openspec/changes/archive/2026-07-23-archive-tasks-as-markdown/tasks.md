## 1. Core Viewport Dispatch

- [x] 1.1 In `internal/ui/viewport.go:23`, remove `ModeViewingArchive` from the tasks case so it becomes `case m.tab == TabTasks && m.mode == ModeNormal:` — archive mode tasks fall through to `loadViewportForArtifact()`

## 2. Viewer Key Guards

- [x] 2.1 In `internal/ui/viewer.go`, guard the `j`/`down` → `TabTasks` branch with `if m.mode == ModeNormal` so archive mode falls through to default scroll
- [x] 2.2 In `internal/ui/viewer.go`, guard the `k`/`up` → `TabTasks` branch with `if m.mode == ModeNormal` so archive mode falls through to default scroll

## 3. Help Bar Unification

- [x] 3.1 In `internal/ui/view.go`, remove the `if m.tab == TabTasks` special case inside the `ModeViewingArchive` help bar block — all archive tabs now show "j/k: scroll"

## 4. Cleanup: Remove Dead `loadTaskItems()` Calls

- [x] 4.1 In `internal/ui/index.go`, remove `m.loadTaskItems()` call added by the previous fix (near archive entry point at line ~762)
- [x] 4.2 In `internal/ui/mouse.go`, remove `m.loadTaskItems()` call added by the previous fix (near archive mouse click handler at line ~152)

## 5. Tests

- [x] 5.1 Update archive+tasks test scenarios in `internal/ui/view_test.go` to expect scroll behavior, not cursor navigation
- [x] 5.2 Add test that archive mode tasks tab uses glamour rendering path (verify `renderCache` hit or glamour output)
- [x] 5.3 Run `make test` and verify all tests pass with race detector

## 6. Verify

- [x] 6.1 Run `make build` to ensure compilation
- [x] 6.2 Run `make lint` to ensure no regressions
- [x] 6.3 Manual smoke test: open an archived change, navigate to tasks tab, press j/k, verify scroll (not cursor movement), verify proposal/design/specs tabs still work in archive mode
