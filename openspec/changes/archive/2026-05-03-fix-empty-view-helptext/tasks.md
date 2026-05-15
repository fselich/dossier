## 1. Fix the empty view help text

- [x] 1.1 In `emptyView()` in `internal/ui/model.go`, replace `helpStyle.Render("\n  q: salir")` with `helpStyle.Render("\n  a/Esc: index  q: quit")`
