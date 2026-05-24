## 1. Core Logic

- [x] 1.1 Add `nextAvailableTab(current Tab, delta int) Tab` method to Model in `internal/ui/model.go` that finds the next/previous available tab, skipping disabled ones and wrapping around, with a guard against infinite loops
- [x] 1.2 Integrate `Tab` (`"tab"`) and `Shift+Tab` (`"shift+tab"`) key handlers in `internal/ui/update.go` under `ModeNormal` and `ModeViewingArchive` cases, calling `nextAvailableTab`
- [x] 1.3 Update help bar strings in `internal/ui/view.go` to show `1-4/Tab` instead of `1-4`

## 2. Verification

- [x] 2.1 Run `go build ./...` to ensure compilation succeeds
- [x] 2.2 Run `go test ./...` to ensure existing tests pass
