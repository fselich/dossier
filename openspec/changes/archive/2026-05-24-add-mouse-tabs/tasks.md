## 1. Enable mouse capture

- [x] 1.1 Add `tea.WithMouseCellMotion()` to `tea.NewProgram` in `cmd/dossier/main.go`

## 2. Mouse event handler

- [x] 2.1 Create `internal/ui/mouse.go` with `handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)` method
- [x] 2.2 Implement wheel scrolling: forward `MouseButtonWheelUp`/`MouseButtonWheelDown` to `m.vp.LineUp(3)` / `m.vp.LineDown(3)`
- [x] 2.3 Implement tab click mapping: calculate X ranges for each tab label, check `tabAvailable`, trigger same logic as keyboard `1`/`2`/`3`/`4`
- [x] 2.4 Gate tab click handling on `ModeNormal` and `ModeViewingArchive` (tabs don't exist in other modes)
- [x] 2.5 Filter to `MouseActionPress` + `MouseButtonLeft` for click handling

## 3. Wire into update loop

- [x] 3.1 Add `case tea.MouseMsg:` in `internal/ui/update.go` that delegates to `handleMouse`

## 4. Verification

- [x] 4.1 Manual test: wheel scrolls content in all modes (normal, index, tasks, spec viewer, config viewer)
- [x] 4.2 Manual test: clicking tabs switches artifacts in normal mode
- [x] 4.3 Manual test: clicking disabled tabs does nothing
- [x] 4.4 Manual test: clicking between tabs does nothing
- [x] 4.5 Manual test: wheel scrolls viewport without moving cursor in index mode
- [x] 4.6 Manual test: wheel scrolls viewport without moving task cursor in tasks tab
- [x] 4.7 Manual test: clicking tabs works in archive viewer mode (`ModeViewingArchive`)
- [x] 4.8 Manual test: clicking outside tab bar area does nothing
