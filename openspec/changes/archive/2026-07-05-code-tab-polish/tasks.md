## 1. Rename tab label

- [x] 1.1 Change `tabLabels[TabGit]` from `"changes"` to `"code"` in `model.go`

## 2. Hide tab in archive mode

- [x] 2.1 In `renderTabBar` in `view.go`, skip TabGit when `m.mode != ModeNormal`
- [x] 2.2 Fix label in `renderTabBar` and mouse click handler to use `"code"` (was `"changes"`)

## 3. Verify

- [x] 3.3 `make build` succeeds
- [x] 3.4 `make lint` passes
- [x] 3.5 `make test` passes
- [x] 3.1 Tab bar shows `code` instead of `changes` in normal mode
- [x] 3.2 Tab bar does NOT show the code tab in archive mode
