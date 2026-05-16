## 1. Style

- [x] 1.1 Add `progressCompleteStyle` variable to `internal/ui/styles.go` with `Foreground(lipgloss.Color("14"))`

## 2. Render Sites

- [x] 2.1 Update `internal/ui/view.go` general progress bar: use `progressCompleteStyle` when `done == total`, otherwise `progressDoneStyle`
- [x] 2.2 Update `internal/ui/index.go` per-change progress bar: same conditional
- [x] 2.3 Update `internal/ui/tasks.go` per-section progress bar: same conditional

## 3. Verification

- [x] 3.1 Run `go build ./...` to confirm no compilation errors
- [x] 3.2 Manually verify: open a change with all tasks done and confirm the bar renders cyan
- [x] 3.3 Manually verify: open a change with partial tasks and confirm the bar remains green
