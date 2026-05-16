## 1. Fix line counting

- [x] 1.1 In `internal/ui/tasks.go`, replace `line++` after `sb.WriteString(rendered + "\n")` (task item case) with `line += lipgloss.Height(rendered)`
- [x] 1.2 In `internal/ui/tasks.go`, replace the two `line++` calls in the section header case with `line += lipgloss.Height(...)` using the actual rendered section string

## 2. Verification

- [x] 2.1 Run `go build ./...` to confirm no compilation errors
- [x] 2.2 Manually verify: open a change with more tasks than fit the terminal height and confirm `j` scrolls through all of them
- [x] 2.3 Manually verify: a task with long text (wrapping to 2+ lines) does not desync scroll when navigating past it
